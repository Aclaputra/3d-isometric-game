package main

import (
	"bytes"
	_ "embed"
	"fmt"
	"image"
	"image/color"
	"log"

	"habitate/assets/images/isometric/forest/tiles"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

const (
	MapWidth  = 30
	MapHeight = 30
	MapLevels = 2
)

const (
	Gravity   = 0.02 // acceleration per frame
	JumpSpeed = -0.5 // initial jump speed
)

type Camera struct {
	X, Y float64
	Zoom float64
}

type Game struct {
	tilesetImg    *ebiten.Image
	tilesetDesc   *TilesetJSON
	characterImg  *ebiten.Image
	characterDesc *CharacterJSON
	tiles         map[string]*Tile
	mapData       [MapLevels][MapHeight][MapWidth]Tile // binary map
	player        Character
	camera        Camera
	debugMode     bool
}

func NewGame() *Game {
	decodeTilesetImg, _, err := image.Decode(bytes.NewReader(tiles.Main_tiles))
	if err != nil {
		log.Fatal(err)
	}
	tilesetImg := ebiten.NewImageFromImage(decodeTilesetImg)

	characterImg, _, err := ebitenutil.NewImageFromFile("assets/images/isometric/character/female.png")
	if err != nil {
		log.Fatal(err)
	}
	tilesetDesc, err := NewTilesetJSON("assets/images/isometric/forest/specs/tiles.json")
	if err != nil {
		log.Fatal(err)
	}
	characterDesc, err := NewCharacterJSON("assets/images/isometric/character/specs.json")
	if err != nil {
		log.Fatal(err)
	}
	tiles, err := tilesetDesc.GenerateTileset(tilesetImg)
	if err != nil {
		log.Fatal(err)
	}
	character, err := characterDesc.GenerateCharacter(characterImg)
	if err != nil {
		log.Fatal(err)
	}
	character.X = 12
	character.Y = 12
	character.Z = 1
	character.Facing = "down"
	character.State = "idle"

	camera := Camera{
		X:    character.X,
		Y:    character.Y,
		Zoom: 1.0,
	}

	noise := NewNoise()
	tilemap := NewTileMap(noise, tiles)

	return &Game{
		tilesetImg:    tilesetImg,
		tilesetDesc:   tilesetDesc,
		characterImg:  characterImg,
		characterDesc: characterDesc,
		tiles:         tiles,
		mapData:       tilemap.mapData,
		player:        character,
		camera:        camera,
		debugMode:     false,
	}
}

func (g *Game) Update() error {
	fps := ebiten.CurrentFPS()
	ebiten.SetWindowTitle(fmt.Sprintf("Habitate v0.01 | FPS: %.0f", fps))

	// Update player movement
	g.updatePlayer()

	// Update camera smoothly
	g.updateCamera()

	return nil
}

func (g *Game) updatePlayer() {
	speed := 0.1
	dx, dy := 0.0, 0.0

	up := ebiten.IsKeyPressed(ebiten.KeyUp) || ebiten.IsKeyPressed(ebiten.KeyW)
	down := ebiten.IsKeyPressed(ebiten.KeyDown) || ebiten.IsKeyPressed(ebiten.KeyS)
	left := ebiten.IsKeyPressed(ebiten.KeyLeft) || ebiten.IsKeyPressed(ebiten.KeyA)
	right := ebiten.IsKeyPressed(ebiten.KeyRight) || ebiten.IsKeyPressed(ebiten.KeyD)

	// --- Movement ---
	if left {
		dx -= speed
	}
	if right {
		dx += speed
	}
	if up {
		dy -= speed
	}
	if down {
		dy += speed
	}

	// Normalize diagonals
	if dx != 0 && dy != 0 {
		dx *= 0.7071
		dy *= 0.7071
	}

	// Facing logic
	switch {
	case up && left:
		g.player.Facing = "up"
	case up && right:
		g.player.Facing = "right"
	case down && left:
		g.player.Facing = "left"
	case down && right:
		g.player.Facing = "down"
	case up:
		g.player.Facing = "up_right"
	case down:
		g.player.Facing = "down_left"
	case left:
		g.player.Facing = "up_left"
	case right:
		g.player.Facing = "down_right"
	}

	// State logic
	if dx == 0 && dy == 0 {
		g.player.State = "idle"
		g.player.AnimTimer = 0
	} else {
		g.player.State = "walk"
	}

	// Animation
	if g.player.State == "walk" {
		g.player.AnimTimer++
		if g.player.AnimTimer > 30 {
			g.player.AnimTimer = 0
			g.player.AnimFrame = (g.player.AnimFrame + 1) % 2
		}
	} else {
		g.player.AnimFrame = 0
	}

	g.player.CurrentImage = g.player.ImageDirections[g.player.Facing][g.player.State][g.player.AnimFrame]

	// -------------------------------------------------------
	//  COLLISION CHECK
	// -------------------------------------------------------

	newX := g.player.X + dx
	newY := g.player.Y + dy

	// Player box at the new location
	newBox := g.player.AABBAt(newX, newY)

	blocked := false

	playerZ := int(g.player.Z)

	for y := 0; y < MapHeight; y++ {
		for x := 0; x < MapWidth; x++ {

			tile := g.mapData[playerZ][y][x]
			if tile.CollisionBox == nil {
				continue
			}

			if AABBIntersects(newBox, *tile.CollisionBox) {
				blocked = true
				break
			}
		}
		if blocked {
			break
		}
	}

	if !blocked {
		g.player.X = newX
		g.player.Y = newY
	}
}

func (g *Game) updateCamera() {
	// Smoothly follow player
	g.camera.X += (g.player.X - g.camera.X) * 0.1
	g.camera.Y += (g.player.Y - g.camera.Y) * 0.1

	// Zoom with mouse wheel
	_, wheelY := ebiten.Wheel()
	if wheelY != 0 {
		g.camera.Zoom += wheelY * 0.1 // adjust zoom speed here
		if g.camera.Zoom < 0.2 {      // minimum zoom
			g.camera.Zoom = 0.2
		}
		if g.camera.Zoom > 3.0 { // maximum zoom
			g.camera.Zoom = 3.0
		}
	}
}

func AABBIntersects(a, b AABB) bool {
	return a.MinX < b.MaxX &&
		a.MaxX > b.MinX &&
		a.MinY < b.MaxY &&
		a.MaxY > b.MinY
}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.Fill(color.RGBA{10, 10, 30, 255})

	tileSize := g.tilesetDesc.Size
	tileWidth := float64(tileSize)
	tileHeight := float64(tileSize / 2)

	zoom := g.camera.Zoom

	// Center of the screen
	screenCenterX := float64(screen.Bounds().Dx()) / 2
	screenCenterY := float64(screen.Bounds().Dy()) / 2

	g.player.CurrentImage = g.player.ImageDirections[g.player.Facing][g.player.State][g.player.AnimFrame]

	// Draw world tiles + cyan collision boxes
	for z := 0; z < MapLevels; z++ {
		for y := 0; y < MapHeight; y++ {

			// --- draw the player at the correct depth ---
			if z == int(g.player.Z) && y == int(g.player.Y) {

				offsetUp := 48 * zoom
				px := screenCenterX
				py := screenCenterY

				op := &ebiten.DrawImageOptions{}
				op.GeoM.Scale(zoom, zoom)
				op.GeoM.Translate(
					px-float64(g.player.CurrentImage.Bounds().Dx())/2*zoom,
					py-float64(g.player.CurrentImage.Bounds().Dy())*zoom+32*zoom-offsetUp,
				)

				screen.DrawImage(g.player.CurrentImage, op)

				// ---- debug box code (unchanged) ----
				if g.debugMode {
					// (keep your debug code here)
				}
			}

			for x := 0; x < MapWidth; x++ {
				tile := g.mapData[z][y][x]
				tileImage := tile.Image
				// for dev test debugging on one block collision
				if tileImage == nil {
					continue
				}

				// World coordinates of the tile
				tileX := float64(x-y) * tileWidth / 2
				tileY := float64(x+y) * tileHeight / 2

				// Offset tile by camera (player-centered)
				screenX := (tileX-g.player.X*tileWidth/2+g.player.Y*tileWidth/2)*zoom + screenCenterX
				screenY := (tileY-g.player.X*tileHeight/2-g.player.Y*tileHeight/2)*zoom + screenCenterY - float64(z)*tileHeight*zoom

				// Draw tile
				options := &ebiten.DrawImageOptions{}
				options.GeoM.Scale(zoom, zoom)
				options.GeoM.Translate(screenX, screenY)

				screen.DrawImage(tileImage, options)
				// --- Draw isometric collision box (cyan diamond) ---
				if g.debugMode && tile.CollisionBox != nil {
					box := tile.CollisionBox

					screenShiftX := tileWidth * -0.01 * zoom
					screenShiftY := tileWidth * -0.55 * zoom

					// Convert each AABB corner into isometric screen cwoords
					convert := func(wx, wy float64) (sx, sy float64) {
						isoX := (wx - wy) * (tileWidth / 2)
						isoY := (wx + wy) * (tileHeight / 2)

						sx = (isoX-g.player.X*tileWidth/2+g.player.Y*tileWidth/2)*zoom + screenCenterX
						sy = (isoY-g.player.X*tileHeight/2-g.player.Y*tileHeight/2)*zoom + screenCenterY - float64(z)*tileHeight*zoom

						sx += screenShiftX
						sy += screenShiftY
						return
					}

					sx0, sy0 := convert(box.MinX, box.MinY)
					sx1, sy1 := convert(box.MaxX, box.MinY)
					sx2, sy2 := convert(box.MaxX, box.MaxY)
					sx3, sy3 := convert(box.MinX, box.MaxY)

					// Draw top diamond (cyan)
					drawQuadFilled(
						screen,
						sx0, sy0,
						sx1, sy1,
						sx2, sy2,
						sx3, sy3,
						color.RGBA{0, 255, 255, 200},
					)

					belowOffset := tileHeight * zoom // natural isometric vertical depth
					drawQuadFilled(
						screen,
						sx0, sy0+belowOffset,
						sx1, sy1+belowOffset,
						sx2, sy2+belowOffset,
						sx3, sy3+belowOffset,
						color.RGBA{100, 100, 255, 200},
					)
				}

			}
		}
	}
}

func drawQuadFilled(screen *ebiten.Image,
	x0, y0, x1, y1, x2, y2, x3, y3 float64,
	col color.RGBA,
) {
	// A white base image (1×1)
	img := ebiten.NewImage(1, 1)
	img.Fill(color.White)

	op := &ebiten.DrawTrianglesOptions{}

	// Convert color to float32
	r := float32(col.R) / 255
	g := float32(col.G) / 255
	b := float32(col.B) / 255
	a := float32(col.A) / 255 // opacity!

	vertices := []ebiten.Vertex{
		{DstX: float32(x0), DstY: float32(y0), SrcX: 0, SrcY: 0, ColorR: r, ColorG: g, ColorB: b, ColorA: a},
		{DstX: float32(x1), DstY: float32(y1), SrcX: 1, SrcY: 0, ColorR: r, ColorG: g, ColorB: b, ColorA: a},
		{DstX: float32(x2), DstY: float32(y2), SrcX: 1, SrcY: 1, ColorR: r, ColorG: g, ColorB: b, ColorA: a},
		{DstX: float32(x3), DstY: float32(y3), SrcX: 0, SrcY: 1, ColorR: r, ColorG: g, ColorB: b, ColorA: a},
	}

	indices := []uint16{
		0, 1, 2, // first triangle
		0, 2, 3, // second triangle
	}

	screen.DrawTriangles(vertices, indices, img, op)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return 1280, 720
}

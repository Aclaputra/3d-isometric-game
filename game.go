package main

import (
	"bytes"
	_ "embed"
	"image"
	"image/color"
	"log"

	"habitate/assets/images/isometric/forest/tiles"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

const (
	MapWidth  = 50
	MapHeight = 50
	MapLevels = 1
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
	mapData       [MapLevels][MapHeight][MapWidth]int // binary map
	player        Character
	camera        Camera
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
	character.X = 25
	character.Y = 5
	character.Facing = "down"
	character.State = "idle"

	camera := Camera{
		X:    character.X,
		Y:    character.Y,
		Zoom: 1.0,
	}

	noise := NewNoise()
	tilemap := NewTileMap(noise)

	return &Game{
		tilesetImg:    tilesetImg,
		tilesetDesc:   tilesetDesc,
		characterImg:  characterImg,
		characterDesc: characterDesc,
		tiles:         tiles,
		mapData:       tilemap.mapData,
		player:        character,
		camera:        camera,
	}
}

func (g *Game) Update() error {
	// Update player movement (if you have one)
	g.updatePlayer()

	// Update camera to follow player smoothly
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

	// --- Normalize diagonals ---
	if dx != 0 && dy != 0 {
		dx *= 0.7071
		dy *= 0.7071
	}

	// --- Determine direction ---
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

	// --- Determine state ---
	if dx == 0 && dy == 0 {
		g.player.State = "idle"
		g.player.AnimTimer = 0 // reset animation when idle
	} else {
		g.player.State = "walk"
	}

	// --- Update animation frame ---
	if g.player.State == "walk" {
		g.player.AnimTimer++
		if g.player.AnimTimer > 30 { // adjust speed (lower = faster)
			g.player.AnimTimer = 0
			g.player.AnimFrame = (g.player.AnimFrame + 1) % 2 // toggle 0–1
		}
	} else {
		g.player.AnimFrame = 0
	}

	// --- Select image ---
	g.player.CurrentImage = g.player.ImageDirections[g.player.Facing][g.player.State][g.player.AnimFrame]

	// --- Apply movement ---
	g.player.X += dx
	g.player.Y += dy
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

func (g *Game) Draw(screen *ebiten.Image) {
	screen.Fill(color.RGBA{10, 10, 30, 255})

	tileSize := g.tilesetDesc.Size
	tileWidth := float64(tileSize)
	tileHeight := float64(tileSize / 2)
	dirtTile := g.tiles["Dirt Block"].Image

	zoom := g.camera.Zoom

	// Center of the screen
	screenCenterX := float64(screen.Bounds().Dx()) / 2
	screenCenterY := float64(screen.Bounds().Dy()) / 2

	// Draw world tiles
	for z := 0; z < MapLevels; z++ {
		for y := 0; y < MapHeight; y++ {
			for x := 0; x < MapWidth; x++ {
				if g.mapData[z][y][x] == 0 {
					continue
				}

				// World coordinates of the tile
				tileX := float64(x-y) * tileWidth / 2
				tileY := float64(x+y) * tileHeight / 2

				// Offset tile by camera (player-centered)
				screenX := (tileX-g.player.X*tileWidth/2+g.player.Y*tileWidth/2)*zoom + screenCenterX
				screenY := (tileY-g.player.X*tileHeight/2-g.player.Y*tileHeight/2)*zoom + screenCenterY - float64(z)*tileHeight*zoom

				options := &ebiten.DrawImageOptions{}
				options.GeoM.Scale(zoom, zoom)
				options.GeoM.Translate(screenX, screenY)
				screen.DrawImage(dirtTile, options)
			}
		}
	}

	// Draw player in center
	g.player.CurrentImage = g.player.ImageDirections[g.player.Facing][g.player.State][g.player.AnimFrame]
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Scale(zoom, zoom)
	op.GeoM.Translate(screenCenterX-float64(g.player.CurrentImage.Bounds().Dx())/2*zoom,
		screenCenterY-float64(g.player.CurrentImage.Bounds().Dy())*zoom+32*zoom)
	screen.DrawImage(g.player.CurrentImage, op)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return 1280, 720
}

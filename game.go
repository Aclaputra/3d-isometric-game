package main

import (
	_ "embed"
	"fmt"
	"image/color"

	"habitate/assets/images/isometric/forest/tiles"
	"habitate/components"
	"habitate/debug_helper"

	"github.com/hajimehoshi/ebiten/v2"
)

const (
	MapWidth  = 30
	MapHeight = 30
	MapLevels = 2
)

type Game struct {
	tiles     map[string]*Tile
	mapData   [MapLevels][MapHeight][MapWidth]Tile
	player    Character
	camera    *components.Camera
	debugMode bool
}

func NewGame() *Game {
	resourceData := NewResourceData(&ResourcePaths{
		TilsetImage: tiles.Main_tiles,
		PlayerImage: "assets/images/isometric/character/female.png",
		TilesetDesc: "assets/images/isometric/forest/specs/tiles.json",
		PlayerDesc:  "assets/images/isometric/character/specs.json",
	})

	resourceData.Player.X = 12
	resourceData.Player.Y = 12
	resourceData.Player.Z = 1
	resourceData.Player.Facing = "down"
	resourceData.Player.State = "idle"

	camera := components.NewCamera(
		resourceData.Player.X,
		resourceData.Player.Y,
		1.0,
	)

	noise := NewNoise()
	tilemap := NewTileMap(noise, resourceData.Tiles)

	return &Game{
		tiles:     resourceData.Tiles,
		mapData:   tilemap.mapData,
		player:    *resourceData.Player,
		camera:    camera,
		debugMode: false,
	}
}

func (g *Game) Update() error {
	fps := ebiten.CurrentFPS()
	ebiten.SetWindowTitle(fmt.Sprintf("Habitate v0.01 | FPS: %.0f", fps))

	// Update player movement
	g.player.Update(g.mapData)

	// Update camera smoothly
	g.camera.Update(g.player.X, g.player.Y)

	return nil
}

func AABBIntersects(a, b AABB) bool {
	return a.MinX < b.MaxX &&
		a.MaxX > b.MinX &&
		a.MinY < b.MaxY &&
		a.MaxY > b.MinY
}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.Fill(color.RGBA{10, 10, 30, 255})

	tileSize := TILE_SIZE
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
					debug_helper.DrawQuadFilled(
						screen,
						sx0, sy0,
						sx1, sy1,
						sx2, sy2,
						sx3, sy3,
						color.RGBA{0, 255, 255, 200},
					)

					belowOffset := tileHeight * zoom // natural isometric vertical depth
					debug_helper.DrawQuadFilled(
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

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return SCREEN_WIDTH, SCREEN_HEIGHT
}

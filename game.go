package main

import (
	_ "embed"
	"image/color"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

const (
	MapWidth  = 50
	MapHeight = 50
	MapLevels = 5
)

type Game struct {
	tilesetImg    *ebiten.Image
	tilesetDesc   *TilesetJSON
	characterImg  *ebiten.Image
	characterDesc *CharacterJSON
	tiles         map[string]*Tile
	mapData       [MapLevels][MapHeight][MapWidth]int // binary map
}

func NewGame() *Game {
	tilesetImg, _, err := ebitenutil.NewImageFromFile("assets/images/isometric/forest/tiles/main.png")
	if err != nil {
		log.Fatal(err)
	}
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

	noise := NewNoise()
	tilemap := NewTileMap(noise)

	return &Game{
		tilesetImg:    tilesetImg,
		tilesetDesc:   tilesetDesc,
		characterImg:  characterImg,
		characterDesc: characterDesc,
		tiles:         tiles,
		mapData:       tilemap.mapData,
	}
}

func (g *Game) Update() error {

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.Fill(color.RGBA{10, 10, 30, 255})

	tileSize := g.tilesetDesc.Size
	tileWidth := float64(tileSize)
	tileHeight := float64(tileSize / 2)
	dirtTile := g.tiles["Dirt Block"].Image

	for z := 0; z < MapLevels; z++ {
		for y := 0; y < MapHeight; y++ {
			for x := 0; x < MapWidth; x++ {
				if g.mapData[z][y][x] == 0 {
					continue
				}

				// Convert to isometric screen coordinates
				screenX := (float64(x-y) * tileWidth / 2)
				screenY := (float64(x+y) * tileHeight / 2)

				options := &ebiten.DrawImageOptions{}
				// Offset horizontally to center the map
				options.GeoM.Translate(screenX+(1280/2), screenY-float64(z)*tileHeight)
				screen.DrawImage(dirtTile, options)
			}
		}
	}

	// // TODO: create a base map into a binary
	// for y := 0; y < 50; y++ {
	// 	for x := 0; x < 50; x++ {
	// 		// Get noise value
	// 		nv := noise.Perlin.Noise2D(float64(x)/10, float64(y)/10)
	// 		// Normalize Perlin value (-1..1 → 0..1)
	// 		nv = (nv + 1) / 2

	// 		// Set level threshold
	// 		level := 0
	// 		if nv > 0.6 {
	// 			level = 1 // raised tile
	// 		}

	// 		screenX := (float64(x-y) * tileWidth / 2)
	// 		screenY := (float64(x+y) * tileHeight / 2)
	// 		options := ebiten.DrawImageOptions{}
	// 		options.GeoM.Translate(screenX+(1280/2), screenY)
	// 		screen.DrawImage(dirtTile, &options)

	// 		if level == 1 {
	// 			// second level
	// 			options2 := ebiten.DrawImageOptions{}
	// 			options2.GeoM.Translate(screenX+(1280/2), screenY-float64(tileHeight))
	// 			screen.DrawImage(dirtTile, &options2)
	// 		}
	// 	}
	// }
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return 1280, 720
}

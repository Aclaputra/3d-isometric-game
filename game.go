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
	MapLevels = 5
)

type Game struct {
	tilesetImg    *ebiten.Image
	tilesetDesc   *TilesetJSON
	characterImg  *ebiten.Image
	characterDesc *CharacterJSON
	tiles         map[string]*Tile
	mapData       [MapLevels][MapHeight][MapWidth]int // binary map
	player        Character
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

	// draw player
	playerScreenX := (float64(5-5) * tileWidth / 2)
	playerScreenY := (float64(5+5) * tileHeight / 2)

	op := &ebiten.DrawImageOptions{}
	playerImage := g.player.ImageDirections["down"]["idle"][0]
	op.GeoM.Translate(playerScreenX+(1280/2), playerScreenY-float64(playerImage.Bounds().Dy())+32)
	screen.DrawImage(playerImage, op)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return 1280, 720
}

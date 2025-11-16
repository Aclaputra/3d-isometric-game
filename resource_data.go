package main

import (
	"bytes"
	"habitate/assets/images/isometric/forest/tiles"
	"image"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

type ResourceData struct {
	Tiles  map[string]*Tile
	Player *Character
}

func NewResourceData(path *ResourcePaths) *ResourceData {
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
	return &ResourceData{
		Tiles:  tiles,
		Player: &character,
	}
}

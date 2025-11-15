package main

import (
	"encoding/json"
	"image"
	"log"
	"os"

	"github.com/hajimehoshi/ebiten/v2"
)

type Axis struct {
	X int `json:"x"`
	Y int `json:"y"`
}

type Dimensions struct {
	Width  int `json:"width"`
	Height int `json:"height"`
}

// Its a tileset because we create the map custom on the game
type TilesetJSON struct {
	Description string           `json:"description"`
	Size        int              `json:"size"`
	TilesetDesc []map[string]any `json:"tileset_descs"`
	ObjectDesc  []map[string]any `json:"object_descs"`
	ItemDesc    []map[string]any `json:"item_descs"`
}

type TilesetDetail struct {
	Name string
	X    int
	Y    int
}

func NewTilesetJSON(filepath string) (*TilesetJSON, error) {
	contents, err := os.ReadFile(filepath)
	if err != nil {
		return nil, err
	}

	var tileObjectJson TilesetJSON
	err = json.Unmarshal(contents, &tileObjectJson)
	if err != nil {
		return nil, err
	}

	return &tileObjectJson, nil
}

func (tj *TilesetJSON) GenerateTileset(image *ebiten.Image) (tiles map[string]*Tile, err error) {
	type PerspectiveDetail struct {
		Ptype any `json:"type"`
		X     any `json:"x"`
		Y     any `json:"y"`
	}

	var tilesets []TilesetDetail

	for _, desc := range tj.TilesetDesc {
		name := desc["name"]
		var perspectives []PerspectiveDetail

		arr, perspectiveExist := desc["perspective"].([]any)
		if perspectiveExist {
			for _, v := range arr {
				var perspective *PerspectiveDetail
				if m, ok := v.(map[string]any); ok {
					perspective = &PerspectiveDetail{
						Ptype: m["type"],
						X:     m["x"],
						Y:     m["y"],
					}
				}
				perspectives = append(perspectives, *perspective)
			}
		}

		if len(perspectives) > 0 {
			prefixName := name.(string)
			for _, perspective := range perspectives {
				tileset := &TilesetDetail{
					Name: prefixName + " - " + perspective.Ptype.(string),
					X:    int(perspective.X.(float64)),
					Y:    int(perspective.Y.(float64)),
				}
				tilesets = append(tilesets, *tileset)
			}
		} else {
			tileset := &TilesetDetail{
				Name: name.(string),
				X:    int(desc["x"].(float64)),
				Y:    int(desc["y"].(float64)),
			}
			tilesets = append(tilesets, *tileset)
		}
	}

	tiles, err = tj.mappingTilesetImage(image, &tilesets)
	if err != nil {
		log.Fatal(err)
	}

	return
}

type Tile struct {
	Name         string
	IsWall       bool // wall collision like
	CollisionBox *AABB
	Image        *ebiten.Image
}

type AABB struct {
	MinX, MinY float64
	MaxX, MaxY float64
}

func (tj *TilesetJSON) mappingTilesetImage(tilesetImage *ebiten.Image, tilesetDetails *[]TilesetDetail) (tiles map[string]*Tile, err error) {
	tileSize := tj.Size
	tileXCount := tilesetImage.Bounds().Dx() / tileSize
	tiles = make(map[string]*Tile, 0)
	for _, desc := range *tilesetDetails {
		tileIndex := desc.Y*tileXCount + desc.X
		sourceX := (tileIndex % tileXCount) * tileSize
		sourceY := (tileIndex / tileXCount) * tileSize
		tileImage := tilesetImage.SubImage(image.Rect(sourceX, sourceY, sourceX+tileSize, sourceY+tileSize)).(*ebiten.Image)

		tiles[desc.Name] = &Tile{
			Name:  desc.Name,
			Image: tileImage,
		}
	}

	return
}

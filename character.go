package main

import (
	"encoding/json"
	"image"
	"os"

	"github.com/hajimehoshi/ebiten/v2"
)

type DetailDirection struct {
	Idle []Axis `json:"idle"`
	Walk []Axis `json:"walk"`
}

type Direction struct {
	Name   string          `json:"name"`
	Detail DetailDirection `json:"detail"`
}

type CharacterJSON struct {
	Description string      `json:"description"`
	Size        int         `json:"size"`
	Directions  []Direction `json:"directions"`
}

func NewCharacterJSON(filepath string) (*CharacterJSON, error) {
	contents, err := os.ReadFile(filepath)
	if err != nil {
		return nil, err
	}

	var characterJson CharacterJSON
	err = json.Unmarshal(contents, &characterJson)
	if err != nil {
		return nil, err
	}

	return &characterJson, nil
}

type Character struct {
	Name            string
	ImageDirections map[string]map[string][]*ebiten.Image
}

func (cj *CharacterJSON) GenerateCharacter(image *ebiten.Image) (character Character, err error) {
	character, err = cj.mappingCharacterImage(image)
	if err != nil {
		return
	}
	character.Name = "Female Character"

	return
}

func (cj *CharacterJSON) mappingCharacterImage(characterImage *ebiten.Image) (character Character, err error) {
	characterSize := cj.Size
	characterXcount := characterImage.Bounds().Dx() / characterSize
	imageDirections := make(map[string]map[string][]*ebiten.Image)

	for _, direction := range cj.Directions {
		imageDirection := make(map[string][]*ebiten.Image)

		imageDirection["idle"] = make([]*ebiten.Image, 0)
		for _, idle := range direction.Detail.Idle {
			characterIndex := idle.Y*characterXcount + idle.X
			sourceX := (characterIndex % characterXcount) * characterSize
			sourceY := (characterIndex / characterXcount) * characterSize
			frameDirectionImage := characterImage.SubImage(image.Rect(sourceX, sourceY, sourceX+characterSize, sourceY+characterSize)).(*ebiten.Image)
			imageDirection["idle"] = append(imageDirection["idle"], frameDirectionImage)
		}

		imageDirection["walk"] = make([]*ebiten.Image, 0)
		for _, walk := range direction.Detail.Walk {
			characterIndex := walk.Y*characterXcount + walk.X
			sourceX := (characterIndex % characterXcount) * characterSize
			sourceY := (characterIndex / characterXcount) * characterSize
			frameDirectionImage := characterImage.SubImage(image.Rect(sourceX, sourceY, sourceX+characterSize, sourceY+characterSize)).(*ebiten.Image)
			imageDirection["walk"] = append(imageDirection["walk"], frameDirectionImage)
		}

		imageDirections[direction.Name] = imageDirection
	}

	character.ImageDirections = imageDirections
	return
}

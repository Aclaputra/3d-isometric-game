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

func (cj *CharacterJSON) GenerateCharacter(image *ebiten.Image) (character Character, err error) {
	character, err = cj.mappingCharacterImage(image)
	if err != nil {
		return
	}
	character.Name = "Female Character"

	return
}

type Character struct {
	X, Y            float64
	Z               float64 // vertical height for gravity
	Name            string
	VelocityZ       float64 // vertical speed
	OnGround        bool    // is player on the ground
	ImageDirections map[string]map[string][]*ebiten.Image
	CurrentImage    *ebiten.Image
	Facing          string
	State           string
	AnimFrame       int
	AnimTimer       int
	FootCollisonBox CharacterCollision
}

type CharacterCollision struct {
	// Player collision box relative to player center
	BoundsMinX float64
	BoundsMinY float64
	BoundsMaxX float64
	BoundsMaxY float64
}

func (p *Character) AABBAt(x, y float64) AABB {
	feetW := 0.30 // same width used in update
	feetH := 0.30 // same height used in update

	return AABB{
		MinX: x - feetW,
		MaxX: x + feetW,
		MinY: y - feetH,
		MaxY: y + feetH,
	}
}

func (p *Character) Update(mapData [MapLevels][MapHeight][MapWidth]Tile) {
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
		p.Facing = "up"
	case up && right:
		p.Facing = "right"
	case down && left:
		p.Facing = "left"
	case down && right:
		p.Facing = "down"
	case up:
		p.Facing = "up_right"
	case down:
		p.Facing = "down_left"
	case left:
		p.Facing = "up_left"
	case right:
		p.Facing = "down_right"
	}

	// State logic
	if dx == 0 && dy == 0 {
		p.State = "idle"
		p.AnimTimer = 0
	} else {
		p.State = "walk"
	}

	// Animation
	if p.State == "walk" {
		p.AnimTimer++
		if p.AnimTimer > 30 {
			p.AnimTimer = 0
			p.AnimFrame = (p.AnimFrame + 1) % 2
		}
	} else {
		p.AnimFrame = 0
	}

	p.CurrentImage = p.ImageDirections[p.Facing][p.State][p.AnimFrame]

	// -------------------------------------------------------
	//  COLLISION CHECK
	// -------------------------------------------------------

	newX := p.X + dx
	newY := p.Y + dy

	// Player box at the new location
	newBox := p.AABBAt(newX, newY)

	blocked := false

	playerZ := int(p.Z)

	for y := 0; y < MapHeight; y++ {
		for x := 0; x < MapWidth; x++ {

			tile := mapData[playerZ][y][x]
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
		p.X = newX
		p.Y = newY
	}
}

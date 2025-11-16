package main

const (
	TILE_SIZE = 32
)

var (
	SCREEN_WIDTH  = 1280
	SCREEN_HEIGHT = 720
)

type (
	ResourcePaths struct {
		PlayerImage string
		PlayerDesc  string
		TilsetImage []byte
		TilesetDesc string
	}
)

type GameState int

const (
	StateStartMenu GameState = iota
	StatePlaying
	StatePaused
)

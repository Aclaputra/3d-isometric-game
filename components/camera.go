package components

import "github.com/hajimehoshi/ebiten/v2"

type Camera struct {
	X, Y float64
	Zoom float64
}

func NewCamera(x, y, zoom float64) *Camera {
	return &Camera{
		X:    x,
		Y:    y,
		Zoom: zoom,
	}
}

func (c *Camera) Update(playerX, playerY float64) {
	// Smoothly follow player
	c.X += (playerX - c.X) * 0.1
	c.Y += (playerY - c.Y) * 0.1

	// Zoom with mouse wheel
	_, wheelY := ebiten.Wheel()
	if wheelY != 0 {
		c.Zoom += wheelY * 0.1 // adjust zoom speed here
		if c.Zoom < 0.2 {      // minimum zoom
			c.Zoom = 0.2
		}
		if c.Zoom > 3.0 { // maximum zoom
			c.Zoom = 3.0
		}
	}
}

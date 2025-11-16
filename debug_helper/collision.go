package debug_helper

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
)

func DrawQuadFilled(screen *ebiten.Image,
	x0, y0, x1, y1, x2, y2, x3, y3 float64,
	col color.RGBA,
) {
	// A white base image (1×1)
	img := ebiten.NewImage(1, 1)
	img.Fill(color.White)

	op := &ebiten.DrawTrianglesOptions{}

	// Convert color to float32
	r := float32(col.R) / 255
	g := float32(col.G) / 255
	b := float32(col.B) / 255
	a := float32(col.A) / 255 // opacity!

	vertices := []ebiten.Vertex{
		{DstX: float32(x0), DstY: float32(y0), SrcX: 0, SrcY: 0, ColorR: r, ColorG: g, ColorB: b, ColorA: a},
		{DstX: float32(x1), DstY: float32(y1), SrcX: 1, SrcY: 0, ColorR: r, ColorG: g, ColorB: b, ColorA: a},
		{DstX: float32(x2), DstY: float32(y2), SrcX: 1, SrcY: 1, ColorR: r, ColorG: g, ColorB: b, ColorA: a},
		{DstX: float32(x3), DstY: float32(y3), SrcX: 0, SrcY: 1, ColorR: r, ColorG: g, ColorB: b, ColorA: a},
	}

	indices := []uint16{
		0, 1, 2, // first triangle
		0, 2, 3, // second triangle
	}

	screen.DrawTriangles(vertices, indices, img, op)
}

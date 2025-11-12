package main

type TileMap struct {
	mapData [MapLevels][MapHeight][MapWidth]int // binary map
}

func NewTileMap(noise *Noise) *TileMap {
	tm := &TileMap{}

	for y := 0; y < MapHeight; y++ {
		for x := 0; x < MapWidth; x++ {

			// Get Perlin noise value (normalized 0..1)
			nv := noise.Perlin.Noise2D(float64(x)/10, float64(y)/10)
			nv = (nv + 1) / 2

			// Convert noise to a "height level"
			height := int(nv * float64(MapLevels))

			// Clamp to range [0, MapLevels-1]
			if height >= MapLevels {
				height = MapLevels - 1
			}

			// Fill all levels below that height (solid stack)
			for z := 0; z <= height; z++ {
				tm.mapData[z][y][x] = 1
			}
		}
	}

	return tm
}

package main

type TileMap struct {
	mapData [MapLevels][MapHeight][MapWidth]Tile // binary map
}

func NewTileMap(noise *Noise, tiles map[string]*Tile) *TileMap {
	tm := &TileMap{}

	dirtTile := tiles["Dirt Block"]
	for y := 0; y < MapHeight; y++ {
		for x := 0; x < MapWidth; x++ {

			// Perlin Noise → value 0..1
			nv := noise.Perlin.Noise2D(float64(x)/10, float64(y)/10)
			nv = (nv + 1) / 2

			// Convert noise to height level
			height := int(nv * float64(MapLevels))
			if height >= MapLevels {
				height = MapLevels - 1
			}

			// Fill stack
			for z := 0; z <= height; z++ {

				tile := &tm.mapData[z][y][x]

				tile.Name = dirtTile.Name
				tile.IsWall = dirtTile.IsWall
				tile.Image = dirtTile.Image

				// if z == 1 && y == 9 && x == 8 {
				offsetX := 1.6 // right in tile-space
				offsetY := 0.6 // up in tile-space

				tile.CollisionBox = &AABB{
					MinX: float64(x) + offsetX,
					MinY: float64(y) + offsetY,
					MaxX: float64(x) + 1 + offsetX,
					MaxY: float64(y) + 1 + offsetY,
				}
				// }

				// CollisionBox is generated in Draw(), not here.
			}
		}
	}

	return tm
}

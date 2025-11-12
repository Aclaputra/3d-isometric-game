package main

import "github.com/aquilax/go-perlin"

type Noise struct {
	Perlin *perlin.Perlin
}

func NewNoise() *Noise {
	alpha := 2.0       // persistence: higher = more jagged, mountain-like
	beta := 1.5        // frequency: controls the horizontal stretch
	n := 6             // number of octaves: more detail
	seed := int64(149) // fixed seed for reproducibility

	return &Noise{
		Perlin: perlin.NewPerlin(alpha, beta, int32(n), seed),
	}
}

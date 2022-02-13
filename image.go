package main

import (
	"image"
	"image/color"
)

type layer map[image.Point]color.Color

type layeredImage struct {
	l layer
	image.Image
}

func (li *layeredImage) At(x, y int) color.Color {
	if c, ok := li.l[image.Pt(x, y)]; ok {
		return c
	}
	return li.Image.At(x, y)
}

func (li *layeredImage) Set(p image.Point, c color.Color) {
	li.l[p] = c
}

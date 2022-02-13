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

func floodFill(m *layeredImage, p image.Point, fromColor, toColor color.Color) {
	m.Set(p, toColor)
	b := m.Bounds()
	for _, pt := range []image.Point{
		image.Pt(p.X-1, p.Y),
		image.Pt(p.X+1, p.Y),
		image.Pt(p.X, p.Y-1),
		image.Pt(p.X, p.Y+1),
	} {
		if pt.In(b) && m.At(pt.X, pt.Y) == fromColor {
			m.Set(pt, toColor)
			floodFill(m, pt, fromColor, toColor)
		}
	}
}

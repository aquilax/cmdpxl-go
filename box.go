package main

import (
	"image"

	"github.com/gdamore/tcell/v2"
)

type drawBox struct {
	borderSize int
	image.Rectangle
}

func newDrawBoxCoord(x1, y1, x2, y2 int) *drawBox {
	return &drawBox{
		borderSize,
		image.Rect(x1, y1, x2, y2).Canon(),
	}
}

func newDrawBox(x, y, width, height int) *drawBox {
	return &drawBox{
		borderSize,
		image.Rect(x, y, x+width-1, y+height-1),
	}
}

func (db *drawBox) draw(s tcell.Screen, style tcell.Style) *drawBox {
	x1 := db.Min.X
	y1 := db.Min.Y
	x2 := db.Max.X
	y2 := db.Max.Y

	for x := x1 + 1; x <= x2-1; x++ {
		// top border
		s.SetContent(x, y1, '─', nil, style)
		// bottom border
		s.SetContent(x, y2, '─', nil, style)
	}
	for y := y1 + 1; y <= y2-1; y++ {
		// top border
		s.SetContent(x1, y, '│', nil, style)
		// bottom border
		s.SetContent(x2, y, '│', nil, style)
	}
	if y1 != y2 && x1 != x2 {
		s.SetContent(x1, y1, '╭', nil, style)
		s.SetContent(x2, y1, '╮', nil, style)
		s.SetContent(x1, y2, '╰', nil, style)
		s.SetContent(x2, y2, '╯', nil, style)
	}
	return db
}

func (db *drawBox) getPoint(x, y int) image.Point {
	return image.Pt(x+db.Min.X+db.borderSize, y+db.Min.Y+db.borderSize)
}

func (db drawBox) getCanvas() image.Rectangle {
	return db.Inset(db.borderSize)
}

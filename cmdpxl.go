package main

import (
	"fmt"
	"image"
	"image/color"

	"github.com/gdamore/tcell/v2"
)

type layer map[image.Point]color.Color

type historyItem struct {
	point image.Point
	color color.Color
}

type CmdPxl struct {
	screenWidth  int
	screenHeight int
	imageWidth   int
	imageHeight  int

	cursorX int
	cursorY int

	paddingX int
	paddingY int

	m              image.Image
	fileName       string
	interfaceStyle tcell.Style
	s              tcell.Screen
	penColor       color.Color
	history        []historyItem
	paintLayer     layer
}

func NewCmdPxl(fileName string, m image.Image) *CmdPxl {
	b := m.Bounds()
	return &CmdPxl{
		interfaceStyle: tcell.StyleDefault.Foreground(tcell.ColorWhite).Background(tcell.ColorReset),
		fileName:       fileName,
		m:              m,
		imageWidth:     b.Max.X,
		imageHeight:    b.Max.Y,
		paddingY:       1,
		cursorX:        0,
		cursorY:        0,
		penColor:       color.White,
		history:        make([]historyItem, 0),
		paintLayer:     make(layer),
	}
}

func (c *CmdPxl) Run() error {
	var err error

	// Initialize screen
	c.s, err = tcell.NewScreen()
	if err != nil {
		return err
	}
	if err := c.s.Init(); err != nil {
		return err
	}

	c.s.SetStyle(c.interfaceStyle)

	defer c.s.Fini()

mainLoop:
	for {
		// Update screen
		c.s.Show()

		// Poll event
		ev := c.s.PollEvent()

		// Process event
		switch ev := ev.(type) {
		case *tcell.EventResize:
			c.screenWidth, c.screenHeight = ev.Size()
			c.paddingX = max(1, (c.screenWidth-max(48, c.imageWidth*2))/2)
			c.s.Sync()
		case *tcell.EventKey:
			// quit
			if ev.Key() == tcell.KeyEscape || ev.Key() == tcell.KeyCtrlC || ev.Rune() == 'x' {
				break mainLoop
			}
			// move cursor
			if ev.Rune() == 'w' {
				c.cursorY = mod(c.cursorY-1, c.imageHeight)
			}
			if ev.Rune() == 's' {
				c.cursorY = mod(c.cursorY+1, c.imageHeight)
			}
			if ev.Rune() == 'a' {
				c.cursorX = mod(c.cursorX-1, c.imageWidth)
			}
			if ev.Rune() == 'd' {
				c.cursorX = mod(c.cursorX+1, c.imageWidth)
			}
			if ev.Rune() == 'e' || ev.Rune() == ' ' {
				c.history = append(c.history, historyItem{image.Point{c.cursorX, c.cursorY}, c.penColor})
				c.paintLayer[image.Point{c.cursorX, c.cursorY}] = c.penColor
			}
			if ev.Rune() == 'z' {
				l := len(c.history)
				if l > 0 {
					c.history = c.history[:l-1]
					c.paintLayer = getLayerFromHistory(c.history)
				}
			}
		}
		c.draw()
	}
	return nil
}

func (c *CmdPxl) draw() {
	c.drawImageBox()
	c.drawImage()
	c.drawInterface()
}

func (c *CmdPxl) drawBox(x1, y1, x2, y2 int, style tcell.Style) {
	if y2 < y1 {
		y1, y2 = y2, y1
	}
	if x2 < x1 {
		x1, x2 = x2, x1
	}
	for x := x1 + 1; x <= x2-1; x++ {
		// top border
		c.s.SetContent(x, y1, '─', nil, style)
		// bottom border
		c.s.SetContent(x, y2, '─', nil, style)
	}
	for y := y1 + 1; y <= y2-1; y++ {
		// top border
		c.s.SetContent(x1, y, '│', nil, style)
		// bottom border
		c.s.SetContent(x2, y, '│', nil, style)
	}
	if y1 != y2 && x1 != x2 {
		c.s.SetContent(x1, y1, '╭', nil, style)
		c.s.SetContent(x2, y1, '╮', nil, style)
		c.s.SetContent(x1, y2, '╰', nil, style)
		c.s.SetContent(x2, y2, '╯', nil, style)
	}
}

func (c *CmdPxl) drawImageBox() {
	offsetY := 5
	x := min(c.imageWidth+1, c.screenWidth/2-2)
	y := min(c.imageHeight+1, c.screenHeight-12)

	x1 := c.paddingX
	y1 := offsetY + c.paddingY
	x2 := x1 + x*2 - 1
	y2 := y1 + y
	c.drawBox(x1, y1, x2, y2, c.interfaceStyle)
}

func (c *CmdPxl) drawImage() {
	offX := c.paddingX + 1
	offY := c.paddingY + 1 + 5
	for y := 0; y < c.imageHeight; y++ {
		for x := 0; x < c.imageWidth; x++ {
			color := c.m.At(x, y)
			if c, ok := c.paintLayer[image.Point{x, y}]; ok {
				color = c
			}
			bgColor := tcell.FromImageColor(color)
			style := tcell.StyleDefault.Background(bgColor)
			if c.cursorX == x && c.cursorY == y {
				style = style.Foreground(tcell.FromImageColor(getFgColor(color)))
				c.s.SetContent(x*2+offX, y+offY, '[', nil, style)
				c.s.SetContent(x*2+offX+1, y+offY, ']', nil, style)
			} else {
				c.s.SetContent(x*2+offX, y+offY, ' ', nil, style)
				c.s.SetContent(x*2+offX+1, y+offY, ' ', nil, style)
			}
		}
	}
}

func (c *CmdPxl) drawInterface() {
	drawText(c.s, c.paddingX, 1, c.interfaceStyle, fmt.Sprintf("CMDPXL-GO: %s (%dx%d) | pos: %02d,%02d", c.fileName, c.imageWidth, c.imageHeight, c.cursorX, c.cursorY))
	drawText(c.s, c.paddingX, c.screenHeight-3, c.interfaceStyle, "[wasd] move | [e] draw | [f] fill | [arrows] pan")
	drawText(c.s, c.paddingX, c.screenHeight-2, c.interfaceStyle, "[z] undo | [t] filters | [x] quit")
}

func drawText(s tcell.Screen, x, y int, style tcell.Style, text string) {
	row := y
	col := x
	for _, r := range text {
		s.SetContent(col, row, r, nil, style)
		col++
	}
}

func min(numbers ...int) int {
	min := numbers[0]

	for _, i := range numbers {
		if min > i {
			min = i
		}
	}

	return min
}

func max(numbers ...int) int {
	max := numbers[0]

	for _, i := range numbers {
		if max < i {
			max = i
		}
	}

	return max
}

func mod(a, b int) int {
	if a < 0 {
		return b - 1
	}
	if a > b-1 {
		return 0
	}
	return a
}

func getFgColor(c color.Color) color.Color {
	// https://socketloop.com/tutorials/golang-find-relative-luminance-or-color-brightness
	red, green, blue, _ := c.RGBA()
	lum := float64(float64(0.299)*float64(red) + float64(0.587)*float64(green) + float64(0.114)*float64(blue))
	if lum < .5 {
		return color.White
	}
	return color.Black
}

func getLayerFromHistory(h []historyItem) layer {
	l := make(layer)
	for _, hi := range h {
		l[hi.point] = hi.color
	}
	return l
}

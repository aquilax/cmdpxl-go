package main

import (
	"fmt"
	"image"
	"image/color"
	"strings"

	"github.com/gdamore/tcell/v2"
)

type direction bool
type state int
type saveImageCallback = func(fileName string, m image.Image) error

const (
	maxHue     = 380
	borderSize = 1

	dirIncrease  direction = true
	dirDecrease  direction = false
	stateDrawing state     = iota
	stateQuit
)

type historyItem struct {
	point image.Point
	color color.Color
}

type CmdPxl struct {
	currentState state

	screenWidth  int
	screenHeight int

	imageWidth  int
	imageHeight int

	maxDrawWidth  int
	maxDrawHeight int

	imageBox *drawBox

	cursorX int
	cursorY int

	paddingX int
	paddingY int

	panX int
	panY int

	paletteSize    int
	m              layeredImage
	fileName       string
	interfaceStyle tcell.Style
	s              tcell.Screen
	penColor       cmdColor
	history        []historyItem

	saveImage saveImageCallback
}

func NewCmdPxl(fileName string, m image.Image, saveImage saveImageCallback) *CmdPxl {
	b := m.Bounds()
	paletteSize := 11

	return &CmdPxl{
		currentState:   stateDrawing,
		interfaceStyle: tcell.StyleDefault.Foreground(tcell.ColorWhite).Background(tcell.ColorReset),
		fileName:       fileName,
		m:              layeredImage{make(layer), m},
		imageWidth:     b.Max.X,
		imageHeight:    b.Max.Y,
		panX:           0,
		panY:           0,
		maxDrawWidth:   0,
		maxDrawHeight:  0,
		paddingY:       1,
		cursorX:        0,
		cursorY:        0,
		paletteSize:    paletteSize,
		penColor:       *NewCmdColor(color.White, paletteSize),
		history:        make([]historyItem, 0),
		saveImage:      saveImage,
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
			c.paddingX = max(0, (c.screenWidth-max(48, c.imageWidth*2))/2)

			c.maxDrawWidth = c.screenWidth - 2*borderSize
			c.maxDrawHeight = c.screenHeight - 13 // chrome

			c.imageBox = c.getImageBox()
			c.s.Sync()
		case *tcell.EventKey:
			if c.currentState == stateDrawing {
				// quit
				if ev.Key() == tcell.KeyEscape || ev.Key() == tcell.KeyCtrlC || ev.Rune() == 'x' {
					// any changes made
					if len(c.history) > 0 {
						c.currentState = stateQuit
					} else {
						// quit directly
						break mainLoop
					}
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
					pt := image.Pt(c.cursorX+c.panX, c.cursorY+c.panY)
					c.history = append(c.history, historyItem{pt, c.penColor.c})
					c.m.Set(pt, c.penColor.c)
				}
				if ev.Rune() == 'z' {
					l := len(c.history)
					if l > 0 {
						c.history = c.history[:l-1]
						c.m.l = getLayerFromHistory(c.history)
					}
				}
				if ev.Rune() == 'D' {
					// debug
					newDrawBox(0, 0, 5, 3).draw(c.s, c.interfaceStyle)
					drawText(c.s, 0, 10, c.interfaceStyle, "box: "+c.imageBox.String())
					drawText(c.s, 0, 11, c.interfaceStyle, "canvas: "+c.imageBox.getCanvas().String())
					drawText(c.s, 0, 12, c.interfaceStyle, fmt.Sprintf("WxH: %dx%d", c.imageBox.getCanvas().Dx(), c.imageBox.getCanvas().Dy()))
					drawText(c.s, 0, 13, c.interfaceStyle, fmt.Sprintf("Pan XxY: %dx%d", c.panX, c.panY))
					drawText(c.s, 0, 14, c.interfaceStyle, fmt.Sprintf("image XxY: %dx%d", c.imageWidth, c.imageHeight))
				}

				// colors

				// hue
				if ev.Rune() == 'j' {
					c.penColor.changeHue(dirIncrease)
				}
				if ev.Rune() == 'u' {
					c.penColor.changeHue(dirDecrease)
				}

				// saturation
				if ev.Rune() == 'k' {
					c.penColor.changeSaturation(dirIncrease)
				}
				if ev.Rune() == 'i' {
					c.penColor.changeSaturation(dirDecrease)
				}

				// value
				if ev.Rune() == 'l' {
					c.penColor.changeValue(dirIncrease)
				}
				if ev.Rune() == 'o' {
					c.penColor.changeValue(dirDecrease)
				}

				// panning
				if ev.Key() == tcell.KeyUp {
					c.panY -= 1
					if c.panY < 0 {
						c.panY = c.imageHeight - (c.imageBox.getCanvas().Dy() + 1)
					}
				}
				if ev.Key() == tcell.KeyDown {
					c.panY += 1
					if c.panY > c.imageHeight-(c.imageBox.getCanvas().Dy()+1) {
						c.panY = 0
					}
				}
				if ev.Key() == tcell.KeyLeft {
					c.panX -= 1
					if c.panX < 0 {
						c.panX = c.imageWidth - ((c.imageBox.getCanvas().Dx() + 1) / 2)
					}
				}
				if ev.Key() == tcell.KeyRight {
					c.panX += 1
					if c.panX > c.imageWidth-((c.imageBox.getCanvas().Dx()+1)/2) {
						c.panX = 0
					}
				}

				if ev.Rune() == 'f' {
					color := c.m.At(c.cursorX, c.cursorY)
					if color != c.penColor.c {
						floodFill(&c.m, image.Pt(c.cursorX, c.cursorY), c.m.At(c.cursorX, c.cursorY), c.penColor.c)
					}
				}

			} else if c.currentState == stateQuit {
				if ev.Rune() == 'y' || ev.Rune() == 'Y' {
					err := c.saveImage(c.fileName, &c.m)
					if err != nil {
						return err
					}
					break mainLoop
				}
				if ev.Rune() == 'n' || ev.Rune() == 'N' || ev.Key() == tcell.KeyEscape {
					c.currentState = stateDrawing
					c.s.Clear()
				}
			}
		}
		c.draw()
	}
	return nil
}

func (c *CmdPxl) draw() {
	c.drawInterface()
	c.drawColorSelect()
	c.imageBox.draw(c.s, c.interfaceStyle)
	c.drawImage(c.imageBox)
	if c.currentState == stateQuit {
		c.drawExitConfirmation()
	}
}

func (c *CmdPxl) drawExitConfirmation() *drawBox {
	confirmation := "Do you want to exit? [y/n]"
	dBox := newDrawBox(0, 0, len(confirmation)+2+borderSize*2, borderSize*2+1).draw(c.s, c.interfaceStyle)
	p := dBox.getPoint(1, 0)
	drawText(c.s, p.X, p.Y, c.interfaceStyle, confirmation)
	return dBox
}

func (c *CmdPxl) getImageBox() *drawBox {
	offsetY := 5
	interfaceRows := 12
	width := min(c.imageWidth*2+2, c.screenWidth)
	height := min(c.imageHeight+2, c.screenHeight-interfaceRows)

	x := c.paddingX
	y := offsetY + c.paddingY
	return newDrawBox(x, y, width, height)
}

func (c *CmdPxl) drawImage(dBox *drawBox) {
	canvas := dBox.getCanvas()
	const pixelWidth = 2
	xBoundary := min(c.imageWidth*2, canvas.Dx()/pixelWidth+1)
	yBoundary := min(c.imageHeight, canvas.Dy()+1)
	for y := 0; y < yBoundary; y++ {
		for x := 0; x < xBoundary; x++ {
			imageColor := c.m.At(x+c.panX, y+c.panY)
			bgColor := tcell.FromImageColor(imageColor)
			style := tcell.StyleDefault.Background(bgColor)
			p := dBox.getPoint(x*2, y)
			if c.cursorX == x && c.cursorY == y {
				style = style.Foreground(tcell.FromImageColor(getFgColor(imageColor)))
				c.s.SetContent(p.X, p.Y, '[', nil, style)
				c.s.SetContent(p.X+1, p.Y, ']', nil, style)
			} else {
				c.s.SetContent(p.X, p.Y, ' ', nil, style)
				c.s.SetContent(p.X+1, p.Y, ' ', nil, style)
			}
		}
	}
}

func (c *CmdPxl) drawInterface() {
	drawText(c.s, c.paddingX, 1, c.interfaceStyle, fmt.Sprintf("CMDPXL-GO: %s (%dx%d) | pos: %03d,%03d", c.fileName, c.imageWidth, c.imageHeight, c.cursorX, c.cursorY))
	p := newDrawBox(c.paddingX, c.screenHeight-4, 100, 4).getPoint(0, 0)
	drawText(c.s, p.X, p.Y, c.interfaceStyle, "[wasd] move | [e] draw | [f] fill | [arrows] pan")
	drawText(c.s, p.X, p.Y+1, c.interfaceStyle, "[z] undo | [t] filters | [x] quit")
}

func (c *CmdPxl) drawColorSelect() {
	sectionWidth := 12
	boxHeight := 4
	numBoxes := 4
	x1 := c.paddingX
	y1 := c.paddingY + 1
	// box
	dBox := newDrawBox(x1, y1, numBoxes*sectionWidth+borderSize, boxHeight).draw(c.s, c.interfaceStyle)
	p := dBox.getPoint(0, 0)
	// instructions
	instructions := "[u/j]: hue  [i/k]: sat  [o/l]: val  current"
	drawText(c.s, p.X, p.Y, c.interfaceStyle, instructions)

	// Color selection
	p = dBox.getPoint(0, 1)
	for offset, cl := range c.penColor.huePalette {
		style := tcell.StyleDefault.Background(tcell.FromImageColor(cl))
		text := ' '
		if offset == c.penColor.huePaletteIndex {
			style.Foreground(tcell.FromImageColor(getFgColor(cl)))
			text = '●'
		}
		c.s.SetContent(p.X+offset, p.Y, text, nil, style)
	}

	// Saturation
	p = dBox.getPoint(sectionWidth*1, 1)
	c.s.SetContent(p.X-1, p.Y-2, '┬', nil, c.interfaceStyle)
	c.s.SetContent(p.X-1, p.Y-1, '│', nil, c.interfaceStyle)
	c.s.SetContent(p.X-1, p.Y+0, '│', nil, c.interfaceStyle)
	c.s.SetContent(p.X-1, p.Y+1, '┴', nil, c.interfaceStyle)
	for offset, cl := range c.penColor.saturationPalette {
		style := tcell.StyleDefault.Background(tcell.FromImageColor(cl))
		text := ' '
		if offset == c.penColor.saturationPaletteIndex {
			style.Foreground(tcell.FromImageColor(getFgColor(cl)))
			text = '●'
		}
		c.s.SetContent(p.X+offset, p.Y, text, nil, style)
	}

	// Value
	p = dBox.getPoint(sectionWidth*2, 1)
	c.s.SetContent(p.X-1, p.Y-2, '┬', nil, c.interfaceStyle)
	c.s.SetContent(p.X-1, p.Y-1, '│', nil, c.interfaceStyle)
	c.s.SetContent(p.X-1, p.Y+0, '│', nil, c.interfaceStyle)
	c.s.SetContent(p.X-1, p.Y+1, '┴', nil, c.interfaceStyle)
	for offset, cl := range c.penColor.valuePalette {
		style := tcell.StyleDefault.Background(tcell.FromImageColor(cl))
		text := ' '
		if offset == c.penColor.valuePaletteIndex {
			style.Foreground(tcell.FromImageColor(getFgColor(cl)))
			text = '●'
		}
		c.s.SetContent(p.X+offset, p.Y, text, nil, style)
	}

	// Current color
	style := tcell.StyleDefault.Background((tcell.FromImageColor(c.penColor.c)))
	p = dBox.getPoint(sectionWidth*3, 1)
	c.s.SetContent(p.X-1, p.Y-2, '┬', nil, c.interfaceStyle)
	c.s.SetContent(p.X-1, p.Y-1, '│', nil, c.interfaceStyle)
	c.s.SetContent(p.X-1, p.Y+0, '│', nil, c.interfaceStyle)
	c.s.SetContent(p.X-1, p.Y+1, '┴', nil, c.interfaceStyle)

	drawText(c.s, p.X, p.Y, style, strings.Repeat(" ", c.paletteSize))
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

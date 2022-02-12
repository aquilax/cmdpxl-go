package main

import (
	"fmt"
	"image"
	"image/color"
	"math"
	"strings"

	"github.com/gdamore/tcell/v2"
	"github.com/lucasb-eyer/go-colorful"
)

const maxHue = 380

type direction bool

const (
	dirIncrease direction = true
	dirDecrease direction = false
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

	paletteSize    int
	m              image.Image
	fileName       string
	interfaceStyle tcell.Style
	s              tcell.Screen
	penColor       cmdColor
	history        []historyItem
	paintLayer     layer
}

func NewCmdPxl(fileName string, m image.Image) *CmdPxl {
	b := m.Bounds()
	paletteSize := 11
	return &CmdPxl{
		interfaceStyle: tcell.StyleDefault.Foreground(tcell.ColorWhite).Background(tcell.ColorReset),
		fileName:       fileName,
		m:              m,
		imageWidth:     b.Max.X,
		imageHeight:    b.Max.Y,
		paddingY:       1,
		cursorX:        0,
		cursorY:        0,
		paletteSize:    paletteSize,
		penColor:       *NewCmdColor(color.RGBA{0, 0, 0, 1}, paletteSize),
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
				c.history = append(c.history, historyItem{image.Point{c.cursorX, c.cursorY}, c.penColor.c})
				c.paintLayer[image.Point{c.cursorX, c.cursorY}] = c.penColor.c
			}
			if ev.Rune() == 'z' {
				l := len(c.history)
				if l > 0 {
					c.history = c.history[:l-1]
					c.paintLayer = getLayerFromHistory(c.history)
				}
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
		}
		c.draw()
	}
	return nil
}

func (c *CmdPxl) draw() {
	c.drawInterface()
	c.drawColorSelect()
	c.drawImageBox()
	c.drawImage()

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

func (c *CmdPxl) drawColorSelect() {
	sectionWidth := 12
	boxHeight := 3
	boxes := 4
	x1 := c.paddingX
	y1 := c.paddingY + 1
	x2 := x1 + boxes*sectionWidth
	y2 := y1 + boxHeight
	// box
	c.drawBox(x1, y1, x2, y2, c.interfaceStyle)
	// instructions
	instructions := "[u/j]: hue  [i/k]: sat  [o/l]: val  current"
	drawText(c.s, x1+1, y1+1, c.interfaceStyle, instructions)

	// Color selection
	for offset, cl := range c.penColor.huePalette {
		style := tcell.StyleDefault.Background(tcell.FromImageColor(cl))
		text := ' '
		if offset == c.penColor.huePaletteIndex {
			style.Foreground(tcell.FromImageColor(getFgColor(cl)))
			text = '●'
		}
		c.s.SetContent(x1+1+offset, y1+2, text, nil, style)
	}

	// Saturation
	for offset, cl := range c.penColor.saturationPalette {
		style := tcell.StyleDefault.Background(tcell.FromImageColor(cl))
		text := ' '
		if offset == c.penColor.saturationPaletteIndex {
			style.Foreground(tcell.FromImageColor(getFgColor(cl)))
			text = '●'
		}
		c.s.SetContent(x1+offset+1+sectionWidth*1, y1+2, text, nil, style)
	}

	// Value
	for offset, cl := range c.penColor.valuePalette {
		style := tcell.StyleDefault.Background(tcell.FromImageColor(cl))
		text := ' '
		if offset == c.penColor.valuePaletteIndex {
			style.Foreground(tcell.FromImageColor(getFgColor(cl)))
			text = '●'
		}
		c.s.SetContent(x1+offset+1+sectionWidth*2, y1+2, text, nil, style)
	}

	// Current color
	style := tcell.StyleDefault.Background((tcell.FromImageColor(c.penColor.c)))
	drawText(c.s, x1+1+sectionWidth*3, y1+2, style, strings.Repeat(" ", c.paletteSize))

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

type cmdColor struct {
	c color.Color

	hue        float64
	saturation float64
	value      float64

	paletteSize int

	huePaletteIndex int
	huePalette      []colorful.Color

	saturationPaletteIndex int
	saturationPalette      []colorful.Color

	valuePaletteIndex int
	valuePalette      []colorful.Color
}

func NewCmdColor(c color.Color, paletteSize int) *cmdColor {
	cl, _ := colorful.MakeColor(c)
	h, s, v := cl.Hsv()
	huePalette := getHuePalette(paletteSize)
	huePaletteIndex := getHuePaletteIndex(h, huePalette)

	saturationPalette := getSaturationPalette(h, paletteSize)
	saturationPaletteIndex := getSaturationPaletteIndex(s, saturationPalette)

	valuePalette := getValuePalette(h, s, paletteSize)
	valuePaletteIndex := getValuePaletteIndex(v, valuePalette)

	return &cmdColor{
		c:                      c,
		hue:                    h,
		saturation:             s,
		value:                  v,
		paletteSize:            paletteSize,
		huePaletteIndex:        huePaletteIndex,
		huePalette:             huePalette,
		saturationPalette:      saturationPalette,
		saturationPaletteIndex: saturationPaletteIndex,
		valuePalette:           valuePalette,
		valuePaletteIndex:      valuePaletteIndex,
	}
}

func getHuePalette(items int) []colorful.Color {
	const (
		saturation = 1.0
		value      = 1.0
	)
	result := make([]colorful.Color, items)
	step := maxHue / float64(items)
	for i := 0; i < items; i++ {
		h := float64(i) * step
		result[i] = colorful.Hsv(h, saturation, value)
	}
	return result
}

func getHuePaletteIndex(hue float64, palette []colorful.Color) int {
	hues := make([]float64, len(palette))
	for i, c := range palette {
		h, _, _ := c.Hsv()
		hues[i] = h
	}
	return getClosestIndex(hue, hues)
}

func getClosestIndex(value float64, list []float64) int {
	prev := 0.0
	for i, v := range list {
		if v > value {
			if math.Abs(prev-value) < math.Abs(v-value) {
				return i - 1
			} else {
				return i
			}
		}
		prev = v
	}
	return len(list) - 1
}

func (cc *cmdColor) changeHue(dir direction) {
	newIndex := -1
	if dir == dirIncrease {
		// up
		newIndex = cc.huePaletteIndex + 1
		if newIndex > cc.paletteSize-1 {
			newIndex = cc.paletteSize - 1
		}
	} else {
		//down
		newIndex = cc.huePaletteIndex - 1
		if newIndex < 0 {
			newIndex = 0
		}
	}
	cl := cc.huePalette[newIndex]
	newHue, _, _ := cl.Hsv()
	cc.c = colorful.Hsv(newHue, cc.saturation, cc.value)
	cc.hue = newHue
	cc.huePaletteIndex = getHuePaletteIndex(newHue, cc.huePalette)
}

func getSaturationPalette(hue float64, items int) []colorful.Color {
	const value = 1.0
	result := make([]colorful.Color, items)
	step := 1.0 / float64(items)
	for i := 0; i < items; i++ {
		saturation := float64(i) * step
		result[i] = colorful.Hsv(hue, saturation, value)
	}
	return result
}

func getSaturationPaletteIndex(saturation float64, palette []colorful.Color) int {
	saturations := make([]float64, len(palette))
	for i, c := range palette {
		_, s, _ := c.Hsv()
		saturations[i] = s
	}
	return getClosestIndex(saturation, saturations)
}

func (cc *cmdColor) changeSaturation(dir direction) {
	newIndex := -1
	if dir == dirIncrease {
		// up
		newIndex = cc.saturationPaletteIndex + 1
		if newIndex > cc.paletteSize-1 {
			newIndex = cc.paletteSize - 1
		}
	} else {
		//down
		newIndex = cc.saturationPaletteIndex - 1
		if newIndex < 0 {
			newIndex = 0
		}
	}
	cl := cc.saturationPalette[newIndex]
	_, newSaturation, _ := cl.Hsv()
	cc.c = colorful.Hsv(cc.hue, newSaturation, cc.value)
	cc.saturation = newSaturation
	cc.saturationPaletteIndex = getSaturationPaletteIndex(newSaturation, cc.saturationPalette)
}

func getValuePalette(hue, saturation float64, items int) []colorful.Color {
	result := make([]colorful.Color, items)
	step := 1.0 / float64(items)
	for i := 0; i < items; i++ {
		value := float64(i) * step
		result[i] = colorful.Hsv(hue, saturation, value)
	}
	return result
}

func getValuePaletteIndex(value float64, palette []colorful.Color) int {
	values := make([]float64, len(palette))
	for i, c := range palette {
		_, _, v := c.Hsv()
		values[i] = v
	}
	return getClosestIndex(value, values)
}

func (cc *cmdColor) changeValue(dir direction) {
	newIndex := -1
	if dir == dirIncrease {
		// up
		newIndex = cc.valuePaletteIndex + 1
		if newIndex > cc.paletteSize-1 {
			newIndex = cc.paletteSize - 1
		}
	} else {
		//down
		newIndex = cc.valuePaletteIndex - 1
		if newIndex < 0 {
			newIndex = 0
		}
	}
	cl := cc.valuePalette[newIndex]
	_, _, newValue := cl.Hsv()
	cc.c = colorful.Hsv(cc.hue, cc.saturation, newValue)
	cc.value = newValue
	cc.valuePaletteIndex = getValuePaletteIndex(newValue, cc.valuePalette)
}

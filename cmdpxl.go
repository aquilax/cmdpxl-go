package main

import (
	"image"

	"github.com/gdamore/tcell/v2"
)

type CmdPxl struct {
	screenWidth  int
	screenHeight int
	imageWidth   int
	imageHeight  int

	paddingX int
	paddingY int

	m              image.Image
	interfaceStyle tcell.Style
	s              tcell.Screen
}

func NewCmdPxl(m image.Image) *CmdPxl {
	b := m.Bounds()
	return &CmdPxl{
		interfaceStyle: tcell.StyleDefault.Foreground(tcell.ColorWhite).Background(tcell.ColorReset),
		m:              m,
		imageWidth:     b.Max.X,
		imageHeight:    b.Max.Y,
		paddingY:       1,
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
			if ev.Key() == tcell.KeyEscape || ev.Key() == tcell.KeyCtrlC {
				break mainLoop
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
	// s.SetContent(x2, y1, tcell.RuneURCorner, nil, style)
	// s.SetContent(x1, y2, tcell.RuneLLCorner, nil, style)
	// s.SetContent(x2, y2, tcell.RuneLRCorner, nil, style)
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
	offsetY := 6
	x := min(c.imageWidth+1, c.screenWidth/2-2)
	y := min(c.imageHeight+1, c.screenHeight-12)

	x1 := 1 + c.paddingX
	y1 := offsetY + c.paddingY
	x2 := x1 + x*2 - 1
	y2 := y1 + y
	c.drawBox(x1, y1, x2, y2, c.interfaceStyle)
}

func (c *CmdPxl) drawImage() {
	offX := c.paddingX + 2
	offY := c.paddingY + 1 + 6
	for j := 0; j < c.imageHeight; j++ {
		for i := 0; i < c.imageWidth; i++ {
			color := c.m.At(i, j)
			c.s.SetContent(i*2+offX, j+offY, ' ', nil, tcell.StyleDefault.Background(tcell.FromImageColor(color)))
			c.s.SetContent(i*2+offX+1, j+offY, ' ', nil, tcell.StyleDefault.Background(tcell.FromImageColor(color)))
		}
	}
}

func (c *CmdPxl) drawInterface() {}

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

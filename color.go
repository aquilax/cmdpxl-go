package main

import (
	"image/color"
	"math"

	"github.com/lucasb-eyer/go-colorful"
)

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
	cc.saturationPalette = getSaturationPalette(newHue, cc.paletteSize)
	cc.valuePalette = getValuePalette(newHue, cc.saturation, cc.paletteSize)
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
	cc.valuePalette = getValuePalette(cc.hue, cc.saturation, cc.paletteSize)
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

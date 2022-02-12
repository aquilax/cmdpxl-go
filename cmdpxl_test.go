package main

import (
	"testing"

	"github.com/lucasb-eyer/go-colorful"
)

func Test_getPaletteIndex(t *testing.T) {
	type args struct {
		color   colorful.Color
		palette []colorful.Color
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{
			"works for green",
			args{
				color:   colorful.Color{R: 0, G: 255, B: 0},
				palette: getPalette(11),
			},
			3,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hue, _, _ := tt.args.color.Hsv()
			if got := getPaletteIndex(hue, tt.args.palette); got != tt.want {
				t.Errorf("getPaletteIndex() = %v, want %v", got, tt.want)
			}
		})
	}
}

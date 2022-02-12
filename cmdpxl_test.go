package main

import (
	"image"
	"reflect"
	"testing"

	"github.com/lucasb-eyer/go-colorful"
)

func Test_getHuePaletteIndex(t *testing.T) {
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
				palette: getHuePalette(11),
			},
			3,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hue, _, _ := tt.args.color.Hsv()
			if got := getHuePaletteIndex(hue, tt.args.palette); got != tt.want {
				t.Errorf("getHuePaletteIndex() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_drawBox_getPoint(t *testing.T) {
	tests := []struct {
		name string
		dBox *drawBox
		args image.Point
		want image.Point
	}{
		{
			"works with one line box",
			newDrawBox(0, 0, 5, 3),
			image.Pt(0, 0),
			image.Pt(1, 1),
		},
		{
			"works with x offest",
			newDrawBox(0, 0, 5, 3),
			image.Pt(1, 0),
			image.Pt(2, 1),
		},
		{
			"works with box offest",
			newDrawBox(10, 20, 5, 3),
			image.Pt(0, 0),
			image.Pt(11, 21),
		},
		{
			"works with box and point offest",
			newDrawBox(10, 20, 5, 3),
			image.Pt(3, 5),
			image.Pt(14, 26),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.dBox.getPoint(tt.args.X, tt.args.Y); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("drawBox.getPoint() = %v, want %v", got, tt.want)
			}
		})
	}
}

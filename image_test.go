package main

import (
	"bytes"
	"encoding/base64"
	"image"
	"image/color"
	"image/png"
	"testing"
)

func Test_Image_Save(t *testing.T) {
	i, _ := createImage("5,5")
	li := layeredImage{make(layer), i}
	li.Set(image.Pt(0, 0), color.White)
	b := new(bytes.Buffer)
	png.Encode(b, &li)
	image1 := base64.StdEncoding.EncodeToString(b.Bytes())

	i2, _ := createImage("5,5")
	li2 := layeredImage{make(layer), i2}
	li2.Set(image.Pt(1, 0), color.White)
	b2 := new(bytes.Buffer)
	png.Encode(b2, &li2)
	image2 := base64.StdEncoding.EncodeToString(b2.Bytes())

	if image1 == image2 {
		t.Errorf("png encoded i1 = %v expected to be different than i2 = %v", image2, image1)
	}
}

func Test_floodFill(t *testing.T) {
	type args struct {
		p       image.Point
		toColor color.Color
	}
	tests := []struct {
		name       string
		args       args
		checkPoint image.Point
	}{
		{
			"works on empty canvas",
			args{
				image.Pt(2, 2),
				color.White,
			},
			image.Pt(0, 0),
		},
	}
	for _, tt := range tests {
		i, _ := createImage("5,5")
		li := layeredImage{make(layer), i}
		fromColor := li.At(tt.args.p.X, tt.args.p.Y)
		t.Run(tt.name, func(t *testing.T) {
			floodFill(&li, tt.args.p, fromColor, tt.args.toColor)
		})
		got := li.At(tt.checkPoint.X, tt.checkPoint.Y)
		if got != tt.args.toColor {
			t.Errorf("Expected to find color %v at %s but found %v instead", tt.args.toColor, tt.checkPoint, got)
		}
	}
}

func Benchmark_floodFill(b *testing.B) {
	i, _ := createImage("100,100")
	li := layeredImage{make(layer), i}
	fromColor := li.At(0, 0)
	for n := 0; n < b.N; n++ {
		li := layeredImage{make(layer), i}
		floodFill(&li, image.Pt(0, 0), fromColor, color.Black)
	}
}

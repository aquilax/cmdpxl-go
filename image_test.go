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

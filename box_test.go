package main

import (
	"image"
	"reflect"
	"testing"

	"github.com/gdamore/tcell/v2"
)

func Test_newDrawBoxCoord(t *testing.T) {
	type args struct {
		x1 int
		y1 int
		x2 int
		y2 int
	}
	tests := []struct {
		name string
		args args
		want *drawBox
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := newDrawBoxCoord(tt.args.x1, tt.args.y1, tt.args.x2, tt.args.y2); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("newDrawBoxCoord() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_newDrawBox(t *testing.T) {
	type args struct {
		x      int
		y      int
		width  int
		height int
	}
	tests := []struct {
		name string
		args args
		want *drawBox
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := newDrawBox(tt.args.x, tt.args.y, tt.args.width, tt.args.height); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("newDrawBox() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_drawBox_draw(t *testing.T) {
	type fields struct {
		borderSize int
		Rectangle  image.Rectangle
	}
	type args struct {
		s     tcell.Screen
		style tcell.Style
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   *drawBox
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db := &drawBox{
				borderSize: tt.fields.borderSize,
				Rectangle:  tt.fields.Rectangle,
			}
			if got := db.draw(tt.args.s, tt.args.style); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("drawBox.draw() = %v, want %v", got, tt.want)
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

func Test_drawBox_getCanvas(t *testing.T) {
	tests := []struct {
		name string
		dBox *drawBox
		want image.Rectangle
	}{
		{
			"7x7 box",
			newDrawBox(0, 0, 7, 7),
			image.Rect(1, 1, 5, 5),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.dBox.getCanvas(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("drawBox.getCanvas() = %v, want %v", got, tt.want)
			}
		})
	}
}

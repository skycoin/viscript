package app

import (
//"fmt"
)

type Rectangle struct {
	Top    float32
	Right  float32
	Bottom float32
	Left   float32
}

func (r *Rectangle) Width() float32 {
	return r.Right - r.Left
}

func (r *Rectangle) Height() float32 {
	return r.Top - r.Bottom
}

func (r *Rectangle) CenterX() float32 {
	return r.Left + r.Width()/2
}

func (r *Rectangle) CenterY() float32 {
	return r.Bottom + r.Height()/2
}

func (r *Rectangle) Contains(x, y float32) bool {
	if x > r.Left && x < r.Right && y > r.Bottom && y < r.Top {
		return true
	}

	return false
}

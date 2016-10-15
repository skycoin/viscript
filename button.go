package main

import (
//"fmt"
)

type Button struct {
	Name      string
	Activated bool
	Rect      *Rectangle
}

func (bu *Button) Init() {
}

func (bu *Button) Draw() {
	rend.DrawQuad(11, 13, bu.Rect)
}

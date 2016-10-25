package main

import (
//"fmt"
//"github.com/go-gl/gl/v2.1/gl"
)

type Button struct {
	Name      string
	Activated bool
	Rect      *Rectangle
}

func (bu *Button) Draw() {
	if bu.Activated {
		rend.Color(green)
	} else {
		rend.Color(white)
	}

	span := bu.Rect.Height() * goldenPercentage // ...of both dimensions of each character
	glTextWidth := float32(len(bu.Name)) * span // in terms of OpenGL/float32 space
	x := bu.Rect.Left + (bu.Rect.Width()-glTextWidth)/2
	verticalLipSpan := (bu.Rect.Height() - span) / 2 // lip or frame edge

	rend.DrawQuad(11, 13, bu.Rect)

	for _, c := range bu.Name {
		rend.DrawCharAtRect(c, &Rectangle{bu.Rect.Top - verticalLipSpan, x + span, bu.Rect.Bottom + verticalLipSpan, x})
		x += span
	}
}

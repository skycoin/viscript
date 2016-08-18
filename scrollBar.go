package main

import (
	"fmt"
	"github.com/go-gl/gl/v2.1/gl"
)

var sb = ScrollBar{0, rectRad, 0, 0}

type ScrollBar struct {
	PosX      float32
	PosY      float32
	LenOfBar  float32
	LenOfVoid float32 // length of the negative space adjacent to bar
}

func (bar ScrollBar) StartFrame() {
	bar.LenOfBar = float32(numYChars) / float32(len(document)) * rectRad * 2

	if /* no need for scrollbar, cuz entire doc less than screen */ len(document) <= numYChars {
		bar.LenOfBar = 0
		view.OffsetY = 0
	} else {
		bar.LenOfVoid = rectRad*2 - bar.LenOfBar
	}
}

func (bar ScrollBar) Scroll(mousePixelDeltaY float64) {
	bar.PosY -= float32(mousePixelDeltaY) * pixelHei

	if bar.PosY < -rectRad+bar.LenOfBar {
		bar.PosY = -rectRad + bar.LenOfBar
	}
	if bar.PosY > rectRad {
		bar.PosY = rectRad
	}

	view.OffsetY = bar.PosY / bar.LenOfVoid * float32(len(document)) * chHei
	fmt.Printf("PosY: %.1f - view.OffsetY: %.1f", bar.PosY, view.OffsetY)
}

func (bar ScrollBar) DrawVertical(atlasX, atlasY float32) {
	gl.Normal3f(0, 0, 1)
	u := float32(atlasX) * uvSpan
	v := float32(atlasY) * uvSpan

	gl.TexCoord2f(u, v+uvSpan) // bl  0, 1
	gl.Vertex3f(rectRad-chWid, bar.PosY-bar.LenOfBar, 0)

	gl.TexCoord2f(u+uvSpan, v+uvSpan) // br  1, 1
	gl.Vertex3f(rectRad, bar.PosY-bar.LenOfBar, 0)

	gl.TexCoord2f(u+uvSpan, v) // tr  1, 0
	gl.Vertex3f(rectRad, bar.PosY, 0)

	gl.TexCoord2f(u, v) // tl  0, 0
	gl.Vertex3f(rectRad-chWid, bar.PosY, 0)

	curX += chWid
}

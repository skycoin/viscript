package main

import (
	"fmt"
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
		viewportOffsetY = 0
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

	viewportOffsetY = bar.PosY / bar.LenOfVoid * float32(len(document)) * chHei
	fmt.Printf("PosY: %.1f - viewportOffsetY: %.1f", bar.PosY, bar.viewportOffsetY)
}

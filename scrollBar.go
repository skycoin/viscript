package main

import (
	"fmt"
)

var sb = scrollBar{0, rectRad, 0, 0}

type scrollBar struct {
	scrollPosX    float32
	scrollPosY    float32
	scrollBarLen  float32
	scrollVoidLen float32
}

func (bar scrollBar) StartFrame() {
	bar.scrollBarLen = float32(numYChars) / float32(len(document)) * rectRad * 2

	if /* no need for scrollbar, cuz entire doc less than screen */ len(document) <= numYChars {
		bar.scrollBarLen = 0
		viewportOffsetY = 0
	} else {
		bar.scrollVoidLen = rectRad*2 - bar.scrollBarLen
	}
}

func (bar scrollBar) Scroll(mousePixelDeltaY float64) {
	bar.scrollPosY -= float32(mousePixelDeltaY) * pixelHei

	if bar.scrollPosY < -rectRad+bar.scrollBarLen {
		bar.scrollPosY = -rectRad + bar.scrollBarLen
	}
	if bar.scrollPosY > rectRad {
		bar.scrollPosY = rectRad
	}

	viewportOffsetY = bar.scrollPosY / bar.scrollVoidLen * float32(len(document)) * chHei
	fmt.Printf("scrollPosY: %.1f - viewportOffsetY: %.1f", bar.scrollPosY, bar.viewportOffsetY)
}

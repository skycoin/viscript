package main

import (
//"fmt"
)

var sb = scrollBar{0, rectRad, 0, 0}

type scrollBar struct {
	scrollPosX    float32
	scrollPosY    float32
	scrollBarLen  float32
	scrollVoidLen float32
}

func scroll(mousePixelDeltaY float64) {
	sb.scrollPosY -= float32(mousePixelDeltaY) * pixelHei

	if sb.scrollPosY < -rectRad+scrollBarLen {
		sb.scrollPosY = -rectRad + scrollBarLen
	}
	if sb.scrollPosY > rectRad {
		sb.scrollPosY = rectRad
	}

	viewportOffsetY = sb.scrollPosY / sb.scrollVoidLen * float32(len(document)) * chHei
	fmt.Printf("scrollPosY: %.1f - viewportOffsetY: %.1f", scrollPosY, viewportOffsetY)
}

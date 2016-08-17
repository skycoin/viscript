package main

import (
//"fmt"
)

var scrollPosX float32
var scrollPosY float32 = rectRad

func scroll(mousePixelDeltaY float64) {
	scrollPosY -= float32(mousePixelDeltaY) * pixelHei

	if scrollPosY < -rectRad+scrollBarLen {
		scrollPosY = -rectRad + scrollBarLen
	}
	if scrollPosY > rectRad { // (rectRad-scrollBarLen)
		scrollPosY = rectRad
	}

	viewportOffsetY = scrollPosY / scrollVoidLen * float32(len(document)) * chHei
	fmt.Printf("scrollPosY: %.1f - viewportOffsetY: %.1f", scrollPosY, viewportOffsetY)
}

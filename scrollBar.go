package main

import (
	"fmt"
	"github.com/go-gl/gl/v2.1/gl"
)

/*
mouse position updates use pixels, so the smallest drag motions will be
a jump of at least 1 pixel height.
the ratio of that height / LenOfBar (bar representing the page size)
gives us the jump sizes in scrolling the text
*/

var sb = ScrollBar{0, rectRad, 0, 0}

type ScrollBar struct {
	PosX      float32
	PosY      float32
	LenOfBar  float32
	LenOfVoid float32 // length of the negative space adjacent to bar
}

func (bar *ScrollBar) StartFrame() {
	if /* content smaller than screen */ len(document) <= numYChars {
		// no need for scrollbar
		bar.LenOfBar = 0
		bar.LenOfVoid = rectRad * 2
	} else {
		bar.LenOfBar = float32(numYChars) / float32(len(document)) * rectRad * 2
		bar.LenOfVoid = rectRad*2 - bar.LenOfBar
	}
}

func (bar *ScrollBar) Scroll(mousePixelDeltaY float64) {
	// y increment in gl space
	yInc := float32(mousePixelDeltaY) * pixelHei
	bar.PosY -= yInc

	if bar.PosY < -rectRad+bar.LenOfBar {
		bar.PosY = -rectRad + bar.LenOfBar
	}
	if bar.PosY > rectRad {
		bar.PosY = rectRad
	}

	view.OffsetY += yInc
	//view.OffsetY = bar.PosY / bar.LenOfVoid * float32(len(document)) * chHei
}

func (bar *ScrollBar) DrawVertical(atlasX, atlasY float32) {
	gl.Normal3f(0, 0, 1)
	u := float32(atlasX) * uvSpan
	v := float32(atlasY) * uvSpan

	top := bar.PosY                 //rectRad - 1
	bott := bar.PosY - bar.LenOfBar //-rectRad + 1

	// bottom left   0, 1
	gl.TexCoord2f(u, v+uvSpan)
	gl.Vertex3f(rectRad-chWid, bott, 0)

	// bottom right   1, 1
	gl.TexCoord2f(u+uvSpan, v+uvSpan)
	gl.Vertex3f(rectRad, bott, 0)

	// top right   1, 0
	gl.TexCoord2f(u+uvSpan, v)
	gl.Vertex3f(rectRad, top, 0)

	// top left   0, 0
	gl.TexCoord2f(u, v)
	gl.Vertex3f(rectRad-chWid, top, 0)

	curX += chWid
}

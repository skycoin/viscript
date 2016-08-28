package main

import (
	//	"fmt"
	"github.com/go-gl/gl/v2.1/gl"
)

/*
mouse position updates use pixels, so the smallest drag motions will be
a jump of at least 1 pixel height.
the ratio of that height / LenOfVoid (bar representing the page size),
compared to the void/offscreen length of the text body,
gives us the jump size in scrolling through the text body
*/

type ScrollBar struct {
	PosX      float32
	PosY      float32
	LenOfBar  float32
	LenOfVoid float32 // length of the negative space adjacent to bar
}

func (bar *ScrollBar) StartFrame(tp TextPanel) {
	size /* of screen */ := textRend.ScreenRad * 2

	if /* content smaller than screen */ len(tp.Body) <= tp.NumCharsY {
		// no need for scrollbar
		bar.LenOfBar = 0
		bar.LenOfVoid = size
		code.LenOfOffscreenY = 0
	} else {
		bar.LenOfBar = float32(tp.NumCharsY) / float32(len(tp.Body)) * size
		bar.LenOfVoid = size - bar.LenOfBar
		code.LenOfOffscreenY = float32(len(tp.Body)-tp.NumCharsY) * textRend.chHei
	}
}

func (bar *ScrollBar) Scroll(mousePixelDeltaY float64) {
	// y increment (for bar) in gl space
	yInc := float32(mousePixelDeltaY) * textRend.pixelHei

	bar.PosY -= yInc

	if bar.PosY < -textRend.ScreenRad+bar.LenOfBar {
		bar.PosY = -textRend.ScreenRad + bar.LenOfBar
	}
	if bar.PosY > textRend.ScreenRad {
		bar.PosY = textRend.ScreenRad
	}

	code.OffsetY -= yInc / bar.LenOfVoid * code.LenOfOffscreenY

	if code.OffsetY > 0 {
		code.OffsetY = 0
	}

	if code.OffsetY < -code.LenOfOffscreenY {
		code.OffsetY = -code.LenOfOffscreenY
	}
}

func (bar *ScrollBar) DrawVertical(atlasX, atlasY float32) {
	rad := textRend.ScreenRad
	sp := textRend.UvSpan
	u := float32(atlasX) * sp
	v := float32(atlasY) * sp

	top := bar.PosY                 //rad - 1
	bott := bar.PosY - bar.LenOfBar //-rad + 1

	gl.Normal3f(0, 0, 1)

	// bottom left   0, 1
	gl.TexCoord2f(u, v+sp)
	gl.Vertex3f(rad-textRend.chWid, bott, 0)

	// bottom right   1, 1
	gl.TexCoord2f(u+sp, v+sp)
	gl.Vertex3f(rad, bott, 0)

	// top right   1, 0
	gl.TexCoord2f(u+sp, v)
	gl.Vertex3f(rad, top, 0)

	// top left   0, 0
	gl.TexCoord2f(u, v)
	gl.Vertex3f(rad-textRend.chWid, top, 0)

	textRend.CurrX += textRend.chWid
}

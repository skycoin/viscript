package main

import (
	//	"fmt"
	"github.com/go-gl/gl/v2.1/gl"
)

/*
mouse position updates use pixels, so the smallest drag motions will be
a jump of at least 1 pixel height.
the ratio of that height / LenOfVoid (bar representing the page size),
compared to the void/offscreen length of the document,
gives us the jump size in scrolling through the document
*/

var sb = ScrollBar{0, textRend.ScreenRad, 0, 0}

type ScrollBar struct {
	PosX      float32
	PosY      float32
	LenOfBar  float32
	LenOfVoid float32 // length of the negative space adjacent to bar
}

func (bar *ScrollBar) StartFrame() {
	size /* of screen */ := textRend.ScreenRad * 2

	if /* content smaller than screen */ len(document) <= textRend.NumCharsY {
		// no need for scrollbar
		bar.LenOfBar = 0
		bar.LenOfVoid = size
		code.LenOfOffscreenY = 0
	} else {
		bar.LenOfBar = float32(textRend.NumCharsY) / float32(len(document)) * size
		bar.LenOfVoid = size - bar.LenOfBar
		code.LenOfOffscreenY = float32(len(document)-textRend.NumCharsY) * chHei
	}
}

func (bar *ScrollBar) Scroll(mousePixelDeltaY float64) {
	// y increment (for bar) in gl space
	yInc := float32(mousePixelDeltaY) * pixelHei

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
	gl.Normal3f(0, 0, 1)
	u := float32(atlasX) * textRend.UvSpan
	v := float32(atlasY) * textRend.UvSpan

	top := bar.PosY                 //textRend.ScreenRad - 1
	bott := bar.PosY - bar.LenOfBar //-textRend.ScreenRad + 1

	// bottom left   0, 1
	gl.TexCoord2f(u, v+textRend.UvSpan)
	gl.Vertex3f(textRend.ScreenRad-chWid, bott, 0)

	// bottom right   1, 1
	gl.TexCoord2f(u+textRend.UvSpan, v+textRend.UvSpan)
	gl.Vertex3f(textRend.ScreenRad, bott, 0)

	// top right   1, 0
	gl.TexCoord2f(u+textRend.UvSpan, v)
	gl.Vertex3f(textRend.ScreenRad, top, 0)

	// top left   0, 0
	gl.TexCoord2f(u, v)
	gl.Vertex3f(textRend.ScreenRad-chWid, top, 0)

	textRend.CurrX += chWid
}

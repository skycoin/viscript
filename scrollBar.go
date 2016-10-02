package main

import (
	//"fmt"
	"github.com/go-gl/gl/v2.1/gl"
)

/* TODO:
* genericize draw name
* remove x/y dupes
 */

type ScrollBar struct {
	IsHorizontal   bool
	PosX           float32 // ... position of the grabbable handle
	PosY           float32
	LenOfBar       float32
	LenOfVoid      float32 // length of the negative space representing the length of entire document
	LenOfOffscreen float32
	ScrollDelta    float32 // distance/offset from home/start of document (negative in Y cuz Y increases upwards)
}

func (bar *ScrollBar) UpdateSize(tp *TextPanel) {
	if bar.IsHorizontal {
		// OPTIMIZEME in the future?  idealistically, the below should only be calculated
		// whenever user changes the size of a line such as by:
		// 		typing/deleting/hitting-enter (could be splitting the biggest line in 2)

		// find....
		numCharsInLongest := 0 // ...line of the document
		for _, line := range tp.Body {
			if numCharsInLongest < len(line) {
				numCharsInLongest = len(line)
			}
		}

		// the rest of this block is an altered copy of the else block
		wid := textRend.CharWid * float32(tp.NumCharsX) /* width of panel */

		if /* content smaller than screen */ numCharsInLongest <= tp.NumCharsX {
			// NO BAR
			bar.LenOfBar = 0
			bar.LenOfVoid = wid
			bar.LenOfOffscreen = 0
		} else {
			bar.LenOfBar = float32(tp.NumCharsX) / float32(numCharsInLongest) * wid
			bar.LenOfVoid = wid - bar.LenOfBar
			bar.LenOfOffscreen = float32(numCharsInLongest-tp.NumCharsX) * textRend.CharWid
		}
	} else {
		hei := textRend.CharHei * float32(tp.NumCharsY) /* height of panel */

		if /* content smaller than screen */ len(tp.Body) <= tp.NumCharsY {
			// NO BAR
			bar.LenOfBar = 0
			bar.LenOfVoid = hei
			bar.LenOfOffscreen = 0
		} else {
			bar.LenOfBar = float32(tp.NumCharsY) / float32(len(tp.Body)) * hei
			bar.LenOfVoid = hei - bar.LenOfBar
			bar.LenOfOffscreen = float32(len(tp.Body)-tp.NumCharsY) * textRend.CharHei
		}
	}
}

func (bar *ScrollBar) DragHandleContainsMouseCursor() bool {
	if bar.IsHorizontal {
		// if in the vertical space
		if curs.MouseGlY <= bar.PosY && curs.MouseGlY >= bar.PosY-textRend.CharHei {
			// if in the horizontal space
			if curs.MouseGlX <= bar.PosX+bar.LenOfBar && curs.MouseGlX >= bar.PosX {
				return true
			}
		}
	} else { // vertical bar
		if curs.MouseGlY <= bar.PosY && curs.MouseGlY >= bar.PosY-bar.LenOfBar {
			if curs.MouseGlX <= bar.PosX+textRend.CharWid && curs.MouseGlX >= bar.PosX {
				return true
			}
		}

	}

	return false
}

func (bar *ScrollBar) ScrollThisMuch(tp *TextPanel, delta float32) {
	/*
		mouse position updates use pixels, so the smallest drag motions will be
		a jump of at least 1 pixel height.
		the ratio of that height / LenOfVoid (bar representing the page size),
		compared to the void/offscreen length of the text body,
		gives us the jump size in scrolling through the text body
	*/

	amount := delta / bar.LenOfVoid * bar.LenOfOffscreen

	if bar.IsHorizontal {
		bar.PosX += delta
		bar.PosX = bar.Clamp(bar.PosX, tp.Left, tp.Right-bar.LenOfBar)
		bar.ScrollDelta += amount
		bar.ScrollDelta = bar.Clamp(bar.ScrollDelta, 0, bar.LenOfOffscreen)
	} else {
		bar.PosY -= delta
		bar.PosY = bar.Clamp(bar.PosY, tp.Bottom+bar.LenOfBar, tp.Top)
		bar.ScrollDelta -= amount
		bar.ScrollDelta = bar.Clamp(bar.ScrollDelta, -bar.LenOfOffscreen, 0)
	}
}

// params: position in relevant dimension, negativemost, & positivemost bounds
func (bar *ScrollBar) Clamp(pos, negBoundary, posBoundary float32) float32 {
	if pos < negBoundary {
		pos = negBoundary
	}
	if pos > posBoundary {
		pos = posBoundary
	}

	return pos
}

func (bar *ScrollBar) Draw(atlasX, atlasY float32, tp TextPanel) {
	// this draws all bars whether you can see them or not (not is when content fits entirely inside panel)
	// not worth optimizing, but mentioning it in case of some weird future bug

	sp /* span */ := textRend.UvSpan
	u := float32(atlasX) * sp
	v := float32(atlasY) * sp

	top := tp.Bottom + textRend.CharHei/2
	bott := tp.Bottom                // bottom
	l := tp.Right - textRend.CharWid // left
	r := tp.Right                    // right

	if bar.IsHorizontal {
		l = bar.PosX
		r = bar.PosX + bar.LenOfBar
	} else {
		top = bar.PosY
		bott = bar.PosY - bar.LenOfBar
	}

	gl.Normal3f(0, 0, 1)

	// bottom left   0, 1
	gl.TexCoord2f(u, v+sp)
	gl.Vertex3f(l, bott, 0)

	// bottom right   1, 1
	gl.TexCoord2f(u+sp, v+sp)
	gl.Vertex3f(r, bott, 0)

	// top right   1, 0
	gl.TexCoord2f(u+sp, v)
	gl.Vertex3f(r, top, 0)

	// top left   0, 0
	gl.TexCoord2f(u, v)
	gl.Vertex3f(l, top, 0)
}

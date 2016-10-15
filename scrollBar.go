package main

import (
//"fmt"
//"github.com/go-gl/gl/v2.1/gl"
)

/* TODO:
* genericize draw name
* remove x/y dupes
 */

type ScrollBar struct {
	IsHorizontal   bool
	Thickness      float32
	PosX           float32 // ... position of the grabbable handle
	PosY           float32
	LenOfBar       float32
	LenOfVoid      float32 // length of the negative space representing the length of entire document
	LenOfOffscreen float32 // total hidden offscreen space (bookending the visible portal)
	ScrollDelta    float32 // distance/offset from home/start of document (negative in Y cuz Y increases upwards)
}

func (bar *ScrollBar) UpdateSize(tp TextPanel) {
	if bar.IsHorizontal {
		// OPTIMIZEME in the future?  idealistically, the below should only be calculated
		// whenever user changes the size of a line such as by:
		// 		typing/deleting/hitting-enter (could be splitting the biggest line in 2).
		// but it probably will never matter.

		// find....
		numCharsInLongest := 0 // ...line of the document
		for _, line := range tp.Body {
			if numCharsInLongest < len(line) {
				numCharsInLongest = len(line)
			}
		}
		numCharsInLongest++ // adding extra space so we have room to show cursor at end of longest

		// the rest of this block is an altered copy of the else block
		panWid := tp.Right - tp.Left - tp.BarVert.Thickness // width of panel (MINUS scrollbar space)

		if /* content smaller than screen */ numCharsInLongest <= tp.NumCharsX {
			// NO BAR
			bar.LenOfBar = 0
			bar.LenOfVoid = panWid
			bar.LenOfOffscreen = 0
		} else {
			totalTextWid := float32(numCharsInLongest) * rend.CharWid
			bar.LenOfOffscreen = totalTextWid - panWid
			bar.LenOfBar = panWid / totalTextWid * panWid
			bar.LenOfVoid = panWid - bar.LenOfBar
		}
	} else { // vertical bar
		panHei := tp.Top - tp.Bottom - tp.BarHori.Thickness // height of panel (MINUS scrollbar space)

		if /* content smaller than screen */ len(tp.Body) <= tp.NumCharsY {
			// NO BAR
			bar.LenOfBar = 0
			bar.LenOfVoid = panHei
			bar.LenOfOffscreen = 0
		} else {
			totalTextHei := float32(len(tp.Body)) * rend.CharHei
			bar.LenOfOffscreen = totalTextHei - panHei
			bar.LenOfBar = panHei / totalTextHei * panHei
			bar.LenOfVoid = panHei - bar.LenOfBar
		}
	}

	if bar.Thickness == 0 {
		bar.Thickness = rend.ScreenRad / 30

		// FIXME (hack, because with current projection a unit span in x does not look the same as in y)
		if bar.IsHorizontal {
			bar.Thickness *= 1.4

			bar.PosY = tp.Bottom + bar.Thickness
		} else {
			bar.PosX = tp.Right - bar.Thickness
		}
	}
}

func (bar *ScrollBar) ContainsMouseCursor(tp *TextPanel) bool {
	if bar.IsHorizontal {
		if curs.MouseGlY <= bar.PosY {
			return true
		}
	} else { // vertical bar
		if curs.MouseGlX >= bar.PosX {
			return true
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
		bar.PosX = Clamp(bar.PosX, tp.Left, tp.Right-bar.LenOfBar-tp.BarVert.Thickness)
		bar.ScrollDelta += amount
		bar.ScrollDelta = Clamp(bar.ScrollDelta, 0, bar.LenOfOffscreen)
	} else {
		bar.PosY -= delta
		bar.PosY = Clamp(bar.PosY, tp.Bottom+bar.LenOfBar+tp.BarHori.Thickness, tp.Top)
		bar.ScrollDelta -= amount
		bar.ScrollDelta = Clamp(bar.ScrollDelta, -bar.LenOfOffscreen, 0)
	}
}

func (bar *ScrollBar) Draw(atlasX, atlasY float32, tp TextPanel) {
	// this draws all bars whether you can see them or not (not is when content fits entirely inside panel)
	// not worth optimizing, but mentioning it in case of some weird future bug

	bar.UpdateSize(tp)

	l := bar.PosX
	t := bar.PosY
	b := tp.Bottom
	r := tp.Right

	// potential future FIXME:
	// with the current "double line" graphic used this doesn't matter,
	// but the .IsVertical case should copy the horizontal case,
	// in order to prevent unwanted texture stretching
	if bar.IsHorizontal {
		th := tp.BarVert.Thickness
		if th < 1 { // at some point it is, which would cause infinite loop
			th = bar.Thickness
		}

		r = l + th
		max := bar.PosX + bar.LenOfBar

		for l < max {
			if r > max {
				r = max
			}

			rend.DrawQuad(atlasX, atlasY, &Rectangle{t, r, b, l})

			l += th
			r += th
		}
	} else {
		b = bar.PosY - bar.LenOfBar

		rend.DrawQuad(atlasX, atlasY, &Rectangle{t, r, b, l})
	}
}

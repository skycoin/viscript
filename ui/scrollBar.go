package ui

import (
	//"fmt"
	"github.com/corpusc/viscript/common"
)

var ScrollBarThickness float32 = 6 / 25 // FUTURE FIXME: base this on screen size, so UHD+ users get thick enough bars

type ScrollBar struct {
	IsHorizontal   bool
	PosX           float32 // ... position of the grabbable handle
	PosY           float32
	LenOfBar       float32
	LenOfVoid      float32 // length of the negative space representing the length of entire document
	LenOfOffscreen float32 // total hidden offscreen space (bookending the visible portal)
	ScrollDelta    float32 // distance/offset from home/start of document (negative in Y cuz Y increases upwards)
	Rect           *common.Rectangle
}

func (bar *ScrollBar) Scroll(delta float32) {
	/*
		mouse position updates use pixels, so the smallest drag motions will be
		a jump of at least 1 pixel height.
		the ratio of that height / LenOfVoid (bar representing the page size),
		compared to the void/offscreen length of the text body,
		gives us the jump size in scrolling through the text body
	*/

	amount := delta / bar.LenOfVoid * bar.LenOfOffscreen

	if bar.IsHorizontal {
		//bar.PosX += delta
		//bar.PosX = Clamp(bar.PosX, panel.Left, panel.Right-bar.LenOfBar-ScrollBarThickness)
		bar.ScrollDelta += amount
		bar.ClampX()
	} else { // vertical
		//bar.PosY -= delta
		//bar.PosY = Clamp(bar.PosY, panel.Bottom+bar.LenOfBar+ScrollBarThickness, panel.Top)
		bar.ScrollDelta -= amount
		bar.ClampY()
	}
}

func (bar *ScrollBar) ClampX() {
	bar.ScrollDelta = common.Clamp(bar.ScrollDelta, 0, bar.LenOfOffscreen)
}

func (bar *ScrollBar) ClampY() {
	bar.ScrollDelta = common.Clamp(bar.ScrollDelta, -bar.LenOfOffscreen, 0)
}

func (bar *ScrollBar) SetSize(panel *common.Rectangle, body []string, charWid, charHei float32) {
	if bar.IsHorizontal {
		// OPTIMIZEME in the future?  idealistically, the below should only be calculated
		// whenever user changes the size of a line, such as by:
		// 		typing/deleting/hitting-enter (could be splitting the biggest line in 2).
		// or changes client window size, or panel size
		// ....but it probably will never really matter.

		// find....
		numCharsInLongest := 0 // ...line of the document
		for _, line := range body {
			if numCharsInLongest < len(line) {
				numCharsInLongest = len(line)
			}
		}
		numCharsInLongest++ // adding extra space so we have room to show cursor at end of longest

		// the rest of this block is an altered copy of the else block
		panWid := panel.Width() - ScrollBarThickness // width of panel (MINUS scrollbar space)

		/* if content smaller than panel width */
		if float32(numCharsInLongest)*charWid <= panel.Width()-ScrollBarThickness {
			// NO BAR
			bar.LenOfBar = 0
			bar.LenOfVoid = panWid
			bar.LenOfOffscreen = 0
		} else {
			totalTextWid := float32(numCharsInLongest) * charWid
			bar.LenOfOffscreen = totalTextWid - panWid
			bar.LenOfBar = panWid / totalTextWid * panWid
			bar.LenOfVoid = panWid - bar.LenOfBar
			bar.PosX = panel.Left + bar.ScrollDelta/bar.LenOfOffscreen*bar.LenOfVoid
			bar.ClampX() // OPTIMIZEME: only do when app resized
		}

		bar.PosY = panel.Bottom + ScrollBarThickness
	} else { // vertical bar
		panHei := panel.Height() - ScrollBarThickness // height of panel (MINUS scrollbar space)

		/* if content smaller than panel height */
		if float32(len(body))*charHei <= panel.Height()-ScrollBarThickness {
			// NO BAR
			bar.LenOfBar = 0
			bar.LenOfVoid = panHei
			bar.LenOfOffscreen = 0
		} else {
			totalTextHei := float32(len(body)) * charHei
			bar.LenOfOffscreen = totalTextHei - panHei
			bar.LenOfBar = panHei / totalTextHei * panHei
			bar.LenOfVoid = panHei - bar.LenOfBar
			bar.PosY = panel.Top + bar.ScrollDelta/bar.LenOfOffscreen*bar.LenOfVoid
			bar.ClampY() // OPTIMIZEME: only do when app resized
		}

		bar.PosX = panel.Right - ScrollBarThickness
	}

	// setup the final rectangle for drawing
	l := bar.PosX
	t := bar.PosY
	b := panel.Bottom
	r := panel.Right

	if bar.IsHorizontal {
		r = l + bar.LenOfBar
	} else {
		b = t - bar.LenOfBar
	}

	bar.Rect = &common.Rectangle{t, r, b, l}
}

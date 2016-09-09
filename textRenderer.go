/* TODO:

* KEY-BASED NAVIGATION (CTRL-HOME/END - PGUP/DN)
* ??DONE?? BACKSPACE/DELETE at the ends of lines
	pulls us up to prev line, or pulls up next line
* HORIZONTAL SCROLLBARS
	horizontal could be charHei thickness
	vertical could easily be a smaller rendering of the first ~40 chars?
		however not if we map the whole vertical space (when scrollspace is taller than screen),
		because this requires scaling the text.  and keeping the aspect ratio means ~40 (max)
		would alter the width of the scrollbar
*/

package main

import (
	//"fmt"
	"github.com/go-gl/gl/v2.1/gl"
)

type TextRenderer struct {
	pixelWid      float32
	pixelHei      float32
	chWid         float32
	chHei         float32
	chWidInPixels int
	chHeiInPixels int
	UvSpan        float32 // looking into 16/16 atlas/grid of character tiles
	ScreenRad     float32 // entire screen radius (distance to edge
	// in the cardinal directions from the center, corners would be farther away)
	MaxCharsX int
	MaxCharsY int
	// current position renderer draws to
	CurrX   float32
	CurrY   float32
	Focused TextPanel
	Panels  []TextPanel
}

func (tr *TextRenderer) Init() {
	tr.UvSpan = float32(1.0) / 16
	tr.ScreenRad = float32(3)
	tr.MaxCharsX = 80
	tr.MaxCharsY = 25
	tr.pixelWid = tr.ScreenRad * 2 / float32(resX)
	tr.pixelHei = tr.ScreenRad * 2 / float32(resY)
	tr.chWid = float32(tr.ScreenRad * 2 / float32(tr.MaxCharsX))
	tr.chHei = float32(tr.ScreenRad * 2 / float32(tr.MaxCharsY))
	tr.chWidInPixels = int(float32(resX) / float32(tr.MaxCharsX))
	tr.chHeiInPixels = int(float32(resY) / float32(tr.MaxCharsY))
}

func (tr *TextRenderer) DrawAll() {
	code.Draw()
	cons.Draw()

	for _, pan := range tr.Panels {
		pan.Draw()
	}

	curs.Draw()
}

func drawCurrentChar(char rune) {
	x := int(char) % 16
	y := int(char) / 16
	w := textRend.chWid // char width
	h := textRend.chHei // char height
	sp := textRend.UvSpan

	gl.Normal3f(0, 0, 1)

	gl.TexCoord2f(float32(x)*sp, float32(y)*sp+sp) // bl  0, 1
	gl.Vertex3f(textRend.CurrX, textRend.CurrY-h, 0)

	gl.TexCoord2f(float32(x)*sp+sp, float32(y)*sp+sp) // br  1, 1
	gl.Vertex3f(textRend.CurrX+w, textRend.CurrY-h, 0)

	gl.TexCoord2f(float32(x)*sp+sp, float32(y)*sp) // tr  1, 0
	gl.Vertex3f(textRend.CurrX+w, textRend.CurrY, 0)

	gl.TexCoord2f(float32(x)*sp, float32(y)*sp) // tl  0, 0
	gl.Vertex3f(textRend.CurrX, textRend.CurrY, 0)

	textRend.CurrX += w
}

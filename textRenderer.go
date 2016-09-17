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

var code = TextPanel{NumCharsY: 14}
var cons = TextPanel{NumCharsY: 10} // console (runtime feedback log)

type TextRenderer struct {
	PixelWid        float32
	PixelHei        float32
	CharWid         float32
	CharHei         float32
	CharWidInPixels int
	CharHeiInPixels int
	UvSpan          float32 // looking into 16/16 atlas/grid of character tiles
	ScreenRad       float32 // entire screen radius (distance to edge
	// in the cardinal directions from the center, corners would be farther away)
	MaxCharsX int
	MaxCharsY int
	// current position renderer draws to
	CurrX   float32
	CurrY   float32
	Focused *TextPanel
	Panels  []TextPanel
}

func (tr *TextRenderer) Init() {
	tr.UvSpan = float32(1.0) / 16
	tr.ScreenRad = float32(3)
	tr.MaxCharsX = 80
	tr.MaxCharsY = 25
	tr.PixelWid = tr.ScreenRad * 2 / float32(resX)
	tr.PixelHei = tr.ScreenRad * 2 / float32(resY)
	tr.CharWid = float32(tr.ScreenRad * 2 / float32(tr.MaxCharsX))
	tr.CharHei = float32(tr.ScreenRad * 2 / float32(tr.MaxCharsY))
	tr.CharWidInPixels = int(float32(resX) / float32(tr.MaxCharsX))
	tr.CharHeiInPixels = int(float32(resY) / float32(tr.MaxCharsY))

	cons.Top = tr.ScreenRad - float32(code.NumCharsY+1)*tr.CharHei
	cons.Init()
	code.Init()
	code.SetupDemoProgram()
	tr.Focused = &code
	tr.Panels = append(tr.Panels, code)
	tr.Panels = append(tr.Panels, cons)
}

func (tr *TextRenderer) DrawAll() {
	for _, pan := range tr.Panels {
		pan.Draw()
	}

	curs.Draw()
}

func (tr *TextRenderer) ScrollFocusedPanel(mousePixelDeltaY float64) {
	for _, pan := range tr.Panels {
		pan.ScrollIfMouseOver(mousePixelDeltaY)
	}
}

func drawCurrentChar(char rune, clipSpan float32) {
	u := float32(int(char) % 16)
	v := float32(int(char) / 16)
	w := textRend.CharWid // char width
	h := textRend.CharHei // char height
	sp := textRend.UvSpan

	gl.Normal3f(0, 0, 1)

	gl.TexCoord2f(u*sp, v*sp+sp) // bl  0, 1
	gl.Vertex3f(textRend.CurrX, textRend.CurrY-h, 0)

	gl.TexCoord2f(u*sp+sp, v*sp+sp) // br  1, 1
	gl.Vertex3f(textRend.CurrX+w, textRend.CurrY-h, 0)

	gl.TexCoord2f(u*sp+sp, v*sp) // tr  1, 0
	gl.Vertex3f(textRend.CurrX+w, textRend.CurrY-clipSpan, 0)

	gl.TexCoord2f(u*sp, v*sp) // tl  0, 0
	gl.Vertex3f(textRend.CurrX, textRend.CurrY-clipSpan, 0)

	textRend.CurrX += w
}

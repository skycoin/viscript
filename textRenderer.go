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
	PixelWid        float32
	PixelHei        float32
	CharWid         float32
	CharHei         float32
	CharWidInPixels int
	CharHeiInPixels int
	UvSpan          float32 // looking into 16/16 atlas/grid of character tiles
	ScreenRad       float32 // entire screen radius (distance to edge
	// in the cardinal directions from the center, corners would be farther away)
	MaxCharsX int // this is used to give us proportions like an 80x25 text console screen, from a 3f by 3f gl space
	MaxCharsY int
	// current position renderer draws to
	CurrX   float32
	CurrY   float32
	Focused *TextPanel
	Panels  []*TextPanel
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

	tr.Panels = append(tr.Panels, &TextPanel{NumCharsY: 14, IsEditable: true})
	tr.Panels = append(tr.Panels, &TextPanel{NumCharsY: 10, IsEditable: true}) // console (runtime feedback log)	// FIXME so its not editable once we're done debugging some things
	tr.Focused = tr.Panels[0]

	tr.Panels[0].Init()
	tr.Panels[0].SetupDemoProgram()
	tr.Panels[1].Top = tr.ScreenRad - float32(tr.Focused.NumCharsY+1)*tr.CharHei
	tr.Panels[1].Init()
}

func (tr *TextRenderer) DrawAll() {
	curs.Update()

	for _, pan := range tr.Panels {
		pan.Draw()
	}
}

func (tr *TextRenderer) ScrollPanelThatIsHoveredOver(mousePixelDeltaX, mousePixelDeltaY float64) {
	for _, pan := range tr.Panels {
		pan.ScrollIfMouseOver(mousePixelDeltaX, mousePixelDeltaY)
	}
}

func (tr *TextRenderer) DrawCharAtRect(char rune, r *Rectangle) {
	u := float32(int(char) % 16)
	v := float32(int(char) / 16)
	sp := textRend.UvSpan

	gl.Normal3f(0, 0, 1)

	gl.TexCoord2f(u*sp, v*sp+sp)
	gl.Vertex3f(r.Left, r.Bottom, 0)

	gl.TexCoord2f(u*sp+sp, v*sp+sp)
	gl.Vertex3f(r.Right, r.Bottom, 0)

	gl.TexCoord2f(u*sp+sp, v*sp)
	gl.Vertex3f(r.Right, r.Top, 0)

	gl.TexCoord2f(u*sp, v*sp)
	gl.Vertex3f(r.Left, r.Top, 0)
}

func (tr *TextRenderer) DrawQuad(atlasX, atlasY float32, r *Rectangle) {
	sp /* span */ := textRend.UvSpan
	u := float32(atlasX) * sp
	v := float32(atlasY) * sp

	gl.Normal3f(0, 0, 1)

	gl.TexCoord2f(u, v+sp)
	gl.Vertex3f(r.Left, r.Bottom, 0)

	gl.TexCoord2f(u+sp, v+sp)
	gl.Vertex3f(r.Right, r.Bottom, 0)

	gl.TexCoord2f(u+sp, v)
	gl.Vertex3f(r.Right, r.Top, 0)

	gl.TexCoord2f(u, v)
	gl.Vertex3f(r.Left, r.Top, 0)
}

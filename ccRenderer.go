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

var rend = CcRenderer{}

type CcRenderer struct {
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

func (cr *CcRenderer) Init() {
	cr.UvSpan = float32(1.0) / 16
	cr.ScreenRad = float32(3)
	cr.MaxCharsX = 80
	cr.MaxCharsY = 25
	cr.PixelWid = cr.ScreenRad * 2 / float32(resX)
	cr.PixelHei = cr.ScreenRad * 2 / float32(resY)
	cr.CharWid = float32(cr.ScreenRad * 2 / float32(cr.MaxCharsX))
	cr.CharHei = float32(cr.ScreenRad * 2 / float32(cr.MaxCharsY))
	cr.CharWidInPixels = int(float32(resX) / float32(cr.MaxCharsX))
	cr.CharHeiInPixels = int(float32(resY) / float32(cr.MaxCharsY))

	cr.Panels = append(cr.Panels, &TextPanel{NumCharsY: 14, IsEditable: true})
	cr.Panels = append(cr.Panels, &TextPanel{NumCharsY: 10, IsEditable: true}) // console (runtime feedback log)	// FIXME so its not editable once we're done debugging some things
	cr.Focused = cr.Panels[0]

	cr.Panels[0].Init()
	cr.Panels[0].SetupDemoProgram()
	cr.Panels[1].Top = cr.ScreenRad - float32(cr.Focused.NumCharsY+1)*cr.CharHei
	cr.Panels[1].Init()
}

func (cr *CcRenderer) DrawAll() {
	curs.Update()
	menu.Draw()

	for _, pan := range cr.Panels {
		pan.Draw()
	}
}

func (cr *CcRenderer) ScrollPanelThatIsHoveredOver(mousePixelDeltaX, mousePixelDeltaY float64) {
	for _, pan := range cr.Panels {
		pan.ScrollIfMouseOver(mousePixelDeltaX, mousePixelDeltaY)
	}
}

func (cr *CcRenderer) DrawCharAtRect(char rune, r *Rectangle) {
	u := float32(int(char) % 16)
	v := float32(int(char) / 16)
	sp := rend.UvSpan

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

func (cr *CcRenderer) DrawQuad(atlasX, atlasY float32, r *Rectangle) {
	sp /* span */ := rend.UvSpan
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

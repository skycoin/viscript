/*
--- TODO: ---

* KEY-BASED NAVIGATION (CTRL-HOME/END - PGUP/DN)
* BACKSPACE/DELETE at the ends of lines
	pulls us up to prev line, or
	pulls up next line


--- OPTIONAL NICETIES: ---

* HORIZONTAL SCROLLBARS
	horizontal could be charHei thickness
	vertical could easily be a smaller rendering of the first ~40 chars?
		however not if we map the whole vertical space (when scrollspace is taller than screen),
		because this requires scaling the text.  and keeping the aspect ratio means ~40 (max)
		would alter the width of the scrollbar
*/

package gfx

import (
	"fmt"
	"github.com/corpusc/viscript/common"
	"github.com/corpusc/viscript/ui"
	"github.com/go-gl/gl/v2.1/gl"
)

var Rend = CcRenderer{}

var goldenRatio = 1.61803398875
var goldenFraction = float32(goldenRatio / (goldenRatio + 1))

var Black = []float32{0, 0, 0, 1}
var Blue = []float32{0, 0, 1, 1}
var Cyan = []float32{0, 0.5, 1, 1}
var Fuschia = []float32{0.6, 0.2, 0.3, 1}
var Gray = []float32{0.25, 0.25, 0.25, 1}
var GrayDark = []float32{0.15, 0.15, 0.15, 1}
var GrayLight = []float32{0.4, 0.4, 0.4, 1}
var Green = []float32{0, 1, 0, 1}
var Magenta = []float32{1, 0, 1, 1}
var Orange = []float32{0.8, 0.35, 0, 1}
var Purple = []float32{0.6, 0, 0.8, 1}
var Red = []float32{1, 0, 0, 1}
var Tan = []float32{0.55, 0.47, 0.37, 1}
var Violet = []float32{0.4, 0.2, 1, 1}
var White = []float32{1, 1, 1, 1}
var Yellow = []float32{1, 1, 0, 1}

// dimensions (in pixel units)
var InitAppWidth int32 = 800
var InitAppHeight int32 = 600
var CurrAppWidth int32 = InitAppWidth
var CurrAppHeight int32 = InitAppHeight
var longerDimension = float32(InitAppWidth) / float32(InitAppHeight)
var InitFrustum = &common.Rectangle{1, longerDimension, -1, -longerDimension}
var PrevFrustum = &common.Rectangle{InitFrustum.Top, InitFrustum.Right, InitFrustum.Bottom, InitFrustum.Left}
var CurrFrustum = &common.Rectangle{InitFrustum.Top, InitFrustum.Right, InitFrustum.Bottom, InitFrustum.Left}

func init() {
	fmt.Println("gfx.init()")

	// one-time setup
	Rend.MaxCharsX = 80
	Rend.MaxCharsY = 25
	Rend.DistanceFromOrigin = 3
	Rend.UvSpan = float32(1.0) / 16 // how much uv a pixel spans
	Rend.RunPanelHeiPerc = 0.4
	Rend.PrevColor = GrayDark
	Rend.CurrColor = GrayDark

	// things to resize later
	Rend.ClientExtentX = Rend.DistanceFromOrigin * longerDimension
	Rend.ClientExtentY = Rend.DistanceFromOrigin
	Rend.CharWid = float32(Rend.ClientExtentX*2) / float32(Rend.MaxCharsX)
	Rend.CharHei = float32(Rend.ClientExtentY*2) / float32(Rend.MaxCharsY)
	Rend.CharWidInPixels = int(float32(CurrAppWidth) / float32(Rend.MaxCharsX))
	Rend.CharHeiInPixels = int(float32(CurrAppHeight) / float32(Rend.MaxCharsY))
	Rend.PixelWid = Rend.ClientExtentX * 2 / float32(CurrAppWidth)
	Rend.PixelHei = Rend.ClientExtentY * 2 / float32(CurrAppHeight)

	// one-time setup of panels
	Rend.Panels = append(Rend.Panels, &TextPanel{BandPercent: 1 - Rend.RunPanelHeiPerc, IsEditable: true})
	Rend.Panels = append(Rend.Panels, &TextPanel{BandPercent: Rend.RunPanelHeiPerc, IsEditable: true}) // console (runtime feedback log)	// FIXME so its not editable once we're done debugging some things
	Rend.Focused = Rend.Panels[0]

	Rend.Panels[0].Init()
	Rend.Panels[0].SetupDemoProgram()
	Rend.Panels[1].Init()

	ui.MainMenu.SetSize(Rend.GetMenuSizedRect())
}

type CcRenderer struct {
	DistanceFromOrigin float32
	RunPanelHeiPerc    float32 // FIXME: hardwired value for a specific use case
	ClientExtentX      float32 // distance from the center to an edge of the app's root/client area
	ClientExtentY      float32
	// ....in the cardinal directions from the center, corners would be farther away)
	PixelWid        float32
	PixelHei        float32
	CharWid         float32
	CharHei         float32
	CharWidInPixels int
	CharHeiInPixels int
	UvSpan          float32 // looking into 16/16 atlas/grid of character tiles
	// FIXME: below is no longer a maximum of what fits on a max-sized panel (taking up the whole app window) anymore.
	// 		but is still used as a guide for sizes
	MaxCharsX int // this is used to give us proportions like an 80x25 text console screen, ....
	MaxCharsY int // ....from a cr.DistanceFromOrigin*2-by-cr.DistanceFromOrigin*2 gl space
	// current position renderer draws to
	CurrX     float32
	CurrY     float32
	PrevColor []float32 // previous
	CurrColor []float32
	Focused   *TextPanel
	Panels    []*TextPanel
}

func (cr *CcRenderer) SetSize() {
	fmt.Printf("CcRenderer.SetSize() - ClientExtentX: %.2f\n", cr.ClientExtentX)
	//cr.ClientExtentX = cr.DistanceFromOrigin * (CurrFrustum.Right / InitFrustum.Right)
	//cr.ClientExtentY = cr.DistanceFromOrigin * (CurrFrustum.Top / InitFrustum.Top)
	*PrevFrustum = *CurrFrustum

	CurrFrustum.Right = float32(CurrAppWidth) / float32(InitAppWidth) * InitFrustum.Right
	CurrFrustum.Left = -CurrFrustum.Right
	CurrFrustum.Top = float32(CurrAppHeight) / float32(InitAppHeight) * InitFrustum.Top
	CurrFrustum.Bottom = -CurrFrustum.Top

	fmt.Printf("CcRenderer.SetSize() - PrevFrustum.Left: %.3f\n", PrevFrustum.Left)
	fmt.Printf("CcRenderer.SetSize() - CurrFrustum.Left: %.3f\n", CurrFrustum.Left)

	cr.ClientExtentX = cr.DistanceFromOrigin * CurrFrustum.Right
	cr.ClientExtentY = cr.DistanceFromOrigin * CurrFrustum.Top

	// things that weren't initialized in this func
	ui.MainMenu.SetSize(cr.GetMenuSizedRect())

	for _, pan := range cr.Panels {
		pan.SetSize()
	}
}

func (cr *CcRenderer) ScrollPanelThatIsHoveredOver(mousePixelDeltaX, mousePixelDeltaY float64) {
	for _, pan := range cr.Panels {
		pan.ScrollIfMouseOver(mousePixelDeltaX, mousePixelDeltaY)
	}
}

func (cr *CcRenderer) GetMenuSizedRect() *common.Rectangle {
	return &common.Rectangle{
		Rend.ClientExtentY,
		Rend.ClientExtentX,
		Rend.ClientExtentY - Rend.CharHei,
		-Rend.ClientExtentX}
}

func (cr *CcRenderer) Color(newColor []float32) {
	cr.PrevColor = cr.CurrColor
	cr.CurrColor = newColor
	gl.Materialfv(gl.FRONT, gl.AMBIENT_AND_DIFFUSE, &newColor[0])
}

func (cr *CcRenderer) DrawAll() {
	Curs.Update()
	cr.DrawMenu()

	for _, pan := range cr.Panels {
		pan.Draw()
	}

	// syntax tree
	if ui.MainMenu.ButtonActivated("Syntax Tree") {
		span := float32(0.8)
		x := -span / 2
		y := ui.MainMenu.Rect.Bottom - 0.2
		r := &common.Rectangle{y, x + span, y - span, x}
		Rend.DrawStretchableRect(11, 13, r)
		Rend.DrawTextInRect("mainBlock", r)
	}

	// // 'crosshair' center indicator
	//var f float32 = Rend.CharHei
	//Rend.DrawCharAtRect('+', &common.Rectangle{f, f, -f, -f})
}

func (cr *CcRenderer) DrawMenu() {
	for _, bu := range ui.MainMenu.Buttons {
		if bu.Activated {
			Rend.Color(Green)
		} else {
			Rend.Color(White)
		}

		Rend.DrawStretchableRect(11, 13, bu.Rect)
		Rend.DrawTextInRect(bu.Name, bu.Rect)
	}
}

func (cr *CcRenderer) DrawTextInRect(s string, r *common.Rectangle) {
	h := r.Height() * goldenFraction   // height of chars
	w := h                             // width of chars (same as height, or else squished to fit rect)
	glTextWidth := float32(len(s)) * w // in terms of OpenGL/float32 space
	lipSpan := (r.Height() - h) / 2    // lip/frame/edge span
	maxW := r.Width() - lipSpan*2      // maximum width for text, which leaves a edge/lip/frame margin

	if glTextWidth > maxW {
		glTextWidth = maxW
		w = maxW / float32(len(s))
	}

	x := r.Left + (r.Width()-glTextWidth)/2

	for _, c := range s {
		Rend.DrawCharAtRect(c, &common.Rectangle{r.Top - lipSpan, x + w, r.Bottom + lipSpan, x})
		x += w
	}
}

func (cr *CcRenderer) DrawCharAtRect(char rune, r *common.Rectangle) {
	u := float32(int(char) % 16)
	v := float32(int(char) / 16)
	sp := Rend.UvSpan

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

func (cr *CcRenderer) DrawQuad(atlasX, atlasY float32, r *common.Rectangle) {
	sp /* span */ := Rend.UvSpan
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

func (cr *CcRenderer) DrawStretchableRect(atlasX, atlasY float32, r *common.Rectangle) {
	// (sometimes called 9 Slicing)
	// draw 9 quads which keep a predictable frame/margin/edge undistorted,
	// while stretching the middle to fit the desired space

	w := r.Width()
	h := r.Height()

	// skip invisible or inverted rects
	if w <= 0 || h <= 0 {
		return
	}

	//var uvEdgeFraction float32 = 0.125 // 1/8
	var uvEdgeFraction float32 = 0.125 / 2 // 1/16
	// we're gonna draw from top to bottom (positivemost to negativemost)

	sp /* span */ := Rend.UvSpan
	u := float32(atlasX) * sp
	v := float32(atlasY) * sp

	gl.Normal3f(0, 0, 1)

	// setup the 4 lines needed (for 3 spanning sections)
	uSpots := []float32{}
	uSpots = append(uSpots, (u))
	uSpots = append(uSpots, (u)+sp*uvEdgeFraction)
	uSpots = append(uSpots, (u+sp)-sp*uvEdgeFraction)
	uSpots = append(uSpots, (u + sp))

	vSpots := []float32{}
	vSpots = append(vSpots, (v))
	vSpots = append(vSpots, (v)+sp*uvEdgeFraction)
	vSpots = append(vSpots, (v+sp)-sp*uvEdgeFraction)
	vSpots = append(vSpots, (v + sp))

	edgeSpan := Rend.PixelWid * 4
	if edgeSpan > w/2 {
		edgeSpan = w / 2
	}

	xSpots := []float32{}
	xSpots = append(xSpots, r.Left)
	xSpots = append(xSpots, r.Left+edgeSpan)
	xSpots = append(xSpots, r.Right-edgeSpan)
	xSpots = append(xSpots, r.Right)

	edgeSpan = Rend.PixelHei * 4
	if edgeSpan > h/2 {
		edgeSpan = h / 2
	}

	ySpots := []float32{}
	ySpots = append(ySpots, r.Top)
	ySpots = append(ySpots, r.Top-edgeSpan)
	ySpots = append(ySpots, r.Bottom+edgeSpan)
	ySpots = append(ySpots, r.Bottom)

	if ySpots[1] > ySpots[0] {
		ySpots[1] = ySpots[0]
	}

	for iX := 0; iX < 3; iX++ {
		for iY := 0; iY < 3; iY++ {
			// draw 1 of 9 rects

			if false { //iX == 1 && iY == 1 {
			} else {
				gl.TexCoord2f(uSpots[iX], vSpots[iY+1]) // left bottom
				gl.Vertex3f(xSpots[iX], ySpots[iY+1], 0)

				gl.TexCoord2f(uSpots[iX+1], vSpots[iY+1]) // right bottom
				gl.Vertex3f(xSpots[iX+1], ySpots[iY+1], 0)

				gl.TexCoord2f(uSpots[iX+1], vSpots[iY]) // right top
				gl.Vertex3f(xSpots[iX+1], ySpots[iY], 0)

				gl.TexCoord2f(uSpots[iX], vSpots[iY]) // left top
				gl.Vertex3f(xSpots[iX], ySpots[iY], 0)
			}
		}
	}
}

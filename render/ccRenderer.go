/*
--- TODO: ---

* KEY-BASED NAVIGATION (CTRL-HOME/END - PGUP/DN)
* BACKSPACE/DELETE at the ends of lines
	pulls us up to prev line, or pulls up next line


--- OPTIONAL NICETIES: ---

* HORIZONTAL SCROLLBARS
	horizontal could be charHei thickness
	vertical could easily be a smaller rendering of the first ~40 chars?
		however not if we map the whole vertical space (when scrollspace is taller than screen),
		because this requires scaling the text.  and keeping the aspect ratio means ~40 (max)
		would alter the width of the scrollbar
*/

package render

import (
	"fmt"
	"github.com/corpusc/viscript/common"
	"github.com/go-gl/gl/v2.1/gl"
)

var Rend = CcRenderer{}

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
	fmt.Println("package render - init()")

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
	MenuInst.SetSize()

	for _, pan := range cr.Panels {
		pan.SetSize()
	}
}

func (cr *CcRenderer) Color(newColor []float32) {
	cr.PrevColor = cr.CurrColor
	cr.CurrColor = newColor
	gl.Materialfv(gl.FRONT, gl.AMBIENT_AND_DIFFUSE, &newColor[0])
}

var bu *Button = &Button{}
var b2 *Button = &Button{}

func (cr *CcRenderer) DrawAll() {
	Curs.Update()
	MenuInst.Draw()

	for _, pan := range cr.Panels {
		pan.Draw()
	}

	// show width of client area
	if bu.Rect == nil {
		var f float32 = 1.3
		bu.Rect = &common.Rectangle{Rend.CharHei + 0.4, f, Rend.CharHei + 0.1, -f} // OPTIMIZEME? is this causing slow GC thrashing?
		b2.Rect = &common.Rectangle{Rend.CharHei + 0.8, f, Rend.CharHei + 0.5, -f} // OPTIMIZEME? is this causing slow GC thrashing?
	}
	bu.Name = fmt.Sprintf("extentX: %.2f", cr.ClientExtentX)
	bu.Draw()
	b2.Name = fmt.Sprintf("extentY: %.2f", cr.ClientExtentY)
	b2.Draw()

	// 'crosshair' center indicator
	var f float32 = Rend.CharHei
	Rend.DrawCharAtRect('+', &common.Rectangle{f, f, -f, -f})
}

func (cr *CcRenderer) ScrollPanelThatIsHoveredOver(mousePixelDeltaX, mousePixelDeltaY float64) {
	for _, pan := range cr.Panels {
		pan.ScrollIfMouseOver(mousePixelDeltaX, mousePixelDeltaY)
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

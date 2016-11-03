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
	"github.com/go-gl/gl/v2.1/gl"
	"viscript/common"
)

var rend = CcRenderer{}
var black = []float32{0, 0, 0, 1}
var blue = []float32{0, 0, 1, 1}
var cyan = []float32{0, 0.5, 1, 1}
var fuschia = []float32{0.6, 0.2, 0.3, 1}
var gray = []float32{0.25, 0.25, 0.25, 1}
var grayDark = []float32{0.15, 0.15, 0.15, 1}
var grayLight = []float32{0.4, 0.4, 0.4, 1}
var green = []float32{0, 1, 0, 1}
var magenta = []float32{1, 0, 1, 1}
var orange = []float32{0.8, 0.35, 0, 1}
var purple = []float32{0.6, 0, 0.8, 1}
var red = []float32{1, 0, 0, 1}
var tan = []float32{0.55, 0.47, 0.37, 1}
var violet = []float32{0.4, 0.2, 1, 1}
var white = []float32{1, 1, 1, 1}
var yellow = []float32{1, 1, 0, 1}

// dimensions (in pixel units)
var initAppWidth int32 = 800
var initAppHeight int32 = 600
var currAppWidth int32 = initAppWidth
var currAppHeight int32 = initAppHeight
var longerDimension = float32(initAppWidth) / float32(initAppHeight)
var initFrustum = &common.Rectangle{1, longerDimension, -1, -longerDimension}
var prevFrustum = &common.Rectangle{initFrustum.Top, initFrustum.Right, initFrustum.Bottom, initFrustum.Left}
var currFrustum = &common.Rectangle{initFrustum.Top, initFrustum.Right, initFrustum.Bottom, initFrustum.Left}

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

func (cr *CcRenderer) Init() { // FIXME: part should be split into something like .Setup/Set___()
	if cr.ClientExtentX == 0.0 || cr.ClientExtentY == 0.0 {
		cr.MaxCharsX = 80
		cr.MaxCharsY = 25
		cr.DistanceFromOrigin = 3
		cr.UvSpan = float32(1.0) / 16 // how much uv a pixel spans
		cr.RunPanelHeiPerc = 0.4
		cr.PrevColor = grayDark
		cr.CurrColor = grayDark

		fmt.Println("CcRenderer.Init(): FIRST TIME")
		cr.ClientExtentX = cr.DistanceFromOrigin * longerDimension
		cr.ClientExtentY = cr.DistanceFromOrigin
		cr.CharWid = float32(cr.ClientExtentX*2) / float32(cr.MaxCharsX)
		cr.CharHei = float32(cr.ClientExtentY*2) / float32(cr.MaxCharsY)
		cr.CharWidInPixels = int(float32(currAppWidth) / float32(cr.MaxCharsX))
		cr.CharHeiInPixels = int(float32(currAppHeight) / float32(cr.MaxCharsY))
		cr.PixelWid = cr.ClientExtentX * 2 / float32(currAppWidth)
		cr.PixelHei = cr.ClientExtentY * 2 / float32(currAppHeight)
	} else {
		fmt.Printf("CcRenderer.Init(): for resize changes - ClientExtentX: %.2f\n", cr.ClientExtentX)
		//cr.ClientExtentX = cr.DistanceFromOrigin * (currFrustum.Right / initFrustum.Right)
		//cr.ClientExtentY = cr.DistanceFromOrigin * (currFrustum.Top / initFrustum.Top)
		*prevFrustum = *currFrustum

		currFrustum.Right = float32(currAppWidth) / float32(initAppWidth) * initFrustum.Right
		currFrustum.Left = -currFrustum.Right
		currFrustum.Top = float32(currAppHeight) / float32(initAppHeight) * initFrustum.Top
		currFrustum.Bottom = -currFrustum.Top

		fmt.Printf("CcRenderer.Init(): for resize changes - prevFrustum.Left: %.3f\n", prevFrustum.Left)
		fmt.Printf("CcRenderer.Init(): for resize changes - currFrustum.Left: %.3f\n", currFrustum.Left)

		cr.ClientExtentX = cr.DistanceFromOrigin * currFrustum.Right
		cr.ClientExtentY = cr.DistanceFromOrigin * currFrustum.Top

		// things that weren't initialized in this func
		menu.SetSize()

		for _, pan := range cr.Panels {
			pan.SetSize()
		}
	}

	if len(cr.Panels) == 0 {
		cr.Panels = append(cr.Panels, &TextPanel{BandPercent: 1 - cr.RunPanelHeiPerc, IsEditable: true})
		cr.Panels = append(cr.Panels, &TextPanel{BandPercent: cr.RunPanelHeiPerc, IsEditable: true}) // console (runtime feedback log)	// FIXME so its not editable once we're done debugging some things
		cr.Focused = cr.Panels[0]

		cr.Panels[0].Init()
		cr.Panels[0].SetupDemoProgram()
		cr.Panels[1].Init()
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
	curs.Update()
	menu.Draw()

	for _, pan := range cr.Panels {
		pan.Draw()
	}

	// show width of client area
	if bu.Rect == nil {
		var f float32 = 1.3
		bu.Rect = &common.Rectangle{rend.CharHei + 0.4, f, rend.CharHei + 0.1, -f} // OPTIMIZEME? is this causing slow GC thrashing?
		b2.Rect = &common.Rectangle{rend.CharHei + 0.8, f, rend.CharHei + 0.5, -f} // OPTIMIZEME? is this causing slow GC thrashing?
	}
	bu.Name = fmt.Sprintf("extentX: %.2f", cr.ClientExtentX)
	bu.Draw()
	b2.Name = fmt.Sprintf("extentY: %.2f", cr.ClientExtentY)
	b2.Draw()

	// 'crosshair' center indicator
	var f float32 = rend.CharHei
	rend.DrawCharAtRect('+', &common.Rectangle{f, f, -f, -f})
}

func (cr *CcRenderer) ScrollPanelThatIsHoveredOver(mousePixelDeltaX, mousePixelDeltaY float64) {
	for _, pan := range cr.Panels {
		pan.ScrollIfMouseOver(mousePixelDeltaX, mousePixelDeltaY)
	}
}

func (cr *CcRenderer) DrawCharAtRect(char rune, r *common.Rectangle) {
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

func (cr *CcRenderer) DrawQuad(atlasX, atlasY float32, r *common.Rectangle) {
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

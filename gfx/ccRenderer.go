package gfx

import (
	"fmt"
	"github.com/corpusc/viscript/app"
	"github.com/corpusc/viscript/ui"
	"github.com/go-gl/gl/v2.1/gl"
)

type CcRenderer struct {
	DistanceFromOrigin float32
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
	CurrX   float32
	CurrY   float32
	Focused *ScrollablePanel
	Panels  []*ScrollablePanel
}

func (cr *CcRenderer) SetSize() {
	fmt.Printf("CcRenderer.SetSize() - ClientExtentX: %.2f\n", cr.ClientExtentX)
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

func (cr *CcRenderer) GetMenuSizedRect() *app.Rectangle {
	return &app.Rectangle{
		Rend.ClientExtentY,
		Rend.ClientExtentX,
		Rend.ClientExtentY - Rend.CharHei,
		-Rend.ClientExtentX}
}

func (cr *CcRenderer) DrawAll() {
	Curs.Update()
	cr.DrawMenu()

	for _, pan := range cr.Panels {
		pan.Draw()
	}
}

func (cr *CcRenderer) DrawMenu() {
	for _, bu := range ui.MainMenu.Buttons {
		if bu.Activated {
			SetColor(Green)
		} else {
			SetColor(White)
		}

		Rend.DrawStretchableRect(11, 13, bu.Rect)
		Rend.DrawTextInRect(bu.Name, bu.Rect)
	}
}

func (cr *CcRenderer) DrawTextInRect(s string, r *app.Rectangle) {
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
		Rend.DrawCharAtRect(c, &app.Rectangle{r.Top - lipSpan, x + w, r.Bottom + lipSpan, x})
		x += w
	}
}

func (cr *CcRenderer) DrawCharAtRect(char rune, r *app.Rectangle) {
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

func (cr *CcRenderer) DrawTriangle(atlasX, atlasY float32, a, b, c app.Vec2) {
	// for convenience, and because drawing some extra triangles
	// (only for flow arrows between tree node blocks ATM) won't matter,
	// we are actually drawing a quad, with the last 2 verts @ the same spot

	sp /* span */ := Rend.UvSpan
	u := float32(atlasX) * sp
	v := float32(atlasY) * sp

	gl.Normal3f(0, 0, 1)

	gl.TexCoord2f(u, v)
	gl.Vertex3f(a.X, a.Y, 0)

	gl.TexCoord2f(u+sp, v)
	gl.Vertex3f(b.X, b.Y, 0)

	gl.TexCoord2f(u+sp/2, v+sp)
	gl.Vertex3f(c.X, c.Y, 0)
	gl.TexCoord2f(u+sp/2, v+sp)
	gl.Vertex3f(c.X, c.Y, 0)
}

func (cr *CcRenderer) DrawQuad(atlasX, atlasY float32, r *app.Rectangle) {
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

func (cr *CcRenderer) DrawStretchableRect(atlasX, atlasY float32, r *app.Rectangle) {
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

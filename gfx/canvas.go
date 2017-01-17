package gfx

import (
	"fmt"
	"github.com/corpusc/viscript/app"
	"github.com/corpusc/viscript/cGfx"
	"github.com/corpusc/viscript/ui"
	"github.com/go-gl/gl/v2.1/gl"
)

var goldenRatio = 1.61803398875
var goldenFraction = float32(goldenRatio / (goldenRatio + 1))

func init() {
	fmt.Println("init() - canvas.go")

	// MORE one-time setup
	ui.MainMenu.SetSize(cGfx.GetMenuSizedRect())
}

func SetSize() {
	// things that weren't initialized in this func
	ui.MainMenu.SetSize(cGfx.GetMenuSizedRect())
}

func DrawAll() {
	DrawMenu()
}

func DrawMenu() {
	for _, bu := range ui.MainMenu.Buttons {
		if bu.Activated {
			//SetColor(cGfx.Green)
		} else {
			//SetColor(cGfx.White)
		}

		Update9SlicedRect(bu.Rect)
		DrawTextInRect(bu.Name, bu.Rect.Rect)
	}
}

func DrawTextInRect(s string, r *app.Rectangle) {
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
		DrawCharAtRect(c, &app.Rectangle{r.Top - lipSpan, x + w, r.Bottom + lipSpan, x})
		x += w
	}
}

func DrawCharAtRect(char rune, r *app.Rectangle) {
	u := float32(int(char) % 16)
	v := float32(int(char) / 16)
	sp := cGfx.UvSpan

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

func DrawTriangle(atlasX, atlasY float32, a, b, c app.Vec2F) {
	// for convenience, and because drawing some extra triangles
	// (only for flow arrows between tree node blocks ATM) won't matter,
	// we are actually drawing a quad, with the last 2 verts @ the same spot

	sp /* span */ := cGfx.UvSpan
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

func DrawQuad(atlasX, atlasY float32, r *app.Rectangle) {
	sp /* span */ := cGfx.UvSpan
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

func Update9SlicedRect(r *app.PicRectangle) {
	// 9 quads (like a tic-tac-toe grid, or a "#", which has 9 cells)
	// which keep a predictable frame/margin/edge undistorted,
	// while stretching the middle to fit the desired space

	w := r.Rect.Width()
	h := r.Rect.Height()

	// skip invisible or inverted rects
	if w <= 0 || h <= 0 {
		return
	}

	//var uvEdgeFraction float32 = 0.125 // 1/8
	var uvEdgeFraction float32 = 0.125 / 2 // 1/16
	// we're gonna draw from top to bottom (positivemost to negativemost)

	sp /* span */ := cGfx.UvSpan
	u := float32(r.AtlasPos.X) * sp
	v := float32(r.AtlasPos.Y) * sp

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

	edgeSpan := cGfx.PixelSize.X * 4
	if edgeSpan > w/2 {
		edgeSpan = w / 2
	}

	xSpots := []float32{}
	xSpots = append(xSpots, r.Rect.Left)
	xSpots = append(xSpots, r.Rect.Left+edgeSpan)
	xSpots = append(xSpots, r.Rect.Right-edgeSpan)
	xSpots = append(xSpots, r.Rect.Right)

	edgeSpan = cGfx.PixelSize.Y * 4
	if edgeSpan > h/2 {
		edgeSpan = h / 2
	}

	ySpots := []float32{}
	ySpots = append(ySpots, r.Rect.Top)
	ySpots = append(ySpots, r.Rect.Top-edgeSpan)
	ySpots = append(ySpots, r.Rect.Bottom+edgeSpan)
	ySpots = append(ySpots, r.Rect.Bottom)

	if ySpots[1] > ySpots[0] {
		ySpots[1] = ySpots[0]
	}

	for iX := 0; iX < 3; iX++ {
		for iY := 0; iY < 3; iY++ {
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

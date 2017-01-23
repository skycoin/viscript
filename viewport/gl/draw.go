package gl

import (
	"github.com/corpusc/viscript/app"
	"github.com/corpusc/viscript/gfx"
	"github.com/corpusc/viscript/viewport/terminal"
	"github.com/go-gl/gl/v2.1/gl"
)

var x float32 = float32(gfx.InitAppWidth) * gfx.PixelSize.X / 2
var Terms terminal.TerminalStack = terminal.TerminalStack{}
var desktop *app.Rectangle = &app.Rectangle{
	gfx.DistanceFromOrigin,
	x,
	-gfx.DistanceFromOrigin,
	-x}

func init() {
	println("gl.init() - draw.go")
	Terms.Init()
	Terms.AddTerminal()
	//Terms.AddTerminal()
}

func SetColor(newColor []float32) {
	//cGfx.PrevColor = cGfx.CurrColor
	//cGfx.CurrColor = newColor
	gl.Materialfv(gl.FRONT, gl.AMBIENT_AND_DIFFUSE, &newColor[0])
}

/*
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
*/

func DrawCharAtRect(char rune, r *app.Rectangle) {
	u := float32(int(char) % 16)
	v := float32(int(char) / 16)
	sp := app.UvSpan

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

func DrawQuad(tile app.Vec2I, r *app.Rectangle) {
	sp /* span */ := app.UvSpan
	u := float32(tile.X) * sp
	v := float32(tile.Y) * sp

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

func DrawTriangle(atlasX, atlasY float32, a, b, c app.Vec2F) { // (so-called tri)
	// for convenience, and because drawing some extra triangles
	// (only for flow arrows between tree node blocks ATM) won't matter,
	// we are actually drawing a quad, with the last 2 verts @ the same spot

	sp /* span */ := app.UvSpan
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

func Draw9Sliced(r *app.PicRectangle) {
	// // skip invisible or inverted rects
	// if w <= 0 || h <= 0 {
	// 	return
	// }

	/*gl.Normal3f(0, 0, 1)

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
	}*/
}

func drawAll() {
	DrawQuad(gfx.Pic_GradientBorder, desktop)

	for _, value := range Terms.Terms {
		//println("drawing terminal --- Key (TermId):", key, "Value:", value)
		DrawQuad(gfx.Pic_GradientBorder, value.Bounds)

		cr := &app.Rectangle{
			value.Bounds.Top,
			value.Bounds.Left + value.SpanX(),
			value.Bounds.Top - value.SpanY(),
			value.Bounds.Left} //current rectangle

		for x := 0; x < value.GridSize.X; x++ {
			for y := 0; y < value.GridSize.Y; y++ {
				if value.Chars[y][x] != 0 {
					DrawCharAtRect(rune(value.Chars[y][x]), cr)
				}

				cr.Top -= value.SpanY()
				cr.Bottom -= value.SpanY()
			}

			cr.Top = value.Bounds.Top
			cr.Bottom = value.Bounds.Top - value.SpanY()

			cr.Left += value.SpanX()
			cr.Right += value.SpanX()
		}
	}

	/*
		cGfx.DrawAll()

		// skip 0 so we can use it as a code for being uninitialized
		for i := 1; i < len(cGfx.Rects); i++ {
			if cGfx.Rects[i].State == app.RectState_Active {
				//gfx.SetColor(gfx.Rects[i].Color)

				if cGfx.Rects[i].Type == app.RectType_9Slice {
					Draw9Sliced(cGfx.Rects[i])
					DrawQuad(cGfx.Pic_GradientBorder.X, cGfx.Pic_GradientBorder.Y, cGfx.Rects[i].Rectangle)
				} else {
					DrawQuad(cGfx.Pic_GradientBorder.X, cGfx.Pic_GradientBorder.Y, cGfx.Rects[i].Rectangle)
				}
			}
		}
	*/
}

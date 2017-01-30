package terminal

import (
	"github.com/corpusc/viscript/app"
	"github.com/corpusc/viscript/viewport/gl"
)

func (self *TerminalStack) Draw() {
	//println("TerminalStack.Draw()")

	for _, value := range self.Terms {
		//println("drawing terminal --- Key (TermId):", key, "Value:", value)
		gl.DrawQuad(gl.Pic_GradientBorder, value.Bounds, value.Depth)

		cr := &app.Rectangle{
			value.Bounds.Top,
			value.Bounds.Left + value.SpanX(),
			value.Bounds.Top - value.SpanY(),
			value.Bounds.Left} //current rectangle

		for x := 0; x < value.GridSize.X; x++ {
			for y := 0; y < value.GridSize.Y; y++ {
				if value.Chars[y][x] != 0 {
					gl.DrawCharAtRect(rune(value.Chars[y][x]), cr, value.Depth)
				}

				if x == int(value.Curs.X) && y == int(value.Curs.Y) {
					// draw cursor
					gl.DrawQuad(
						gl.Pic_GradientBorder,
						gl.Curs.GetAnimationModifiedRect(*cr), value.Depth)
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
}

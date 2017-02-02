package terminal

import (
	"github.com/corpusc/viscript/app"
	"github.com/corpusc/viscript/viewport/gl"
)

func (self *TerminalStack) Draw() {
	//println("TerminalStack.Draw()")

	for _, value := range self.DrawOrder {
		term := self.Terms[value]
		// println("Drawing: ", term.TerminalId)
		//println("drawing terminal --- Key (TermId):", key, "term:", term)
		gl.DrawQuad(gl.Pic_GradientBorder, term.Bounds, term.Depth)

		cr := &app.Rectangle{
			term.Bounds.Top,
			term.Bounds.Left + term.SpanX(),
			term.Bounds.Top - term.SpanY(),
			term.Bounds.Left} //current rectangle

		for x := 0; x < term.GridSize.X; x++ {
			for y := 0; y < term.GridSize.Y; y++ {
				if term.Chars[y][x] != 0 {
					gl.DrawCharAtRect(rune(term.Chars[y][x]), cr, term.Depth)
				}

				if x == int(term.Curs.X) && y == int(term.Curs.Y) {
					// draw cursor
					gl.DrawQuad(
						gl.Pic_GradientBorder,
						gl.Curs.GetAnimationModifiedRect(*cr), term.Depth*2)
				}

				cr.Top -= term.SpanY()
				cr.Bottom -= term.SpanY()
			}

			cr.Top = term.Bounds.Top
			cr.Bottom = term.Bounds.Top - term.SpanY()

			cr.Left += term.SpanX()
			cr.Right += term.SpanX()
		}
	}
}

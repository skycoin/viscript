package terminal

import (
	"github.com/corpusc/viscript/app"
	"github.com/corpusc/viscript/viewport/gl"
)

func (self *TerminalStack) Draw() {
	for _, value := range self.Terms {
		z := value.Depth

		if value == self.Focused {
			z = 10

			gl.SetColor(gl.White)
		} else {
			gl.SetColor(gl.Gray)
		}

		gl.Draw9SlicedRect(gl.Pic_GradientBorder, value.Bounds, z)

		cr := &app.Rectangle{ //current rect
			value.Bounds.Top,
			value.Bounds.Left + value.CharSize.X,
			value.Bounds.Top - value.CharSize.Y,
			value.Bounds.Left}

		cr.Left += value.BorderSize //start with the initial rect being offset by the border margin
		cr.Right += value.BorderSize
		cr.Top -= value.BorderSize
		cr.Bottom -= value.BorderSize

		for x := 0; x < value.GridSize.X; x++ {
			for y := 0; y < value.GridSize.Y; y++ {
				if value.Chars[y][x] != 0 {
					gl.DrawCharAtRect(rune(value.Chars[y][x]), cr, z)
				}

				//draw cursor (if it's here)
				if x == int(value.Curs.X) &&
					y == int(value.Curs.Y) {
					gl.DrawQuad(
						gl.Pic_GradientBorder,
						gl.Curs.GetCurrentFrame(*cr), z)
				}

				cr.Top -= value.CharSize.Y
				cr.Bottom -= value.CharSize.Y
			}

			cr.Top = value.Bounds.Top - value.BorderSize
			cr.Bottom = value.Bounds.Top - value.BorderSize - value.CharSize.Y

			cr.Left += value.CharSize.X
			cr.Right += value.CharSize.X
		}
	}
}

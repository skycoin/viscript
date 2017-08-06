package terminal

import (
	"github.com/skycoin/viscript/app"
	"github.com/skycoin/viscript/viewport/gl"
	"strconv"
)

func (ts *TerminalStack) Draw() {
	for _, t := range ts.Terms {
		z := t.Depth

		if t == ts.Focused {
			gl.SetColor(gl.White)
		} else {
			gl.SetColor(gl.Gray)
		}

		drawIdTab(t, z)

		//main window background
		gl.Draw9SlicedRect(gl.Pic_GradientBorder, t.Bounds, z)

		//current rect (in character grid of main window)
		cr := &app.Rectangle{
			t.Bounds.Top,
			t.Bounds.Left + t.CharSize.X,
			t.Bounds.Top - t.CharSize.Y,
			t.Bounds.Left}

		cr.Left += t.BorderSize //start with the initial character grid rect being offset by the border margin
		cr.Right += t.BorderSize
		cr.Top -= t.BorderSize
		cr.Bottom -= t.BorderSize

		for x := 0; x < t.GridSize.X; x++ {
			for y := 0; y < t.GridSize.Y; y++ {
				if t.Chars[y][x] != 0 {
					gl.DrawCharAtRect(rune(t.Chars[y][x]), cr, z)
				}

				//draw cursor (if it's here)
				if x == int(t.Cursor.X) &&
					y == int(t.Cursor.Y) {

					gl.DrawQuad(
						gl.Pic_GradientBorder,
						gl.Curs.GetCurrentFrame(*cr), z)
				}

				cr.Top -= t.CharSize.Y
				cr.Bottom -= t.CharSize.Y
			}

			cr.Top = t.Bounds.Top - t.BorderSize
			cr.Bottom = t.Bounds.Top - t.BorderSize - t.CharSize.Y

			cr.Left += t.CharSize.X
			cr.Right += t.CharSize.X
		}
	}
}

//
//
//private
func drawIdTab(t *Terminal, z float32) {
	id := strconv.Itoa(int(t.TerminalId))

	idText := &app.Rectangle{ //rectangle initially used to draw tab background
		t.Bounds.Top + t.BorderSize*2 + t.CharSize.Y,
		t.Bounds.Left + t.BorderSize*2 + t.CharSize.X*float32(len(id)),
		t.Bounds.Top,
		t.Bounds.Left}

	//id tab background
	gl.Draw9SlicedRect(gl.Pic_GradientBorder, idText, z)

	//push in edges to leave a border visible....
	idText.Top -= t.BorderSize
	idText.Bottom += t.BorderSize
	idText.Left += t.BorderSize
	idText.Right = idText.Left + t.CharSize.X //....and shrink width to char size

	//draw the id #
	for i := 0; i < len(id); i++ {
		gl.DrawCharAtRect(rune(id[i]), idText, z)
		idText.Left += t.CharSize.X
		idText.Right += t.CharSize.X
	}
}

package terminal

import (
	"github.com/skycoin/viscript/app"
	"github.com/skycoin/viscript/viewport/gl"
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
	////handle arg conversion errors
	// id, err := strconv.Atoi(t.TerminalId)
	// if err != nil {
	// 	st.PrintError("Unable to convert passed index.")
	// 	s := "err.Error(): \"" + err.Error() + "\""
	// 	st.PrintError(s)
	// 	return
	// }

	// //handle index range errors
	// if id < 0 ||
	// 	id >= len(st.storedTerminalIds) {
	// 	st.PrintError("Index not in range.")
	// 	return
	// }

	// string{strconv.Itoa(int(st.storedTerminalIds[id]))})
	//
	//
	//
	//
	stringer := "" + string(t.TerminalId)

	idText := &app.Rectangle{ //rectangle initially used to draw tab background (FIXME for multidigit ids)
		t.Bounds.Top + t.BorderSize*2 + t.CharSize.Y,
		t.Bounds.Left + t.BorderSize*2 + t.CharSize.X,
		t.Bounds.Top,
		t.Bounds.Left}

	//id tab background
	gl.Draw9SlicedRect(gl.Pic_GradientBorder, idText, z)

	//id rectangle now pushes in to leave a border visible
	idText.Left += t.BorderSize
	idText.Right -= t.BorderSize
	idText.Top -= t.BorderSize
	idText.Bottom += t.BorderSize

	//draw the id #
	//
	//(single A atm)
	gl.DrawCharAtRect(rune(stringer[0]), idText, z)
}

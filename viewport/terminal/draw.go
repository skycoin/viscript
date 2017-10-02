package terminal

import (
	"github.com/skycoin/viscript/app"
	"github.com/skycoin/viscript/viewport/gl"
	"strconv"
)

func (ts *TerminalStack) DrawTextMode() { //for running headless mode
	for _, t := range ts.TermMap {
		if t == ts.TermMap[ts.FocusedId] { //in text/headless mode, only draw focused
			println(app.GetBarOfChars("_", t.GridSize.X-1))

			for y := 0; y < t.GridSize.Y; y++ {
				line := ""
				blank := true

				for x := 0; x < t.GridSize.X; x++ {
					if t.Chars[y][x] == 0 {
						line += " "
					} else {
						line += string(t.Chars[y][x])
						blank = false
					}
				}

				if blank {
					line = ""
				}

				println(line)
			}
		}
	}
}

func (ts *TerminalStack) Draw() {
	for _, t := range ts.TermMap {
		z := t.Depth

		if t.TerminalId == ts.FocusedId {
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
//
//

func drawIdTab(t *Terminal, z float32) {
	//...with a rectangle whose bottom lip/edge will be covered by main window
	id := strconv.Itoa(int(t.TerminalId))

	tr := &app.Rectangle{ //text rectangle (initially used to draw whole tab background)
		t.Bounds.Top + t.BorderSize + t.CharSize.Y,
		t.Bounds.Left + t.BorderSize*2 + t.CharSize.X*float32(len(id)),
		t.Bounds.Top - t.BorderSize,
		t.Bounds.Left}

	//id tab background
	gl.Draw9SlicedRect(gl.Pic_GradientBorder, tr, z)

	//push in edges to encompass ONLY the text (leaving a border visible)
	tr.Top -= t.BorderSize
	tr.Bottom += t.BorderSize
	tr.Left += t.BorderSize
	tr.Right = tr.Left + t.CharSize.X //....and shrink width to char size

	//draw the id #
	for i := 0; i < len(id); i++ {
		gl.DrawCharAtRect(rune(id[i]), tr, z)
		tr.Left += t.CharSize.X
		tr.Right += t.CharSize.X
	}
}

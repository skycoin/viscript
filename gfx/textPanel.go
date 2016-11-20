package gfx

import (
	"fmt"
	"github.com/corpusc/viscript/common"
	"github.com/corpusc/viscript/ui"
	"github.com/go-gl/gl/v2.1/gl"
)

type TextPanel struct {
	BandPercent float32 // what percentage (in the relevant dimension) of the entire parent PanelBand do we occupy?
	CursX       int     // current cursor/insert position (in character grid cells/units)
	CursY       int
	MouseX      int // current mouse position in character grid space (units/cells)
	MouseY      int
	IsEditable  bool
	Rect        *common.Rectangle
	Selection   *ui.SelectionRange
	BarHori     *ui.ScrollBar // horizontal
	BarVert     *ui.ScrollBar // vertical
	Body        []string
}

func (tp *TextPanel) Init() {
	fmt.Printf("TextPanel.Init()\n")

	tp.Selection = &ui.SelectionRange{}
	tp.Selection.Init()

	// scrollbars
	tp.BarHori = &ui.ScrollBar{IsHorizontal: true}
	tp.BarVert = &ui.ScrollBar{}
	tp.BarHori.Rect = &common.Rectangle{}
	tp.BarVert.Rect = &common.Rectangle{}

	tp.SetSize()
}

func (tp *TextPanel) SetSize() {
	fmt.Printf("TextPanel.SetSize()\n")

	tp.Rect = &common.Rectangle{
		Rend.ClientExtentY - Rend.CharHei,
		Rend.ClientExtentX,
		-Rend.ClientExtentY,
		-Rend.ClientExtentX}

	if tp.BandPercent == Rend.RunPanelHeiPerc { // FIXME: this is hardwired for one use case for now
		tp.Rect.Top = tp.Rect.Bottom + tp.Rect.Height()*tp.BandPercent
	} else {
		tp.Rect.Bottom = tp.Rect.Bottom + tp.Rect.Height()*Rend.RunPanelHeiPerc
	}

	// set scrollbars' upper left corners
	tp.BarHori.Rect.Left = tp.Rect.Left
	tp.BarHori.Rect.Top = tp.Rect.Bottom + ui.ScrollBarThickness
	tp.BarVert.Rect.Left = tp.Rect.Right - ui.ScrollBarThickness
	tp.BarVert.Rect.Top = tp.Rect.Top
}

func (tp *TextPanel) RespondToMouseClick() {
	Rend.Focused = tp

	// diffs/deltas from home position of panel (top left corner)
	glDeltaXFromHome := Curs.MouseGlX - tp.Rect.Left
	glDeltaYFromHome := Curs.MouseGlY - tp.Rect.Top
	tp.MouseX = int((glDeltaXFromHome + tp.BarHori.ScrollDelta) / Rend.CharWid)
	tp.MouseY = int(-(glDeltaYFromHome + tp.BarVert.ScrollDelta) / Rend.CharHei)

	if tp.MouseY < 0 {
		tp.MouseY = 0
	}

	if tp.MouseY >= len(tp.Body) {
		tp.MouseY = len(tp.Body) - 1
	}
}

func (tp *TextPanel) GoToTopEdge() {
	Rend.CurrY = tp.Rect.Top - tp.BarVert.ScrollDelta
}
func (tp *TextPanel) GoToLeftEdge() float32 {
	Rend.CurrX = tp.Rect.Left - tp.BarHori.ScrollDelta
	return Rend.CurrX
}
func (tp *TextPanel) GoToTopLeftCorner() {
	tp.GoToTopEdge()
	tp.GoToLeftEdge()
}

func (tp *TextPanel) Draw() {
	tp.GoToTopLeftCorner()
	tp.DrawBackground(11, 13)

	cX := Rend.CurrX // current (internal/logic cursor) drawing position
	cY := Rend.CurrY
	cW := Rend.CharWid
	cH := Rend.CharHei
	b := tp.BarHori.Rect.Top // bottom of text area

	// body of text
	for y, line := range tp.Body {
		// if line visible
		if cY <= tp.Rect.Top+cH && cY >= b {
			r := &common.Rectangle{cY, cX + cW, cY - cH, cX} // t, r, b, l

			// if line needs vertical adjustment
			if cY > tp.Rect.Top {
				r.Top = tp.Rect.Top
			}
			if cY-cH < b {
				r.Bottom = b
			}

			//parseLine(y, line, true)
			Rend.Color(Gray)

			// process line of text
			for x, c := range line {
				// if char visible
				if cX >= tp.Rect.Left-cW && cX < tp.BarVert.Rect.Left {
					common.ClampLeftAndRightOf(r, tp.Rect.Left, tp.BarVert.Rect.Left)
					Rend.DrawCharAtRect(c, r)

					if tp.IsEditable && Curs.Visible == true {
						if x == tp.CursX && y == tp.CursY {
							Rend.Color(White)
							Rend.DrawCharAtRect('_', r)
							Rend.Color(Rend.PrevColor)
						}
					}
				}

				cX += cW
				r.Left = cX
				r.Right = cX + cW
			}

			// draw cursor at the end of line if needed
			if cX < tp.BarVert.Rect.Left && y == tp.CursY && tp.CursX == len(line) {
				if tp.IsEditable && Curs.Visible == true {
					Rend.Color(White)
					common.ClampLeftAndRightOf(r, tp.Rect.Left, tp.BarVert.Rect.Left)
					Rend.DrawCharAtRect('_', r)
				}
			}

			cX = tp.GoToLeftEdge()
		}

		cY -= cH // go down a line height
	}

	Rend.Color(GrayDark)
	tp.DrawScrollbarChrome(10, 11, tp.Rect.Right-ui.ScrollBarThickness, tp.Rect.Top)                          // vertical bar background
	tp.DrawScrollbarChrome(13, 12, tp.Rect.Left, tp.Rect.Bottom+ui.ScrollBarThickness)                        // horizontal bar background
	tp.DrawScrollbarChrome(12, 11, tp.Rect.Right-ui.ScrollBarThickness, tp.Rect.Bottom+ui.ScrollBarThickness) // corner elbow piece
	Rend.Color(Gray)
	//tp.BarHori.SetSize(tp.Rect, tp.Body, cW, cH)
	tp.BarVert.SetSize(tp.Rect, tp.Body, cW, cH)
	//	Rend.DrawQuad(11, 13, tp.BarHori.Rect) // 2,11 (pixel checkerboard)    // 14, 15 (square in the middle)
	Rend.DrawQuad(11, 13, tp.BarVert.Rect) // 13, 12 (double horizontal lines)    // 10, 11 (double vertical lines)
	Rend.Color(White)
}

// ATM the only different between the 2 funcs below is the top left corner (involving 3 vertices)
func (tp *TextPanel) DrawScrollbarChrome(atlasCellX, atlasCellY, l, t float32) { // left, top
	sp := Rend.UvSpan
	u := float32(atlasCellX) * sp
	v := float32(atlasCellY) * sp

	gl.Normal3f(0, 0, 1)

	// bottom left   0, 1
	gl.TexCoord2f(u, v+sp)
	gl.Vertex3f(l, tp.Rect.Bottom, 0)

	// bottom right   1, 1
	gl.TexCoord2f(u+sp, v+sp)
	gl.Vertex3f(tp.Rect.Right, tp.Rect.Bottom, 0)

	// top right   1, 0
	gl.TexCoord2f(u+sp, v)
	gl.Vertex3f(tp.Rect.Right, t, 0)

	// top left   0, 0
	gl.TexCoord2f(u, v)
	gl.Vertex3f(l, t, 0)
}

func (tp *TextPanel) DrawBackground(atlasCellX, atlasCellY float32) {
	sp := Rend.UvSpan
	u := float32(atlasCellX) * sp
	v := float32(atlasCellY) * sp

	Rend.Color(GrayDark)
	gl.Normal3f(0, 0, 1)

	// bottom left   0, 1
	gl.TexCoord2f(u, v+sp)
	gl.Vertex3f(tp.Rect.Left, tp.Rect.Bottom+ui.ScrollBarThickness, 0)

	// bottom right   1, 1
	gl.TexCoord2f(u+sp, v+sp)
	gl.Vertex3f(tp.Rect.Right-ui.ScrollBarThickness, tp.Rect.Bottom+ui.ScrollBarThickness, 0)

	// top right   1, 0
	gl.TexCoord2f(u+sp, v)
	gl.Vertex3f(tp.Rect.Right-ui.ScrollBarThickness, tp.Rect.Top, 0)

	// top left   0, 0
	gl.TexCoord2f(u, v)
	gl.Vertex3f(tp.Rect.Left, tp.Rect.Top, 0)
}

func (tp *TextPanel) ScrollIfMouseOver(mousePixelDeltaX, mousePixelDeltaY float64) {
	if tp.ContainsMouseCursor() {
		// position increments in gl space
		xInc := float32(mousePixelDeltaX) * Rend.PixelWid
		yInc := float32(mousePixelDeltaY) * Rend.PixelHei
		tp.BarHori.Scroll(xInc)
		tp.BarVert.Scroll(yInc)
	}
}

func (tp *TextPanel) ContainsMouseCursor() bool {
	return MouseCursorIsInside(tp.Rect)
}

func (tp *TextPanel) ContainsMouseCursorInsideOfScrollBars() bool {
	return MouseCursorIsInside(&common.Rectangle{
		tp.Rect.Top, tp.Rect.Right - ui.ScrollBarThickness, tp.Rect.Bottom + ui.ScrollBarThickness, tp.Rect.Left})
}

func (tp *TextPanel) RemoveCharacter(fromUnderCursor bool) {
	if fromUnderCursor {
		if len(tp.Body[tp.CursY]) > tp.CursX {
			tp.Body[tp.CursY] = tp.Body[tp.CursY][:tp.CursX] + tp.Body[tp.CursY][tp.CursX+1:len(tp.Body[tp.CursY])]
		}
	} else {
		if tp.CursX > 0 {
			tp.Body[tp.CursY] = tp.Body[tp.CursY][:tp.CursX-1] + tp.Body[tp.CursY][tp.CursX:len(tp.Body[tp.CursY])]
			tp.CursX--
		}
	}
}

func (tp *TextPanel) SetupDemoProgram() {
	tp.Body = append(tp.Body, "// ------- variable declarations ------- -------")
	//tp.Body = append(tp.Body, "var myVar int32")
	tp.Body = append(tp.Body, "var a int32 = 42")
	tp.Body = append(tp.Body, "var b int32 = 58")
	tp.Body = append(tp.Body, "")
	tp.Body = append(tp.Body, "// ------- builtin function calls ------- ------- ------- ------- ------- ------- ------- end")
	tp.Body = append(tp.Body, "//    sub32(7, 9)")
	//tp.Body = append(tp.Body, "sub32(4,8)")
	//tp.Body = append(tp.Body, "mult32(7, 7)")
	//tp.Body = append(tp.Body, "mult32(3,5)")
	//tp.Body = append(tp.Body, "div32(8,2)")
	//tp.Body = append(tp.Body, "div32(15,  3)")
	//tp.Body = append(tp.Body, "add32(2,3)")
	//tp.Body = append(tp.Body, "add32(a, b)")
	tp.Body = append(tp.Body, "")
	tp.Body = append(tp.Body, "// ------- user function calls -------")
	tp.Body = append(tp.Body, "myFunc(a, b)")
	tp.Body = append(tp.Body, "")
	tp.Body = append(tp.Body, "// ------- function declarations -------")
	tp.Body = append(tp.Body, "func myFunc(a int32, b int32){")
	tp.Body = append(tp.Body, "")
	tp.Body = append(tp.Body, "        div32(6, 2)")
	tp.Body = append(tp.Body, "        innerFunc(a,b)")
	tp.Body = append(tp.Body, "}")
	tp.Body = append(tp.Body, "")
	tp.Body = append(tp.Body, "func innerFunc (a, b int32) {")
	tp.Body = append(tp.Body, "        var locA int32 = 71")
	tp.Body = append(tp.Body, "        var locB int32 = 29")
	tp.Body = append(tp.Body, "        sub32(locA, locB)")
	tp.Body = append(tp.Body, "    }    ")

	/*
		for i := 0; i < 22; i++ {
			tp.Body = append(tp.Body, fmt.Sprintf("%d: put lots of text on screen", i))
		}
	*/
}

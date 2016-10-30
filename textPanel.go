package main

import (
	"fmt"
	"github.com/go-gl/gl/v2.1/gl"
)

type TextPanel struct {
	BandPercent float32 // what percentage (in the relevant dimension) of the entire parent PanelBand do we occupy?
	CursX       int     // current cursor/insert position (in character grid cells/units)
	CursY       int
	MouseX      int // current mouse position in character grid space (units/cells)
	MouseY      int
	IsEditable  bool
	Rect        *Rectangle
	Selection   *SelectionRange
	BarHori     *ScrollBar // horizontal
	BarVert     *ScrollBar // vertical
	Body        []string
}

func (tp *TextPanel) Init() {
	fmt.Printf("TextPanel.Init()\n")

	tp.Selection = &SelectionRange{}
	tp.Selection.Init()

	// scrollbar
	tp.BarHori = &ScrollBar{IsHorizontal: true}
	tp.BarVert = &ScrollBar{}

	tp.SetSize()
}

func (tp *TextPanel) SetSize() {
	fmt.Printf("TextPanel.SetSize()\n")

	tp.Rect = &Rectangle{
		rend.ClientExtentY - rend.CharHei,
		rend.ClientExtentX,
		-rend.ClientExtentY,
		-rend.ClientExtentX}

	if tp.BandPercent == rend.RunPanelHeiPerc { // FIXME: this is hardwired for one use case for now
		tp.Rect.Top = tp.Rect.Bottom + tp.Rect.Height()*tp.BandPercent
	} else {
		tp.Rect.Bottom = tp.Rect.Bottom + tp.Rect.Height()*rend.RunPanelHeiPerc
	}

	// scrollbar
	tp.BarHori.PosX = tp.Rect.Left
	tp.BarHori.PosY = tp.Rect.Bottom + tp.BarHori.Thickness
	tp.BarVert.PosX = tp.Rect.Right - tp.BarVert.Thickness
	tp.BarVert.PosY = tp.Rect.Top
}

func (tp *TextPanel) RespondToMouseClick() {
	rend.Focused = tp

	// diffs/deltas from home position of panel (top left corner)
	glDeltaXFromHome := curs.MouseGlX - tp.Rect.Left
	glDeltaYFromHome := curs.MouseGlY - tp.Rect.Top
	tp.MouseX = int((glDeltaXFromHome + tp.BarHori.ScrollDelta) / rend.CharWid)
	tp.MouseY = int(-(glDeltaYFromHome + tp.BarVert.ScrollDelta) / rend.CharHei)

	if tp.MouseY < 0 {
		tp.MouseY = 0
	}

	if tp.MouseY >= len(tp.Body) {
		tp.MouseY = len(tp.Body) - 1
	}
}

func (tp *TextPanel) GoToTopEdge() {
	rend.CurrY = tp.Rect.Top - tp.BarVert.ScrollDelta
}
func (tp *TextPanel) GoToLeftEdge() float32 {
	rend.CurrX = tp.Rect.Left - tp.BarHori.ScrollDelta
	return rend.CurrX
}
func (tp *TextPanel) GoToTopLeftCorner() {
	tp.GoToTopEdge()
	tp.GoToLeftEdge()
}

func (tp *TextPanel) Draw() {
	tp.GoToTopLeftCorner()
	tp.DrawBackground(11, 13)

	cX := rend.CurrX // current (internal/logic cursor) drawing position
	cY := rend.CurrY
	cW := rend.CharWid
	cH := rend.CharHei
	b := tp.BarHori.PosY // bottom of text area

	// body of text
	for y, line := range tp.Body {
		// if line visible
		if cY <= tp.Rect.Top+cH && cY >= b {
			r := &Rectangle{cY, cX + cW, cY - cH, cX} // t, r, b, l

			// if line needs vertical adjustment
			if cY > tp.Rect.Top {
				r.Top = tp.Rect.Top
			}
			if cY-cH < b {
				r.Bottom = b
			}

			parseLine(y, line, true)

			// process line of text
			for x, c := range line {
				// if char visible
				if cX >= tp.Rect.Left-cW && cX < tp.BarVert.PosX {
					ClampLeftAndRightOf(r, tp.Rect.Left, tp.BarVert.PosX)
					rend.DrawCharAtRect(c, r)

					if tp.IsEditable && curs.Visible == true {
						if x == tp.CursX && y == tp.CursY {
							rend.Color(white)
							rend.DrawCharAtRect('_', r)
							rend.Color(rend.PrevColor)
						}
					}
				}

				cX += cW
				r.Left = cX
				r.Right = cX + cW
			}

			// draw cursor at the end of line if needed
			if cX < tp.BarVert.PosX && y == tp.CursY && tp.CursX == len(line) {
				if tp.IsEditable && curs.Visible == true {
					rend.Color(white)
					ClampLeftAndRightOf(r, tp.Rect.Left, tp.BarVert.PosX)
					rend.DrawCharAtRect('_', r)
				}
			}

			cX = tp.GoToLeftEdge()
		}

		cY -= cH // go down a line height
	}

	rend.Color(grayDark)
	tp.DrawScrollbarChrome(10, 11, tp.Rect.Right-tp.BarVert.Thickness, tp.Rect.Top)                         // vertical bar background
	tp.DrawScrollbarChrome(13, 12, tp.Rect.Left, tp.Rect.Bottom+tp.BarHori.Thickness)                       // horizontal bar background
	tp.DrawScrollbarChrome(12, 11, tp.Rect.Right-tp.BarVert.Thickness, tp.Rect.Bottom+tp.BarHori.Thickness) // corner elbow piece
	rend.Color(gray)
	tp.BarHori.Draw(11, 13, *tp) // 2,11 (pixel checkerboard)    // 14, 15 (square in the middle)
	tp.BarVert.Draw(11, 13, *tp) // 13, 12 (double horizontal lines)    // 10, 11 (double vertical lines)
	rend.Color(white)
}

// ATM the only different between the 2 funcs below is the top left corner (involving 3 vertices)
func (tp *TextPanel) DrawScrollbarChrome(atlasCellX, atlasCellY, l, t float32) { // left, top
	sp := rend.UvSpan
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
	sp := rend.UvSpan
	u := float32(atlasCellX) * sp
	v := float32(atlasCellY) * sp

	rend.Color(grayDark)
	gl.Normal3f(0, 0, 1)

	// bottom left   0, 1
	gl.TexCoord2f(u, v+sp)
	gl.Vertex3f(tp.Rect.Left, tp.Rect.Bottom+tp.BarHori.Thickness, 0)

	// bottom right   1, 1
	gl.TexCoord2f(u+sp, v+sp)
	gl.Vertex3f(tp.Rect.Right-tp.BarVert.Thickness, tp.Rect.Bottom+tp.BarHori.Thickness, 0)

	// top right   1, 0
	gl.TexCoord2f(u+sp, v)
	gl.Vertex3f(tp.Rect.Right-tp.BarVert.Thickness, tp.Rect.Top, 0)

	// top left   0, 0
	gl.TexCoord2f(u, v)
	gl.Vertex3f(tp.Rect.Left, tp.Rect.Top, 0)
}

func (tp *TextPanel) ScrollIfMouseOver(mousePixelDeltaX, mousePixelDeltaY float64) {
	if tp.ContainsMouseCursor() {
		// position increments in gl space
		xInc := float32(mousePixelDeltaX) * rend.PixelWid
		yInc := float32(mousePixelDeltaY) * rend.PixelHei
		tp.BarHori.Scroll(tp, xInc)
		tp.BarVert.Scroll(tp, yInc)
	}
}

func (tp *TextPanel) ContainsMouseCursor() bool {
	return MouseCursorIsInside(tp.Rect)
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

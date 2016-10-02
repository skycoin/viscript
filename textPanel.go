package main

import (
	"fmt"
	"github.com/go-gl/gl/v2.1/gl"
)

type TextPanel struct {
	Top        float32
	Bottom     float32
	Left       float32
	Right      float32
	NumCharsX  int
	NumCharsY  int
	CursX      int // current cursor/insert position (in character grid cells/units)
	CursY      int
	MouseX     int // current mouse position in character grid space (units/cells)
	MouseY     int
	IsEditable bool
	Selection  *SelectionRange
	BarHori    *ScrollBar // horizontal
	BarVert    *ScrollBar // vertical
	Body       []string
}

func (tp *TextPanel) Init() {
	tp.Selection = &SelectionRange{}
	tp.Selection.Init()

	tp.Left = -textRend.ScreenRad
	tp.Right = textRend.ScreenRad

	if tp.Top == 0 {
		tp.Top = textRend.ScreenRad - textRend.CharHei
	}

	tp.Bottom = tp.Top - float32(tp.NumCharsY)*textRend.CharHei

	tp.BarHori = &ScrollBar{IsHorizontal: true}
	tp.BarVert = &ScrollBar{}
	tp.BarHori.PosX = tp.Left
	tp.BarHori.PosY = tp.Bottom //- textRend.CharHei
	tp.BarVert.PosX = tp.Right  //- textRend.CharWid
	tp.BarVert.PosY = tp.Top

	if tp.NumCharsX == 0 {
		tp.NumCharsX = textRend.MaxCharsX
	}
	if tp.NumCharsY == 0 {
		tp.NumCharsY = textRend.MaxCharsY
	}

	fmt.Printf("TextPanel.Init()    t: %.2f, b: %.2f\n", tp.Top, tp.Bottom)
}

func (tp *TextPanel) RespondToMouseClick() {
	textRend.Focused = tp

	// diffs/deltas from home position of panel (top left corner)
	glDeltaXFromHome := curs.MouseGlX - tp.Left
	glDeltaYFromHome := curs.MouseGlY - tp.Top
	tp.MouseX = int((glDeltaXFromHome + tp.BarHori.ScrollDelta) / textRend.CharWid)
	tp.MouseY = int(-(glDeltaYFromHome + tp.BarVert.ScrollDelta) / textRend.CharHei)

	if tp.MouseY < 0 {
		tp.MouseY = 0
	}

	if tp.MouseY >= len(tp.Body) {
		tp.MouseY = len(tp.Body) - 1
	}
}

func (tp *TextPanel) GoToTopEdge() {
	textRend.CurrY = tp.Top - tp.BarVert.ScrollDelta
}
func (tp *TextPanel) GoToLeftEdge() {
	textRend.CurrX = tp.Left - tp.BarHori.ScrollDelta
}
func (tp *TextPanel) GoToTopLeftCorner() {
	tp.GoToTopEdge()
	tp.GoToLeftEdge()
}

func (tp *TextPanel) Draw() {
	tp.BarHori.UpdateSize(tp)
	tp.BarVert.UpdateSize(tp)
	tp.GoToTopLeftCorner()
	tp.DrawBackground(11, 13)

	// body of text
	for y, line := range tp.Body {
		// if line visible
		if textRend.CurrY <= tp.Top+textRend.CharHei && textRend.CurrY >= tp.Bottom {
			clipSpan := &Rectangle{}

			// if line needs clipping
			if textRend.CurrY > tp.Top {
				clipSpan.Top = textRend.CurrY - tp.Top
			}
			if textRend.CurrY-textRend.CharHei < tp.Bottom {
				clipSpan.Bottom = (textRend.CurrY - textRend.CharHei) - tp.Bottom
			}

			// draw line of text
			for x, c := range line {
				textRend.DrawCharAtCurrentPosition(c, clipSpan)

				if tp.IsEditable && curs.Visible == true {
					//fmt.Printf("tp.CursX, tp.CursY: %d,%d\n", tp.CursX, tp.CursY)

					if x == tp.CursX && y == tp.CursY {
						textRend.CurrX -= textRend.CharWid
						textRend.DrawCharAtCurrentPosition('_', clipSpan)
					}
				}
			}

			// draw cursor at the end of line if needed
			if y == tp.CursY && tp.CursX == len(line) {
				if tp.IsEditable && curs.Visible == true {
					textRend.DrawCharAtCurrentPosition('_', clipSpan)
				}
			}

			tp.GoToLeftEdge()
		}

		textRend.CurrY -= textRend.CharHei // go down a line height
	}

	tp.BarHori.Draw(2, 11, *tp)
	tp.BarVert.Draw(2, 11, *tp)
	tp.DrawScrollbarCorner(12, 11)
}

// ATM the only different between the 2 funcs below is the top left corner (involving 3 vertices)
func (tp *TextPanel) DrawScrollbarCorner(atlasCellX, atlasCellY float32) {
	sp := textRend.UvSpan
	u := float32(atlasCellX) * sp
	v := float32(atlasCellY) * sp

	gl.Normal3f(0, 0, 1)

	// bottom left   0, 1
	gl.TexCoord2f(u, v+sp)
	gl.Vertex3f(tp.BarVert.PosX, tp.Bottom, 0)

	// bottom right   1, 1
	gl.TexCoord2f(u+sp, v+sp)
	gl.Vertex3f(tp.Right, tp.Bottom, 0)

	// top right   1, 0
	gl.TexCoord2f(u+sp, v)
	gl.Vertex3f(tp.Right, tp.BarHori.PosY, 0)

	// top left   0, 0
	gl.TexCoord2f(u, v)
	gl.Vertex3f(tp.BarVert.PosX, tp.BarHori.PosY, 0)
}

func (tp *TextPanel) DrawBackground(atlasCellX, atlasCellY float32) {
	sp := textRend.UvSpan
	u := float32(atlasCellX) * sp
	v := float32(atlasCellY) * sp

	gl.Normal3f(0, 0, 1)

	// bottom left   0, 1
	gl.TexCoord2f(u, v+sp)
	gl.Vertex3f(tp.Left, tp.Bottom, 0)

	// bottom right   1, 1
	gl.TexCoord2f(u+sp, v+sp)
	gl.Vertex3f(tp.Right, tp.Bottom, 0)

	// top right   1, 0
	gl.TexCoord2f(u+sp, v)
	gl.Vertex3f(tp.Right, tp.Top, 0)

	// top left   0, 0
	gl.TexCoord2f(u, v)
	gl.Vertex3f(tp.Left, tp.Top, 0)
}

func (tp *TextPanel) ScrollIfMouseOver(mousePixelDeltaX, mousePixelDeltaY float64) {
	if tp.ContainsMouseCursor() {
		// position increments in gl space
		xInc := float32(mousePixelDeltaX) * textRend.PixelWid
		yInc := float32(mousePixelDeltaY) * textRend.PixelHei
		tp.BarHori.ScrollThisMuch(tp, xInc)
		tp.BarVert.ScrollThisMuch(tp, yInc)
	}
}

func (tp *TextPanel) ContainsMouseCursor() bool {
	if curs.MouseGlY < tp.Top && curs.MouseGlY > tp.Bottom {
		if curs.MouseGlX < tp.Right && curs.MouseGlX > tp.Left {
			return true
		}
	}

	return false
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
	tp.Body = append(tp.Body, "// ------- variable declarations ------- ------- ------- ------- ------- ------- ------- ------- end")
	tp.Body = append(tp.Body, "var myVar int32")
	tp.Body = append(tp.Body, "var a int32 = 42")
	tp.Body = append(tp.Body, "var b int32 = 58")
	tp.Body = append(tp.Body, "")
	tp.Body = append(tp.Body, "// ------- builtin function calls -------")
	tp.Body = append(tp.Body, "    sub32(7, 9)")
	tp.Body = append(tp.Body, "sub32(4,8)")
	tp.Body = append(tp.Body, "mult32(7, 7)")
	tp.Body = append(tp.Body, "mult32(3,5)")
	tp.Body = append(tp.Body, "div32(8,2)")
	tp.Body = append(tp.Body, "div32(15,  3)")
	tp.Body = append(tp.Body, "add32(2,3)")
	tp.Body = append(tp.Body, "add32(a, b)")
	tp.Body = append(tp.Body, "")
	tp.Body = append(tp.Body, "// ------- user function calls -------")
	tp.Body = append(tp.Body, "myFunc(a, b)")
	tp.Body = append(tp.Body, "")
	tp.Body = append(tp.Body, "// ------- function declarations -------")
	tp.Body = append(tp.Body, "func myFunc(a int32, b int32){")
	tp.Body = append(tp.Body, "")
	tp.Body = append(tp.Body, "        innerFunc(a,b)")
	tp.Body = append(tp.Body, "}")
	tp.Body = append(tp.Body, "func innerFunc (a, b int32) {")
	tp.Body = append(tp.Body, "        var myLocal int32")
	tp.Body = append(tp.Body, "    }    ")

	/*
		for i := 0; i < 22; i++ {
			tp.Body = append(tp.Body, fmt.Sprintf("%d: put lots of text on screen", i))
		}
	*/
}

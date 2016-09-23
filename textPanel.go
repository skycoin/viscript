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
	Selection  SelectionRange
	Bar        *ScrollBar
	Body       []string
}

func (tp *TextPanel) Init() {
	tp.Selection = SelectionRange{}
	tp.Selection.Init()

	tp.Left = -textRend.ScreenRad
	tp.Right = textRend.ScreenRad

	if tp.Top == 0 {
		tp.Top = textRend.ScreenRad - textRend.CharHei
	}

	tp.Bottom = tp.Top - float32(tp.NumCharsY)*textRend.CharHei

	tp.Bar = &ScrollBar{}
	tp.Bar.PosX = tp.Right - textRend.CharWid
	tp.Bar.PosY = tp.Top

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
	//fmt.Printf("glDeltaYFromHome: %.2f\n", glDeltaYFromHome)
	tp.MouseX = int((glDeltaXFromHome + tp.Bar.ScrollDeltaX) / textRend.CharWid)
	tp.MouseY = int(-(glDeltaYFromHome + tp.Bar.ScrollDeltaY) / textRend.CharHei)
	//fmt.Printf("tp.MouseY: %d\n", tp.MouseY)

	if tp.MouseY < 0 {
		tp.MouseY = 0
	}

	if tp.MouseY >= len(tp.Body) {
		tp.MouseY = len(tp.Body) - 1
	}
}

func (tp *TextPanel) GoToTopEdge() {
	//fmt.Printf("GoToTopEdge() tp.ScrollDeltaY: %.2f\n", tp.ScrollDeltaY)
	textRend.CurrY = tp.Top - tp.Bar.ScrollDeltaY //tp.Bar.OffsetY
}
func (tp *TextPanel) GoToLeftEdge() {
	textRend.CurrX = tp.Left
}
func (tp *TextPanel) GoToTopLeftCorner() {
	tp.GoToTopEdge()
	tp.GoToLeftEdge()
}

func (tp *TextPanel) Draw() {
	tp.Bar.UpdateSize(tp)
	tp.GoToTopLeftCorner()
	tp.DrawBackground(11, 13)

	// body of text
	for y, line := range tp.Body {
		if textRend.CurrY <= tp.Top+textRend.CharHei { // if line visible
			var clipSpan float32

			if textRend.CurrY > tp.Top { // if line needs its top clipped
				clipSpan = textRend.CurrY - tp.Top
			}

			// draw line of text
			for x, c := range line {
				drawCurrentChar(c, clipSpan)

				if tp.IsEditable && curs.Visible == true {
					//fmt.Printf("tp.CursX, tp.CursY: %d,%d\n", tp.CursX, tp.CursY)

					if x == tp.CursX && y == tp.CursY {
						textRend.CurrX -= textRend.CharWid
						drawCurrentChar('_', clipSpan)
					}
				}
			}

			// draw at the end of line if needed
			if y == tp.CursY && tp.CursX == len(line) {
				if tp.IsEditable && curs.Visible == true {
					drawCurrentChar('_', clipSpan)
				}
			}

			tp.GoToLeftEdge()
		}

		textRend.CurrY -= textRend.CharHei // go down a line height
	}

	tp.Bar.DrawVertical(2, 11)
}

func (tp *TextPanel) DrawBackground(atlasCellX, atlasCellY float32) {
	rad := textRend.ScreenRad
	sp := textRend.UvSpan
	u := float32(atlasCellX) * sp
	v := float32(atlasCellY) * sp

	gl.Normal3f(0, 0, 1)

	// bottom left   0, 1
	gl.TexCoord2f(u, v+sp)
	gl.Vertex3f(-rad, tp.Bottom, 0)

	// bottom right   1, 1
	gl.TexCoord2f(u+sp, v+sp)
	gl.Vertex3f(rad, tp.Bottom, 0)

	// top right   1, 0
	gl.TexCoord2f(u+sp, v)
	gl.Vertex3f(rad, tp.Top, 0)

	// top left   0, 0
	gl.TexCoord2f(u, v)
	gl.Vertex3f(-rad, tp.Top, 0)
}

func (tp *TextPanel) ScrollIfMouseOver(mousePixelDeltaY float64) {
	if tp.ContainsMouseCursor() {
		// y position increment (for bar) in gl space
		yInc := float32(mousePixelDeltaY) * textRend.PixelHei
		tp.Bar.ScrollThisMuch(tp, yInc)
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
	tp.Body = append(tp.Body, "// ------- variable declarations -------")
	tp.Body = append(tp.Body, "var myVar int32")
	tp.Body = append(tp.Body, "var a int32 = 42")
	tp.Body = append(tp.Body, "var b int32 = 58")
	tp.Body = append(tp.Body, "")
	tp.Body = append(tp.Body, "// ------- function declarations -------")
	tp.Body = append(tp.Body, "func myFunc(a,b){")
	tp.Body = append(tp.Body, "}")
	tp.Body = append(tp.Body, "func nuthaFunc (a, b) {")
	tp.Body = append(tp.Body, "        var myLocal int32")
	tp.Body = append(tp.Body, "    }    ")
	tp.Body = append(tp.Body, "")
	tp.Body = append(tp.Body, "// ------- function calls -------")
	tp.Body = append(tp.Body, "    sub32(7, 9)")
	tp.Body = append(tp.Body, "sub32(4,8)")
	tp.Body = append(tp.Body, "mult32(7, 7)")
	tp.Body = append(tp.Body, "mult32(3,5)")
	tp.Body = append(tp.Body, "div32(8,2)")
	tp.Body = append(tp.Body, "div32(15,  3)")
	tp.Body = append(tp.Body, "add32(2,3)")
	tp.Body = append(tp.Body, "add32(a, b)")
	tp.Body = append(tp.Body, "")

	for i := 0; i < 22; i++ {
		tp.Body = append(tp.Body, fmt.Sprintf("%d: put lots of text on screen", i))
	}
}

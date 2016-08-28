package main

import (
	"fmt"
	"github.com/go-gl/gl/v2.1/gl"
)

type TextPanel struct {
	Top             float32
	Bottom          float32
	Left            float32
	Right           float32
	NumCharsX       int
	NumCharsY       int
	OffsetY         float32
	LenOfOffscreenY float32
	Selection       SelectionRange
	Bar             ScrollBar
	Body            []string
}

func (tp *TextPanel) Init() {
	tp.Selection = SelectionRange{}
	tp.Selection.Init()

	tp.Left = -textRend.ScreenRad
	tp.Right = textRend.ScreenRad

	if tp.Top == 0 {
		tp.Top = textRend.ScreenRad
	}

	tp.Bottom = tp.Top - float32(tp.NumCharsY)*textRend.chHei

	fmt.Printf("TextPanel.Init()    t: %.2f, b: %.2f\n", tp.Top, tp.Bottom)
}

func (tp *TextPanel) Draw() {
	tp.DrawBackground(11, 13)

	for _, line := range tp.Body {
		for _, c := range line {
			drawCurrentChar(c)
		}

		textRend.CurrX = -textRend.ScreenRad
		textRend.CurrY -= textRend.chHei
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

	textRend.CurrX += textRend.chWid
}

func (tp *TextPanel) ScrollIfMouseOver(mousePixelDeltaY float64) {
	if curs.MouseGlY < tp.Top && curs.MouseGlY > tp.Bottom {
		if curs.MouseGlX < tp.Right && curs.MouseGlX > tp.Left {
			tp.Bar.Scroll(mousePixelDeltaY) // FIXME to work with any TextPanel
		}
	}
}

func (tp *TextPanel) RemoveCharacter(fromUnderCursor bool) {
	if fromUnderCursor {
		if len(tp.Body[curs.Y]) > curs.X {
			tp.Body[curs.Y] = tp.Body[curs.Y][:curs.X] + tp.Body[curs.Y][curs.X+1:len(tp.Body[curs.Y])]
		}
	} else {
		if curs.X > 0 {
			tp.Body[curs.Y] = tp.Body[curs.Y][:curs.X-1] + tp.Body[curs.Y][curs.X:len(tp.Body[curs.Y])]
			curs.X--
		}
	}
}

func (tp *TextPanel) SetupDemoProgram() {
	tp.Body = append(tp.Body, "PRESS CTRL-R to RUN your program")
	tp.Body = append(tp.Body, "")
	tp.Body = append(tp.Body, "------- variable declarations -------")
	tp.Body = append(tp.Body, "var myVar int32")
	tp.Body = append(tp.Body, "var a int32 = 42")
	tp.Body = append(tp.Body, "var b int32 = 58")
	tp.Body = append(tp.Body, "")
	tp.Body = append(tp.Body, "------- function declarations -------")
	tp.Body = append(tp.Body, "func myFunc(a,b){")
	tp.Body = append(tp.Body, "}")
	tp.Body = append(tp.Body, "func nuthaFunc (a, b) {")
	tp.Body = append(tp.Body, "        var myLocal int32")
	tp.Body = append(tp.Body, "    }    ")
	tp.Body = append(tp.Body, "")
	tp.Body = append(tp.Body, "------- function calls -------")
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

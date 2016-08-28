package main

import (
	"fmt"
)

type TextPanel struct {
	NumCharsX       int
	NumCharsY       int
	OffsetY         float32
	LenOfOffscreenY float32
	Selection       SelectionRange
	Bar             ScrollBar
	Body            []string
}

func (tp *TextPanel) Init() {
	textRend.Init()

	tp.Selection = SelectionRange{}
	tp.Selection.Init()
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

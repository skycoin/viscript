/* TODO:
CTRL-ARROW moves by whole word
select range with mouse
" " " arrow keys
CTRL-HOME/END
PGUP/DN
RIGHT at end of line goes to next line
BACKSPACE/DELETE at the ends of lines
	pulls us up to prev line, or pulls up next line
H/V SCROLLBARS
	horizontal could be charHei thickness
	vertical could easily be a smaller rendering of the first ~40 chars?
		however not if we map the whole vertical space (when scrollspace is taller than screen),
		because this requires scaling the text.  and keeping the aspect ratio means ~40 (max)
		would alter the width of the scrollbar
*/

package main

import (
	"fmt"
	"github.com/go-gl/gl/v2.1/gl"
	"math"
)

// character
var chWid = float32(textRend.ScreenRad * 2 / float32(textRend.NumCharsX))
var chHei = float32(textRend.ScreenRad * 2 / float32(textRend.NumCharsY))
var chWidInPixels = int(float32(resX) / float32(textRend.NumCharsX))
var chHeiInPixels = int(float32(resY) / float32(textRend.NumCharsY))
var pixelWid = textRend.ScreenRad * 2 / float32(resX)
var pixelHei = textRend.ScreenRad * 2 / float32(resY)
var document = make([]string, 0)

// selection
// future consideration/fixme:
// need to sanitize start/end positions.
// since they may be beyond the last line character of the line.
// also, in addition to backspace/delete, typing any visible character should delete marked text.
// complication:   if they start or end on invalid characters (of the line string),
// the forward or backwards direction from there determines where the visible selection
// range starts/ends....
// will an end position always be defined (when value is NOT math.MaxUint32),
// when a START is?  because that determines where the first VISIBLY marked
// character starts
var selectionStartX = math.MaxUint32
var selectionStartY = math.MaxUint32
var selectionEndX = math.MaxUint32
var selectionEndY = math.MaxUint32
var selectingRangeOfText = false

var textRend = TextRenderer{}

func initDoc() {
	textRend.Init()

	document = append(document, "PRESS CTRL-R to RUN your program")
	document = append(document, "")
	document = append(document, "------- variable declarations -------")
	document = append(document, "var myVar int32")
	document = append(document, "var a int32 = 42")
	document = append(document, "var b int32 = 58")
	document = append(document, "")
	document = append(document, "------- function declarations -------")
	document = append(document, "func myFunc(a,b){")
	document = append(document, "}")
	document = append(document, "func nuthaFunc (a, b) {")
	document = append(document, "        var myLocal int32")
	document = append(document, "    }    ")
	document = append(document, "")
	document = append(document, "------- function calls -------")
	document = append(document, "    sub32(7, 9)")
	document = append(document, "sub32(4,8)")
	document = append(document, "mult32(7, 7)")
	document = append(document, "mult32(3,5)")
	document = append(document, "div32(8,2)")
	document = append(document, "div32(15,  3)")
	document = append(document, "add32(2,3)")
	document = append(document, "add32(a, b)")
	document = append(document, "")

	for i := 0; i < 22; i++ {
		document = append(document, fmt.Sprintf("%d: put lots of text on screen", i))
	}
}

type TextRenderer struct {
	UvSpan    float32
	ScreenRad float32 // entire screen radius (distance to edge
	// in the cardinal directions from the center, corners would be farther away)
	NumCharsX int
	NumCharsY int
	// current place renderer draws to
	CurrX float32
	CurrY float32
}

func (tr *TextRenderer) Init() {
	tr.UvSpan = float32(1.0) / 16
	tr.ScreenRad = float32(3)
	tr.NumCharsX = 80
	tr.NumCharsY = 25
	tr.CurrX = -tr.ScreenRad
	tr.CurrY = tr.ScreenRad
}

func (tr *TextRenderer) DrawAll() {
	for _, line := range document {
		for _, c := range line {
			drawCharAtCurrentPos(c)
		}

		tr.CurrX = -tr.ScreenRad
		tr.CurrY -= chHei
	}

	sb.DrawVertical(2, 11)
	tr.DrawCharAt('#', curs.MouseX, curs.MouseY)
	curs.Draw()

	tr.CurrX = -tr.ScreenRad
	tr.CurrY = tr.ScreenRad - code.OffsetY
}

func commonMovementKeyHandling() {
	if selectingRangeOfText {
		selectionEndX = curs.X
		selectionEndY = curs.Y
	} else { // arrow keys without shift gets rid selection
		selectionStartX = math.MaxUint32
		selectionStartY = math.MaxUint32
		selectionEndX = math.MaxUint32
		selectionEndY = math.MaxUint32
	}
}

func removeCharacter(fromUnderCursor bool) {
	if fromUnderCursor {
		if len(document[curs.Y]) > curs.X {
			document[curs.Y] = document[curs.Y][:curs.X] + document[curs.Y][curs.X+1:len(document[curs.Y])]
		}
	} else {
		if curs.X > 0 {
			document[curs.Y] = document[curs.Y][:curs.X-1] + document[curs.Y][curs.X:len(document[curs.Y])]
			curs.X--
		}
	}
}

func (tr *TextRenderer) DrawCharAt(letter rune, posX int, posY int) {
	x := int(letter) % 16
	y := int(letter) / 16
	sp := tr.UvSpan

	gl.Normal3f(0, 0, 1)

	gl.TexCoord2f(float32(x)*sp, float32(y)*sp+sp) // bl  0, 1
	gl.Vertex3f(-tr.ScreenRad+float32(posX)*chWid, tr.ScreenRad-float32(posY)*chHei-chHei, 0)

	gl.TexCoord2f(float32(x)*sp+sp, float32(y)*sp+sp) // br  1, 1
	gl.Vertex3f(-tr.ScreenRad+float32(posX)*chWid+chWid, tr.ScreenRad-float32(posY)*chHei-chHei, 0)

	gl.TexCoord2f(float32(x)*sp+sp, float32(y)*sp) // tr  1, 0
	gl.Vertex3f(-tr.ScreenRad+float32(posX)*chWid+chWid, tr.ScreenRad-float32(posY)*chHei, 0)

	gl.TexCoord2f(float32(x)*sp, float32(y)*sp) // tl  0, 0
	gl.Vertex3f(-tr.ScreenRad+float32(posX)*chWid, tr.ScreenRad-float32(posY)*chHei, 0)

	tr.CurrX += chWid
}

func drawCharAtCurrentPos(letter rune) {
	x := int(letter) % 16
	y := int(letter) / 16
	sp := textRend.UvSpan

	gl.Normal3f(0, 0, 1)

	gl.TexCoord2f(float32(x)*sp, float32(y)*sp+sp) // bl  0, 1
	gl.Vertex3f(textRend.CurrX, textRend.CurrY-chHei, 0)

	gl.TexCoord2f(float32(x)*sp+sp, float32(y)*sp+sp) // br  1, 1
	gl.Vertex3f(textRend.CurrX+chWid, textRend.CurrY-chHei, 0)

	gl.TexCoord2f(float32(x)*sp+sp, float32(y)*sp) // tr  1, 0
	gl.Vertex3f(textRend.CurrX+chWid, textRend.CurrY, 0)

	gl.TexCoord2f(float32(x)*sp, float32(y)*sp) // tl  0, 0
	gl.Vertex3f(textRend.CurrX, textRend.CurrY, 0)

	textRend.CurrX += chWid
}

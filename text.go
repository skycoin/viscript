/* TODO:
CTRL-ARROW moves by whole word
select range with mouse
" " " arrow keys
CTRL-HOME/END
PGUP/DN
RIGHT at end of line goes to next line
BACKSPACE/DELETE at the ends of lines
	pulls us up to prev line, or pulls up next line
*/

package main

import (
	"github.com/go-gl/gl/v2.1/gl"
	"math"
	"time"
)

// character
var uvSpan = float32(1.0) / 16
var rectRad = float32(3) // rectangular radius (distance to edge of app window
// in the cardinal directions from the center, corners would be farther away)
var curX = -rectRad // the current pos of the DRAWing cursor
var curY = rectRad
var chWid = float32(rectRad * 2 / 80) // char w
var chHei = float32(rectRad * 2 / 25) // char h
var pixelWid = int(float32(resX) / 80)
var pixelHei = int(float32(resY) / 25)
var mouseX int = 0 // char position of mouse pointer
var mouseY int = 0
var document = make([]string, 0)

// cursor
var nextBlinkChange = time.Now()
var cursVisible = true
var cursX = 0
var cursY = 0

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

func initDoc() {
	document = append(document, "testing init line")
}

func makeChars() {
	for _, line := range document {
		for _, c := range line {
			drawCharAtCurrentPos(c)
		}

		curX = -rectRad
		curY -= chHei
	}

	drawCharAt('#', mouseX, mouseY)
	drawCursorMaybe()

	curX = -rectRad
	curY = rectRad
}

func commonMovementKeyHandling() {
	if selectingRangeOfText {
		selectionEndX = cursX
		selectionEndY = cursY
	} else { // arrow keys without shift gets rid selection
		selectionStartX = math.MaxUint32
		selectionStartY = math.MaxUint32
		selectionEndX = math.MaxUint32
		selectionEndY = math.MaxUint32
	}
}

func removeCharacter(fromUnderCursor bool) {
	if fromUnderCursor {
		if len(document[cursY]) > cursX {
			document[cursY] = document[cursY][:cursX] + document[cursY][cursX+1:len(document[cursY])]
		}
	} else {
		if cursX > 0 {
			document[cursY] = document[cursY][:cursX-1] + document[cursY][cursX:len(document[cursY])]
			cursX--
		}
	}
}

func drawCursorMaybe() {
	if nextBlinkChange.Before(time.Now()) {
		nextBlinkChange = time.Now().Add(time.Millisecond * 170)
		cursVisible = !cursVisible
	}

	if cursVisible == true {
		drawCharAt('_', cursX, cursY)
	}
}

func drawCharAt(letter rune, posX int, posY int) {
	x := int(letter) % 16
	y := int(letter) / 16

	gl.Normal3f(0, 0, 1)

	gl.TexCoord2f(float32(x)*uvSpan, float32(y)*uvSpan+uvSpan) // bl  0, 1
	gl.Vertex3f(-rectRad+float32(posX)*chWid, rectRad-float32(posY)*chHei-chHei, 0)

	gl.TexCoord2f(float32(x)*uvSpan+uvSpan, float32(y)*uvSpan+uvSpan) // br  1, 1
	gl.Vertex3f(-rectRad+float32(posX)*chWid+chWid, rectRad-float32(posY)*chHei-chHei, 0)

	gl.TexCoord2f(float32(x)*uvSpan+uvSpan, float32(y)*uvSpan) // tr  1, 0
	gl.Vertex3f(-rectRad+float32(posX)*chWid+chWid, rectRad-float32(posY)*chHei, 0)

	gl.TexCoord2f(float32(x)*uvSpan, float32(y)*uvSpan) // tl  0, 0
	gl.Vertex3f(-rectRad+float32(posX)*chWid, rectRad-float32(posY)*chHei, 0)

	curX += chWid

	if curX >= rectRad {
		curX = -rectRad
		curY -= chHei
	}
}

func drawCharAtCurrentPos(letter rune) {
	x := int(letter) % 16
	y := int(letter) / 16

	gl.Normal3f(0, 0, 1)

	gl.TexCoord2f(float32(x)*uvSpan, float32(y)*uvSpan+uvSpan) // bl  0, 1
	gl.Vertex3f(curX, curY-chHei, 0)

	gl.TexCoord2f(float32(x)*uvSpan+uvSpan, float32(y)*uvSpan+uvSpan) // br  1, 1
	gl.Vertex3f(curX+chWid, curY-chHei, 0)

	gl.TexCoord2f(float32(x)*uvSpan+uvSpan, float32(y)*uvSpan) // tr  1, 0
	gl.Vertex3f(curX+chWid, curY, 0)

	gl.TexCoord2f(float32(x)*uvSpan, float32(y)*uvSpan) // tl  0, 0
	gl.Vertex3f(curX, curY, 0)

	curX += chWid

	if curX >= rectRad {
		curX = -rectRad
		curY -= chHei
	}
}

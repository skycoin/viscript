package main

import (
	//"fmt"
	"github.com/go-gl/gl/v2.1/gl"
	"github.com/go-gl/glfw/v3.1/glfw"
	"time"
)

var curs Cursors = Cursors{NextBlinkChange: time.Now(), Visible: true}

/*					2 cursors:

* mouse:           keeps updated to current pointer position
* text insert:     keyboard typing will be inserted here

 */

type Cursors struct {
	NextBlinkChange time.Time
	Visible         bool
	X               int // current text insert position (in character grid space)
	Y               int
	MouseX          int // current mouse position (in character grid space)
	MouseY          int
	MouseGlX        float32 // current mouse position in OpenGL space
	MouseGlY        float32
}

func (c *Cursors) Draw() {
	// mouse
	c.DrawCharAt('#', c.MouseX, c.MouseY)

	// text/char insert
	if c.NextBlinkChange.Before(time.Now()) {
		c.NextBlinkChange = time.Now().Add(time.Millisecond * 170)
		c.Visible = !c.Visible
	}

	if c.Visible == true {
		c.DrawCharAt('_', c.X, c.Y)
	}
}

func (c *Cursors) DrawCharAt(char rune, posX int, posY int) {
	x := int(char) % 16
	y := int(char) / 16
	w := textRend.chWid // char width
	h := textRend.chHei // char height
	sp := textRend.UvSpan
	rad := textRend.ScreenRad

	gl.Normal3f(0, 0, 1)

	gl.TexCoord2f(float32(x)*sp, float32(y)*sp+sp) // bl  0, 1
	gl.Vertex3f(-rad+float32(posX)*w, rad-float32(posY)*h-h, 0)

	gl.TexCoord2f(float32(x)*sp+sp, float32(y)*sp+sp) // br  1, 1
	gl.Vertex3f(-rad+float32(posX)*w+w, rad-float32(posY)*h-h, 0)

	gl.TexCoord2f(float32(x)*sp+sp, float32(y)*sp) // tr  1, 0
	gl.Vertex3f(-rad+float32(posX)*w+w, rad-float32(posY)*h, 0)

	gl.TexCoord2f(float32(x)*sp, float32(y)*sp) // tl  0, 0
	gl.Vertex3f(-rad+float32(posX)*w, rad-float32(posY)*h, 0)

	textRend.CurrX += w
}

func (c *Cursors) ConvertMouseClickToTextCursorPosition(button uint8, action uint8) {
	if glfw.MouseButton(button) == glfw.MouseButtonLeft &&
		glfw.Action(action) == glfw.Press {

		if c.MouseY < len(code.Body) {
			curs.Y = c.MouseY

			if c.MouseX <= len(code.Body[curs.Y]) {
				curs.X = c.MouseX
			} else {
				curs.X = len(code.Body[curs.Y])
			}
		} else {
			curs.Y = len(code.Body) - 1
		}
	}
}

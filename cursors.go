package main

import (
	"fmt"
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
	MouseX          int // current position in character grid space (units/cells)
	MouseY          int
	MouseGlX        float32 // current mouse position in "analog" OpenGL space
	MouseGlY        float32
}

func (c *Cursors) Update() {
	// mouse
	// (don't know that we ever need this, but it only worked properly
	// when we had just 1 panel that took up the whole client area of this app)
	//c.DrawCharAt('#', c.MouseX, c.MouseY)

	// text/char insert
	if c.NextBlinkChange.Before(time.Now()) {
		c.NextBlinkChange = time.Now().Add(time.Millisecond * 170)
		c.Visible = !c.Visible
	}
}

func (c *Cursors) UpdatePosition(x, y float32) {
	c.MouseGlX = -textRend.ScreenRad + x*textRend.PixelWid
	c.MouseGlY = textRend.ScreenRad - y*textRend.PixelHei
	foc := textRend.Focused
	panY := c.MouseGlY - foc.Top
	fmt.Printf("panY: %.2f\n", panY)
	c.MouseY = int(( /**/ -panY /*+ -foc.Bar.ScrollDistY*/) / textRend.CharHei)
	c.MouseX = int(x) / textRend.CharWidInPixels
	fmt.Printf("c.MouseY: %d\n", c.MouseY)

	if c.MouseY < 0 {
		c.MouseY = 0
	}

	if c.MouseY >= len(foc.Body) {
		c.MouseY = len(foc.Body) - 1
	}
}

// this function was designed for a single panel, and before scrolling was added
func (c *Cursors) DrawCharAt(char rune, posX, posY int) {
	rad := textRend.ScreenRad
	sp := textRend.UvSpan
	u := sp * float32(int(char)%16)
	v := sp * float32(int(char)/16)
	w := textRend.CharWid // char width
	h := textRend.CharHei // char height
	x := -rad + float32(posX)*w
	y := rad - float32(posY)*h

	gl.Normal3f(0, 0, 1)

	// bottom left
	gl.TexCoord2f(u, v+sp)
	gl.Vertex3f(x, y-h, 0)

	// bottom right
	gl.TexCoord2f(u+sp, v+sp)
	gl.Vertex3f(x+w, y-h, 0)

	// top right
	gl.TexCoord2f(u+sp, v)
	gl.Vertex3f(x+w, y, 0)

	// top left
	gl.TexCoord2f(u, v)
	gl.Vertex3f(x, y, 0)
}

func (c *Cursors) ConvertMouseClickToTextCursorPosition(button, action uint8) {
	if glfw.MouseButton(button) == glfw.MouseButtonLeft &&
		glfw.Action(action) == glfw.Press {

		foc := textRend.Focused

		if foc.IsEditable && foc.ContainsMouseCursor() {
			if !foc.Bar.DragHandleContainsMouseCursor() {
				fmt.Printf("ConvertMouseClickToTextCursorPosition --- c.MouseY: %d\n", c.MouseY)

				if c.MouseY < len(foc.Body) {
					foc.CursY = c.MouseY

					if c.MouseX <= len(foc.Body[foc.CursY]) {
						foc.CursX = c.MouseX
					} else {
						foc.CursX = len(foc.Body[foc.CursY])
					}
				} else {
					foc.CursY = len(foc.Body) - 1
				}
			}
		}
	}
}

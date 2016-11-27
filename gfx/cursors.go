package gfx

import (
	//"fmt"
	"github.com/corpusc/viscript/common"
	"github.com/go-gl/gl/v2.1/gl"
	"github.com/go-gl/glfw/v3.2/glfw"
	"time"
)

var Curs Cursors = Cursors{NextBlinkChange: time.Now(), Visible: true}

func MouseCursorIsInside(r *common.Rectangle) bool {
	if Curs.MouseGlY < r.Top && Curs.MouseGlY > r.Bottom {
		if Curs.MouseGlX < r.Right && Curs.MouseGlX > r.Left {
			return true
		}
	}

	return false
}

type Cursors struct {
	NextBlinkChange time.Time
	Visible         bool
	MouseGlX        float32 // current mouse position in OpenGL space
	MouseGlY        float32
}

func (c *Cursors) Update() {
	if c.NextBlinkChange.Before(time.Now()) {
		c.NextBlinkChange = time.Now().Add(time.Millisecond * 170)
		c.Visible = !c.Visible
	}
}

func (c *Cursors) UpdatePosition(x, y float32) {
	c.MouseGlX = -Rend.ClientExtentX + x*Rend.PixelWid
	c.MouseGlY = Rend.ClientExtentY - y*Rend.PixelHei
}

// this function was designed for a single panel, and before scrolling was added
func (c *Cursors) DEPRECATED_DrawCharAt(char rune, posX, posY int) {
	sp := Rend.UvSpan
	u := sp * float32(int(char)%16)
	v := sp * float32(int(char)/16)
	w := Rend.CharWid // char width
	h := Rend.CharHei // char height
	x := -Rend.ClientExtentX + float32(posX)*w
	y := Rend.ClientExtentY - float32(posY)*h

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

		foc := Rend.Focused

		if foc.IsEditable && foc.ContainsMouseCursorInsideOfScrollBars() {
			if foc.MouseY < len(foc.Body) {
				foc.CursY = foc.MouseY

				if foc.MouseX <= len(foc.Body[foc.CursY]) {
					foc.CursX = foc.MouseX
				} else {
					foc.CursX = len(foc.Body[foc.CursY])
				}
			} else {
				foc.CursY = len(foc.Body) - 1
			}
		}
	}
}

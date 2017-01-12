package mouse

import (
	//"fmt"
	"github.com/corpusc/viscript/app"
)

var GlX float32 // current mouse position in OpenGL space
var GlY float32
var PixelDelta app.Vec2F

// private
var prevPixelPos app.Vec2F

func CursorIsInside(r *app.Rectangle) bool {
	if GlY < r.Top && GlY > r.Bottom {
		if GlX < r.Right && GlX > r.Left {
			return true
		}
	}

	return false
}

func UpdatePosition(pos, extents, pixSize app.Vec2F) {
	GlX = -extents.X + pos.X*pixSize.X
	GlY = extents.Y - pos.Y*pixSize.Y
	PixelDelta.X = pos.X - prevPixelPos.X
	PixelDelta.Y = pos.Y - prevPixelPos.Y
	prevPixelPos.X = pos.X
	prevPixelPos.Y = pos.Y
}

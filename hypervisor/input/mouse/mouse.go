package mouse

import (
	//"fmt"
	"github.com/corpusc/viscript/app"
)

var (
	GlX               float32 // current mouse position in OpenGL space
	GlY               float32
	PixelDelta        app.Vec2F
	HoldingLeftButton bool

	// private
	pixelSize_    app.Vec2F
	prevPixelPos  app.Vec2F
	canvasExtents app.Vec2F
)

func SetSizes(extents, pixelSize app.Vec2F) {
	canvasExtents = extents
	pixelSize_ = pixelSize
}

func CursorIsInside(r *app.Rectangle) bool {
	if GlY < r.Top && GlY > r.Bottom {
		if GlX < r.Right && GlX > r.Left {
			return true
		}
	}

	return false
}

func UpdatePosition(pos app.Vec2F) {
	GlX = -canvasExtents.X + pos.X*pixelSize_.X
	GlY = canvasExtents.Y - pos.Y*pixelSize_.Y
	PixelDelta.X = pos.X - prevPixelPos.X
	PixelDelta.Y = pos.Y - prevPixelPos.Y
	prevPixelPos.X = pos.X
	prevPixelPos.Y = pos.Y
}

package mouse

import (
	"github.com/corpusc/viscript/app"
	"math"
)

var (
	GlX         float32 //current mouse position in OpenGL space
	GlY         float32
	PrevGlX     float32
	PrevGlY     float32
	PixelDelta  app.Vec2F
	HoldingLeft bool

	// private
	pixelSize_    app.Vec2F
	prevPixelPos  app.Vec2F
	canvasExtents app.Vec2F
	nearThresh    float64 //nearness threshold (how close pointer should be to the edge)
)

func Update(pos app.Vec2F) {
	PrevGlX = GlX
	PrevGlY = GlY
	GlX = -canvasExtents.X + pos.X*pixelSize_.X
	GlY = canvasExtents.Y - pos.Y*pixelSize_.Y
	PixelDelta.X = pos.X - prevPixelPos.X
	PixelDelta.Y = pos.Y - prevPixelPos.Y
	prevPixelPos.X = pos.X
	prevPixelPos.Y = pos.Y
}

func NearRight(bounds *app.Rectangle) bool {
	return math.Abs(float64(GlX-bounds.Right)) <= nearThresh
}

func NearBottom(bounds *app.Rectangle) bool {
	return math.Abs(float64(GlY-bounds.Bottom)) <= nearThresh
}

func IncreaseNearnessThreshold() {
	nearThresh = 10.0
}

func DecreaseNearnessThreshold() {
	nearThresh = 0.05
}

func PointerIsInside(r *app.Rectangle) bool {
	if GlY <= r.Top && GlY >= r.Bottom {
		if GlX <= r.Right && GlX >= r.Left {
			return true
		}
	}

	return false
}

func SetSizes(extents, pixelSize app.Vec2F) {
	canvasExtents = extents
	pixelSize_ = pixelSize
}

func GetScrollDeltaX() float32 {
	return PixelDelta.X * pixelSize_.X
}

func GetScrollDeltaY() float32 {
	return PixelDelta.Y * pixelSize_.Y
}

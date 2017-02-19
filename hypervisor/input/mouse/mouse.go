package mouse

import (
	//"fmt"
	"github.com/corpusc/viscript/app"
	"math"
)

var (
	GlX              float32 // current mouse position in OpenGL space
	GlY              float32
	PrevGlX          float32
	PrevGlY          float32
	PixelDelta       app.Vec2F
	HoldingLeft      bool
	IsNearRight      bool
	IsNearBottom     bool
	IsInsideTerminal bool
	Bounds           *app.Rectangle
	EdgeGlMaxAbs     float64 // how close cursor should be to the edge

	// private
	pixelSize_    app.Vec2F
	prevPixelPos  app.Vec2F
	canvasExtents app.Vec2F
)

func Update(pos app.Vec2F) {
	UpdatePosition(pos)
	IsNearRight = CursorIsNearRightEdge()
	IsNearBottom = CursorIsNearBottomEdge()
	IsInsideTerminal = CursorIsInside(Bounds)
}

func CursorIsNearRightEdge() bool {
	return math.Abs(float64(GlX-Bounds.Right)) <= EdgeGlMaxAbs
}

func CursorIsNearBottomEdge() bool {
	return math.Abs(float64(GlY-Bounds.Bottom)) <= EdgeGlMaxAbs
}

func IncreaseEdgeGlMaxAbs() {
	EdgeGlMaxAbs = 10.0
}

func DecreaseEdgeGlMaxAbs() {
	EdgeGlMaxAbs = 0.05
}

func CursorIsInside(r *app.Rectangle) bool {
	if GlY <= r.Top && GlY >= r.Bottom {
		if GlX <= r.Right && GlX >= r.Left {
			return true
		}
	}

	return false
}

func UpdatePosition(pos app.Vec2F) {
	PrevGlX = GlX
	PrevGlY = GlY
	GlX = -canvasExtents.X + pos.X*pixelSize_.X
	GlY = canvasExtents.Y - pos.Y*pixelSize_.Y
	PixelDelta.X = pos.X - prevPixelPos.X
	PixelDelta.Y = pos.Y - prevPixelPos.Y
	prevPixelPos.X = pos.X
	prevPixelPos.Y = pos.Y
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

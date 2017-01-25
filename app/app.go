package app

const Name = "V I S C R I P T"
const UvSpan = float32(1.0) / 16 // span of a tile/cell in uv space

// params: float value, negativemost, & positivemost bounds
func Clamp(f, negBoundary, posBoundary float32) float32 {
	if f < negBoundary {
		f = negBoundary
	}
	if f > posBoundary {
		f = posBoundary
	}

	return f
}

// params: Rectangle, negativemost, & positivemost bounds
func ClampLeftAndRightOf(r *Rectangle, negBoundary, posBoundary float32) *Rectangle {
	if r.Left < negBoundary {
		r.Left = negBoundary
	}
	if r.Right > posBoundary {
		r.Right = posBoundary
	}

	return r
}

type Vec2I struct {
	X int
	Y int
}

type Vec2UI32 struct {
	X uint32
	Y uint32
}

type Vec2F struct {
	X float32
	Y float32
}

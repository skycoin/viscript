package app

var Name = "V I S C R I P T"

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

type Vec2 struct {
	X float32
	Y float32
}

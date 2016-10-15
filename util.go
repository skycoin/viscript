package main

import (
//"fmt"
)

func MouseCursorIsInside(r *Rectangle) bool {
	if curs.MouseGlY < r.Top && curs.MouseGlY > r.Bottom {
		if curs.MouseGlX < r.Right && curs.MouseGlX > r.Left {
			return true
		}
	}

	return false
}

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

package main

import (
	//"fmt"
	"viscript/common"
)

var goldenRatio = 1.61803398875
var goldenPercentage = float32(goldenRatio / (goldenRatio + 1))

func MouseCursorIsInside(r *common.Rectangle) bool {
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
func ClampLeftAndRightOf(r *common.Rectangle, negBoundary, posBoundary float32) *common.Rectangle {
	if r.Left < negBoundary {
		r.Left = negBoundary
	}
	if r.Right > posBoundary {
		r.Right = posBoundary
	}

	return r
}

package gfx

import (
	//"fmt"
	"github.com/corpusc/viscript/common"
)

var goldenRatio = 1.61803398875
var goldenPercentage = float32(goldenRatio / (goldenRatio + 1))

func MouseCursorIsInside(r *common.Rectangle) bool {
	if Curs.MouseGlY < r.Top && Curs.MouseGlY > r.Bottom {
		if Curs.MouseGlX < r.Right && Curs.MouseGlX > r.Left {
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

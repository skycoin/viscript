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

package gl

import (
	"github.com/corpusc/viscript/cGfx"
	"github.com/go-gl/gl/v2.1/gl"
)

func SetColor(newColor []float32) {
	cGfx.PrevColor = cGfx.CurrColor
	cGfx.CurrColor = newColor
	gl.Materialfv(gl.FRONT, gl.AMBIENT_AND_DIFFUSE, &newColor[0])
}

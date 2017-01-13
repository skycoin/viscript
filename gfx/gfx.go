package gfx

import (
	"fmt"
	"github.com/corpusc/viscript/app"
	"github.com/corpusc/viscript/ui"
	"github.com/go-gl/gl/v2.1/gl"
)

var PrevColor []float32 // previous
var CurrColor []float32

var goldenRatio = 1.61803398875
var goldenFraction = float32(goldenRatio / (goldenRatio + 1))

// dimensions (in pixel units)
var InitAppWidth int32 = 800
var InitAppHeight int32 = 600
var CurrAppWidth int32 = InitAppWidth
var CurrAppHeight int32 = InitAppHeight
var longerDimension = float32(InitAppWidth) / float32(InitAppHeight)
var InitFrustum = &app.Rectangle{1, longerDimension, -1, -longerDimension}
var PrevFrustum = &app.Rectangle{InitFrustum.Top, InitFrustum.Right, InitFrustum.Bottom, InitFrustum.Left}
var CurrFrustum = &app.Rectangle{InitFrustum.Top, InitFrustum.Right, InitFrustum.Bottom, InitFrustum.Left}

// colors
var Black = []float32{0, 0, 0, 1}
var Blue = []float32{0, 0, 1, 1}
var Cyan = []float32{0, 0.5, 1, 1}
var Fuschia = []float32{0.6, 0.2, 0.3, 1}
var Gray = []float32{0.25, 0.25, 0.25, 1}
var GrayDark = []float32{0.15, 0.15, 0.15, 1}
var GrayLight = []float32{0.4, 0.4, 0.4, 1}
var Green = []float32{0, 1, 0, 1}
var Magenta = []float32{1, 0, 1, 1}
var Maroon = []float32{0.5, 0.03, 0.207, 1}
var MaroonDark = []float32{0.24, 0.014, 0.1035, 1}
var Orange = []float32{0.8, 0.35, 0, 1}
var Purple = []float32{0.6, 0, 0.8, 1}
var Red = []float32{1, 0, 0, 1}
var Tan = []float32{0.55, 0.47, 0.37, 1}
var Violet = []float32{0.4, 0.2, 1, 1}
var White = []float32{1, 1, 1, 1}
var Yellow = []float32{1, 1, 0, 1}

// ^^^
// as above, so below   (keep these synchronized)
// VVV

func SetColorFromText(s string) {
	switch s {
	case "<color=Black":
		SetColor(Black)
	case "<color=Blue":
		SetColor(Blue)
	case "<color=Cyan":
		SetColor(Cyan)
	case "<color=Fuschia":
		SetColor(Fuschia)
	case "<color=Gray":
		SetColor(Gray)
	case "<color=GrayDark":
		SetColor(GrayDark)
	case "<color=GrayLight":
		SetColor(GrayLight)
	case "<color=Green":
		SetColor(Green)
	case "<color=Magenta":
		SetColor(Magenta)
	case "<color=Maroon":
		SetColor(Maroon)
	case "<color=MaroonDark":
		SetColor(MaroonDark)
	case "<color=Orange":
		SetColor(Orange)
	case "<color=Purple":
		SetColor(Purple)
	case "<color=Red":
		SetColor(Red)
	case "<color=Tan":
		SetColor(Tan)
	case "<color=Violet":
		SetColor(Violet)
	case "<color=White":
		SetColor(White)
	case "<color=Yellow":
		SetColor(Yellow)
	}
}

func init() {
	fmt.Println("gfx.init()")

	// one-time setup
	PrevColor = GrayDark
	CurrColor = GrayDark

	MaxCharsX = 80
	MaxCharsY = 25
	DistanceFromOrigin = 3
	UvSpan = float32(1.0) / 16 // how much uv a pixel spans

	// things that are resized later
	CanvasExtents.X = DistanceFromOrigin * longerDimension
	CanvasExtents.Y = DistanceFromOrigin
	CharWid = float32(CanvasExtents.X*2) / float32(MaxCharsX)
	CharHei = float32(CanvasExtents.Y*2) / float32(MaxCharsY)
	CharWidInPixels = int(float32(CurrAppWidth) / float32(MaxCharsX))
	CharHeiInPixels = int(float32(CurrAppHeight) / float32(MaxCharsY))
	PixelSize.X = CanvasExtents.X * 2 / float32(CurrAppWidth)
	PixelSize.Y = CanvasExtents.Y * 2 / float32(CurrAppHeight)

	// MORE one-time setup
	ui.MainMenu.SetSize(GetMenuSizedRect())
}

func SetColor(newColor []float32) {
	PrevColor = CurrColor
	CurrColor = newColor
	gl.Materialfv(gl.FRONT, gl.AMBIENT_AND_DIFFUSE, &newColor[0])
}

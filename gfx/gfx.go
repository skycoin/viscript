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

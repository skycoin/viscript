package cGfx

import (
	"fmt"
	"github.com/corpusc/viscript/app"
)

/*
------------- STARTING UP A NEW *CLEAN* VERSION OF GFX -------------

.....which (ATM) has only "app" as a dependency
("app" has NONE)

migrate things over here until we can delete/replace
the current "gfx" package with THIS

--------------------------------------------------------------------
*/

// pics (tile positions in atlas)
var Pic_GradientBorder = app.Vec2I{11, 13}
var Pic_PixelCheckerboard = app.Vec2I{2, 11}
var Pic_SquareInTheMiddle = app.Vec2I{14, 15}
var Pic_DoubleLinesHorizontal = app.Vec2I{13, 12}
var Pic_DoubleLinesVertical = app.Vec2I{10, 11}

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

var PrevColor []float32 // previous
var CurrColor []float32

// dimensions (in pixel units)
var InitAppWidth int32 = 800
var InitAppHeight int32 = 600
var CurrAppWidth int32 = InitAppWidth
var CurrAppHeight int32 = InitAppHeight
var longerDimension = float32(InitAppWidth) / float32(InitAppHeight)
var InitFrustum = &app.Rectangle{1, longerDimension, -1, -longerDimension}
var PrevFrustum = &app.Rectangle{InitFrustum.Top, InitFrustum.Right, InitFrustum.Bottom, InitFrustum.Left}
var CurrFrustum = &app.Rectangle{InitFrustum.Top, InitFrustum.Right, InitFrustum.Bottom, InitFrustum.Left}

var (
	// distance from the center to an edge of the app's root/client area
	// ....in the cardinal directions from the center, corners would be farther away)
	CanvasExtents      app.Vec2F
	PixelSize          app.Vec2F
	DistanceFromOrigin float32
	CharWid            float32
	CharHei            float32
	CharWidInPixels    int
	CharHeiInPixels    int
	UvSpan             float32 // looking into 16/16 atlas/grid of character tiles
	// FIXME: below is no longer a maximum of what fits on a max-sized panel (taking up the whole app window) anymore.
	// 		but is still used as a guide for sizes
	MaxCharsX int // this is used to give us proportions like an 80x25 text console screen, ....
	MaxCharsY int // ....from a DistanceFromOrigin*2-by-DistanceFromOrigin*2 gl space
	// current position renderer draws to
	CurrX float32
	CurrY float32
)

func init() {
	fmt.Println("cGfx.init()")

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
}

func SetSize() {
	fmt.Printf("cGfx.SetSize() - CanvasExtents.X: %.2f\n", CanvasExtents.X)
	*PrevFrustum = *CurrFrustum

	CurrFrustum.Right = float32(CurrAppWidth) / float32(InitAppWidth) * InitFrustum.Right
	CurrFrustum.Left = -CurrFrustum.Right
	CurrFrustum.Top = float32(CurrAppHeight) / float32(InitAppHeight) * InitFrustum.Top
	CurrFrustum.Bottom = -CurrFrustum.Top

	fmt.Printf("cGfx.SetSize() - PrevFrustum.Left: %.3f\n", PrevFrustum.Left)
	fmt.Printf("cGfx.SetSize() - CurrFrustum.Left: %.3f\n", CurrFrustum.Left)

	CanvasExtents.X = DistanceFromOrigin * CurrFrustum.Right
	CanvasExtents.Y = DistanceFromOrigin * CurrFrustum.Top

	// things that weren't initialized in this func
}

func GetMenuSizedRect() *app.Rectangle {
	return &app.Rectangle{
		CanvasExtents.Y,
		CanvasExtents.X,
		CanvasExtents.Y - CharHei,
		-CanvasExtents.X}
}

func DrawAll() {
	Curs.Update()
}

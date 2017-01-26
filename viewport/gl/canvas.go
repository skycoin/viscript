// canvas == the whole "client" area of the graphical OpenGL window (of root app)
package gl

import (
	"github.com/corpusc/viscript/app"
)

// dimensions (in pixel units)
var InitAppWidth int = 800 // initial/startup size (when resizing, compare against this)
var InitAppHeight int = 600
var CurrAppWidth = int32(InitAppWidth) // current
var CurrAppHeight = int32(InitAppHeight)
var longerDimension = float32(InitAppWidth) / float32(InitAppHeight)
var InitFrustum = &app.Rectangle{1, longerDimension, -1, -longerDimension}
var PrevFrustum = &app.Rectangle{InitFrustum.Top, InitFrustum.Right, InitFrustum.Bottom, InitFrustum.Left}
var CurrFrustum = &app.Rectangle{InitFrustum.Top, InitFrustum.Right, InitFrustum.Bottom, InitFrustum.Left}

var (
	// distance from the center to top & bottom edges of the canvas
	DistanceFromOrigin float32 = 1
)

var (
	CanvasExtents   app.Vec2F
	PixelSize       app.Vec2F
	CharWid         float32
	CharHei         float32
	CharWidInPixels int
	CharHeiInPixels int
	// FIXME: below is no longer a maximum of what fits on a max-sized panel (taking up the whole app window) anymore.
	// 		but is still used as a guide for sizes
	MaxCharsX int // this is used to give us proportions similar to a 80x25 text console screen
	MaxCharsY int
	// current position renderer draws to
	CurrX float32
	CurrY float32
)

func init() {
	println("canvas.init()")
	// one-time setup
	PrevColor = GrayDark
	CurrColor = GrayDark

	// FIXME: these are NO LONGER used as maximums, but more as guidelines for text size
	MaxCharsX = 80
	MaxCharsY = 25

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
	MainMenu.SetSize(GetMenuSizedRect())
}

func GetMenuSizedRect() *app.Rectangle {
	return &app.Rectangle{
		CanvasExtents.Y,
		CanvasExtents.X,
		CanvasExtents.Y - CharHei,
		-CanvasExtents.X}
}

func SetSize(x, y int32) {
	println("canvas.SetSize()")
	*PrevFrustum = *CurrFrustum
	CurrAppWidth = x
	CurrAppHeight = y

	CurrFrustum.Right = float32(CurrAppWidth) / float32(InitAppWidth) * InitFrustum.Right
	CurrFrustum.Left = -CurrFrustum.Right
	CurrFrustum.Top = float32(CurrAppHeight) / float32(InitAppHeight) * InitFrustum.Top
	CurrFrustum.Bottom = -CurrFrustum.Top

	CanvasExtents.X = DistanceFromOrigin * CurrFrustum.Right
	CanvasExtents.Y = DistanceFromOrigin * CurrFrustum.Top

	// things that weren't initialized in this func
	MainMenu.SetSize(GetMenuSizedRect())
}

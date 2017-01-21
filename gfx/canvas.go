// canvas == the whole "client" area of the graphical OpenGL window
package gfx

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
	DistanceFromOrigin float32 = 3
)

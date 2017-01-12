package gl

import (
	"fmt"
	//"log"

	//"bytes"
	//"math"
	//"strconv"

	//"encoding/binary"

	"github.com/corpusc/viscript/gfx"
	//"github.com/corpusc/viscript/msg"
	//"github.com/corpusc/viscript/script"
	"github.com/go-gl/glfw/v3.2/glfw"
)

func InitMiscEvents(w *glfw.Window) {
	w.SetFramebufferSizeCallback(onFramebufferSize)
}

//direct, not wrapped
//NOT AN INPUT EVENT
func onFramebufferSize(w *glfw.Window, width, height int) {
	fmt.Printf("onFramebufferSize() - width, height: %d, %d\n", width, height)
	gfx.CurrAppWidth = int32(width)
	gfx.CurrAppHeight = int32(height)
	gfx.SetSize()
}

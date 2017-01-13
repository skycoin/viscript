package gl

import (
	//"fmt"
	//"log"

	//"bytes"
	//"math"
	//"strconv"

	"github.com/corpusc/viscript/msg"
	//"github.com/corpusc/viscript/script"
	"github.com/go-gl/glfw/v3.2/glfw"
)

func InitMiscEvents(w *glfw.Window) {
	w.SetFramebufferSizeCallback(onFrameBufferSize)
}

func onFrameBufferSize(w *glfw.Window, width, height int) {
	m := msg.MessageFrameBufferSize{uint32(width), uint32(height)}
	InputEvents <- msg.Serialize(msg.TypeFrameBufferSize, m)
}

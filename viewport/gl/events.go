package gl

import (
	"github.com/corpusc/viscript/msg"
	"github.com/go-gl/glfw/v3.2/glfw"
)

func InitMiscEvents(w *glfw.Window) {
	w.SetFramebufferSizeCallback(onFrameBufferSize)
}

func onFrameBufferSize(w *glfw.Window, width, height int) {
	msg.SerializeAndDispatch(
		InputEvents,
		msg.TypeFrameBufferSize,
		msg.MessageFrameBufferSize{X: uint32(width), Y: uint32(height)})
}

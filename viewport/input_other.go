package viewport

import (
	"github.com/corpusc/viscript/msg"
	"github.com/corpusc/viscript/viewport/gl"
)

func onFrameBufferSize(m msg.MessageFrameBufferSize) {
	if DebugPrintInputEvents {
		print("TypeFrameBufferSize")
		showUInt32("X", m.X)
		showUInt32("Y", m.Y)
		println()
	}

	gl.SetSize(int32(m.X), int32(m.Y))
}

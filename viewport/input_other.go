package viewport

import (
	"github.com/corpusc/viscript/msg"
	"github.com/corpusc/viscript/viewport/gl"
)

func onFrameBufferSize(m msg.MessageFrameBufferSize) {
	gl.SetSize(int32(m.X), int32(m.Y))
}

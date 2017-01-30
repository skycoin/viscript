package viewport

import (
	"fmt"
	"github.com/corpusc/viscript/msg"
	"github.com/corpusc/viscript/viewport/gl"
)

func onFrameBufferSize(m msg.MessageFrameBufferSize) {
	fmt.Printf("onFrameBufferSize() - x, y: %d, %d\n", m.X, m.Y)
	gl.SetSize(int32(m.X), int32(m.Y))
}

package hypervisor

import (
	"bytes"
	"fmt"
	"github.com/corpusc/viscript/app"
	"github.com/corpusc/viscript/cGfx"
	"github.com/corpusc/viscript/gl"
	"github.com/corpusc/viscript/hypervisor/input/mouse"
	"github.com/corpusc/viscript/msg"
	"github.com/go-gl/glfw/v3.2/glfw"
	"math"
)

func onChar(m msg.MessageOnCharacter) {
	InsertRuneIntoDocument("Rune", m.Rune)
	//script.Process(false)
}

// triggered both by moving **AND*** by pressing buttons
func onMouseCursorPos(m msg.MessageMousePos) {
	x := float32(m.X)
	y := float32(m.Y)

	mouse.UpdatePosition(
		app.Vec2F{x, y},
		cGfx.CanvasExtents,
		cGfx.PixelSize) // state update

	// rendering update
	if /* LMB held */ gl.GlfwWindow.GetMouseButton(glfw.MouseButtonLeft) == glfw.Press {
		ScrollTermThatHasMousePointer(mouse.PixelDelta.X, mouse.PixelDelta.Y)
	}
}

func onMouseScroll(m msg.MessageMouseScroll) {
	var delta float32 = 30

	if eitherControlKeyHeld() { // horizontal ability from 1D scrolling
		ScrollTermThatHasMousePointer(float32(m.Y)*-delta, 0)
	} else { // can handle both x & y for 2D scrolling
		ScrollTermThatHasMousePointer(float32(m.X)*delta, float32(m.Y)*-delta)
	}
}

func onFrameBufferSize(m msg.MessageFrameBufferSize) {
	fmt.Printf("onFrameBufferSize() - x, y: %d, %d\n", m.X, m.Y)
	cGfx.CurrAppWidth = int32(m.X)
	cGfx.CurrAppHeight = int32(m.Y)
	cGfx.SetSize()
	SetSize()
}

func eitherControlKeyHeld() bool { // FIXME: bake into msg when serialized (or use bundled Mod field for 1s that have it)
	if gl.GlfwWindow.GetKey(glfw.KeyLeftControl) == glfw.Press || gl.GlfwWindow.GetKey(glfw.KeyRightControl) == glfw.Press {
		return true
	} else {
		return false
	}
}

// WARNING: given arguments must be in range
func insert(slice []string, index int, value string) []string {
	slice = slice[0 : len(slice)+1]      // grow the slice by one element
	copy(slice[index+1:], slice[index:]) // move the upper part of the slice out of the way and open a hole
	slice[index] = value
	return slice
}

// similar to insert method, instead moves current slice element and appends to one above
func remove(slice []string, index int, value string) []string {
	slice = append(slice[:index], slice[index+1:]...)
	slice[index-1] = slice[index-1] + value
	return slice
}

func movedCursorSoUpdateDependents() {
	foc := Focused

	// autoscroll to keep cursor visible
	ls := float32(foc.CursX) * cGfx.CharWid // left side (of cursor, in virtual space)
	rs := ls + cGfx.CharWid                 // right side ^

	if ls < foc.BarHori.ScrollDelta {
		foc.BarHori.ScrollDelta = ls
	}

	if rs > foc.BarHori.ScrollDelta+foc.Content.Width() {
		foc.BarHori.ScrollDelta = rs - foc.Content.Width()
	}

	// --- Selection Marking ---
	//
	// when SM is made functional,
	// we should probably detect whether cursor
	// position should update Start_ or End_ at this point.
	// rather than always making that the "end".
	// i doubt marking forwards or backwards will ever alter what is
	// done with the selection

	if foc.Selection.CurrentlySelecting {
		foc.Selection.EndX = foc.CursX
		foc.Selection.EndY = foc.CursY
	} else { // moving cursor without shift gets rid of selection
		foc.Selection.StartX = math.MaxUint32
		foc.Selection.StartY = math.MaxUint32
		foc.Selection.EndX = math.MaxUint32
		foc.Selection.EndY = math.MaxUint32
	}
}

func getSlice(wBuf *bytes.Buffer, err error) (data []byte) {
	data = make([]byte, 0)

	if err != nil {
		fmt.Println("binary.Write failed:", err)
	} else {
		b := wBuf.Bytes()

		for i := 0; i < wBuf.Len(); i++ {
			data = append(data, b[i])
		}
	}

	return
}

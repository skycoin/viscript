package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"github.com/go-gl/glfw/v3.1/glfw"
	"math"
)

const PREFIX_SIZE = 5 // guaranteed minimum size of every message (4 for length & 1 for type)
var events = make(chan []byte, 256)
var prevMousePixelX float64
var prevMousePixelY float64
var mousePixelDeltaX float64
var mousePixelDeltaY float64

func initInputEvents(w *glfw.Window) {
	w.SetCharCallback(onChar)
	w.SetKeyCallback(onKey)
	w.SetMouseButtonCallback(onMouseButton)
	w.SetScrollCallback(onMouseScroll)
	w.SetCursorPosCallback(onMouseCursorPos)
}

func pollEventsAndHandleAnInput(w *glfw.Window) {
	glfw.PollEvents()

	// (at the moment we have no reason to poll/detect keys being held)
	//if w.GetKey(glfw.KeyEscape) == glfw.Press {
	//	fmt.Println("PRESSED ESCape")
	//}
}

func onMouseCursorPos(w *glfw.Window, x float64, y float64) {
	curs.MouseGlX = -textRend.ScreenRad + float32(x)*textRend.pixelWid
	curs.MouseGlY = textRend.ScreenRad - float32(y)*textRend.pixelHei
	curs.MouseX = int(x) / textRend.chWidInPixels
	curs.MouseY = int(y) / textRend.chHeiInPixels
	mousePixelDeltaX = x - prevMousePixelX
	mousePixelDeltaY = y - prevMousePixelY
	prevMousePixelX = x
	prevMousePixelY = y

	if /* LMB held */ w.GetMouseButton(glfw.MouseButtonLeft) == glfw.Press {
		code.ScrollIfMouseOver(mousePixelDeltaY)
		cons.ScrollIfMouseOver(mousePixelDeltaY)
	}

	// build message
	content := append(getBytesOfFloat64(x), getBytesOfFloat64(y)...)
	dispatchWithPrefix(content, MessageMousePos)
}

func onMouseScroll(w *glfw.Window, xOff float64, yOff float64) {
	code.ScrollIfMouseOver(yOff * -30)
	cons.ScrollIfMouseOver(yOff * -30)

	// build message
	content := append(getBytesOfFloat64(xOff), getBytesOfFloat64(yOff)...)
	dispatchWithPrefix(content, MessageMouseScroll)
}

// apparently every time this is fired, a mouse position event is ALSO fired
func onMouseButton(
	window *glfw.Window,
	b glfw.MouseButton,
	action glfw.Action,
	mod glfw.ModifierKey) {

	if action != glfw.Press {
		switch glfw.MouseButton(b) {
		case glfw.MouseButtonLeft:
			for _, pan := range textRend.Panels {
				if pan.ContainsMouseCursor() {
					textRend.Focused = pan
				}
			}
		default:
		}
	}

	// build message
	content := append(getByteOfUInt8(uint8(b)), getByteOfUInt8(uint8(action))...)
	content = append(content, getByteOfUInt8(uint8(mod))...)
	dispatchWithPrefix(content, MessageMouseButton)
}

func onChar(w *glfw.Window, char rune) {
	dispatchWithPrefix(getBytesOfRune(char), MessageCharacter)
}

// WEIRD BEHAVIOUR OF KEY EVENTS.... for a PRESS, you can get a
// shift/alt/ctrl/super event through the "mod" variable,
// (see the top of "action == glfw.Press" section for an example)
// regardless of left/right key used.
// BUT for RELEASE, the "mod" variable will NOT tell you what key it is!
// you will have to find the specific left or right key that was released via
// the "action" variable!
func onKey(
	w *glfw.Window,
	key glfw.Key,
	scancode int,
	action glfw.Action,
	mod glfw.ModifierKey) {

	if action == glfw.Release {
		switch key {

		case glfw.KeyEscape:
			w.SetShouldClose(true)

		case glfw.KeyLeftShift:
			fallthrough
		case glfw.KeyRightShift:
			fmt.Println("done selecting")
			code.Selection.CurrentlySelecting = false // FIXME to work with whichever has focus
			// TODO?  possibly flip around if selectionStart comes after selectionEnd in the page flow?

		case glfw.KeyLeftControl:
			fallthrough
		case glfw.KeyRightControl:
			fmt.Println("Control RELEASED")

		case glfw.KeyLeftAlt:
			fallthrough
		case glfw.KeyRightAlt:
			fmt.Println("Alt RELEASED")

		case glfw.KeyLeftSuper:
			fallthrough
		case glfw.KeyRightSuper:
			fmt.Println("'Super' modifier key RELEASED")
		}
	} else { // glfw.Repeat   or   glfw.Press
		switch mod {
		case glfw.ModShift:
			fmt.Println("start selecting")
			code.Selection.CurrentlySelecting = true // FIXME to work on any TextPanel
			code.Selection.StartX = curs.X
			code.Selection.StartY = curs.Y
		}

		switch key {
		case glfw.KeyEnter:
			startOfLine := code.Body[curs.Y][:curs.X]
			restOfLine := code.Body[curs.Y][curs.X:len(code.Body[curs.Y])]
			//restOfDoc := code.Body[curs.Y+1 : len(code.Body)]
			//startOfDoc := code.Body[:curs.Y]

			newDoc := make([]string, 0)

			for i := 0; i < len(code.Body); i++ {
				//fmt.Printf("___[%d]: %s\n", i, code.Body[i])

				if i == curs.Y {
					newDoc = append(newDoc, startOfLine)
					newDoc = append(newDoc, restOfLine)
				} else {
					newDoc = append(newDoc, code.Body[i])
				}
			}

			code.Body = newDoc
			curs.X = 0
			curs.Y++

		case glfw.KeyHome:
			commonMovementKeyHandling()
			curs.X = 0
		case glfw.KeyEnd:
			commonMovementKeyHandling()
			curs.X = len(code.Body[curs.Y])
		case glfw.KeyUp:
			commonMovementKeyHandling()

			if curs.Y > 0 {
				curs.Y--

				if curs.X > len(code.Body[curs.Y]) {
					curs.X = len(code.Body[curs.Y])
				}
			}
		case glfw.KeyDown:
			commonMovementKeyHandling()

			if curs.Y < len(code.Body)-1 {
				curs.Y++

				if curs.X > len(code.Body[curs.Y]) {
					curs.X = len(code.Body[curs.Y])
				}
			}
		case glfw.KeyLeft:
			commonMovementKeyHandling()

			if curs.X == 0 {
				if curs.Y > 0 {
					curs.Y--
					curs.X = len(code.Body[curs.Y])
				}
			} else {
				if mod == glfw.ModControl {
					curs.X = getWordSkipPos(curs.X, -1)
				} else {
					curs.X--
				}
			}
		case glfw.KeyRight:
			commonMovementKeyHandling()

			if curs.X < len(code.Body[curs.Y]) {
				if mod == glfw.ModControl {
					curs.X = getWordSkipPos(curs.X, 1)
				} else {
					curs.X++
				}
			}
		case glfw.KeyBackspace:
			code.RemoveCharacter(false)
		case glfw.KeyDelete:
			code.RemoveCharacter(true)
		}
	}

	// build message
	content := getByteOfUInt8(uint8(key))
	content = append(content, getBytesOfSInt32(int32(scancode))...)
	content = append(content, getByteOfUInt8(uint8(action))...)
	content = append(content, getByteOfUInt8(uint8(mod))...)
	dispatchWithPrefix(content, MessageKey)
}

func dispatchWithPrefix(content []byte, msgType uint8) {
	//prefix := make([]byte, PREFIX_SIZE)
	prefix := append(
		getBytesOfUInt32(uint32(len(content))+PREFIX_SIZE),
		getByteOfUInt8(msgType)...)

	events <- append(prefix, content...)
}

func getWordSkipPos(xIn int, change int) int {
	peekPos := xIn

	for {
		peekPos += change

		if peekPos < 0 {
			return 0
		}

		if peekPos >= len(code.Body[curs.Y]) {
			return len(code.Body[curs.Y])
		}

		if string(code.Body[curs.Y][peekPos]) == " " {
			return peekPos
		}
	}
}

func commonMovementKeyHandling() {
	if code.Selection.CurrentlySelecting {
		code.Selection.EndX = curs.X
		code.Selection.EndY = curs.Y
	} else { // arrow keys without shift gets rid selection
		code.Selection.StartX = math.MaxUint32
		code.Selection.StartY = math.MaxUint32
		code.Selection.EndX = math.MaxUint32
		code.Selection.EndY = math.MaxUint32
	}
}

// the rest of these getBytesOfType() funcs are identical except for the value type
func getBytesOfRune(value rune) (data []byte) {
	wBuf := new(bytes.Buffer)
	err := binary.Write(wBuf, binary.LittleEndian, value)
	data = getSlice(wBuf, err)
	return
}

func getByteOfUInt8(value uint8) (data []byte) {
	wBuf := new(bytes.Buffer)
	err := binary.Write(wBuf, binary.LittleEndian, value)
	data = getSlice(wBuf, err)
	return
}

func getBytesOfSInt32(value int32) (data []byte) {
	wBuf := new(bytes.Buffer)
	err := binary.Write(wBuf, binary.LittleEndian, value)
	data = getSlice(wBuf, err)
	return
}

func getBytesOfUInt32(value uint32) (data []byte) {
	wBuf := new(bytes.Buffer)
	err := binary.Write(wBuf, binary.LittleEndian, value)
	data = getSlice(wBuf, err)
	return
}

func getBytesOfFloat64(value float64) (data []byte) {
	wBuf := new(bytes.Buffer)
	err := binary.Write(wBuf, binary.LittleEndian, value)
	data = getSlice(wBuf, err)
	return
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

package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"github.com/go-gl/glfw/v3.1/glfw"
)

const PREFIX_SIZE = 5 // guaranteed minimum size of every message (4 for length & 1 for type)
var events = make(chan []byte, 256)

func initInputEvents(w *glfw.Window) {
	w.SetCharCallback(onChar)
	w.SetKeyCallback(onKey)
	w.SetMouseButtonCallback(onMouseButton)
	w.SetScrollCallback(onMouseScroll)
	w.SetCursorPosCallback(onMouseCursorPos)
}

func pollEventsAndHandleAnInput(window *glfw.Window) {
	glfw.PollEvents()

	// (at the moment we have no reason to poll/detect keys being held)
	//if window.GetKey(glfw.KeyEscape) == glfw.Press {
	//	fmt.Println("PRESSED ESCape")
	//}
}

var prevMousePixelX;
var prevMousePixelY;
var mousePixelDeltaX;
var mousePixelDeltaY;
func onMouseCursorPos(w *glfw.Window, x float64, y float64) {
	mouseX = int(x) / chWidInPixels
	mouseY = int(y) / chHeiInPixels
	mousePixelDeltaX = x - prevMousePixelX
	mousePixelDeltaY = y - prevMousePixelY
	prevMousePixelX = x;
	prevMousePixelY = y;

	// build message
	content := append(getBytesOfFloat64(x), getBytesOfFloat64(y)...)
	dispatchWithPrefix(content, MessageMousePos)
}

func onMouseScroll(w *glfw.Window, xOff float64, yOff float64) {
	//fmt.Println("onScroll()")

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

	//fmt.Println("onMouseButton()")

	if action != glfw.Press {
		switch glfw.MouseButton(b) {
		case glfw.MouseButtonLeft:
		default:
		}
	}

	// build message
	content := append(getBytesOfUInt8(uint8(b)), getBytesOfUInt8(uint8(action))...)
	content = append(content, getBytesOfUInt8(uint8(mod))...)
	dispatchWithPrefix(content, MessageMouseButton)
}

func dispatchWithPrefix(content []byte, msgType uint8) {
	//prefix := make([]byte, PREFIX_SIZE)
	prefix := append(
		getBytesOfUInt32(uint32(len(content))+PREFIX_SIZE),
		getBytesOfUInt8(msgType)...)

	events <- append(prefix, content...)
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
			selectingRangeOfText = false
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
			selectingRangeOfText = true
			selectionStartX = cursX
			selectionStartY = cursY
		}

		switch key {
		case glfw.KeyEnter:
			startOfLine := document[cursY][:cursX]
			restOfLine := document[cursY][cursX:len(document[cursY])]
			//restOfDoc := document[cursY+1 : len(document)]
			//startOfDoc := document[:cursY]

			newDoc := make([]string, 0)

			for i := 0; i < len(document); i++ {
				//fmt.Printf("___[%d]: %s\n", i, document[i])

				if i == cursY {
					newDoc = append(newDoc, startOfLine)
					newDoc = append(newDoc, restOfLine)
				} else {
					newDoc = append(newDoc, document[i])
				}
			}

			document = newDoc
			cursX = 0
			cursY++

		case glfw.KeyHome:
			commonMovementKeyHandling()
			cursX = 0
		case glfw.KeyEnd:
			commonMovementKeyHandling()
			cursX = len(document[cursY])
		case glfw.KeyUp:
			commonMovementKeyHandling()

			if cursY > 0 {
				cursY--

				if cursX > len(document[cursY]) {
					cursX = len(document[cursY])
				}
			}
		case glfw.KeyDown:
			commonMovementKeyHandling()

			if cursY < len(document)-1 {
				cursY++

				if cursX > len(document[cursY]) {
					cursX = len(document[cursY])
				}
			}
		case glfw.KeyLeft:
			commonMovementKeyHandling()

			if cursX == 0 {
				if cursY > 0 {
					cursY--
					cursX = len(document[cursY])
				}
			} else {
				if mod == glfw.ModControl {
					cursX = getWordSkipPos(cursX, -1)
				} else {
					cursX--
				}
			}
		case glfw.KeyRight:
			commonMovementKeyHandling()

			if cursX < len(document[cursY]) {
				if mod == glfw.ModControl {
					cursX = getWordSkipPos(cursX, 1)
				} else {
					cursX++
				}
			}
		case glfw.KeyBackspace:
			removeCharacter(false)
		case glfw.KeyDelete:
			removeCharacter(true)
		}
	}

	// build message
	content := getBytesOfUInt8(uint8(key))
	content = append(content, getBytesOfSInt32(int32(scancode))...)
	content = append(content, getBytesOfUInt8(uint8(action))...)
	content = append(content, getBytesOfUInt8(uint8(mod))...)
	dispatchWithPrefix(content, MessageKey)
}

func getWordSkipPos(xIn int, change int) (cursX int) {
	peekPos := xIn

	for {
		peekPos += change

		if peekPos < 0 {
			cursX = 0
			return
		}

		if peekPos >= len(document[cursY]) {
			cursX = len(document[cursY])
			return
		}

		if string(document[cursY][peekPos]) == " " {
			cursX = peekPos
			return
		}
	}
}

func onChar(w *glfw.Window, char rune) {
	dispatchWithPrefix(getBytesOfRune(char), MessageCharacter)
}

// the rest of these getBytesOfType() funcs are identical except for the value type
func getBytesOfRune(value rune) (data []byte) {
	wBuf := new(bytes.Buffer)
	err := binary.Write(wBuf, binary.LittleEndian, value)
	data = getSlice(wBuf, err)
	return
}

func getBytesOfUInt8(value uint8) (data []byte) {
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

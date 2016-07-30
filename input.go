package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"github.com/go-gl/glfw/v3.1/glfw"
)

func initInputEvents(w *glfw.Window) {
	w.SetCharCallback(onChar)
	w.SetKeyCallback(onKey)
	w.SetMouseButtonCallback(onMouseButton)
	w.SetScrollCallback(onMouseScroll)
	w.SetCursorPosCallback(onMouseCursorPos)
}

func pollEventsAndHandleAnInput(window *glfw.Window) {
	glfw.PollEvents()

	// (at the moment we have no reason to detect keys being held)
	// poll a particular key state
	//if window.GetKey(glfw.KeyEscape) == glfw.Press {
	//	fmt.Println("PRESSED ESCape")
	//	window.SetShouldClose(true)
	//}
}

const PREFIX_SIZE = 5 // guaranteed minimum size of every message (4 for length & 1 for type)

func onMouseCursorPos(w *glfw.Window, x float64, y float64) {
	//fmt.Println("onMouseCursorPos()")
	mouseX = int(x) / pixelWid
	mouseY = int(y) / pixelHei

	// build message
	content := append(getBytesOfFloat64(x), getBytesOfFloat64(y)...)
	buildPrefix(content, MessageMousePos)
}

func onMouseScroll(w *glfw.Window, xOff float64, yOff float64) {
	//fmt.Println("onScroll()")

	// build message
	content := append(getBytesOfFloat64(xOff), getBytesOfFloat64(yOff)...)
	buildPrefix(content, MessageMouseScroll)
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
	buildPrefix(content, MessageMouseButton)
}

func buildPrefix(content []byte, msgType uint8) {
	//prefix := make([]byte, PREFIX_SIZE)
	prefix := append(
		getBytesOfUInt32(uint32(len(content))+PREFIX_SIZE),
		getBytesOfUInt8(msgType)...)

	processMessage(append(prefix, content...))
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
			cursX = 0
			cursY++
			document = append(document, "")
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
				cursX--
			}
		case glfw.KeyRight:
			commonMovementKeyHandling()

			if cursX < len(document[cursY]) {
				cursX++
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
	buildPrefix(content, MessageKey)
}

func onChar(w *glfw.Window, char rune) {
	// when a Unicode character is input.
	//fmt.Printf("onChar(): %c\n", char)
	document[cursY] = document[cursY][:cursX] + string(char) + document[cursY][cursX:len(document[cursY])]
	cursX++

	// build message
	buildPrefix(getBytesOfRune(char), MessageCharacter)
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

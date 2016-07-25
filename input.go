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

var prefixLen uint32 = 0
var contentLen uint32 = 0 // current send message CONTENT index, for iterating across funcs

func onMouseCursorPos(w *glfw.Window, x float64, y float64) {
	//fmt.Println("onMouseCursorPos()")
	mouseX = int(x) / pixelWid
	mouseY = int(y) / pixelHei

	// build message
	content := make([]byte, 0) // dynamic size
	addFloat64(x, content)
	addFloat64(y, content)

	buildPrefix(content, MessageMousePos)
}

func onMouseScroll(w *glfw.Window, xOff float64, yOff float64) {
	//fmt.Println("onScroll()")

	// build message
	content := make([]byte, 0) // dynamic size
	addFloat64(xOff, content)
	addFloat64(yOff, content)

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
	content := make([]byte, 0) // dynamic size
	addMouseButton(b, content)
	addAction(action, content)
	addModifierKey(mod, content)

	buildPrefix(content, MessageMouseButton)
}

func buildPrefix(content []byte, msgType uint8) {
	prefix := make([]byte, PREFIX_SIZE)
	addPREfixUInt32(uint32(len(content))+PREFIX_SIZE, prefix)
	contentLen = 0
	addPREfixUInt8(msgType, prefix)
	prefixLen = 0

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
	mods glfw.ModifierKey) {

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
	} else { // action might be glfw.Press, glfw.? (repeated key (held) ), or ?
		switch mods {
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
	var length uint32 = 5
	message := make([]byte, length)
	addPREfixUInt32(length, message)
	addPREfixUInt8(MessageKey, message)
	contentLen = 0

	processMessage(message)
}

func onChar(w *glfw.Window, char rune) {
	// when a Unicode character is input.
	//fmt.Printf("onChar(): %c\n", char)
	document[cursY] = document[cursY][:cursX] + string(char) + document[cursY][cursX:len(document[cursY])]
	cursX++

	// build message
	var length uint32 = 5
	message := make([]byte, length)
	addPREfixUInt32(length, message)
	addPREfixUInt8(MessageCharacter, message)
	contentLen = 0

	processMessage(message)
}

func addPREfixUInt32(value uint32, data []byte) {
	wBuf := new(bytes.Buffer)
	err := binary.Write(wBuf, binary.LittleEndian, value)

	if err != nil {
		fmt.Println("binary.Write failed:", err)
	} else {
		for i := 0; i < wBuf.Len(); i++ {
			data[prefixLen] = wBuf.Bytes()[i] // optimizeme?  converts->[]byte each iteration
			prefixLen++
		}
	}
}

func addPREfixUInt8(value uint8, data []byte) {
	wBuf := new(bytes.Buffer)
	err := binary.Write(wBuf, binary.LittleEndian, value)

	if err != nil {
		fmt.Println("binary.Write failed:", err)
	} else {
		for i := 0; i < wBuf.Len(); i++ {
			data[prefixLen] = wBuf.Bytes()[i]
			prefixLen++
		}
	}
}

// the rest of these add[type]() functions are identical except for the value type
func addFloat64(value float64, data []byte) {
	wBuf := new(bytes.Buffer)
	err := binary.Write(wBuf, binary.LittleEndian, value)
	copyBytes(wBuf, data, err)
}

func addMouseButton(value glfw.MouseButton, data []byte) {
	wBuf := new(bytes.Buffer)
	err := binary.Write(wBuf, binary.LittleEndian, value)
	copyBytes(wBuf, data, err)
}

func addAction(value glfw.Action, data []byte) {
	wBuf := new(bytes.Buffer)
	err := binary.Write(wBuf, binary.LittleEndian, value)
	copyBytes(wBuf, data, err)
}

func addModifierKey(value glfw.ModifierKey, data []byte) {
	wBuf := new(bytes.Buffer)
	err := binary.Write(wBuf, binary.LittleEndian, value)
	copyBytes(wBuf, data, err)
}

func copyBytes(wBuf *bytes.Buffer, data []byte, err error) {
	if err != nil {
		fmt.Println("binary.Write failed:", err)
	} else {
		b := wBuf.Bytes()

		for i := 0; i < wBuf.Len(); i++ {
			data = append(data, b[i])
		}
	}
}

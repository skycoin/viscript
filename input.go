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
	curs.UpdatePosition(float32(x), float32(y))
	mousePixelDeltaX = x - prevMousePixelX
	mousePixelDeltaY = y - prevMousePixelY
	prevMousePixelX = x
	prevMousePixelY = y

	if /* LMB held */ w.GetMouseButton(glfw.MouseButtonLeft) == glfw.Press {
		textRend.ScrollFocusedPanel(mousePixelDeltaY)
	}

	// build message
	content := append(getBytesOfFloat64(x), getBytesOfFloat64(y)...)
	dispatchWithPrefix(content, MessageMousePos)
}

func onMouseScroll(w *glfw.Window, xOff float64, yOff float64) {
	textRend.ScrollFocusedPanel(yOff * -30)

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

	if action == glfw.Press {
		switch glfw.MouseButton(b) {
		case glfw.MouseButtonLeft:
			for _, pan := range textRend.Panels {
				if pan.ContainsMouseCursor() {
					textRend.Focused = &pan
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
			code.Selection.StartX = curs.TextX
			code.Selection.StartY = curs.TextY
		}

		switch key {

		case glfw.KeyEnter:
			startOfLine := code.Body[curs.TextY][:curs.TextX]
			restOfLine := code.Body[curs.TextY][curs.TextX:len(code.Body[curs.TextY])]
			fmt.Printf("startOfLine: \"%s\"\n", startOfLine)
			fmt.Printf(" restOfLine: \"%s\"\n", restOfLine)
			code.Body[curs.TextY] = startOfLine
			code.Body = insert(code.Body, curs.TextY+1, restOfLine)

			/*
				newDoc := make([]string, 0)

				for i := 0; i < len(code.Body); i++ {
					//fmt.Printf("___[%d]: %s\n", i, code.Body[i])

					if i == curs.TextY {
						newDoc = append(newDoc, startOfLine)
						newDoc = append(newDoc, restOfLine)
					} else {
						newDoc = append(newDoc, code.Body[i])
					}
				}

				code.Body = newDoc
			*/
			curs.TextX = 0
			curs.TextY++

		case glfw.KeyHome:
			commonMovementKeyHandling()
			curs.TextX = 0
		case glfw.KeyEnd:
			commonMovementKeyHandling()
			curs.TextX = len(code.Body[curs.TextY])
		case glfw.KeyUp:
			commonMovementKeyHandling()

			if curs.TextY > 0 {
				curs.TextY--

				if curs.TextX > len(code.Body[curs.TextY]) {
					curs.TextX = len(code.Body[curs.TextY])
				}
			}
		case glfw.KeyDown:
			commonMovementKeyHandling()

			if curs.TextY < len(code.Body)-1 {
				curs.TextY++

				if curs.TextX > len(code.Body[curs.TextY]) {
					curs.TextX = len(code.Body[curs.TextY])
				}
			}
		case glfw.KeyLeft:
			commonMovementKeyHandling()

			if curs.TextX == 0 {
				if curs.TextY > 0 {
					curs.TextY--
					curs.TextX = len(code.Body[curs.TextY])
				}
			} else {
				if mod == glfw.ModControl {
					curs.TextX = getWordSkipPos(curs.TextX, -1)
				} else {
					curs.TextX--
				}
			}
		case glfw.KeyRight:
			commonMovementKeyHandling()

			if curs.TextX < len(code.Body[curs.TextY]) {
				if mod == glfw.ModControl {
					curs.TextX = getWordSkipPos(curs.TextX, 1)
				} else {
					curs.TextX++
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

// must be in range
func insert(slice []string, index int, value string) []string {
	slice = slice[0 : len(slice)+1]      // grow the slice by one element
	copy(slice[index+1:], slice[index:]) // move the upper part of the slice out of the way and open a hole
	slice[index] = value
	return slice
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

		if peekPos >= len(code.Body[curs.TextY]) {
			return len(code.Body[curs.TextY])
		}

		if string(code.Body[curs.TextY][peekPos]) == " " {
			return peekPos
		}
	}
}

func commonMovementKeyHandling() {
	if code.Selection.CurrentlySelecting {
		code.Selection.EndX = curs.TextX
		code.Selection.EndY = curs.TextY
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

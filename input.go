package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"github.com/corpusc/viscript/gfx"
	"github.com/corpusc/viscript/parser"
	"github.com/corpusc/viscript/ui"
	"github.com/go-gl/glfw/v3.2/glfw"
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
	w.SetFramebufferSizeCallback(onFramebufferSize)
}

func onFramebufferSize(w *glfw.Window, width, height int) {
	fmt.Printf("onFramebufferSize() - width, height: %d, %d\n", width, height)
	gfx.CurrAppWidth = int32(width)
	gfx.CurrAppHeight = int32(height)
	gfx.Rend.SetSize()
}

func onMouseCursorPos(w *glfw.Window, x float64, y float64) {
	gfx.Curs.UpdatePosition(float32(x), float32(y))
	mousePixelDeltaX = x - prevMousePixelX
	mousePixelDeltaY = y - prevMousePixelY
	prevMousePixelX = x
	prevMousePixelY = y

	if /* LMB held */ w.GetMouseButton(glfw.MouseButtonLeft) == glfw.Press {
		gfx.Rend.ScrollPanelThatIsHoveredOver(mousePixelDeltaX, mousePixelDeltaY)
	}

	// build message
	content := append(getBytesOfFloat64(x), getBytesOfFloat64(y)...)
	dispatchWithPrefix(content, MessageMousePos)
}

func onMouseScroll(w *glfw.Window, xOff, yOff float64) {
	var delta float64 = 30

	// if horizontal
	if w.GetKey(glfw.KeyLeftShift) == glfw.Press || w.GetKey(glfw.KeyRightShift) == glfw.Press {
		gfx.Rend.ScrollPanelThatIsHoveredOver(yOff*-delta, 0)
	} else {
		gfx.Rend.ScrollPanelThatIsHoveredOver(xOff*delta, yOff*-delta)
	}

	// build message
	content := append(getBytesOfFloat64(xOff), getBytesOfFloat64(yOff)...)
	dispatchWithPrefix(content, MessageMouseScroll)
}

// apparently every time this is fired, a mouse position event is ALSO fired
func onMouseButton(
	w *glfw.Window,
	b glfw.MouseButton,
	action glfw.Action,
	mod glfw.ModifierKey) {

	if action == glfw.Press {
		switch glfw.MouseButton(b) {
		case glfw.MouseButtonLeft:
			// respond to button push
			if gfx.MouseCursorIsInside(ui.MainMenu.Rect) {
				for _, bu := range ui.MainMenu.Buttons {
					if gfx.MouseCursorIsInside(bu.Rect) {
						bu.Activated = !bu.Activated

						if bu.Activated && bu.Name == "Run" {
							parser.Parse()
						}

						gfx.Con.Add(fmt.Sprintf("%s toggled\n", bu.Name))
					}
				}
			} else {
				// respond to click in text panel
				for _, pan := range gfx.Rend.Panels {
					if pan.ContainsMouseCursor() {
						pan.RespondToMouseClick()
					}
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

// WEIRD BEHAVIOUR OF KEY EVENTS.... for a PRESS, you can detect a
// shift/alt/ctrl/super key through the "mod" variable,
// (see the top of "action == glfw.Press" section for an example)
// regardless of left/right key used.
// BUT for RELEASE, the "mod" variable will NOT tell you what key it is!
// so you will have to handle both left & right mod keys via the "action" variable!
func onKey(
	w *glfw.Window,
	key glfw.Key,
	scancode int,
	action glfw.Action,
	mod glfw.ModifierKey) {

	foc := gfx.Rend.Focused

	if action == glfw.Release {
		switch key {

		case glfw.KeyEscape:
			w.SetShouldClose(true)

		case glfw.KeyLeftShift:
			fallthrough
		case glfw.KeyRightShift:
			fmt.Println("done selecting")
			foc.Selection.CurrentlySelecting = false // TODO?  possibly flip around if selectionStart comes after selectionEnd in the page flow?

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
			foc.Selection.CurrentlySelecting = true
			foc.Selection.StartX = foc.CursX
			foc.Selection.StartY = foc.CursY
		}

		switch key {

		case glfw.KeyEnter:
			b := foc.Body
			startOfLine := b[foc.CursY][:foc.CursX]
			restOfLine := b[foc.CursY][foc.CursX:len(b[foc.CursY])]
			//fmt.Printf("startOfLine: \"%s\"\n", startOfLine)
			//fmt.Printf(" restOfLine: \"%s\"\n", restOfLine)
			b[foc.CursY] = startOfLine
			//fmt.Printf("foc.CursX: \"%d\"  -  foc.CursY: \"%d\"\n", foc.CursX, foc.CursY)
			foc.Body = insert(b, foc.CursY+1, restOfLine)

			foc.CursX = 0
			foc.CursY++

			if foc.CursY >= len(b) {
				foc.CursY = len(b) - 1
			}

		case glfw.KeyHome:
			commonMovementKeyHandling()
			foc.CursX = 0
		case glfw.KeyEnd:
			commonMovementKeyHandling()
			foc.CursX = len(foc.Body[foc.CursY])
		case glfw.KeyUp:
			commonMovementKeyHandling()

			if foc.CursY > 0 {
				foc.CursY--

				if foc.CursX > len(foc.Body[foc.CursY]) {
					foc.CursX = len(foc.Body[foc.CursY])
				}
			}
		case glfw.KeyDown:
			commonMovementKeyHandling()

			if foc.CursY < len(foc.Body)-1 {
				foc.CursY++

				if foc.CursX > len(foc.Body[foc.CursY]) {
					foc.CursX = len(foc.Body[foc.CursY])
				}
			}
		case glfw.KeyLeft:
			commonMovementKeyHandling()

			if foc.CursX == 0 {
				if foc.CursY > 0 {
					foc.CursY--
					foc.CursX = len(foc.Body[foc.CursY])
				}
			} else {
				if mod == glfw.ModControl {
					foc.CursX = getWordSkipPos(foc.CursX, -1)
				} else {
					foc.CursX--
				}
			}
		case glfw.KeyRight:
			commonMovementKeyHandling()

			if foc.CursX < len(foc.Body[foc.CursY]) {
				if mod == glfw.ModControl {
					foc.CursX = getWordSkipPos(foc.CursX, 1)
				} else {
					foc.CursX++
				}
			}
		case glfw.KeyBackspace:
			foc.RemoveCharacter(false)
		case glfw.KeyDelete:
			foc.RemoveCharacter(true)

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
	foc := gfx.Rend.Focused

	for {
		peekPos += change

		if peekPos < 0 {
			return 0
		}

		if peekPos >= len(foc.Body[foc.CursY]) {
			return len(foc.Body[foc.CursY])
		}

		if string(foc.Body[foc.CursY][peekPos]) == " " {
			return peekPos
		}
	}
}

func commonMovementKeyHandling() {
	foc := gfx.Rend.Focused

	if foc.Selection.CurrentlySelecting {
		foc.Selection.EndX = foc.CursX
		foc.Selection.EndY = foc.CursY
	} else { // arrow keys without shift gets rid selection
		foc.Selection.StartX = math.MaxUint32
		foc.Selection.StartY = math.MaxUint32
		foc.Selection.EndX = math.MaxUint32
		foc.Selection.EndY = math.MaxUint32
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

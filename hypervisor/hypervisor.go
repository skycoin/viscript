package hypervisor

import (
	"fmt"
	_ "image/png"
	/*
		"go/build"
		"runtime"
	*/
	"bytes"
	"encoding/binary"
	"github.com/corpusc/viscript/gfx"
	"github.com/corpusc/viscript/msg"
	"github.com/corpusc/viscript/script"
	"github.com/corpusc/viscript/ui"
	"github.com/go-gl/glfw/v3.2/glfw"
	"math"
	"strconv"
)

var Events = make(chan []byte, 256)
var prevMousePixelX float64
var prevMousePixelY float64
var mousePixelDeltaX float64
var mousePixelDeltaY float64

func InitInputEvents(w *glfw.Window) {
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

	//content := append(getBytesOfFloat64(x), getBytesOfFloat64(y)...)
	//dispatchWithPrefix(content, msg.TypeMousePos)
	// build message

	var m msg.MessageMousePos
	m.X = x
	m.Y = y
	DispatchEvent(msg.TypeMousePos, m)
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
	//content := append(getBytesOfFloat64(xOff), getBytesOfFloat64(yOff)...)
	//dispatchWithPrefix(content, msg.TypeMouseScroll)

	var m msg.MessageMouseScroll
	m.X = xOff
	m.Y = yOff
	DispatchEvent(msg.TypeMouseScroll, m)

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

						switch bu.Name {
						case "Run":
							if bu.Activated {
								script.Process(true)
							}
							break
						case "Testing Tree":
							if bu.Activated {
								script.Process(true)
								script.MakeTree()
							} else { // deactivated
								// remove all panels with trees
								b := gfx.Rend.Panels[:0]
								for _, pan := range gfx.Rend.Panels {
									if len(pan.Trees) < 1 {
										b = append(b, pan)
									}
								}
								gfx.Rend.Panels = b
								//fmt.Printf("len of b (from gfx.Rend.Panels) after removing ones with trees: %d\n", len(b))
								//fmt.Printf("len of gfx.Rend.Panels: %d\n", len(gfx.Rend.Panels))
							}
							break
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
	//content := append(getByteOfUInt8(uint8(b)), getByteOfUInt8(uint8(action))...)
	//content = append(content, getByteOfUInt8(uint8(mod))...)
	//dispatchWithPrefix(content, msg.TypeMouseButton)

	//MessageMouseButton
	var m msg.MessageMouseButton
	m.Button = uint8(b)
	m.Action = uint8(action)
	m.Mod = uint8(mod)
	DispatchEvent(msg.TypeMouseButton, m)
}

//FIX
func onChar(w *glfw.Window, char rune) {
	//dispatchWithPrefix(getBytesOfRune(char), msg.TypeCharacter)

	var m msg.MessageCharacter
	_ = m
	DispatchEvent(msg.TypeCharacter, m)
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
		//fmt.Println("Control RELEASED")
		case glfw.KeyLeftSuper:
			fallthrough
		case glfw.KeyRightSuper:
			fmt.Println("'Super' modifier key RELEASED")
		}
	} else { // glfw.Repeat   or   glfw.Press
		b := foc.TextBodies[0]

		CharWid := int32(gfx.Rend.CharWidInPixels)
		CharHei := int32(gfx.Rend.CharHeiInPixels)
		numOfCharsV := gfx.CurrAppHeight / CharHei
		numOfCharsH := gfx.CurrAppWidth / CharWid

		s := strconv.Itoa(int(numOfCharsV))

		fmt.Printf("Rectangle Right %s\n\n\n", s)
		
		
		switch mod {
		case glfw.ModShift:
			fmt.Println("started selecting")
			foc.Selection.CurrentlySelecting = true
			foc.Selection.StartX = foc.CursX
			foc.Selection.StartY = foc.CursY
		}

		switch key {

		case glfw.KeyEnter:
			startOfLine := b[foc.CursY][:foc.CursX]
			restOfLine := b[foc.CursY][foc.CursX:len(b[foc.CursY])]
			//fmt.Printf("startOfLine: \"%s\"\n", startOfLine)
			//fmt.Printf(" restOfLine: \"%s\"\n", restOfLine)
			b[foc.CursY] = startOfLine
			//fmt.Printf("foc.CursX: \"%d\"  -  foc.CursY: \"%d\"\n", foc.CursX, foc.CursY)
			b = insert(b, foc.CursY+1, restOfLine)

			foc.CursX = 0
			foc.CursY++
			foc.TextBodies[0] = b

			if foc.CursY >= len(b) {
				foc.CursY = len(b) - 1
			}

		case glfw.KeyHome:
			if w.GetKey(glfw.KeyLeftControl) == glfw.Press || w.GetKey(glfw.KeyRightControl) == glfw.Press {

			} else {
				commonMovementKeyHandling()
				foc.CursX = 0
			}
		case glfw.KeyEnd:
			commonMovementKeyHandling()
			foc.CursX = len(b[foc.CursY])
		case glfw.KeyUp:
			commonMovementKeyHandling()

			if foc.CursY > 0 {
				foc.CursY--

				if foc.CursX > len(b[foc.CursY]) {
					foc.CursX = len(b[foc.CursY])
				}
			}
		case glfw.KeyDown:
			commonMovementKeyHandling()

			if foc.CursY < len(b)-1 {
				if numOfCharsV < (int32(foc.CursY) + 1) {
					gfx.Rend.ScrollPanelThatIsHoveredOver(0, float64(CharHei))
				}
				foc.CursY++

				if foc.CursX > len(b[foc.CursY]) {
					foc.CursX = len(b[foc.CursY])
				}
			}
		case glfw.KeyLeft:
			commonMovementKeyHandling()

			if foc.CursX == 0 {
				if foc.CursY > 0 {
					foc.CursY--
					foc.CursX = len(b[foc.CursY])
				}
			} else {
				if mod == glfw.ModControl {
					foc.CursX = getWordSkipPos(foc.CursX, -1)
				} else {
					if (numOfCharsH - int32(foc.CursX)) > (int32(foc.CursX) + 4) {
						gfx.Rend.ScrollPanelThatIsHoveredOver(float64(-CharWid), 0)
					}
					foc.CursX--
				}
			}
		case glfw.KeyRight:

			commonMovementKeyHandling()

			if foc.CursX < len(b[foc.CursY]) {
				if mod == glfw.ModControl {
					foc.CursX = getWordSkipPos(foc.CursX, 1)
				} else {
					fmt.Println(numOfCharsH)
					if numOfCharsH < (int32(foc.CursX) + 4){
						gfx.Rend.ScrollPanelThatIsHoveredOver(float64(CharWid), 0)
					}
					foc.CursX++
				}
			}
		case glfw.KeyBackspace:
			if foc.CursX == 0 {
				b = remove(b, foc.CursY, b[foc.CursY])
				foc.TextBodies[0] = b
				foc.CursY--
				foc.CursX = len(b[foc.CursY])

			} else {
				foc.RemoveCharacter(false)
			}

		case glfw.KeyDelete:
			foc.RemoveCharacter(true)
			fmt.Println("Key Deleted")

		}

		script.Process(false)
	}

	// build message
	//content := getByteOfUInt8(uint8(key))
	//content = append(content, getBytesOfSInt32(int32(scancode))...)
	//content = append(content, getByteOfUInt8(uint8(action))...)
	//content = append(content, getByteOfUInt8(uint8(mod))...)
	//dispatchWithPrefix(content, msg.TypeKey)

	var m msg.MessageKey
	m.Key = uint8(key)
	m.Scan = uint32(scancode)
	m.Action = uint8(action)
	m.Mod = uint8(mod)
	DispatchEvent(msg.TypeKey, m)
}

// must be in range
func insert(slice []string, index int, value string) []string {
	slice = slice[0 : len(slice)+1]      // grow the slice by one element
	copy(slice[index+1:], slice[index:]) // move the upper part of the slice out of the way and open a hole
	slice[index] = value
	return slice
}

<<<<<<< HEAD
// similar to insert method, instead moves current slice element and appends to one above
func remove(slice []string, index int, value string) []string {
	slice = append(slice[:index], slice[index+1:]...)
	slice[index-1] = slice[index-1] + value
	return slice
}

=======
/*
>>>>>>> corpusc/master
func dispatchWithPrefix(content []byte, msgType uint16) {
	prefix := append(
		getBytesOfUInt32(uint32(len(content))+msg.PREFIX_SIZE),
		getBytesOfUInt16(msgType)...)

	Events <- append(prefix, content...)
}
*/

func DispatchEvent(msgType uint16, event interface{}) {
	b := msg.Serialize(msgType, event)
	Events <- b
}

func getWordSkipPos(xIn int, change int) int {
	peekPos := xIn
	foc := gfx.Rend.Focused
	b := foc.TextBodies[0]

	for {
		peekPos += change

		if peekPos < 0 {
			return 0
		}

		if peekPos >= len(b[foc.CursY]) {
			return len(b[foc.CursY])
		}

		if string(b[foc.CursY][peekPos]) == " " {
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

/*
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

func getBytesOfUInt16(value uint16) (data []byte) {
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
*/

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

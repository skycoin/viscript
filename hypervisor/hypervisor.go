package hypervisor

import (
	"fmt"
	"image"
	"image/draw"
	_ "image/png"
	"log"
	"os"
	/*
		"go/build"
		"runtime"
	*/
	"bytes"
	"encoding/binary"
	"github.com/corpusc/viscript/app"
	"github.com/corpusc/viscript/gfx"
	"github.com/corpusc/viscript/msg"
	"github.com/corpusc/viscript/script"
	"github.com/corpusc/viscript/ui"
	"github.com/go-gl/gl/v2.1/gl"
	"github.com/go-gl/glfw/v3.2/glfw"
	"math"




	//"github.com/corpusc/viscript/msg" //message structs are here

)

const PREFIX_SIZE = 5 // guaranteed minimum size of every message (4 for length & 1 for type)
const (
	MessageMousePos    = iota // 0
	MessageMouseScroll        // 1
	MessageMouseButton        // 2
	MessageCharacter
	MessageKey
)



var Events = make(chan []byte, 256)
var prevMousePixelX float64
var prevMousePixelY float64
var mousePixelDeltaX float64
var mousePixelDeltaY float64

var (
	Texture   uint32
	rotationX float32
	rotationY float32
)

func InitRenderer() {
	fmt.Println("initRenderer()")

	gl.Enable(gl.DEPTH_TEST)
	gl.Enable(gl.LIGHTING)
	//gl.Enable(gl.ALPHA_TEST)

	gl.ClearColor(0.5, 0.5, 0.5, 0.0)
	gl.ClearDepth(1)
	gl.DepthFunc(gl.LEQUAL)

	ambient := []float32{0.5, 0.5, 0.5, 1}
	diffuse := []float32{1, 1, 1, 1}
	lightPosition := []float32{-5, 5, 10, 0}
	gl.Lightfv(gl.LIGHT0, gl.AMBIENT, &ambient[0])
	gl.Lightfv(gl.LIGHT0, gl.DIFFUSE, &diffuse[0])
	gl.Lightfv(gl.LIGHT0, gl.POSITION, &lightPosition[0])
	gl.Enable(gl.LIGHT0)

	gl.MatrixMode(gl.PROJECTION)
	gl.LoadIdentity()
	setFrustum(gfx.InitFrustum)
	//gl.Frustum(-1, 1, -1, 1, 1.0, 10.0)
	//gl.Frustum(left, right, bottom, top, zNear, zFar)
	gl.MatrixMode(gl.MODELVIEW)
	gl.LoadIdentity()

	// future FIXME: finished app would not have a demo program loaded on startup?
	script.Process(false)
}

func setFrustum(r *app.Rectangle) {
	gl.Frustum(
		float64(r.Left),
		float64(r.Right),
		float64(r.Bottom),
		float64(r.Top), 1.0, 10.0)
}

func DrawScene() {
	//rotationX += 0.5
	//rotationY += 0.5
	gl.Viewport(0, 0, gfx.CurrAppWidth, gfx.CurrAppHeight) // OPTIMIZEME?  could set flag upon frame buffer size change event
	if *gfx.PrevFrustum != *gfx.CurrFrustum {
		*gfx.PrevFrustum = *gfx.CurrFrustum
		gl.MatrixMode(gl.PROJECTION)
		gl.LoadIdentity()
		setFrustum(gfx.CurrFrustum)
		fmt.Println("CHANGE OF FRUSTUM")
		fmt.Printf(".Panels[0].Rect.Right: %.2f\n", gfx.Rend.Panels[0].Rect.Right)
		fmt.Printf(".Panels[0].Rect.Top: %.2f\n", gfx.Rend.Panels[0].Rect.Top)
	}
	gl.MatrixMode(gl.MODELVIEW) //.PROJECTION)                   //.MODELVIEW)
	gl.LoadIdentity()
	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
	gl.Translatef(0, 0, -gfx.Rend.DistanceFromOrigin)
	//gl.Rotatef(rotationX, 1, 0, 0)
	//gl.Rotatef(rotationY, 0, 1, 0)

	gl.BindTexture(gl.TEXTURE_2D, Texture)

	gl.Begin(gl.QUADS)
	gfx.Rend.DrawAll()
	gl.End()
}

func NewTexture(file string) uint32 {
	imgFile, err := os.Open(file)
	if err != nil {
		log.Fatalf("texture %q not found on disk: %v\n", file, err)
	}

	img, _, err := image.Decode(imgFile)
	if err != nil {
		panic(err)
	}

	rgba := image.NewRGBA(img.Bounds())
	if rgba.Stride != rgba.Rect.Size().X*4 {
		panic("unsupported stride")
	}

	draw.Draw(rgba, rgba.Bounds(), img, image.Point{0, 0}, draw.Src)

	var texture uint32
	gl.Enable(gl.TEXTURE_2D)
	gl.GenTextures(1, &texture)
	gl.BindTexture(gl.TEXTURE_2D, texture)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.LINEAR)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.LINEAR)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_S, gl.CLAMP_TO_EDGE)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_T, gl.CLAMP_TO_EDGE)
	gl.TexImage2D(
		gl.TEXTURE_2D,
		0,
		gl.RGBA,
		int32(rgba.Rect.Size().X),
		int32(rgba.Rect.Size().Y),
		0,
		gl.RGBA,
		gl.UNSIGNED_BYTE,
		gl.Ptr(rgba.Pix))

	return texture
}

func destroyScene() {
}

func makeCube() {
	gl.Normal3f(0, 0, 1)
	gl.TexCoord2f(0, 0)
	gl.Vertex3f(-1, -1, 1)
	gl.TexCoord2f(1, 0)
	gl.Vertex3f(1, -1, 1)
	gl.TexCoord2f(1, 1)
	gl.Vertex3f(1, 1, 1)
	gl.TexCoord2f(0, 1)
	gl.Vertex3f(-1, 1, 1)

	gl.Normal3f(0, 0, -1)
	gl.TexCoord2f(1, 0)
	gl.Vertex3f(-1, -1, -1)
	gl.TexCoord2f(1, 1)
	gl.Vertex3f(-1, 1, -1)
	gl.TexCoord2f(0, 1)
	gl.Vertex3f(1, 1, -1)
	gl.TexCoord2f(0, 0)
	gl.Vertex3f(1, -1, -1)

	gl.Normal3f(0, 1, 0)
	gl.TexCoord2f(0, 1)
	gl.Vertex3f(-1, 1, -1)
	gl.TexCoord2f(0, 0)
	gl.Vertex3f(-1, 1, 1)
	gl.TexCoord2f(1, 0)
	gl.Vertex3f(1, 1, 1)
	gl.TexCoord2f(1, 1)
	gl.Vertex3f(1, 1, -1)

	gl.Normal3f(0, -1, 0)
	gl.TexCoord2f(1, 1)
	gl.Vertex3f(-1, -1, -1)
	gl.TexCoord2f(0, 1)
	gl.Vertex3f(1, -1, -1)
	gl.TexCoord2f(0, 0)
	gl.Vertex3f(1, -1, 1)
	gl.TexCoord2f(1, 0)
	gl.Vertex3f(-1, -1, 1)

	gl.Normal3f(1, 0, 0)
	gl.TexCoord2f(1, 0)
	gl.Vertex3f(1, -1, -1)
	gl.TexCoord2f(1, 1)
	gl.Vertex3f(1, 1, -1)
	gl.TexCoord2f(0, 1)
	gl.Vertex3f(1, 1, 1)
	gl.TexCoord2f(0, 0)
	gl.Vertex3f(1, -1, 1)

	gl.Normal3f(-1, 0, 0)
	gl.TexCoord2f(0, 0)
	gl.Vertex3f(-1, -1, -1)
	gl.TexCoord2f(1, 0)
	gl.Vertex3f(-1, -1, 1)
	gl.TexCoord2f(1, 1)
	gl.Vertex3f(-1, 1, 1)
	gl.TexCoord2f(0, 1)
	gl.Vertex3f(-1, 1, -1)
}

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
	dispatchWithPrefix(content, msg.TypeMousePos)
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
	dispatchWithPrefix(content, msg.TypeMouseScroll)
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
	content := append(getByteOfUInt8(uint8(b)), getByteOfUInt8(uint8(action))...)
	content = append(content, getByteOfUInt8(uint8(mod))...)
	dispatchWithPrefix(content, msg.TypeMouseButton)
}

func onChar(w *glfw.Window, char rune) {
	dispatchWithPrefix(getBytesOfRune(char), msg.TypeCharacter)
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
		//fmt.Println("Control RELEASED")

		case glfw.KeyLeftAlt:
			fallthrough
		case glfw.KeyRightAlt:
		//fmt.Println("Alt RELEASED")

		case glfw.KeyLeftSuper:
			fallthrough
		case glfw.KeyRightSuper:
			//fmt.Println("'Super' modifier key RELEASED")
		}
	} else { // glfw.Repeat   or   glfw.Press
		b := foc.TextBodies[0]

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

			if foc.CursY >= len(b) {
				foc.CursY = len(b) - 1
			}

		case glfw.KeyHome:
			commonMovementKeyHandling()
			foc.CursX = 0
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
					foc.CursX--
				}
			}
		case glfw.KeyRight:
			commonMovementKeyHandling()

			if foc.CursX < len(b[foc.CursY]) {
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

		script.Process(false)
	}

	// build message
	content := getByteOfUInt8(uint8(key))
	content = append(content, getBytesOfSInt32(int32(scancode))...)
	content = append(content, getByteOfUInt8(uint8(action))...)
	content = append(content, getByteOfUInt8(uint8(mod))...)
	dispatchWithPrefix(content, msg.TypeKey)
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

	Events <- append(prefix, content...)
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



var curRecByte = 0 // current receive message index

func monitorEvents(ch chan []byte) {
	select {
	case v := <-ch:
		processMessage(v)
	default:
	//fmt.Println("monitorEvents() default")
	}
}

/* message processing example */

/*
func ProcessIncomingMessages() {
	//have a channel for incoming messages
	for msg := range self.IncomingChannel {
		switch msg.GetMessageType(msg) {
		//InRouteMessage is the only message coming in to node from transports
		case msg.MsgMessageMousePos:
			var m1 msg.MessageMousePos
			msg.Deserialize(msg, m1)
			//self.HandleInRouteMessage(m1)
			fmt.Printf("MessageMousePos: X= %f, Y= %f \n", m1.X, m2.X)
		case msg.MsgMessageMouseScroll:
			var m2 msg.MessageMouseScroll
			msg.Deserialize(msg, m1)

		case msg.MsgMessageMouseButton:
			var m3 msg.MessageMouseButton
			mesg.Deserialize(msg, m1)

		case msg.MsgMessageCharacter:
			var m4 msg.MessageCharacter
			msg.Deserialize(msg, m1)

		case msg.MsgMessageKey:
			var m5 msg.MessageKey
			msg.Deserialize(msg, m1)

		default:
			fmt.Println("UNKNOWN MESSAGE TYPE!")

		}
	}
}
*/

func processMessage(message []byte) {
	switch getMessageType(".", message) {

	case MessageMousePos:
		s("MessageMousePos", message)
		showFloat64("X", message)
		showFloat64("Y", message)

	case MessageMouseScroll:
		s("MessageMouseScroll", message)
		showFloat64("X Offset", message)
		showFloat64("Y Offset", message)

	case MessageMouseButton:
		s("MessageMouseButton", message)
		gfx.Curs.ConvertMouseClickToTextCursorPosition(
			getAndShowUInt8("Button", message),
			getAndShowUInt8("Action", message))
		getAndShowUInt8("Mod", message)

	case MessageCharacter:
		s("MessageCharacter", message)
		insertRuneIntoDocument("Rune", message)
		script.Process(false)

	case MessageKey:
		s("MessageKey", message)
		getAndShowUInt8("Key", message)
		showSInt32("Scan", message)
		getAndShowUInt8("Action", message)
		getAndShowUInt8("Mod", message)

	default:
		fmt.Println("UNKNOWN MESSAGE TYPE!")
	}

	fmt.Println()
	curRecByte = 0
}

func s(s string, message []byte) {
	fmt.Print(s)
	showUInt32("Len", message)
	curRecByte++ // skipping message type's space
}

func getMessageType(s string, message []byte) (value uint8) {
	rBuf := bytes.NewReader(message[4:5])
	err := binary.Read(rBuf, binary.LittleEndian, &value)

	if err != nil {
		fmt.Println("binary.Read failed: ", err)
	} else {
		//fmt.Printf("from byte buffer, %s: %d\n", s, value)
	}

	return
}

func insertRuneIntoDocument(s string, message []byte) {
	var value rune
	var size = 4

	rBuf := bytes.NewReader(message[curRecByte : curRecByte+size])
	err := binary.Read(rBuf, binary.LittleEndian, &value)
	curRecByte += size

	if err != nil {
		fmt.Println("binary.Read failed: ", err)
	} else {
		fmt.Printf("   [%s: %s]", s, string(value))

		f := gfx.Rend.Focused
		b := f.TextBodies[0]
		b[f.CursY] = b[f.CursY][:f.CursX] + string(value) + b[f.CursY][f.CursX:len(b[f.CursY])]
		f.CursX++
	}
}

func getAndShowUInt8(s string, message []byte) (value uint8) {
	var size = 1

	rBuf := bytes.NewReader(message[curRecByte : curRecByte+size])
	err := binary.Read(rBuf, binary.LittleEndian, &value)
	curRecByte += size

	if err != nil {
		fmt.Println("binary.Read failed: ", err)
	} else {
		fmt.Printf("   [%s: %d]", s, value)
	}

	return
}

// the rest of these funcs are almost identical, just top 2 vars customized (and string format)
func showSInt32(s string, message []byte) {
	var value int32
	var size = 4

	rBuf := bytes.NewReader(message[curRecByte : curRecByte+size])
	err := binary.Read(rBuf, binary.LittleEndian, &value)
	curRecByte += size

	if err != nil {
		fmt.Println("binary.Read failed: ", err)
	} else {
		fmt.Printf("   [%s: %d]", s, value)
	}
}

func showUInt32(s string, message []byte) {
	var value uint32
	var size = 4

	rBuf := bytes.NewReader(message[curRecByte : curRecByte+size])
	err := binary.Read(rBuf, binary.LittleEndian, &value)
	curRecByte += size

	if err != nil {
		fmt.Println("binary.Read failed: ", err)
	} else {
		fmt.Printf("   [%s: %d]", s, value)
	}
}

func showFloat64(s string, message []byte) {
	var value float64
	var size = 8

	rBuf := bytes.NewReader(message[curRecByte : curRecByte+size])
	err := binary.Read(rBuf, binary.LittleEndian, &value)
	curRecByte += size

	if err != nil {
		fmt.Println("binary.Read failed: ", err)
	} else {
		fmt.Printf("   [%s: %.1f]", s, value)
	}
}

// Copyright 2014 The go-gl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Renders a textured spinning cube using GLFW 3.1 and OpenGL 2.1.
package main

import (
	"fmt"
	"go/build"
	"image"
	"image/draw"
	_ "image/png"
	"log"
	"math"
	"os"
	"runtime"
	"time"

	"github.com/go-gl/gl/v2.1/gl"
	"github.com/go-gl/glfw/v3.1/glfw"
)

var (
	// gl graphics
	texture   uint32
	rotationX float32
	rotationY float32
)

func init() {
	// GLFW event handling must run on the main OS thread
	// See documentation for functions that are only allowed to be called from the main thread.
	runtime.LockOSThread()
}

// character rendering
var uvSpan = float32(1.0) / 16
var rectRad = float32(3) // rectangular radius (distance to edge in the cardinal directions from the center, corners would be farther away)
var curX = -rectRad
var curY = rectRad
var lnHei = float32(0.25)
var chWid = lnHei / 2
var runes = []rune("aaa AAA Blah Blah fluffernutter sandwiches Blah Blah fluffernutter sandwiches")
var document = make([]string, 0)

// cursor
var nextBlinkChange = time.Now()
var cursVisible = true
var cursX = 0
var cursY = 0

// selection
// future consideration/fixme:
// need to sanitize start/end positions.
// since they may be beyond the last line character of the line.
// also, in addition to backspace/delete, typing any visible character should delete marked text.
// complication:   if they start or end on invalid characters (of the line string),
// the forward or backwards direction from there determines where the visible selection
// range starts/ends....
// will an end position always be defined (when value is NOT math.MaxUint32),
// when a START is?  because that determines where the first VISIBLY marked
// character starts
var selectionStartX = math.MaxUint32
var selectionStartY = math.MaxUint32
var selectionEndX = math.MaxUint32
var selectionEndY = math.MaxUint32
var selectingRangeOfText = false

func main() {
	if err := glfw.Init(); err != nil {
		log.Fatalln("failed to initialize glfw:", err)
	}
	defer glfw.Terminate()

	glfw.WindowHint(glfw.Resizable, glfw.False)
	glfw.WindowHint(glfw.ContextVersionMajor, 2)
	glfw.WindowHint(glfw.ContextVersionMinor, 1)
	window, err := glfw.CreateWindow(800, 600, "V I S", nil, nil)
	if err != nil {
		panic(err)
	}
	window.MakeContextCurrent()

	if err := gl.Init(); err != nil {
		panic(err)
	}

	texture = newTexture("Bisasam_24x24_Shadowed.png")
	defer gl.DeleteTextures(1, &texture)

	// init
	setupScene()
	initDoc()
	initInputEvents(window)

	for !window.ShouldClose() {
		pollEventsAndHandleAnInput(window)
		drawScene()
		window.SwapBuffers()
	}
}

func initDoc() {
	var line = make([]rune, 0)

	for _, c := range runes {
		line = append(line, c)
	}

	document = append(document, "testing init line")
}

func initInputEvents(w *glfw.Window) {
	w.SetCharCallback(onChar)
	w.SetKeyCallback(onKey)
	w.SetMouseButtonCallback(onMouseBtn)
	w.SetScrollCallback(onMouseScroll)
	w.SetCursorPosCallback(onMouseCursorPos)
}

func pollEventsAndHandleAnInput(window *glfw.Window) {
	glfw.PollEvents()

	// poll a particular key state
	//if window.GetKey(glfw.KeyEscape) /* Action */ == glfw.Press {
	//	fmt.Println("PRESSED ESCape")
	//	window.SetShouldClose(true)
	//}
}

func onMouseCursorPos(w *glfw.Window, x float64, y float64) {
	fmt.Printf("onMouseCursorMove() x: %.1f   y: %.1f\n", x, y)
}

func onMouseScroll(window *glfw.Window, xOff float64, yOff float64) {
	fmt.Println("onScroll()")
}

func onMouseBtn(
	window *glfw.Window,
	b glfw.MouseButton,
	action glfw.Action,
	mods glfw.ModifierKey) {

	fmt.Println("onMouseBtn()")

	if action != glfw.Press {
		switch glfw.MouseButton(b) {
		case glfw.MouseButtonLeft:
		default:
		}
	}
}

// WEIRD BEHAVIOUR OF KEY EVENTS.... for a PRESS, you can get a
// shift/alt/ctrl/super event through the "mods" variable,
// (see the top of "action == glfw.Press" section for an example)
// regardless of left/right key used.
// BUT for RELEASE, the "mods" variable will NOT tell you what key it is!
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
	}

	if action == glfw.Press {
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

	if key == glfw.KeyEscape && action == glfw.Release {
		w.SetShouldClose(true)
	}
}

func commonMovementKeyHandling() {
	if selectingRangeOfText {
		selectionEndX = cursX
		selectionEndY = cursY
	} else { // arrow keys without shift gets rid selection
		selectionStartX = math.MaxUint32
		selectionStartY = math.MaxUint32
		selectionEndX = math.MaxUint32
		selectionEndY = math.MaxUint32
	}
}

func onChar(w *glfw.Window, char rune) {
	// when a Unicode character is input.
	//func (w *Window) SetCharacterCallback(cbfun func(w *Window, char uint)) {}

	fmt.Printf("onChar(): %c\n", char)
	//document[len(document)-1] += string(char)
	document[cursY] = document[cursY][:cursX] + string(char) + document[cursY][cursX:len(document[cursY])]
	cursX++
}

func removeCharacter(fromUnderCursor bool) {
	if fromUnderCursor {
		if len(document[cursY]) > cursX {
			document[cursY] = document[cursY][:cursX] + document[cursY][cursX+1:len(document[cursY])]
		}
	} else {
		if cursX > 0 {
			document[cursY] = document[cursY][:cursX-1] + document[cursY][cursX:len(document[cursY])]
			cursX--
		}
	}
}

func newTexture(file string) uint32 {
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

func setupScene() {
	fmt.Println("setupScene()")

	gl.Enable(gl.DEPTH_TEST)
	gl.Enable(gl.LIGHTING)

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
	gl.Frustum(-1, 1, -1, 1, 1.0, 10.0)
	gl.MatrixMode(gl.MODELVIEW)
	gl.LoadIdentity()
}

func destroyScene() {

}

func drawScene() {
	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

	gl.MatrixMode(gl.MODELVIEW)
	gl.LoadIdentity()
	gl.Translatef(0, 0, -3.0)
	//gl.Rotatef(rotationX, 1, 0, 0)
	//gl.Rotatef(rotationY, 0, 1, 0)

	//rotationX += 0.5
	//rotationY += 0.5

	gl.BindTexture(gl.TEXTURE_2D, texture)

	gl.Color4f(1, 1, 1, 1)

	gl.Begin(gl.QUADS)

	//makeCube()
	makeChars()

	gl.End()
}

func makeChars() {
	for _, line := range document {
		for _, c := range line {
			drawCharAtCurrentPos(c)
		}

		curX = -rectRad
		curY -= lnHei
	}

	drawCursorMaybe()

	curX = -rectRad
	curY = rectRad
}

func drawCursorMaybe() {
	if nextBlinkChange.Before(time.Now()) {
		nextBlinkChange = time.Now().Add(time.Millisecond * 170)
		cursVisible = !cursVisible
	}

	if cursVisible == true {
		drawCharAt('_', cursX, cursY)
	}
}

func drawCharAt(letter rune, posX int, posY int) {
	x := int(letter) % 16
	y := int(letter) / 16

	gl.Normal3f(0, 0, 1)

	gl.TexCoord2f(float32(x)*uvSpan, float32(y)*uvSpan+uvSpan) // bl  0, 1
	gl.Vertex3f(-rectRad+float32(posX)*chWid, rectRad-float32(posY)*lnHei-lnHei, 0)

	gl.TexCoord2f(float32(x)*uvSpan+uvSpan, float32(y)*uvSpan+uvSpan) // br  1, 1
	gl.Vertex3f(-rectRad+float32(posX)*chWid+chWid, rectRad-float32(posY)*lnHei-lnHei, 0)

	gl.TexCoord2f(float32(x)*uvSpan+uvSpan, float32(y)*uvSpan) // tr  1, 0
	gl.Vertex3f(-rectRad+float32(posX)*chWid+chWid, rectRad-float32(posY)*lnHei, 0)

	gl.TexCoord2f(float32(x)*uvSpan, float32(y)*uvSpan) // tl  0, 0
	gl.Vertex3f(-rectRad+float32(posX)*chWid, rectRad-float32(posY)*lnHei, 0)

	curX += lnHei / 2

	if curX >= rectRad {
		curX = -rectRad
		curY -= lnHei
	}
}

func drawCharAtCurrentPos(letter rune) {
	x := int(letter) % 16
	y := int(letter) / 16

	gl.Normal3f(0, 0, 1)

	gl.TexCoord2f(float32(x)*uvSpan, float32(y)*uvSpan+uvSpan) // bl  0, 1
	gl.Vertex3f(curX, curY-lnHei, 0)

	gl.TexCoord2f(float32(x)*uvSpan+uvSpan, float32(y)*uvSpan+uvSpan) // br  1, 1
	gl.Vertex3f(curX+lnHei/2, curY-lnHei, 0)

	gl.TexCoord2f(float32(x)*uvSpan+uvSpan, float32(y)*uvSpan) // tr  1, 0
	gl.Vertex3f(curX+lnHei/2, curY, 0)

	gl.TexCoord2f(float32(x)*uvSpan, float32(y)*uvSpan) // tl  0, 0
	gl.Vertex3f(curX, curY, 0)

	curX += lnHei / 2

	if curX >= rectRad {
		curX = -rectRad
		curY -= lnHei
	}
}

// Set the working directory to the root of Go package, so that its assets can be accessed.
func init() {
	//dir, err := importPathToDir("github.com/go-gl/examples/glfw31-gl21-cube")
	dir, err := importPathToDir("c21/")
	if err != nil {
		log.Fatalln("Unable to find Go package in your GOPATH, it's needed to load assets:", err)
	}
	err = os.Chdir(dir)
	if err != nil {
		log.Panicln("os.Chdir:", err)
	}
}

// importPathToDir resolves the absolute path from importPath.
// There doesn't need to be a valid Go package inside that import path,
// but the directory must exist.
func importPathToDir(importPath string) (string, error) {
	p, err := build.Import(importPath, "", build.FindOnly)
	if err != nil {
		return "", err
	}
	return p.Dir, nil
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

/*


accepted
According to the Go specification:

For an expression x of interface type and a type T, the primary expression x.(T) asserts that x is not nil and that the value stored in x is of type T.
A "type assertion" allows you to declare an interface value contains a certain concrete type or that its concrete type satisfies another interface.

In your example, you were asserting data (type interface{}) has the concrete type string. If you are wrong, the program will panic at runtime. You do not need to worry about efficiency, checking just requires comparing two pointer values.

If you were unsure if it was a string or not, you could test using the two return syntax.

str, ok := data.(string)
If data is not a string, ok will be false. It is then common to wrap such a statement into an if statement like so:

if str, ok := data.(string); ok {
    // act on str
} else {
    // not string
}






    s := make([]string, 3)
    fmt.Println("emp:", s)
We can set and get just like with arrays.
    s[0] = "a"
    s[1] = "b"
    s[2] = "c"
    fmt.Println("set:", s)
    fmt.Println("get:", s[2])
len returns the length of the slice as expected.
    fmt.Println("len:", len(s))
In addition to these basic operations, slices support several more that make them richer than arrays. One is the
builtin append, which returns a slice containing one or more new values. Note that we need to accept a return value
from append as we may get a new slice value.
    s = append(s, "d")
    s = append(s, "e", "f")
    fmt.Println("apd:", s)


    c := make([]string, len(s))
    copy(c, s)
    fmt.Println("cpy:", c)
Slices support a “slice” operator with the syntax slice[low:high]. For example, this gets a slice of the elements s[2], s[3], and s[4].
    l := s[2:5]
    fmt.Println("sl1:", l)
This slices up to (but excluding) s[5].
    l = s[:5]
    fmt.Println("sl2:", l)
And this slices up from (and including) s[2].
    l = s[2:]
    fmt.Println("sl3:", l)
We can declare and initialize a variable for slice in a single line as well.
    t := []string{"g", "h", "i"}
    fmt.Println("dcl:", t)
Slices can be composed into multi-dimensional data structures. The length of the inner slices can vary, unlike with multi-dimensional arrays.
    twoD := make([][]int, 3)
    for i := 0; i < 3; i++ {
        innerLen := i + 1
        twoD[i] = make([]int, innerLen)
        for j := 0; j < innerLen; j++ {
            twoD[i][j] = i + j
        }
    }
    fmt.Println("2d: ", twoD)
}
*/

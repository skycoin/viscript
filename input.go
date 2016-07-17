package main

import (
	"fmt"
	/*
		"go/build"
		"image"
		"image/draw"
		_ "image/png"
		"log"
		"math"
		"os"
		"runtime"
		"time"
	*/

	//"github.com/go-gl/gl/v2.1/gl"
	"github.com/go-gl/glfw/v3.1/glfw"
)

func initInputEvents(w *glfw.Window) {
	w.SetCharCallback(onChar)
	w.SetKeyCallback(onKey)
	w.SetMouseButtonCallback(onMouseBtn)
	w.SetScrollCallback(onMouseScroll)
	w.SetCursorPosCallback(onMouseCursorPos)
}

func onMouseCursorPos(w *glfw.Window, x float64, y float64) {
	mouseX = int(x) / pixelWid
	mouseY = int(y) / pixelHei
	fmt.Printf("onMouseCursorPos() X: %.1f   Y: %.1f\n", x, y)
}

func pollEventsAndHandleAnInput(window *glfw.Window) {
	glfw.PollEvents()

	// poll a particular key state
	//if window.GetKey(glfw.KeyEscape) /* Action */ == glfw.Press {
	//	fmt.Println("PRESSED ESCape")
	//	window.SetShouldClose(true)
	//}
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

func onChar(w *glfw.Window, char rune) {
	// when a Unicode character is input.
	//func (w *Window) SetCharacterCallback(cbfun func(w *Window, char uint)) {}

	fmt.Printf("onChar(): %c\n", char)
	//document[len(document)-1] += string(char)
	document[cursY] = document[cursY][:cursX] + string(char) + document[cursY][cursX:len(document[cursY])]
	cursX++
}

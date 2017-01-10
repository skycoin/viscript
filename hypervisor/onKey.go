package hypervisor

import (
	"fmt"
	"github.com/corpusc/viscript/gfx"
	"github.com/corpusc/viscript/msg"
	//"github.com/corpusc/viscript/script"
	"github.com/go-gl/glfw/v3.2/glfw"
)

// WEIRD BEHAVIOUR OF KEY EVENTS.... for a PRESS, you can detect a
// shift/alt/ctrl/super key through the "mod" variable,
// (see the top of "action == glfw.Press" section for an example)
// regardless of left/right key used.
// BUT for RELEASE, the "mod" variable will NOT tell you what key it is!
// so you will have to handle both left & right mod keys via the "action" variable!

// func onKey(
// 	w *glfw.Window,
// 	key glfw.Key,
// 	scancode int,
// 	action glfw.Action,
// 	mod glfw.ModifierKey) {
func onKey(m msg.MessageKey) {
	foc := gfx.Rend.Focused

	if glfw.Action(m.Action) == glfw.Release {
		fmt.Println("release", m.Key)

		switch glfw.Key(m.Key) {

		case glfw.KeyEscape:
			//w.SetShouldClose(true)
			fmt.Println("case glfw.KeyEscape:")
			CloseWindow <- 1
			HypervisorScreenTeardown()

		case glfw.KeyLeftShift:
			fallthrough
		case glfw.KeyRightShift:
			fmt.Println("Done selecting")
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
	} else { // glfw.Press   or   glfw.Repeat
		fmt.Println("press")

		b := foc.TextBodies[0]

		switch glfw.ModifierKey(m.Mod) {
		case glfw.ModShift:
			fmt.Println("Started selecting")
			foc.Selection.CurrentlySelecting = true
			foc.Selection.StartX = foc.CursX
			foc.Selection.StartY = foc.CursY
		}

		switch glfw.Key(m.Key) {
		case glfw.KeyEnter:
			startOfLine := b[foc.CursY][:foc.CursX]
			restOfLine := b[foc.CursY][foc.CursX:len(b[foc.CursY])]
			b[foc.CursY] = startOfLine
			b = insert(b, foc.CursY+1, restOfLine)

			foc.CursX = 0
			foc.CursY++
			foc.TextBodies[0] = b

			if foc.CursY >= len(b) {
				foc.CursY = len(b) - 1
			}
		case glfw.KeyHome:
			if eitherControlKeyHeld() {
				foc.CursY = 0
			}

			foc.CursX = 0
			movedCursorSoUpdateDependents()
		case glfw.KeyEnd:
			if eitherControlKeyHeld() {
				foc.CursY = len(b) - 1
			}

			foc.CursX = len(b[foc.CursY])
			movedCursorSoUpdateDependents()
		case glfw.KeyUp:
			if foc.CursY > 0 {
				foc.CursY--

				if foc.CursX > len(b[foc.CursY]) {
					foc.CursX = len(b[foc.CursY])
				}
			}

			movedCursorSoUpdateDependents()
		case glfw.KeyDown:
			if foc.CursY < len(b)-1 {
				foc.CursY++

				if foc.CursX > len(b[foc.CursY]) {
					foc.CursX = len(b[foc.CursY])
				}
			}

			movedCursorSoUpdateDependents()
		case glfw.KeyLeft:
			if foc.CursX == 0 {
				if foc.CursY > 0 {
					foc.CursY--
					foc.CursX = len(b[foc.CursY])
				}
			} else {
				if glfw.ModifierKey(m.Mod) == glfw.ModControl {
					foc.CursX = getWordSkipPos(foc.CursX, -1)
				} else {
					foc.CursX--
				}
			}

			movedCursorSoUpdateDependents()
		case glfw.KeyRight:
			if foc.CursX < len(b[foc.CursY]) {
				if glfw.ModifierKey(m.Mod) == glfw.ModControl {
					foc.CursX = getWordSkipPos(foc.CursX, 1)
				} else {
					foc.CursX++
				}
			}

			movedCursorSoUpdateDependents()
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

		//script.Process(false)
	}
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

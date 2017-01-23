package hypervisor

import (
	"fmt"
	"github.com/corpusc/viscript/msg"
	//"github.com/corpusc/viscript/script"
)

// WEIRD BEHAVIOUR OF KEY EVENTS.... for a PRESS, you can detect a
// shift/alt/ctrl/super key through the "mod" variable,
// (see the top of "action == glfw.Press" section for an example)
// regardless of left/right key used.
// BUT for RELEASE, the "mod" variable will NOT tell you what key it is!
// so you will have to handle both left & right modifier keys via the "action" variable!

func onKey(m msg.MessageKey) {
	//foc := Focused

	if msg.Action(m.Action) == msg.Release {
		fmt.Println("\nrelease --------- ", m.Key)

		switch m.Key {

		case msg.KeyEscape:
			fmt.Println("case msg.KeyEscape:")
			fmt.Printf("CLOSE OPENGL WINDOW \n")
			CloseWindow = true

		case msg.KeyLeftShift:
			fallthrough
		case msg.KeyRightShift:
			fmt.Println("Done selecting")
			//foc.Selection.CurrentlySelecting = false // TODO?  possibly flip around if selectionStart comes after selectionEnd in the page flow?

		case msg.KeyLeftControl:
			fallthrough
		case msg.KeyRightControl:
			fmt.Println("Control RELEASED")

		case msg.KeyLeftAlt:
			fallthrough
		case msg.KeyRightAlt:
			fmt.Println("Alt RELEASED")

		case msg.KeyLeftSuper:
			fallthrough
		case msg.KeyRightSuper:
			fmt.Println("'Super' modifier key RELEASED")
		}
	} else { // glfw.Press   or   glfw.Repeat
		fmt.Println("\npress --------- ", m.Key)

		/*
			b := foc.TextBodies[0]

			switch glfw.ModifierKey(m.Mod) {
			case glfw.ModShift:
				fmt.Println("Started selecting")
				foc.Selection.CurrentlySelecting = true
				foc.Selection.StartX = foc.CursX
				foc.Selection.StartY = foc.CursY
			case glfw.ModAlt:
				fmt.Println("glfw.ModAlt")
			case glfw.ModControl:
				fmt.Println("glfw.ModControl")
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
		*/

		//script.Process(false)
	}
}

// func getWordSkipPos(xIn int, change int) int {

// 	peekPos := xIn
// 	foc := Focused
// 	b := foc.TextBodies[0]

// 	for {
// 		peekPos += change

// 		if peekPos < 0 {
// 			return 0
// 		}

// 		if peekPos >= len(b[foc.CursY]) {
// 			return len(b[foc.CursY])
// 		}

// 		if string(b[foc.CursY][peekPos]) == " " {
// 			return peekPos
// 		}
// 	}
// }

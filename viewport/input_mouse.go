package viewport

import (
	"fmt"
	"math/rand"

	"time"

	"github.com/corpusc/viscript/msg"
)

var lastMousePos msg.MessageMousePos = msg.MessageMousePos{0, 0}

// triggered both by moving **AND*** by pressing buttons
func onMouseCursorPos(m msg.MessageMousePos) {
	lastMousePos = m
	/*
		x := float32(m.X)
		y := float32(m.Y)

		mouse.UpdatePosition(
			app.Vec2F{x, y},
			gl.CanvasExtents,
			gl.PixelSize) // state update

		// rendering update
		if gl.GlfwWindow.GetMouseButton(glfw.MouseButtonLeft) == glfw.Press {
			ScrollTermThatHasMousePointer(mouse.PixelDelta.X, mouse.PixelDelta.Y)
		}
	*/
}

func onMouseScroll(m msg.MessageMouseScroll) {
	/*
		var delta float32 = 30

		if eitherControlKeyHeld() { // horizontal ability from 1D scrolling
			ScrollTermThatHasMousePointer(float32(m.Y)*-delta, 0)
		} else { // can handle both x & y for 2D scrolling
			ScrollTermThatHasMousePointer(float32(m.X)*delta, float32(m.Y)*-delta)
		}
	*/
}

// apparently every time this is fired, a mouse position event is ALSO fired
func onMouseButton(m msg.MessageMouseButton) {
	// convertClickToTextCursorPosition(m.Button, m.Action)
	if msg.Action(m.Action) == msg.Press {

		keys := make([]msg.TerminalId, 0, len(Terms.Terms))

		for k, _ := range Terms.Terms {
			keys = append(keys, k)
		}

		rand.Seed(time.Now().Unix())
		randKey := keys[rand.Intn(len(Terms.Terms))]

		println("\nRandom Term Key from terminal stack:", randKey)

		for terminalKeyId, eachTerminal := range Terms.Terms {
			fmt.Println("\n")
			fmt.Println("TerminalId In Stack:", terminalKeyId, "Terminal ID:", eachTerminal.TerminalId)
			fmt.Println("Mouse X:", lastMousePos.X, "Mouse Y:", lastMousePos.Y)
			eachTerminal.Bounds.Print()
			// if eachTerminal.Bounds.Contains(float32(lastMousePos.X), float32(lastMousePos.Y)) {
			// 	fmt.Println("Cursor is indside", terminalKeyId)
			// 	// Terms.FocusedId = terminalKeyId
			// 	// Terms.Focused = eachTerminal

			// }
		}

		keys = nil
		keys = make([]msg.TerminalId, 0, len(Terms.Terms))

		for k, _ := range Terms.Terms {
			if k != randKey {
				keys = append(keys, k)
			}
		}

		keys = append(keys, randKey)

		// TODO
		Terms.DrawOrder = keys
		println("\nNew Keys")
		for k := range keys {
			print(keys[k])
			print(" ")
		}

	}

	// TODO contains doesn't work because of normalized Top, Left, Bottom, Right and pixel coodinate comparison

	// switch glfw.MouseButton(m.Button) {
	// case glfw.MouseButtonLeft:
	// 	// respond to clicks in ui rectangles
	// 	if mouse.CursorIsInside(ui.MainMenu.Rect) {
	// 		respondToAnyMenuButtonClicks()
	// 	} else { // respond to any panel clicks outside of menu
	// 		for _, t := range Terms {
	// 			if t.ContainsMouseCursor() {
	// 				t.RespondToMouseClick()
	// 			}
	// 		}
	// 	}
	// }
}

func convertClickToTextCursorPosition(button, action uint8) {
	// if glfw.MouseButton(button) == glfw.MouseButtonLeft &&
	// 	glfw.Action(action) == glfw.Press {

	// 	foc := Focused

	// 	if foc.IsEditable && foc.Content.Contains(mouse.GlX, mouse.GlY) {
	// 		if foc.MouseY < len(foc.TextBodies[0]) {
	// 			foc.CursY = foc.MouseY

	// 			if foc.MouseX <= len(foc.TextBodies[0][foc.CursY]) {
	// 				foc.CursX = foc.MouseX
	// 			} else {
	// 				foc.CursX = len(foc.TextBodies[0][foc.CursY])
	// 			}
	// 		} else {
	// 			foc.CursY = len(foc.TextBodies[0]) - 1
	// 		}
	// 	}
	// }
}

func respondToAnyMenuButtonClicks() {
	// for _, bu := range ui.MainMenu.Buttons {
	// 	if mouse.CursorIsInside(bu.Rect.Rectangle) {
	// 		bu.Activated = !bu.Activated

	// 		switch bu.Name {
	// 		case "Run":
	// 			if bu.Activated {
	// 				//script.Process(true)
	// 			}
	// 			break
	// 		case "Testing Tree":
	// 			if bu.Activated {
	// 				//script.Process(true)
	// 				//script.MakeTree()
	// 			} else { // deactivated
	// 				// remove all terminals with trees
	// 				b := Terms[:0]
	// 				for _, t := range Terms {
	// 					if len(t.Trees) < 1 {
	// 						b = append(b, t)
	// 					}
	// 				}
	// 				Terms = b
	// 				//fmt.Printf("len of b (from Terms) after removing ones with trees: %d\n", len(b))
	// 				//fmt.Printf("len of Terms: %d\n", len(Terms))
	// 			}
	// 			break
	// 		}

	// 		app.Con.Add(fmt.Sprintf("%s toggled\n", bu.Name))
	// 	}
	// }
}

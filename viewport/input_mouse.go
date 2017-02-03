package viewport

import (
	"github.com/corpusc/viscript/app"
	"github.com/corpusc/viscript/hypervisor/input/mouse"
	"github.com/corpusc/viscript/msg"
)

// triggered both by moving **AND*** by pressing buttons
func onMouseCursorPos(m msg.MessageMousePos) {
	x := float32(m.X)
	y := float32(m.Y)

	mouse.UpdatePosition(
		app.Vec2F{x, y},
		gl.CanvasExtents,
		gl.PixelSize) // state update

	if mouse.HoldingLeftButton {
		//ScrollTermThatHasMousePointer(mouse.PixelDelta.X, mouse.PixelDelta.Y)
	}
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
	convertClickToTextCursorPosition(m.Button, m.Action)

	if msg.Action(m.Action) == msg.Press {
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

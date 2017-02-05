package viewport

import (
	"github.com/corpusc/viscript/app"
	"github.com/corpusc/viscript/hypervisor/input/mouse"
	"github.com/corpusc/viscript/msg"
)

// triggered both by moving **AND*** by pressing buttons
func onMouseCursorPos(m msg.MessageMousePos) {
	mouse.UpdatePosition(app.Vec2F{float32(m.X), float32(m.Y)}) // state update

	if mouse.HoldingLeftButton {
		println("TODO: implement 'ScrollFocusedTerm()'")
		//old logic is in ScrollTermThatHasMousePointer(mouse.PixelDelta.X, mouse.PixelDelta.Y),
		//which was a janky way to do it

		// taken from ScrollIfMouseOver looks like an old implementation
		// of terminal had BarHori and BarVert it doesn't appear like that
		// now. Should we add bars or no? also here's kind of an abstract logic
		// logic taken from old scrollablepanel and divided in two if user can
		// scroll in any terminal under the cursor and if user can scroll only
		// the one that is currently focused.
		//
		// If scroll only focused parts of the
		// for _, t := range Terms.Terms {
		// 	if mouse.CursorIsInside(t.Bounds) {
		// 		xInc := mouse.GetScrollDeltaX()
		// 		yInc := mouse.GetScrollDeltaY()
		// 		t.BarHori.Scroll(xInc)
		// 		t.BarVert.Scroll(yInc)
		// 	}
		// }
		//
		// If scroll only focused
		// if mouse.CursorIsInside(Terms.Focused.Bounds){
		// 	xInc := mouse.GetScrollDeltaX()
		// 	yInc := mouse.GetScrollDeltaY()
		// 	Terms.Focused.BarHori.Scroll(xInc)
		// 	Terms.Focused.BarVert.Scroll(yInc)
		// }
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
		switch msg.MouseButton(m.Button) {
		case msg.MouseButtonLeft:
			mouse.HoldingLeftButton = true

			// // detect clicks in rects
			// if mouse.CursorIsInside(ui.MainMenu.Rect) {
			// 	respondToAnyMenuButtonClicks()
			// } else { // respond to any panel clicks outside of menu
			FocusOnTopmostRectThatContainsPointer()
			// }
		}
	} else if msg.Action(m.Action) == msg.Release {
		switch msg.MouseButton(m.Button) {
		case msg.MouseButtonLeft:
			mouse.HoldingLeftButton = false
		}
	}
}

func FocusOnTopmostRectThatContainsPointer() {
	var topmostZ float32
	var topmostId msg.TerminalId

	for id, t := range Terms.Terms {
		if mouse.CursorIsInside(t.Bounds) {
			if topmostZ < t.Depth {
				topmostZ = t.Depth
				topmostId = id
			}
		}
	}

	if topmostZ > 0 {
		Terms.FocusedId = topmostId
		Terms.Focused = Terms.Terms[topmostId]
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

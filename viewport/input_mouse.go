package viewport

import (
	"fmt"
	"github.com/corpusc/viscript/app"
	"github.com/corpusc/viscript/hypervisor/input/mouse"
	"github.com/corpusc/viscript/msg"
	"github.com/corpusc/viscript/viewport/gl"
)

// triggered both by moving **AND*** by pressing buttons
func onMouseCursorPos(m msg.MessageMousePos) {
	// if DebugPrintInputEvents {
	// 	fmt.Print("TypeMousePos")
	// 	showFloat64("X", m.X)
	// 	showFloat64("Y", m.Y)
	// 	println()
	// }
	focused := Terms.Focused
	mouse.Bounds = focused.Bounds
	mouse.Update(app.Vec2F{float32(m.X), float32(m.Y)})

	// set cursor appropriately
	if mouse.IsNearRight && !focused.ResizingBottom && !mouse.HoldingLeft {
		gl.SetHResizeCursor()
	} else if mouse.IsNearBottom && !focused.ResizingRight && !mouse.HoldingLeft {
		gl.SetVResizeCursor()
	} else if mouse.IsInsideTerminal {
		gl.SetIBeamCursor()
	} else {
		gl.SetArrowCursor()
	}

	if mouse.HoldingLeft {
		println("TODO EVENTUALLY: implement something like 'ScrollFocusedTerm()'")
		//old logic is in ScrollTermThatHasMousePointer(mouse.PixelDelta.X, mouse.PixelDelta.Y),
		//which was a janky way to do it

		// Determination should be here if the mouse is over scrollbar or over the
		// area where terminal can be moved. Moving windows happens in GL space
		// coordinates because I thought pixel delta was used for scrollbar scrolling

		// REFACTORME: cause I made it messy i guess
		// FIXME: Also the context in this case text is left there and
		// allowed to right not only the bounds
		// should resize or it should be using characters as kind of measures

		if mouse.IsNearRight && !focused.ResizingBottom {
			gl.SetHResizeCursor()
			mouse.IncreaseEdgeGlMaxAbs()
			Terms.ResizeFocusedTerminalRight(mouse.GlX)
		} else if mouse.IsNearBottom && !focused.ResizingRight {
			gl.SetVResizeCursor()
			mouse.IncreaseEdgeGlMaxAbs()
			Terms.ResizeFocusedTerminalBottom(mouse.GlY)
		}

		if mouse.CursorIsInside(focused.Bounds) && !focused.IsResizing() {

			deltaVec := app.Vec2F{mouse.GlX - mouse.PrevGlX,
				mouse.GlY - mouse.PrevGlY}
			Terms.MoveFocusedTerminal(deltaVec)
			gl.SetHandCursor()

			if DebugPrintInputEvents {
				println("\nTerminal Id:", focused.TerminalId,
					"\nTop", focused.Bounds.Top,
					"\nLeft", focused.Bounds.Left,
					"\nRight", focused.Bounds.Right,
					"\nBottom", focused.Bounds.Bottom,
					"\n\n GL MouseX:", mouse.GlX,
					"\n GL MouseY:", mouse.GlY,
					"\n\n Previous GL MouseX:", mouse.PrevGlX,
					"\n Previous GL MouseY:", mouse.PrevGlY,
					"\n\n DeltaVecX:", deltaVec.X,
					"\n DeltaVecY:", deltaVec.Y,
					"\n\n Rect Center X:", focused.Bounds.CenterX(),
					"\n Rect Center Y:", focused.Bounds.CenterY())
			}
		}
	} else {
		focused.SetResizingOff()
		mouse.DecreaseEdgeGlMaxAbs()
	}
}

func onMouseScroll(m msg.MessageMouseScroll) {
	if DebugPrintInputEvents {
		print("TypeMouseScroll")
		showFloat64("X Offset", m.X)
		showFloat64("Y Offset", m.Y)
		showBool("HoldingControl", m.HoldingControl)
		println()
	}
}

// apparently every time this is fired, a mouse position event is ALSO fired
func onMouseButton(m msg.MessageMouseButton) {
	if DebugPrintInputEvents {
		fmt.Print("TypeMouseButton")
		showUInt8("Button", m.Button)
		showUInt8("Action", m.Action)
		showUInt8("Mod", m.Mod)
		println()
	}

	convertClickToTextCursorPosition(m.Button, m.Action)

	if msg.Action(m.Action) == msg.Press {
		switch msg.MouseButton(m.Button) {
		case msg.MouseButtonLeft:
			mouse.HoldingLeft = true

			// // detect clicks in rects
			// if mouse.CursorIsInside(ui.MainMenu.Rect) {
			// 	respondToAnyMenuButtonClicks()
			// } else { // respond to any panel clicks outside of menu
			focusOnTopmostRectThatContainsPointer()
			// }
		}
	} else if msg.Action(m.Action) == msg.Release {
		switch msg.MouseButton(m.Button) {
		case msg.MouseButtonLeft:
			mouse.HoldingLeft = false
		}
	}
}

func focusOnTopmostRectThatContainsPointer() {
	if mouse.CursorIsInside(Terms.Focused.Bounds) {
		return // because it's focused. can't just click through like that
	}
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

// the rest of these funcs are almost identical, just top 2 vars customized (and string format)
func showBool(s string, x bool) {
	fmt.Printf("   [%s: %t]", s, x)
}

func showUInt8(s string, x uint8) {
	fmt.Printf("   [%s: %d]", s, x)
}

func showSInt32(s string, x int32) {
	fmt.Printf("   [%s: %d]", s, x)
}

func showUInt32(s string, x uint32) {
	fmt.Printf("   [%s: %d]", s, x)
}

func showFloat64(s string, f float64) {
	fmt.Printf("   [%s: %.1f]", s, f)
}

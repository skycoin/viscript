package process

import (
	"fmt"
	"github.com/corpusc/viscript/msg"
)

func (self *State) ProcessInputEvents(msgType uint16, message []byte) []byte {

	switch msgType {

	case msg.TypeMousePos:
		var msgMousePos msg.MessageMousePos
		msg.MustDeserialize(message, &msgMousePos)

		if self.DebugPrintInputEvents {
			fmt.Print("TypeMousePos")
			showFloat64("X", msgMousePos.X)
			showFloat64("Y", msgMousePos.Y)
		}

		onMouseCursorPos(msgMousePos)

	case msg.TypeMouseScroll:
		var msgScroll msg.MessageMouseScroll
		msg.MustDeserialize(message, &msgScroll)

		if self.DebugPrintInputEvents {
			fmt.Print("TypeMouseScroll")
			showFloat64("X Offset", msgScroll.X)
			showFloat64("Y Offset", msgScroll.Y)
		}

		onMouseScroll(msgScroll)

	case msg.TypeMouseButton:
		var msgBtn msg.MessageMouseButton
		msg.MustDeserialize(message, &msgBtn)

		if self.DebugPrintInputEvents {
			fmt.Print("TypeMouseButton")
			showUInt8("Button", msgBtn.Button)
			showUInt8("Action", msgBtn.Action)
			showUInt8("Mod", msgBtn.Mod)
		}

		onMouseButton(msgBtn)

	case msg.TypeChar:
		var msgChar msg.MessageOnCharacter
		msg.MustDeserialize(message, &msgChar)

		if self.DebugPrintInputEvents {
			fmt.Print("TypeChar")
		}

		onChar(msgChar)

	case msg.TypeKey:
		var keyMsg msg.MessageKey
		msg.MustDeserialize(message, &keyMsg)

		if self.DebugPrintInputEvents {
			fmt.Print("TypeKey")
			showUInt8("Key", keyMsg.Key)
			showUInt32("Scan", keyMsg.Scan)
			showUInt8("Action", keyMsg.Action)
			showUInt8("Mod", keyMsg.Mod)
		}

		onKey(keyMsg)

	case msg.TypeFrameBufferSize:
		// FIXME: BRAD SAYS THIS IS NOT INPUT
		var m msg.MessageFrameBufferSize
		msg.MustDeserialize(message, &m)

		if self.DebugPrintInputEvents {
			fmt.Print("TypeFrameBufferSize")
			showUInt32("X", m.X)
			showUInt32("Y", m.Y)
		}

		onFrameBufferSize(m)
	default:
		fmt.Println("UNKNOWN MESSAGE TYPE!")
	}

	if self.DebugPrintInputEvents {
		fmt.Println()
	}

	//curRecByte = 0
	return message
}

func showUInt8(s string, x uint8) uint8 {
	fmt.Printf("   [%s: %d]", s, x)
	return x
}

// the rest of these funcs are almost identical, just top 2 vars customized (and string format)
func showSInt32(s string, x int32) {
	fmt.Printf("   [%s: %d]", s, x)
}

func showUInt32(s string, x uint32) {
	fmt.Printf("   [%s: %d]", s, x)
}

func showFloat64(s string, f float64) float64 {
	fmt.Printf("   [%s: %.1f]", s, f)
	return f
}

//
//EVENT HANDLERS
//

// triggered both by moving **AND*** by pressing buttons
func onMouseCursorPos(m msg.MessageMousePos) {
	/*
		x := float32(m.X)
		y := float32(m.Y)

		mouse.UpdatePosition(
			app.Vec2F{x, y},
			gfx.CanvasExtents,
			gfx.PixelSize) // state update

		// rendering update
		//if LMB held
		gl.GlfwWindow.GetMouseButton(glfw.MouseButtonLeft) == glfw.Press {
			ScrollPanelThatIsHoveredOver(mouse.PixelDelta.X, mouse.PixelDelta.Y)
		}
	*/
}

func onMouseScroll(m msg.MessageMouseScroll) {
	/*
		var delta float32 = 30

		if eitherControlKeyHeld() { // horizontal ability from 1D scrolling
			ScrollPanelThatIsHoveredOver(float32(m.Y)*-delta, 0)
		} else { // can handle both x & y for 2D scrolling
			ScrollPanelThatIsHoveredOver(float32(m.X)*delta, float32(m.Y)*-delta)
		}
	*/
}

func onFrameBufferSize(m msg.MessageFrameBufferSize) {
	/*
		fmt.Printf("onFrameBufferSize() - x, y: %d, %d\n", m.X, m.Y)
		gfx.CurrAppWidth = int32(m.X)
		gfx.CurrAppHeight = int32(m.Y)
		gfx.SetSize()
		SetSize()
	*/
}

func onChar(m msg.MessageOnCharacter) {
	//InsertRuneIntoDocument("Rune", m.Rune)
	//script.Process(false)
}

func onKey(m msg.MessageKey) {
	/*
		foc := Focused

		if glfw.Action(m.Action) == glfw.Release {
			fmt.Println("release --------- ", m.Key)

			switch glfw.Key(m.Key) {

			case glfw.KeyEscape:
				fmt.Println("case glfw.KeyEscape:")
				gl.GlfwWindow.SetShouldClose(true)
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
			fmt.Println("press --------- ", m.Key)

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

			//script.Process(false)
		}
	*/
}

func onMouseButton(m msg.MessageMouseButton) {
	/*
		convertClickToTextCursorPosition(m.Button, m.Action)

		if glfw.Action(m.Action) == glfw.Press {
			switch glfw.MouseButton(m.Button) {
			case glfw.MouseButtonLeft:
				// respond to clicks in ui rectangles
				if mouse.CursorIsInside(ui.MainMenu.Rect) {
					respondToAnyMenuButtonClicks()
				} else { // respond to any panel clicks outside of menu
					for _, pan := range Panels {
						if pan.ContainsMouseCursor() {
							pan.RespondToMouseClick()
						}
					}
				}
			}
		}
	*/
}

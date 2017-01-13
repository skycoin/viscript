package hypervisor

import (
	"fmt"
	"github.com/corpusc/viscript/app"
	"github.com/corpusc/viscript/hypervisor/input/mouse"
	"github.com/corpusc/viscript/msg"
	//"github.com/corpusc/viscript/script"
	"github.com/corpusc/viscript/ui"
	"github.com/go-gl/glfw/v3.2/glfw"
)

// apparently every time this is fired, a mouse position event is ALSO fired
// func onMouseButton(
// 	w *glfw.Window,
// 	bt glfw.MouseButton,
// 	action glfw.Action,
// 	mod glfw.ModifierKey) {

func onMouseButton(m msg.MessageMouseButton) {
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
}

func convertClickToTextCursorPosition(button, action uint8) {
	if glfw.MouseButton(button) == glfw.MouseButtonLeft &&
		glfw.Action(action) == glfw.Press {

		foc := Focused

		if foc.IsEditable && foc.Content.Contains(mouse.GlX, mouse.GlY) {
			if foc.MouseY < len(foc.TextBodies[0]) {
				foc.CursY = foc.MouseY

				if foc.MouseX <= len(foc.TextBodies[0][foc.CursY]) {
					foc.CursX = foc.MouseX
				} else {
					foc.CursX = len(foc.TextBodies[0][foc.CursY])
				}
			} else {
				foc.CursY = len(foc.TextBodies[0]) - 1
			}
		}
	}
}

func respondToAnyMenuButtonClicks() {
	for _, bu := range ui.MainMenu.Buttons {
		if mouse.CursorIsInside(bu.Rect) {
			bu.Activated = !bu.Activated

			switch bu.Name {
			case "Run":
				if bu.Activated {
					//script.Process(true)
				}
				break
			case "Testing Tree":
				if bu.Activated {
					//script.Process(true)
					//script.MakeTree()
				} else { // deactivated
					// remove all panels with trees
					b := Panels[:0]
					for _, pan := range Panels {
						if len(pan.Trees) < 1 {
							b = append(b, pan)
						}
					}
					Panels = b
					//fmt.Printf("len of b (from gfx.Panels) after removing ones with trees: %d\n", len(b))
					//fmt.Printf("len of gfx.Panels: %d\n", len(gfx.Panels))
				}
				break
			}

			app.Con.Add(fmt.Sprintf("%s toggled\n", bu.Name))
		}
	}
}

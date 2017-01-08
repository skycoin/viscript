package hypervisor

import (
	"fmt"
	"github.com/corpusc/viscript/gfx"
	"github.com/corpusc/viscript/msg"
	"github.com/corpusc/viscript/script"
	"github.com/corpusc/viscript/ui"
	//"github.com/go-gl/glfw/v3.2/glfw"
)

// apparently every time this is fired, a mouse position event is ALSO fired
// func onMouseButton(
// 	w *glfw.Window,
// 	bt glfw.MouseButton,
// 	action glfw.Action,
// 	mod glfw.ModifierKey) {

func onMouseButton(m msg.MessageMouseButton) {

	gfx.Curs.ConvertMouseClickToTextCursorPosition(m.Button, m.Action)

	/*
		if action == glfw.Press {
			switch glfw.MouseButton(bt) {
			case glfw.MouseButtonLeft:
				// respond to clicks in ui rectangles
				if gfx.MouseCursorIsInside(ui.MainMenu.Rect) {
					respondToAnyMenuButtonClicks()
				} else { // respond to any panel clicks outside of menu
					for _, pan := range gfx.Rend.Panels {
						if pan.ContainsMouseCursor() {
							pan.RespondToMouseClick()
						}
					}
				}
			}
		}
	*/

	//MessageMouseButton
	/*
		var m msg.MessageMouseButton
		m.Button = uint8(bt)
		m.Action = uint8(action)
		m.Mod = uint8(mod)
		//DispatchEvent(msg.TypeMouseButton, m)
		b := msg.Serialize(msg.TypeMouseButton, m)
		InputEvents <- b
	*/
}

func respondToAnyMenuButtonClicks() {
	for _, bu := range ui.MainMenu.Buttons {
		if gfx.MouseCursorIsInside(bu.Rect) {
			bu.Activated = !bu.Activated

			switch bu.Name {
			case "Run":
				if bu.Activated {
					script.Process(true)
				}
				break
			case "Testing Tree":
				if bu.Activated {
					script.Process(true)
					script.MakeTree()
				} else { // deactivated
					// remove all panels with trees
					b := gfx.Rend.Panels[:0]
					for _, pan := range gfx.Rend.Panels {
						if len(pan.Trees) < 1 {
							b = append(b, pan)
						}
					}
					gfx.Rend.Panels = b
					//fmt.Printf("len of b (from gfx.Rend.Panels) after removing ones with trees: %d\n", len(b))
					//fmt.Printf("len of gfx.Rend.Panels: %d\n", len(gfx.Rend.Panels))
				}
				break
			}

			gfx.Con.Add(fmt.Sprintf("%s toggled\n", bu.Name))
		}
	}
}

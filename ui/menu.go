package ui

import (
	"fmt"
	"github.com/corpusc/viscript/common"
)

var MainMenu = Menu{}

func init() {
	fmt.Println("ui.init() in menu.go")
	MainMenu.Buttons = append(MainMenu.Buttons, &Button{Name: "Run"})
	MainMenu.Buttons = append(MainMenu.Buttons, &Button{Name: "Syntax Tree"})
	MainMenu.Buttons = append(MainMenu.Buttons, &Button{Name: "Item 3"})
	MainMenu.Buttons = append(MainMenu.Buttons, &Button{Name: "Item 4"})
}

type Menu struct {
	IsVertical bool // controls which dimension gets divided up for button sizes
	Rect       *common.Rectangle
	Buttons    []*Button
}

func (m *Menu) SetSize(rect *common.Rectangle) {
	m.Rect = rect

	// depending on vertical or horizontal layout, only 1 dimension (of the below 4 variables) is used
	x := m.Rect.Left
	y := m.Rect.Top
	wid := m.Rect.Width() / float32(len(m.Buttons))  // width of buttons
	hei := m.Rect.Height() / float32(len(m.Buttons)) // height of buttons

	for _, b := range m.Buttons {
		nr := &common.Rectangle{rect.Top, rect.Right, rect.Bottom, rect.Left}

		if m.IsVertical {
			nr.Top = y
			nr.Bottom = y - hei
		} else { // horizontally laid out
			nr.Left = x
			nr.Right = x + wid
		}

		b.Rect = nr

		x += wid
		y -= hei
	}
}

package render

import (
	"fmt"
	"github.com/corpusc/viscript/common"
	"github.com/corpusc/viscript/ui"
)

var MenuInst = Menu{}

type Menu struct {
	IsVertical bool // controls which dimension gets divided up for button sizes
	Rect       *common.Rectangle
	Buttons    []*ui.Button
}

func (m *Menu) Init() {
	fmt.Println("menu.Init()")
	m.Buttons = append(m.Buttons, &ui.Button{Name: "Run"})
	m.Buttons = append(m.Buttons, &ui.Button{Name: "Syntax Tree"})
	m.Buttons = append(m.Buttons, &ui.Button{Name: "Menu Item 3"})
	m.SetSize()
}

func (m *Menu) SetSize() {
	m.Rect = m.GetMenuSizedRect()
	fmt.Printf("m.Rect: %s\n", m.Rect)

	// depending on vertical or horizontal layout, only 1 dimension (of the below 4 variables) is used
	x := m.Rect.Left
	y := m.Rect.Top
	wid := m.Rect.Width() / float32(len(m.Buttons))  // width of buttons
	hei := m.Rect.Height() / float32(len(m.Buttons)) // height of buttons

	for _, b := range m.Buttons {
		nr := m.GetMenuSizedRect()

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

func (m *Menu) Draw() {
	for _, bu := range m.Buttons {
		if bu.Activated {
			Rend.Color(Green)
		} else {
			Rend.Color(White)
		}

		span := bu.Rect.Height() * goldenPercentage // ...of both dimensions of each character
		glTextWidth := float32(len(bu.Name)) * span // in terms of OpenGL/float32 space
		x := bu.Rect.Left + (bu.Rect.Width()-glTextWidth)/2
		verticalLipSpan := (bu.Rect.Height() - span) / 2 // lip or frame edge

		Rend.DrawQuad(11, 13, bu.Rect)

		for _, c := range bu.Name {
			Rend.DrawCharAtRect(c, &common.Rectangle{bu.Rect.Top - verticalLipSpan, x + span, bu.Rect.Bottom + verticalLipSpan, x})
			x += span
		}
	}
}

func (m *Menu) GetMenuSizedRect() *common.Rectangle {
	return &common.Rectangle{
		Rend.ClientExtentY,
		Rend.ClientExtentX,
		Rend.ClientExtentY - Rend.CharHei,
		-Rend.ClientExtentX}
}

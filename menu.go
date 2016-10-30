package main

import (
	"fmt"
)

var menu = Menu{}

type Menu struct {
	IsVertical bool // controls which dimension gets divided up for button sizes
	Rect       *Rectangle
	Buttons    []*Button
}

func (m *Menu) Init() {
	fmt.Println("menu.Init()")
	m.Buttons = append(m.Buttons, &Button{Name: "Run"})
	m.Buttons = append(m.Buttons, &Button{Name: "Syntax Tree"})
	m.Buttons = append(m.Buttons, &Button{Name: "Menu Item 3"})
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
	for _, b := range m.Buttons {
		b.Draw()
	}
}

func (m *Menu) GetMenuSizedRect() *Rectangle {
	return &Rectangle{
		rend.ClientExtentY,
		rend.ClientExtentX,
		rend.ClientExtentY - rend.CharHei,
		-rend.ClientExtentX}
}

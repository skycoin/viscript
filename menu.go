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
	m.Rect = &Rectangle{
		rend.ScreenRad,
		rend.ScreenRad,
		rend.ScreenRad - rend.CharHei,
		-rend.ScreenRad}

	fmt.Println("menu.Init()")
	fmt.Printf("rend.CharHei: %.2f\n", rend.CharHei)
	fmt.Printf("m.Rect.Bottom: %.2f\n", m.Rect.Bottom)
	fmt.Printf("m.Rect: %s\n", m.Rect)
	m.Buttons = append(m.Buttons, &Button{Name: "Syntax Tree"})
	m.Buttons = append(m.Buttons, &Button{Name: "Run"})
	nr := m.Rect //&Rectangle{} // new rectangle

	// depending on vertical or horizontal layout, only 1 dimension (of the below 4 variables) is used
	x := nr.Left
	y := nr.Top
	wid := nr.Width() / float32(len(m.Buttons))  // width of buttons
	hei := nr.Height() / float32(len(m.Buttons)) // height of buttons

	for _, b := range m.Buttons {
		if m.IsVertical {
			nr.Top = y
			nr.Bottom = y - hei
		} else { // horizontally laid out
			nr.Left = x
			nr.Right = x + wid
		}

		b.Rect = nr
		fmt.Printf("nr: %s\n", nr)

		x += wid
		y -= hei
	}
}

func (m *Menu) Draw() {
	for _, b := range m.Buttons {
		b.Draw()
	}
}

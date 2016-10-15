package main

import (
//"fmt"
)

var menu = &Menu{}

type Menu struct {
	Rect       *Rectangle
	IsVertical bool // controls which dimension gets divided up for button sizes
	Buttons    []*Button
}

func (m *Menu) Init() {
	m.Buttons = append(m.Buttons, &Button{Name: "Syntax Tree"})
	m.Buttons = append(m.Buttons, &Button{Name: "Run"})
	nr := m.Rect //&Rectangle{} // new rectangle

	if m.IsVertical {
		nr.Left = 0
		nr.Right = 0
	} else { // horizontally laid out
		nr.Top = 0
		nr.Bottom = 0
	}

	for _, b := range m.Buttons {
		b.Draw()
	}
}

func (m *Menu) Draw() {
	for _, b := range m.Buttons {
		b.Draw()
	}
}

package hypervisor

import (
	"fmt"
)

var Focused *Terminal
var Terms []*Terminal

var runOutputTerminalFrac = float32(0.4) // TEMPORARY fraction of vertical strip height which is dedicated to running code

func initTerminals() {
	Terms = append(Terms, &Terminal{FractionOfStrip: 1 - runOutputTerminalFrac, IsEditable: true})
	Terms = append(Terms, &Terminal{FractionOfStrip: runOutputTerminalFrac, IsEditable: false}) // console (runtime feedback log)	// FIXME so its not editable once we're done debugging some things
	Focused = Terms[0]

	Terms[0].Init()
	Terms[0].SetupDemoProgram()
	Terms[1].Init()
}

func ScrollTermThatHasMousePointer(mousePixelDeltaX, mousePixelDeltaY float32) {
	for _, t := range Terms {
		t.ScrollIfMouseOver(mousePixelDeltaX, mousePixelDeltaY)
	}
}

func InsertRuneIntoDocument(s string, message uint32) string {
	f := Focused
	b := f.TextBodies[0]
	resultsDif := f.CursX - len(b[f.CursY])
	fmt.Printf("Rune   [%s: %s]", s, string(message))

	if f.CursX > len(b[f.CursY]) {
		b[f.CursY] = b[f.CursY][:f.CursX-resultsDif] + b[f.CursY][:len(b[f.CursY])] + string(message)
		fmt.Printf("line is %s\n", b[f.CursY])
		f.CursX++
	} else {
		b[f.CursY] = b[f.CursY][:f.CursX] + string(message) + b[f.CursY][f.CursX:len(b[f.CursY])]
		f.CursX++
	}

	return string(message)
}

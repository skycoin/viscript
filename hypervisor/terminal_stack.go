package hypervisor

import (
	"fmt"
)

var Focused *Terminal
var Panels []*Terminal

var runPanelHeiFrac = float32(0.4) // TEMPORARY fraction of vertical strip height which is dedicated to running code

func initPanels() {
	Panels = append(Panels, &Terminal{FractionOfStrip: 1 - runPanelHeiFrac, IsEditable: true})
	Panels = append(Panels, &Terminal{FractionOfStrip: runPanelHeiFrac, IsEditable: true}) // console (runtime feedback log)	// FIXME so its not editable once we're done debugging some things
	Focused = Panels[0]

	Panels[0].Init()
	Panels[0].SetupDemoProgram()
	Panels[1].Init()
}

// refactoring (possibly termporary) additions
func SetSize() {
	for _, pan := range Panels {
		pan.SetSize()
	}
}

func ScrollPanelThatIsHoveredOver(mousePixelDeltaX, mousePixelDeltaY float32) {
	for _, pan := range Panels {
		pan.ScrollIfMouseOver(mousePixelDeltaX, mousePixelDeltaY)
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

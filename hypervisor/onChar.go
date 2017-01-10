package hypervisor

import (
	"fmt"
	"github.com/corpusc/viscript/gfx"
	"github.com/corpusc/viscript/msg"
)

//Working
//func onChar(w *glfw.Window, char rune) {
func onChar(m msg.MessageOnCharacter) {
	InsertRuneIntoDocument("Rune", m.Rune)
	//script.Process(false)
}

func InsertRuneIntoDocument(s string, message uint32) string {
	f := gfx.Rend.Focused
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

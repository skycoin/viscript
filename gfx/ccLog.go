package gfx

import (
	"fmt"
	"github.com/corpusc/viscript/app"
	"strings"
)

var Con = CcLog{} // console log, displays runtime feedback

type CcLog struct {
	Name  string
	Lines []string
}

func (log CcLog) Add(s string) {
	fmt.Printf(s)
	s = strings.Replace(s, "\n", "", -1)
	log.Lines = append(log.Lines, s)

	if len(Panels) > 1 {
		Panels[1].TextBodies[0] = append(Panels[1].TextBodies[0], s)
	}
}

// numLines: use odd number for an exact middle point
func MakeHighlyVisibleLogHeader(s string, numLines int) {
	s = " " + s + " "
	fillChar := "#"
	osOnly := s == app.Name

	var bar string
	for i := 0; i < 79; i++ {
		bar += fillChar
	}

	var spaces string
	for i := 0; i < len(s); i++ {
		spaces += " "
	}

	var bookend string
	for i := 0; i < (79-len(s))/2; i++ {
		bookend += fillChar
	}

	middle := numLines / 2
	for i := 0; i < numLines; i++ {
		switch {
		case i == middle:
			predPrint(osOnly, bookend+s+bookend)
		case i == middle-1 || i == middle+1:
			predPrint(osOnly, bookend+spaces+bookend)
		default:
			predPrint(osOnly, bar)
		}
	}
}

// prints only to OS console window if it's for the app's name
func predPrint(osOnly bool, s string) {
	if osOnly {
		fmt.Println(s)
	} else {
		Con.Add(fmt.Sprintf("%s\n", s))
	}
}

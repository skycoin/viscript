package app

import (
	"fmt"
	//"github.com/skycoin/viscript/hypervisor"
	"strings"
)

var Con = CcLog{} //console log, displays runtime feedback
//private
var fillChar = "*"         //character used to surround/highlight text
const numOSBoxColumns = 79 //assumes 80 column lines.  but to still look fine
//....with more columns, we make each a new line & reserve the last column

type CcLog struct {
	Name  string
	Lines []string
}

func (log CcLog) Add(s string) {
	fmt.Printf(s)
	s = strings.Replace(s, "\n", "", -1)
	log.Lines = append(log.Lines, s)

	/*
		if len(Terms) > 1 {
			Terms[1].TextBodies[0] = append(Terms[1].TextBodies[0], s)
		}
	*/
}

func GetLabeledBarOfChars(label, chars string, totalChars uint32) string {
	bar := label

	modifyLeftEdge := true
	for len(bar) < int(totalChars) {
		if modifyLeftEdge {
			bar = chars + bar
		} else {
			bar += chars
		}

		modifyLeftEdge = !modifyLeftEdge
	}

	modifyLeftEdge = true
	for len(bar) > int(totalChars) {
		if modifyLeftEdge {
			bar = bar[1:]
		} else {
			bar = bar[:len(bar)-1]
		}

		modifyLeftEdge = !modifyLeftEdge
	}

	return bar
}

func GetBarOfChars(chars string, howMany int) string {
	bar := ""

	for i := 0; i < howMany; i++ {
		bar += chars
	}

	return bar
}

// numLines: use odd number for an exact middle point
func MakeHighlyVisibleLogEntry(s string, numLines int) {
	s = " " + s + " "
	osBoxOnly := s == Name

	bar := GetBarOfChars(fillChar, numOSBoxColumns)

	var spaces string
	for i := 0; i < len(s); i++ {
		spaces += " "
	}

	bookend := ""
	for i := 0; i < (numOSBoxColumns-len(s))/2; i++ {
		bookend += fillChar
	}

	vMid := numLines / 2 //vertical middle
	for i := 0; i < numLines; i++ {
		switch {
		case i == vMid:
			maybePrint(osBoxOnly, bookend+s+bookend)
		case i == vMid-1 || i == vMid+1:
			maybePrint(osBoxOnly, bookend+spaces+bookend)
		default:
			maybePrint(osBoxOnly, bar)
		}
	}
}

// prints only to OS console window if it's for the app's name
func maybePrint(osBoxOnly bool, s string) {
	if len(s) < numOSBoxColumns {
		s += fillChar
	}

	if osBoxOnly {
		println(s)
	} else {
		Con.Add(fmt.Sprintf("%s\n", s))
	}
}

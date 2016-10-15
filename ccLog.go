package main

import (
	"fmt"
	"strings"
)

var con = CcLog{} // console log, displays runtime feedback (including parsing errors)

type CcLog struct {
	Name  string
	Lines []string
}

func (log CcLog) Add(s string) {
	fmt.Printf(s)
	s = strings.Replace(s, "\n", "", -1)
	log.Lines = append(log.Lines, s)

	if len(rend.Panels) > 1 {
		rend.Panels[1].Body = append(rend.Panels[1].Body, s)
	}
}

func (log CcLog) Draw() {
	// TODO: add the graphical drawing
}

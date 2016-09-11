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
	s = strings.Replace(s, "\n", "", -1)
	log.Lines = append(log.Lines, s)
	textRend.Panels[1].Body = append(textRend.Panels[1].Body, s)
	fmt.Printf(s)
}

func (log CcLog) Draw() {
	// TODO: add the graphical drawing
}

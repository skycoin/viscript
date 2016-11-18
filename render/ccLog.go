package render

import (
	"fmt"
	"strings"
)

var Con = CcLog{} // console log, displays runtime feedback (including parsing errors)

type CcLog struct {
	Name  string
	Lines []string
}

func (log CcLog) Add(s string) {
	fmt.Printf(s)
	s = strings.Replace(s, "\n", "", -1)
	log.Lines = append(log.Lines, s)

	if len(Rend.Panels) > 1 {
		Rend.Panels[1].Body = append(Rend.Panels[1].Body, s)
	}
}

func (log CcLog) Draw() {
	// TODO: add the graphical drawing
}

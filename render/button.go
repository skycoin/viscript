package render

import (
	//"fmt"
	"github.com/corpusc/viscript/common"
)

type Button struct {
	Name      string
	Activated bool
	Rect      *common.Rectangle
}

func (bu *Button) Draw() {
	if bu.Activated {
		Rend.Color(green)
	} else {
		Rend.Color(White)
	}

	span := bu.Rect.Height() * goldenPercentage // ...of both dimensions of each character
	glTextWidth := float32(len(bu.Name)) * span // in terms of OpenGL/float32 space
	x := bu.Rect.Left + (bu.Rect.Width()-glTextWidth)/2
	verticalLipSpan := (bu.Rect.Height() - span) / 2 // lip or frame edge

	Rend.DrawQuad(11, 13, bu.Rect)

	for _, c := range bu.Name {
		Rend.DrawCharAtRect(c, &common.Rectangle{bu.Rect.Top - verticalLipSpan, x + span, bu.Rect.Bottom + verticalLipSpan, x})
		x += span
	}
}

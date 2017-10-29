package viewport

import (
	"github.com/skycoin/viscript/app"
	"github.com/skycoin/viscript/viewport/gl"
	"github.com/skycoin/viscript/viewport/terminal"
)

var taskBarDepth float32 = 10.1

func drawTaskBar() {
	drawTaskBarBackground()
	drawStartButton()
	drawTerminalButtons()
}

func drawTaskBarBackground() {
	hmm := float32(0.3)

	gl.Draw9SlicedRect(
		gl.Pic_GradientBorder,
		&app.Rectangle{
			-gl.CanvasExtents.Y + hmm,
			gl.CanvasExtents.X,
			-gl.CanvasExtents.Y,
			-gl.CanvasExtents.X},
		taskBarDepth)
}

func drawStartButton() {
}

func drawTerminalButtons() {
	terminal.Terms.DrawTaskBarButtons()
}

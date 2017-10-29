package viewport

import (
	"github.com/skycoin/viscript/app"
	"github.com/skycoin/viscript/viewport/gl"
	"github.com/skycoin/viscript/viewport/terminal"
)

var (
	taskBarDepth      float32 = 10.1
	taskBarHeight     float32 = 0.1
	taskBarBorderSpan float32 = taskBarHeight / 12
	taskBarCurrX      float32
	taskBarBounds     *app.Rectangle
)

func drawTaskBar() {
	drawTaskBarBackground()
	drawStartButton()
	drawTerminalButtons()
}

func drawTaskBarBackground() {
	taskBarBounds = &app.Rectangle{
		-gl.CanvasExtents.Y + taskBarHeight,
		gl.CanvasExtents.X,
		-gl.CanvasExtents.Y,
		-gl.CanvasExtents.X}

	gl.Draw9SlicedRect(
		gl.Pic_GradientBorder,
		taskBarBounds,
		taskBarDepth)
}

func drawStartButton() {
	taskBarBounds.Top -= taskBarBorderSpan
	taskBarBounds.Bottom += taskBarBorderSpan
	taskBarBounds.Left += taskBarBorderSpan
	taskBarBounds.Right = taskBarBounds.Left + taskBarBounds.Height()

	taskBarCurrX = taskBarBounds.Right

	gl.Draw9SlicedRect(
		gl.Pic_GradientBorder,
		taskBarBounds,
		taskBarDepth)

	taskBarBounds.Left += taskBarBorderSpan
	taskBarBounds.Right -= taskBarBorderSpan
	taskBarBounds.Top -= taskBarBorderSpan
	taskBarBounds.Bottom += taskBarBorderSpan

	gl.Draw9SlicedRect(
		gl.Pic_TriangleUp,
		taskBarBounds,
		taskBarDepth)
}

func drawTerminalButtons() {
	for _, term := range terminal.Terms.TermMap {
		if term.TerminalId == terminal.Terms.FocusedId {
			gl.SetColor(gl.White)
		} else {
			gl.SetColor(gl.Gray)
		}

		//drawIdTab(term, taskBarDepth)
	}
}

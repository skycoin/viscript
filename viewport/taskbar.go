package viewport

import (
	"github.com/skycoin/viscript/app"
	"github.com/skycoin/viscript/viewport/gl"
	"github.com/skycoin/viscript/viewport/terminal"
)

var (
	taskBarDepth      float32 = 11
	taskBarHeight     float32 = 0.1
	taskBarBorderSpan float32 = taskBarHeight / 12
	taskBarCharWid    float32 = taskBarHeight/2 - taskBarBorderSpan*2
	buttonBounds      *app.Rectangle
	charBounds        *app.Rectangle
)

func drawTaskBar() {
	gl.SetColor(gl.Gray)
	drawTaskBarBackground()
	drawStartButton()
	drawTerminalButtons()
}

func drawTaskBarBackground() {
	buttonBounds = &app.Rectangle{
		-gl.CanvasExtents.Y + taskBarHeight,
		gl.CanvasExtents.X,
		-gl.CanvasExtents.Y,
		-gl.CanvasExtents.X}

	gl.Draw9SlicedRect(
		gl.Pic_GradientBorder,
		buttonBounds,
		taskBarDepth)
}

func drawStartButton() {
	//now make buttons inset from task bar
	buttonBounds.Top -= taskBarBorderSpan
	buttonBounds.Bottom += taskBarBorderSpan
	buttonBounds.Left += taskBarBorderSpan
	buttonBounds.Right = buttonBounds.Left + buttonBounds.Height()

	gl.Draw9SlicedRect(
		gl.Pic_GradientBorder,
		buttonBounds,
		taskBarDepth)

	charBounds = &app.Rectangle{
		buttonBounds.Top - taskBarBorderSpan,
		buttonBounds.Right - taskBarBorderSpan,
		buttonBounds.Bottom + taskBarBorderSpan,
		buttonBounds.Left + taskBarBorderSpan}

	buttonBounds.Left += buttonBounds.Width()
	buttonBounds.Right += buttonBounds.Width()

	gl.Draw9SlicedRect(
		gl.Pic_TriangleUp,
		charBounds,
		taskBarDepth)

	charBounds.Left += taskBarCharWid
	charBounds.Right += taskBarCharWid
}

func drawTerminalButtons() {
	for _, term := range terminal.Terms.TermMap {
		charBounds.Left = buttonBounds.Left + taskBarBorderSpan
		charBounds.Right = charBounds.Left + taskBarCharWid

		if term.TerminalId == terminal.Terms.FocusedId {
			gl.SetColor(gl.White)
		} else {
			gl.SetColor(gl.Gray)
		}

		max := len(term.TabText)
		currTextWid := taskBarCharWid*float32(max) + taskBarBorderSpan*2
		buttonBounds.Right = buttonBounds.Left + currTextWid

		//draw button background
		gl.Draw9SlicedRect(
			gl.Pic_GradientBorder,
			buttonBounds,
			taskBarDepth)

		//draw the id # text
		for i := 0; i < max; i++ {
			gl.DrawCharAtRect(rune(term.TabText[i]), charBounds, taskBarDepth)
			charBounds.Left += taskBarCharWid
			charBounds.Right += taskBarCharWid
		}

		buttonBounds.Left += currTextWid
		buttonBounds.Right += currTextWid
	}
}

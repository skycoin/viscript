package viewport

import (
	"github.com/skycoin/viscript/app"
	"github.com/skycoin/viscript/viewport/gl"
	"github.com/skycoin/viscript/viewport/terminal"
)

var (
	buttonBounds *app.Rectangle
	charBounds   *app.Rectangle
)

func drawTaskBar() {
	gl.SetColor(gl.Gray)
	drawTaskBarBackground()
	drawStartButton()
	drawTerminalButtons()
}

func drawTaskBarBackground() {
	buttonBounds = &app.Rectangle{
		-gl.CanvasExtents.Y + app.TaskBarHeight,
		gl.CanvasExtents.X,
		-gl.CanvasExtents.Y,
		-gl.CanvasExtents.X}

	gl.Draw9SlicedRect(
		gl.Pic_GradientBorder,
		buttonBounds,
		app.TaskBarDepth)
}

func drawStartButton() {
	//now make buttons inset from task bar
	buttonBounds.Top -= app.TaskBarBorderSpan
	buttonBounds.Bottom += app.TaskBarBorderSpan
	buttonBounds.Left += app.TaskBarBorderSpan
	buttonBounds.Right = buttonBounds.Left + buttonBounds.Height()

	gl.Draw9SlicedRect(
		gl.Pic_GradientBorder,
		buttonBounds,
		app.TaskBarDepth)

	charBounds = &app.Rectangle{
		buttonBounds.Top - app.TaskBarBorderSpan,
		buttonBounds.Right - app.TaskBarBorderSpan,
		buttonBounds.Bottom + app.TaskBarBorderSpan,
		buttonBounds.Left + app.TaskBarBorderSpan}

	buttonBounds.Left += buttonBounds.Width()
	buttonBounds.Right += buttonBounds.Width()

	gl.Draw9SlicedRect(
		gl.Pic_TriangleUp,
		charBounds,
		app.TaskBarDepth)

	charBounds.Left += app.TaskBarCharWid
	charBounds.Right += app.TaskBarCharWid
}

func drawTerminalButtons() {
	for _, term := range terminal.Terms.TermMap {
		charBounds.Left = term.TaskBarButton.Left + app.TaskBarBorderSpan
		charBounds.Right = charBounds.Left + app.TaskBarCharWid

		if term.TerminalId == terminal.Terms.FocusedId {
			gl.SetColor(gl.White)
		} else {
			gl.SetColor(gl.Gray)
		}

		max := len(term.TabText)
		term.TaskBarButton.Right = charBounds.Left + float32(max)*app.TaskBarCharWid + app.TaskBarBorderSpan

		//draw button background
		gl.Draw9SlicedRect(
			gl.Pic_GradientBorder,
			term.TaskBarButton,
			app.TaskBarDepth)

		//draw the id # text
		for i := 0; i < max; i++ {
			gl.DrawCharAtRect(rune(term.TabText[i]), charBounds, app.TaskBarDepth)
			charBounds.Left += app.TaskBarCharWid
			charBounds.Right += app.TaskBarCharWid
		}
	}
}

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
	taskBarCurrX      float32
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
	buttonBounds.Top -= taskBarBorderSpan
	buttonBounds.Bottom += taskBarBorderSpan
	buttonBounds.Left += taskBarBorderSpan
	buttonBounds.Right = buttonBounds.Left + buttonBounds.Height()

	taskBarCurrX = buttonBounds.Right

	gl.Draw9SlicedRect(
		gl.Pic_GradientBorder,
		buttonBounds,
		taskBarDepth)

	charBounds = &app.Rectangle{
		buttonBounds.Top - taskBarBorderSpan,
		buttonBounds.Right - taskBarBorderSpan,
		buttonBounds.Bottom + taskBarBorderSpan,
		buttonBounds.Left + taskBarBorderSpan}

	buttonBounds.Left += buttonBounds.Height()
	buttonBounds.Right += buttonBounds.Height()

	gl.Draw9SlicedRect(
		gl.Pic_TriangleUp,
		charBounds,
		taskBarDepth)

	charBounds.Left += charBounds.Height()
	charBounds.Right += charBounds.Height()
}

func drawTerminalButtons() {
	for _, term := range terminal.Terms.TermMap {
		charBounds.Left += taskBarBorderSpan * 10  //2
		charBounds.Right += taskBarBorderSpan * 10 //2

		if term.TerminalId == terminal.Terms.FocusedId {
			gl.SetColor(gl.White)
		} else {
			gl.SetColor(gl.Gray)
		}

		max := len(term.TabText)
		currTextWid := charBounds.Width()*float32(max) + taskBarBorderSpan*2
		buttonBounds.Right = buttonBounds.Left + currTextWid

		//draw button background
		gl.Draw9SlicedRect(
			gl.Pic_GradientBorder,
			buttonBounds,
			taskBarDepth)

		//draw the id # text
		// for i := 0; i < max; i++ {
		// 	gl.DrawCharAtRect(rune(term.TabText[i]), charBounds, taskBarDepth)
		// 	charBounds.Left += term.CharSize.X
		// 	charBounds.Right += term.CharSize.X
		// }

		buttonBounds.Left += currTextWid
		buttonBounds.Right += currTextWid
	}
}

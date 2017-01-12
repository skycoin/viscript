package gfx

import (
	"fmt"
	"github.com/corpusc/viscript/app"
	"github.com/corpusc/viscript/hypervisor/input/mouse"
	"github.com/corpusc/viscript/tree"
	"github.com/corpusc/viscript/ui"
	"github.com/go-gl/gl/v2.1/gl"
	"math"
)

type ScrollablePanel struct {
	FractionOfStrip float32 // fraction of the parent PanelStrip (in 1 dimension)
	CursX           int     // current cursor/insert position (in character grid cells/units)
	CursY           int
	MouseX          int // current mouse position in character grid space (units/cells)
	MouseY          int
	IsEditable      bool // editing is hardwired to TextBodies[0], but we probably never want
	// to edit text unless the whole panel is dedicated to just one TextBody (& no graphical trees)
	Whole      *app.Rectangle // the whole panel, including chrome (title bar & scroll bars)
	Content    *app.Rectangle // viewport into virtual space, bordered by above rect
	Selection  *ui.SelectionRange
	BarHori    *ui.ScrollBar // horizontal
	BarVert    *ui.ScrollBar // vertical
	TextBodies [][]string
	TextColors []*ColorSpot
	Trees      []*tree.Tree
}

func (sp *ScrollablePanel) Init() {
	fmt.Printf("ScrollablePanel.Init()\n")

	sp.TextBodies = append(sp.TextBodies, []string{})

	sp.Selection = &ui.SelectionRange{}
	sp.Selection.Init()

	// scrollbars
	sp.BarHori = &ui.ScrollBar{IsHorizontal: true}
	sp.BarVert = &ui.ScrollBar{}
	sp.BarHori.Rect = &app.Rectangle{}
	sp.BarVert.Rect = &app.Rectangle{}

	sp.SetSize()
}

func (sp *ScrollablePanel) SetSize() {
	fmt.Printf("ScrollablePanel.SetSize()\n")

	sp.Whole = &app.Rectangle{
		CanvasExtents.Y - CharHei,
		CanvasExtents.X,
		-CanvasExtents.Y,
		-CanvasExtents.X}

	if sp.FractionOfStrip == runPanelHeiFrac { // FIXME: this is hardwired for one use case for now
		sp.Whole.Top = sp.Whole.Bottom + sp.Whole.Height()*sp.FractionOfStrip
	} else {
		sp.Whole.Bottom = sp.Whole.Bottom + sp.Whole.Height()*runPanelHeiFrac
	}

	sp.Content = &app.Rectangle{}
	sp.Content.Top = sp.Whole.Top
	sp.Content.Right = sp.Whole.Right - ui.ScrollBarThickness
	sp.Content.Bottom = sp.Whole.Bottom + ui.ScrollBarThickness
	sp.Content.Left = sp.Whole.Left

	// set scrollbars' upper left corners
	sp.BarHori.Rect.Left = sp.Whole.Left
	sp.BarHori.Rect.Top = sp.Content.Bottom
	sp.BarVert.Rect.Left = sp.Content.Right
	sp.BarVert.Rect.Top = sp.Whole.Top
}

func (sp *ScrollablePanel) RespondToMouseClick() {
	Focused = sp

	// diffs/deltas from home position of panel (top left corner)
	glDeltaFromHome := app.Vec2F{
		mouse.GlX - sp.Whole.Left,
		mouse.GlY - sp.Whole.Top}

	sp.MouseX = int((glDeltaFromHome.X + sp.BarHori.ScrollDelta) / CharWid)
	sp.MouseY = int(-(glDeltaFromHome.Y + sp.BarVert.ScrollDelta) / CharHei)

	if sp.MouseY < 0 {
		sp.MouseY = 0
	}

	if sp.MouseY >= len(sp.TextBodies[0]) {
		sp.MouseY = len(sp.TextBodies[0]) - 1
	}
}

func (sp *ScrollablePanel) GoToTopEdge() {
	CurrY = sp.Whole.Top - sp.BarVert.ScrollDelta
}
func (sp *ScrollablePanel) GoToLeftEdge() float32 {
	CurrX = sp.Whole.Left - sp.BarHori.ScrollDelta
	return CurrX
}
func (sp *ScrollablePanel) GoToTopLeftCorner() {
	sp.GoToTopEdge()
	sp.GoToLeftEdge()
}

func (sp *ScrollablePanel) Draw() {
	sp.GoToTopLeftCorner()
	sp.DrawBackground(11, 13)
	sp.DrawText()
	SetColor(GrayDark)
	sp.DrawScrollbarChrome(10, 11, sp.Whole.Right-ui.ScrollBarThickness, sp.Whole.Top)                          // vertical bar background
	sp.DrawScrollbarChrome(13, 12, sp.Whole.Left, sp.Whole.Bottom+ui.ScrollBarThickness)                        // horizontal bar background
	sp.DrawScrollbarChrome(12, 11, sp.Whole.Right-ui.ScrollBarThickness, sp.Whole.Bottom+ui.ScrollBarThickness) // corner elbow piece
	SetColor(Gray)
	sp.BarHori.SetSize(sp.Whole, sp.TextBodies[0], CharWid, CharHei) // FIXME? (to consider multiple bodies & multiple trees)
	sp.BarVert.SetSize(sp.Whole, sp.TextBodies[0], CharWid, CharHei)
	DrawStretchableRect(11, 13, sp.BarHori.Rect) // 2,11 (pixel checkerboard)    // 14, 15 (square in the middle)
	DrawStretchableRect(11, 13, sp.BarVert.Rect) // 13, 12 (double horizontal lines)    // 10, 11 (double vertical lines)
	SetColor(White)
	sp.DrawTree()
}

func (sp *ScrollablePanel) DrawText() {
	cX := CurrX // current drawing position
	cY := CurrY
	cW := CharWid
	cH := CharHei
	b := sp.BarHori.Rect.Top // bottom of text area

	// setup for colored text
	ncId := 0         // next color
	var nc *ColorSpot // ^
	if /* colors exist */ len(sp.TextColors) > 0 {
		nc = sp.TextColors[ncId]
	}

	// iterate over lines
	for y, line := range sp.TextBodies[0] {
		lineVisible := cY <= sp.Whole.Top+cH && cY >= b

		if lineVisible {
			r := &app.Rectangle{cY, cX + cW, cY - cH, cX} // t, r, b, l

			// if line needs vertical adjustment
			if cY > sp.Whole.Top {
				r.Top = sp.Whole.Top
			}
			if cY-cH < b {
				r.Bottom = b
			}

			// iterate over runes
			SetColor(Gray)
			for x, c := range line {
				ncId, nc = sp.changeColorIfCodeAt(x, y, ncId, nc)

				// drawing
				if /* char visible */ cX >= sp.Whole.Left-cW && cX < sp.BarVert.Rect.Left {
					app.ClampLeftAndRightOf(r, sp.Whole.Left, sp.BarVert.Rect.Left)
					DrawCharAtRect(c, r)

					if sp.IsEditable { //&& Curs.Visible == true {
						if x == sp.CursX && y == sp.CursY {
							SetColor(White)
							//DrawCharAtRect('_', r)
							DrawStretchableRect(11, 13, Curs.GetAnimationModifiedRect(*r))
							SetColor(PrevColor)
						}
					}
				}

				cX += cW
				r.Left = cX
				r.Right = cX + cW
			}

			// draw cursor at the end of line if needed
			if cX < sp.BarVert.Rect.Left && y == sp.CursY && sp.CursX == len(line) {
				if sp.IsEditable { //&& Curs.Visible == true {
					SetColor(White)
					app.ClampLeftAndRightOf(r, sp.Whole.Left, sp.BarVert.Rect.Left)
					//DrawCharAtRect('_', r)
					DrawStretchableRect(11, 13, Curs.GetAnimationModifiedRect(*r))
				}
			}

			cX = sp.GoToLeftEdge()
		} else { // line not visible
			for x := range line {
				ncId, nc = sp.changeColorIfCodeAt(x, y, ncId, nc)
			}
		}

		cY -= cH // go down a line height
	}
}

func (sp *ScrollablePanel) changeColorIfCodeAt(x, y, ncId int, nc *ColorSpot) (int, *ColorSpot) {
	if /* colors exist */ len(sp.TextColors) > 0 {
		if x == nc.Pos.X &&
			y == nc.Pos.Y {
			SetColor(nc.Color)
			//fmt.Println("-------- nc-------, then 3rd():", nc, sp.TextColors[2])
			ncId++

			if ncId < len(sp.TextColors) {
				nc = sp.TextColors[ncId]
			}
		}
	}

	return ncId, nc
}

// ATM the only different between the 2 funcs below is the top left corner (involving 3 vertices)
func (sp *ScrollablePanel) DrawScrollbarChrome(atlasCellX, atlasCellY, l, t float32) { // left, top
	span := UvSpan
	u := float32(atlasCellX) * span
	v := float32(atlasCellY) * span

	gl.Normal3f(0, 0, 1)

	// bottom left   0, 1
	gl.TexCoord2f(u, v+span)
	gl.Vertex3f(l, sp.Whole.Bottom, 0)

	// bottom right   1, 1
	gl.TexCoord2f(u+span, v+span)
	gl.Vertex3f(sp.Whole.Right, sp.Whole.Bottom, 0)

	// top right   1, 0
	gl.TexCoord2f(u+span, v)
	gl.Vertex3f(sp.Whole.Right, t, 0)

	// top left   0, 0
	gl.TexCoord2f(u, v)
	gl.Vertex3f(l, t, 0)
}

func (sp *ScrollablePanel) DrawBackground(atlasCellX, atlasCellY float32) {
	SetColor(GrayDark)
	DrawStretchableRect(atlasCellX, atlasCellY,
		&app.Rectangle{
			sp.Whole.Top,
			sp.Whole.Right - ui.ScrollBarThickness,
			sp.Whole.Bottom + ui.ScrollBarThickness,
			sp.Whole.Left})
}

func (sp *ScrollablePanel) ScrollIfMouseOver(mousePixelDeltaX, mousePixelDeltaY float32) {
	if sp.ContainsMouseCursor() {
		// position increments in gl space
		xInc := mousePixelDeltaX * PixelSize.X
		yInc := mousePixelDeltaY * PixelSize.Y
		sp.BarHori.Scroll(xInc)
		sp.BarVert.Scroll(yInc)
	}
}

func (sp *ScrollablePanel) ContainsMouseCursor() bool {
	return mouse.CursorIsInside(sp.Whole)
}

func (sp *ScrollablePanel) RemoveCharacter(fromUnderCursor bool) {
	txt := sp.TextBodies[0]

	if fromUnderCursor {
		if len(txt[sp.CursY]) > sp.CursX {
			txt[sp.CursY] = txt[sp.CursY][:sp.CursX] + txt[sp.CursY][sp.CursX+1:len(txt[sp.CursY])]
		}
	} else {
		if sp.CursX > 0 {
			txt[sp.CursY] = txt[sp.CursY][:sp.CursX-1] + txt[sp.CursY][sp.CursX:len(txt[sp.CursY])]
			sp.CursX--
		}
	}
}

func (sp *ScrollablePanel) DrawTree() {
	if len(sp.Trees) > 0 {
		// setup main rect
		span := float32(1.3)
		x := -span / 2
		y := sp.Whole.Top - 0.1
		r := &app.Rectangle{y, x + span, y - span, x}

		sp.drawNodeAndDescendants(r, 0)
	}
}

func (sp *ScrollablePanel) drawNodeAndDescendants(r *app.Rectangle, nodeId int) {
	//fmt.Println("drawNode(r *app.Rectangle)")
	nameBar := &app.Rectangle{r.Top, r.Right, r.Top - 0.2*r.Height(), r.Left}
	DrawStretchableRect(11, 13, r)
	SetColor(Blue)
	DrawStretchableRect(11, 13, nameBar)
	DrawTextInRect(sp.Trees[0].Nodes[nodeId].Text, nameBar)
	SetColor(White)

	cX := r.CenterX()
	rSp := r.Width() // rect span (height & width are the same)
	t := r.Bottom - rSp*0.5
	b := r.Bottom - rSp*1.5

	node := sp.Trees[0].Nodes[nodeId] // FIXME? .....
	// find sp.Trees[0].Nodes[i].....
	// ......(if we ever use multiple trees per panel)
	// ......(also update DrawTree to use range)

	if /* left child exists */ node.ChildIdL != math.MaxInt32 {
		x := cX - rSp*1.5
		sp.drawArrowAndChild(r, &app.Rectangle{t, x + rSp, b, x}, node.ChildIdL)
	}

	if /* right child exists */ node.ChildIdR != math.MaxInt32 {
		x := cX + rSp*0.5
		sp.drawArrowAndChild(r, &app.Rectangle{t, x + rSp, b, x}, node.ChildIdR)
	}
}

func (sp *ScrollablePanel) drawArrowAndChild(parent, child *app.Rectangle, childId int) {
	latExt := child.Width() * 0.15 // lateral extent of arrow's triangle top
	DrawTriangle(9, 1,
		app.Vec2F{parent.CenterX() - latExt, parent.Bottom},
		app.Vec2F{parent.CenterX() + latExt, parent.Bottom},
		app.Vec2F{child.CenterX(), child.Top})
	sp.drawNodeAndDescendants(child, childId)
}

func (sp *ScrollablePanel) SetupDemoProgram() {
	txt := []string{}

	txt = append(txt, "// ------- variable declarations ------- -------")
	//txt = append(txt, "var myVar int32")
	txt = append(txt, "var a int32 = 42 // end-of-line comment")
	txt = append(txt, "var b int32 = 58")
	txt = append(txt, "")
	txt = append(txt, "// ------- builtin function calls ------- ------- ------- ------- ------- ------- ------- end")
	txt = append(txt, "//    sub32(7, 9)")
	//txt = append(txt, "sub32(4,8)")
	//txt = append(txt, "mult32(7, 7)")
	//txt = append(txt, "mult32(3,5)")
	//txt = append(txt, "div32(8,2)")
	//txt = append(txt, "div32(15,  3)")
	//txt = append(txt, "add32(2,3)")
	//txt = append(txt, "add32(a, b)")
	txt = append(txt, "")
	txt = append(txt, "// ------- user function calls -------")
	txt = append(txt, "myFunc(a, b)")
	txt = append(txt, "")
	txt = append(txt, "// ------- function declarations -------")
	txt = append(txt, "func myFunc(a int32, b int32){")
	txt = append(txt, "")
	txt = append(txt, "        div32(6, 2)")
	txt = append(txt, "        innerFunc(a,b)")
	txt = append(txt, "}")
	txt = append(txt, "")
	txt = append(txt, "func innerFunc (a, b int32) {")
	txt = append(txt, "        var locA int32 = 71")
	txt = append(txt, "        var locB int32 = 29")
	txt = append(txt, "        sub32(locA, locB)")
	txt = append(txt, "}")

	/*
		for i := 0; i < 22; i++ {
			txt = append(txt, fmt.Sprintf("%d: put lots of text on screen", i))
		}
	*/

	sp.TextBodies[0] = txt
}

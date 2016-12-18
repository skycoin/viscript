package gfx

import (
	"fmt"
	"github.com/corpusc/viscript/app"
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
	Rect       *app.Rectangle
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
	fmt.Printf("TextPanel.SetSize()\n")

	sp.Rect = &app.Rectangle{
		Rend.ClientExtentY - Rend.CharHei,
		Rend.ClientExtentX,
		-Rend.ClientExtentY,
		-Rend.ClientExtentX}

	if sp.FractionOfStrip == runPanelHeiFrac { // FIXME: this is hardwired for one use case for now
		sp.Rect.Top = sp.Rect.Bottom + sp.Rect.Height()*sp.FractionOfStrip
	} else {
		sp.Rect.Bottom = sp.Rect.Bottom + sp.Rect.Height()*runPanelHeiFrac
	}

	// set scrollbars' upper left corners
	sp.BarHori.Rect.Left = sp.Rect.Left
	sp.BarHori.Rect.Top = sp.Rect.Bottom + ui.ScrollBarThickness
	sp.BarVert.Rect.Left = sp.Rect.Right - ui.ScrollBarThickness
	sp.BarVert.Rect.Top = sp.Rect.Top
}

func (sp *ScrollablePanel) RespondToMouseClick() {
	Rend.Focused = sp

	// diffs/deltas from home position of panel (top left corner)
	glDeltaXFromHome := Curs.MouseGlX - sp.Rect.Left
	glDeltaYFromHome := Curs.MouseGlY - sp.Rect.Top
	sp.MouseX = int((glDeltaXFromHome + sp.BarHori.ScrollDelta) / Rend.CharWid)
	sp.MouseY = int(-(glDeltaYFromHome + sp.BarVert.ScrollDelta) / Rend.CharHei)

	if sp.MouseY < 0 {
		sp.MouseY = 0
	}

	if sp.MouseY >= len(sp.TextBodies[0]) {
		sp.MouseY = len(sp.TextBodies[0]) - 1
	}
}

func (sp *ScrollablePanel) GoToTopEdge() {
	Rend.CurrY = sp.Rect.Top - sp.BarVert.ScrollDelta
}
func (sp *ScrollablePanel) GoToLeftEdge() float32 {
	Rend.CurrX = sp.Rect.Left - sp.BarHori.ScrollDelta
	return Rend.CurrX
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
	sp.DrawScrollbarChrome(10, 11, sp.Rect.Right-ui.ScrollBarThickness, sp.Rect.Top)                          // vertical bar background
	sp.DrawScrollbarChrome(13, 12, sp.Rect.Left, sp.Rect.Bottom+ui.ScrollBarThickness)                        // horizontal bar background
	sp.DrawScrollbarChrome(12, 11, sp.Rect.Right-ui.ScrollBarThickness, sp.Rect.Bottom+ui.ScrollBarThickness) // corner elbow piece
	SetColor(Gray)
	sp.BarHori.SetSize(sp.Rect, sp.TextBodies[0], Rend.CharWid, Rend.CharHei) // FIXME? (to consider multiple bodies & multiple trees)
	sp.BarVert.SetSize(sp.Rect, sp.TextBodies[0], Rend.CharWid, Rend.CharHei)
	Rend.DrawStretchableRect(11, 13, sp.BarHori.Rect) // 2,11 (pixel checkerboard)    // 14, 15 (square in the middle)
	Rend.DrawStretchableRect(11, 13, sp.BarVert.Rect) // 13, 12 (double horizontal lines)    // 10, 11 (double vertical lines)
	SetColor(White)
	sp.DrawTree()
}

func (sp *ScrollablePanel) DrawText() {
	var colorText = "" // when non-zero length, we're building this in place of drawing
	cX := Rend.CurrX   // current (internal/logic cursor) drawing position
	cY := Rend.CurrY
	cW := Rend.CharWid
	cH := Rend.CharHei
	b := sp.BarHori.Rect.Top // bottom of text area

	text := sp.TextBodies[0]
	// setup for colored text
	if /* it's been generated */ len(sp.TextBodies) > 1 {
		text = sp.TextBodies[1]
	} else { // otherwise generate, so demo program starts out colorized
		//script.Process(false)    FIXME: can't import "script", cuz dupes "gfx"
	}

	// iterate all runes
	for y, line := range text {
		// if line visible
		if cY <= sp.Rect.Top+cH && cY >= b {
			r := &app.Rectangle{cY, cX + cW, cY - cH, cX} // t, r, b, l

			// if line needs vertical adjustment
			if cY > sp.Rect.Top {
				r.Top = sp.Rect.Top
			}
			if cY-cH < b {
				r.Bottom = b
			}

			// process line of text
			SetColor(Gray)
			for x, c := range line {
				if /* ending color building */ c == '>' {
					SetColorFromText(colorText)
					colorText = ""
				} else if /* building color */ c == '<' || len(colorText) > 0 {
					colorText += string(c)
				} else { // drawing
					if /* visible */ cX >= sp.Rect.Left-cW && cX < sp.BarVert.Rect.Left {
						app.ClampLeftAndRightOf(r, sp.Rect.Left, sp.BarVert.Rect.Left)
						Rend.DrawCharAtRect(c, r)

						if sp.IsEditable { //&& Curs.Visible == true {
							if x == sp.CursX && y == sp.CursY {
								SetColor(White)
								//Rend.DrawCharAtRect('_', r)
								Rend.DrawStretchableRect(11, 13, Curs.GetAnimationModifiedRect(*r))
								SetColor(PrevColor)
							}
						}
					}

					cX += cW
					r.Left = cX
					r.Right = cX + cW
				}
			}

			// draw cursor at the end of line if needed
			if cX < sp.BarVert.Rect.Left && y == sp.CursY && sp.CursX == len(line) {
				if sp.IsEditable { //&& Curs.Visible == true {
					SetColor(White)
					app.ClampLeftAndRightOf(r, sp.Rect.Left, sp.BarVert.Rect.Left)
					//Rend.DrawCharAtRect('_', r)
					Rend.DrawStretchableRect(11, 13, Curs.GetAnimationModifiedRect(*r))
				}
			}

			cX = sp.GoToLeftEdge()
		}

		cY -= cH // go down a line height
	}
}

// ATM the only different between the 2 funcs below is the top left corner (involving 3 vertices)
func (sp *ScrollablePanel) DrawScrollbarChrome(atlasCellX, atlasCellY, l, t float32) { // left, top
	span := Rend.UvSpan
	u := float32(atlasCellX) * span
	v := float32(atlasCellY) * span

	gl.Normal3f(0, 0, 1)

	// bottom left   0, 1
	gl.TexCoord2f(u, v+span)
	gl.Vertex3f(l, sp.Rect.Bottom, 0)

	// bottom right   1, 1
	gl.TexCoord2f(u+span, v+span)
	gl.Vertex3f(sp.Rect.Right, sp.Rect.Bottom, 0)

	// top right   1, 0
	gl.TexCoord2f(u+span, v)
	gl.Vertex3f(sp.Rect.Right, t, 0)

	// top left   0, 0
	gl.TexCoord2f(u, v)
	gl.Vertex3f(l, t, 0)
}

func (sp *ScrollablePanel) DrawBackground(atlasCellX, atlasCellY float32) {
	SetColor(GrayDark)
	Rend.DrawStretchableRect(atlasCellX, atlasCellY,
		&app.Rectangle{
			sp.Rect.Top,
			sp.Rect.Right - ui.ScrollBarThickness,
			sp.Rect.Bottom + ui.ScrollBarThickness,
			sp.Rect.Left})
}

func (sp *ScrollablePanel) ScrollIfMouseOver(mousePixelDeltaX, mousePixelDeltaY float64) {
	if sp.ContainsMouseCursor() {
		// position increments in gl space
		xInc := float32(mousePixelDeltaX) * Rend.PixelWid
		yInc := float32(mousePixelDeltaY) * Rend.PixelHei
		sp.BarHori.Scroll(xInc)
		sp.BarVert.Scroll(yInc)
	}
}

func (sp *ScrollablePanel) ContainsMouseCursor() bool {
	return MouseCursorIsInside(sp.Rect)
}

func (sp *ScrollablePanel) ContainsMouseCursorInsideOfScrollBars() bool {
	return MouseCursorIsInside(&app.Rectangle{
		sp.Rect.Top, sp.Rect.Right - ui.ScrollBarThickness, sp.Rect.Bottom + ui.ScrollBarThickness, sp.Rect.Left})
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
		y := sp.Rect.Top - 0.1
		r := &app.Rectangle{y, x + span, y - span, x}

		sp.drawNodeAndDescendants(r, 0)
	}
}

func (sp *ScrollablePanel) drawNodeAndDescendants(r *app.Rectangle, nodeId int) {
	//fmt.Println("drawNode(r *app.Rectangle)")
	nameBar := &app.Rectangle{r.Top, r.Right, r.Top - 0.2*r.Height(), r.Left}
	Rend.DrawStretchableRect(11, 13, r)
	SetColor(Blue)
	Rend.DrawStretchableRect(11, 13, nameBar)
	Rend.DrawTextInRect(sp.Trees[0].Nodes[nodeId].Text, nameBar)
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
	Rend.DrawTriangle(9, 1,
		app.Vec2{parent.CenterX() - latExt, parent.Bottom},
		app.Vec2{parent.CenterX() + latExt, parent.Bottom},
		app.Vec2{child.CenterX(), child.Top})
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

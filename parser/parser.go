package parser

import (
	"fmt"
	"github.com/corpusc/viscript/app"
	"github.com/corpusc/viscript/gfx"
	"github.com/corpusc/viscript/tree"
	"github.com/corpusc/viscript/ui"
	"math"
	"regexp"
	"strconv"
)

// ISSUES?

/*
only looks for 1 expression per line

no attempt is made to get anything inside a function
which may come after the opening curly brace, on the same line

closing curly brace of function only recognized as a "}" amidst spaces
*/

// NOTES

/*

TODO:
* make sure names are unique within a partic scope
* allow // comments at any position

*/

var types = []string{"bool", "int32", "string"} // FIXME: should allow [] and [42] prefixes
var builtinFuncs = []string{"add32", "sub32", "mult32", "div32"}
var mainBlock = &CodeBlock{Name: "main"} // the root/entry/top/alpha level of the program
var currBlock = mainBlock

// REGEX (raw strings to avoid having to quote backslashes)
var declaredVar = regexp.MustCompile(`^( +)?var( +)?([a-zA-Z]\w*)( +)?int32(( +)?=( +)?([0-9]+))?$`)
var declFuncStart = regexp.MustCompile(`^func ([a-zA-Z]\w*)( +)?\((.*)\)( +)?\{$`)
var declFuncEnd = regexp.MustCompile(`^( +)?\}( +)?$`)
var calledFunc = regexp.MustCompile(`^( +)?([a-zA-Z]\w*)\(([0-9]+|[a-zA-Z]\w*),( +)?([0-9]+|[a-zA-Z]\w*)\)$`)
var comment = regexp.MustCompile(`^//.*`)

type VarBool struct {
	name  string
	value bool
}

type VarInt32 struct {
	name  string
	value int32
}

type VarString struct {
	name  string
	value string
}

type CodeBlock struct {
	Name        string
	VarBools    []*VarBool
	VarInt32s   []*VarInt32
	VarStrings  []*VarString
	SubBlocks   []*CodeBlock
	Expressions []string
	Parameters  []string // unused atm
}

func Parse() {
	// clear script
	mainBlock = &CodeBlock{Name: "main"}

	// clear the OS and graphical consoles
	gfx.Con.Lines = []string{}
	gfx.Rend.Panels[1].TextBodies[0] = []string{}

	gfx.MakeHighlyVisibleLogHeader(`PARSING`, 5)
	parseAll()

	// setup trees & expressions in new panel
	i := len(gfx.Rend.Panels)
	gfx.Rend.Panels = append(gfx.Rend.Panels, &gfx.ScrollablePanel{FractionOfStrip: 1})
	gfx.Rend.Panels[i].Init()
	gfx.Rend.Panels[i].Trees = append(gfx.Rend.Panels[i].Trees, &tree.Tree{PanelId: i})

	gfx.MakeHighlyVisibleLogHeader(`RUNNING`, 5)
	run(mainBlock)
}

func parseAll() {
	for i, line := range gfx.Rend.Panels[0].TextBodies[0] {
		ParseLine(i, line, false)
	}
}

func Draw() {
	// setup main rect
	span := float32(1.8)
	x := -span / 2
	y := ui.MainMenu.Rect.Bottom - 0.1
	r := &app.Rectangle{y, x + span, y - span, x}

	drawCodeBlock(mainBlock, r)
}

func drawCodeBlock(cb *CodeBlock, r *app.Rectangle) {
	nameLabel := &app.Rectangle{r.Top, r.Right, r.Top - 0.2*r.Height(), r.Left}
	gfx.Rend.DrawStretchableRect(11, 13, r)
	gfx.SetColor(gfx.Blue)
	gfx.Rend.DrawStretchableRect(11, 13, nameLabel)
	gfx.Rend.DrawTextInRect(cb.Name, nameLabel)
	gfx.SetColor(gfx.White)

	cX := r.CenterX()
	rW := r.Width() // rect width
	num := float32(len(cb.SubBlocks))
	rowW := num * rW
	latExt := rW * 0.15 // lateral extent of arrow's triangle top

	if num > 1 {
		rowW += (num - 1) * rW / 2
	}

	x := cX - rowW/2

	// subblock row
	t := r.Bottom - r.Height()/2
	b := r.Bottom - r.Height()*1.5

	for _, curr := range cb.SubBlocks {
		gfx.Rend.DrawTriangle(9, 1,
			app.Vec2{cX - latExt, r.Bottom},
			app.Vec2{cX + latExt, r.Bottom},
			app.Vec2{x + rW/2, t})
		drawCodeBlock(curr, &app.Rectangle{t, x + rW, b, x})

		x += r.Width() * 1.5
	}
}

func ParseLine(i int, line string, coloring bool) {
	switch {
	case declaredVar.MatchString(line):
		result := declaredVar.FindStringSubmatch(line)

		if coloring {
			gfx.SetColor(gfx.Violet)
		} else {
			var s = fmt.Sprintf("%d: var (%s) declared", i, result[3])
			//printIntsFrom(currBlock)

			if result[8] == "" {
				currBlock.VarInt32s = append(currBlock.VarInt32s, &VarInt32{result[3], 0})
			} else {
				value, err := strconv.Atoi(result[8])
				if err != nil {
					s = fmt.Sprintf("%s... BUT COULDN'T CONVERT ASSIGNMENT (%s) TO A NUMBER!", s, result[8])
				} else {
					currBlock.VarInt32s = append(currBlock.VarInt32s, &VarInt32{result[3], int32(value)})
					s = fmt.Sprintf("%s & assigned: %d", s, value)
				}
			}

			gfx.Con.Add(fmt.Sprintf("%s\n", s))
		}
	case declFuncStart.MatchString(line):
		result := declFuncStart.FindStringSubmatch(line)

		if coloring {
			gfx.SetColor(gfx.Fuschia)
		} else {
			gfx.Con.Add(fmt.Sprintf("%d: func (%s) declared, with params: %s\n", i, result[1], result[3]))

			if currBlock.Name == "main" {
				currBlock = &CodeBlock{Name: result[1]}
				mainBlock.SubBlocks = append(mainBlock.SubBlocks, currBlock) // FUTURE FIXME: methods in structs shouldn't be on main/root func
			} else {
				gfx.Con.Add("Func'y func-ception! CAN'T PUT A FUNC INSIDE A FUNC!\n")
			}
		}
	case declFuncEnd.MatchString(line):
		if coloring {
			gfx.SetColor(gfx.Fuschia)
		} else {
			gfx.Con.Add(fmt.Sprintf("func close...\n"))
			//printIntsFrom(mainBlock)
			//printIntsFrom(currBlock)

			if currBlock.Name == "main" {
				gfx.Con.Add(fmt.Sprintf("ERROR! Main\\Root level function doesn't need enclosure!\n"))
			} else {
				currBlock = mainBlock
			}
		}
	case calledFunc.MatchString(line): // FIXME: hardwired for 2 params each
		result := calledFunc.FindStringSubmatch(line)

		if coloring {
			gfx.SetColor(gfx.Fuschia)
		} else {
			gfx.Con.Add(fmt.Sprintf("%d: func call (%s) expressed\n", i, result[2]))
			gfx.Con.Add(fmt.Sprintf("currBlock: %s\n", currBlock))
			currBlock.Expressions = append(currBlock.Expressions, line)
			/*
				currBlock.Expressions = append(currBlock.Expressions, result[2])
				currBlock.Parameters = append(currBlock.Parameters, result[3])
				currBlock.Parameters = append(currBlock.Parameters, result[5])
			*/
			//printIntsFrom(currBlock)

			/*
				// prints out all captures
				for i, v := range result {
					gfx.Con.Add(fmt.Sprintf("%d. %s\n", i, v))
				}
			*/
		}
	case comment.MatchString(line): // allow "//" comments    FIXME to allow this at any later point in the line
		if coloring {
			gfx.SetColor(gfx.GrayDark)
		}
	case line == "":
		// just ignore
	default:
		if coloring {
			gfx.SetColor(gfx.White)
		} else {
			gfx.Con.Add(fmt.Sprintf("SYNTAX ERROR on line %d: \"%s\"\n", i, line))
		}
	}
}

func run(pb *CodeBlock) { // passed block of code
	gfx.Con.Add(fmt.Sprintf("running function: '%s'\n", pb.Name))

	for i, line := range pb.Expressions {
		gfx.Con.Add(fmt.Sprintf("------evaluating expression: '%s\n", line))

		switch {
		case calledFunc.MatchString(line): // FIXME: hardwired for 2 params each
			result := calledFunc.FindStringSubmatch(line)
			gfx.Con.Add(fmt.Sprintf("%d: calling func (%s) with params: %s, %s\n", i, result[2], result[3], result[5]))

			a := getInt32(result[3])
			if /* not legit num */ a == math.MaxInt32 {
				return
			}
			b := getInt32(result[5])
			if /* not legit num */ b == math.MaxInt32 {
				return
			}

			switch result[2] {
			case "add32":
				gfx.Con.Add(fmt.Sprintf("%d + %d = %d\n", a, b, a+b))
			case "sub32":
				gfx.Con.Add(fmt.Sprintf("%d - %d = %d\n", a, b, a-b))
			case "mult32":
				gfx.Con.Add(fmt.Sprintf("%d * %d = %d\n", a, b, a*b))
			case "div32":
				gfx.Con.Add(fmt.Sprintf("%d / %d = %d\n", a, b, a/b))
			default:
				for _, cb := range pb.SubBlocks {
					gfx.Con.Add((fmt.Sprintf("CodeBlock.Name considered: %s   switching on: %s\n", cb.Name, result[2])))

					if cb.Name == result[2] {
						gfx.Con.Add((fmt.Sprintf("'%s' matched '%s'\n", cb.Name, result[2])))
						run(cb)
					}
				}
			}
		}
	}
}

func getInt32(s string) int32 {
	value, err := strconv.Atoi(s)

	if err != nil {
		for _, v := range currBlock.VarInt32s {
			if s == v.name {
				return v.value
			}
		}

		if currBlock.Name != "main" {
			for _, v := range mainBlock.VarInt32s {
				if s == v.name {
					return v.value
				}
			}
		}

		gfx.Con.Add(fmt.Sprintf("ERROR!  '%s' IS NOT A VALID VARIABLE/FUNCTION!\n", s))
		return math.MaxInt32
	}

	return int32(value)
}

func printIntsFrom(f *CodeBlock) {
	if len(f.VarInt32s) == 0 {
		gfx.Con.Add(fmt.Sprintf("%s has no elements!\n", f.Name))
	} else {
		for i, v := range f.VarInt32s {
			gfx.Con.Add(fmt.Sprintf("%s.VarInt32s[%d]: %s = %d\n", f.Name, i, v.name, v.value))
		}
	}
}

/*
The FindAllStringSubmatch-function will, for each match, return an array with the
entire match in the first field and the
content of the groups in the remaining fields.
The arrays for all the matches are then captured in a container array.

the number of fields in the resulting array always matches the number of groups plus one.
*/

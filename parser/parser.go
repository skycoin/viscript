package parser

import (
	"fmt"
	"strings"
	//"github.com/corpusc/viscript/app"
	"github.com/corpusc/viscript/gfx"
	"github.com/corpusc/viscript/tree"
	//"github.com/corpusc/viscript/ui"
	"math"
	"regexp"
	"strconv"
)

/*

KNOWN ISSUES:

* only looks for 1 expression per line

* no attempt is made to get anything inside any block
	which may come after the opening curly brace, on the same line

* closing curly brace of a block only recognized as a "}" amidst spaces



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

func MakeTree() {
	// setup trees & expressions in new panel
	pI := len(gfx.Rend.Panels) // panel id

	// new panel
	gfx.Rend.Panels = append(gfx.Rend.Panels, &gfx.ScrollablePanel{FractionOfStrip: 1})
	gfx.Rend.Panels[pI].Init()

	// new tree
	gfx.Rend.Panels[pI].Trees = append(
		gfx.Rend.Panels[pI].Trees, &tree.Tree{pI, []*tree.Node{}})
	nI := 0 // node id
	makeNode(pI, math.MaxInt32, math.MaxInt32, math.MaxInt32, "top")
	makeNode(pI, math.MaxInt32, math.MaxInt32, nI, "1st left")
	makeNode(pI, math.MaxInt32, math.MaxInt32, nI, "1st right")
	nI++
	makeNode(pI, math.MaxInt32, math.MaxInt32, nI, "level 3")
	nI++
	makeNode(pI, math.MaxInt32, math.MaxInt32, nI, "level 4")
	nI++
	makeNode(pI, math.MaxInt32, math.MaxInt32, nI, "level 5")
}
func makeTree(tree tree.Tree) {
	addNode(mainBlock)
}

func makeNode(panelId, childIdL, childIdR, parentId int, s string) {
	gfx.Rend.Panels[panelId].Trees[0].Nodes = append(
		gfx.Rend.Panels[panelId].Trees[0].Nodes, &tree.Node{s, childIdL, childIdR, parentId})

	if parentId != math.MaxInt32 {
		// set pointer to child
		gfx.Rend.Panels[panelId].Trees[0].Nodes[parentId].ChildIdR = len(gfx.Rend.Panels[panelId].Trees[0].Nodes)
	}
}
func addNode(cb *CodeBlock) {
	for _, curr := range cb.SubBlocks {
		addNode(curr)
	}
}

func Parse() {
	// clear script
	mainBlock = &CodeBlock{Name: "main"}

	// clear the OS and graphical consoles
	gfx.Con.Lines = []string{}
	gfx.Rend.Panels[1].TextBodies[0] = []string{}

	gfx.MakeHighlyVisibleLogHeader(`PARSING`, 5)
	parseAll()

	gfx.MakeHighlyVisibleLogHeader(`RUNNING`, 5)
	run(mainBlock)
}

func parseAll() {
	bods := gfx.Rend.Panels[0].TextBodies
	// make a 2nd copy of code, but with inserted color markup
	bods = append(bods, []string{})

	for i, line := range bods[0] {
		processLine(i, line, false)
		bods[1] = append(bods[1], lexAndColorMarkupLine(line))
	}
}

func lexAndColorMarkupLine(line string) string {
	s := line // the dynamic/processed offshoot

	// strip any comments
	i := strings.Index(line, "//")
	if /* there is a comment */ i != -1 {
		s = line[:i]
		line = line[:i] + "<color=GrayDark>" + line[i: /*len(line)*/]
		fmt.Println("line:", line)
	}

	s = strings.TrimSpace(s)

	if len(s) > 0 {
		fmt.Println("s:", s)

		// tokenize
		lex := strings.Split(s, " ")

		for i := range lex {
			fmt.Printf("lex: %d '%s'\n", i, lex[i])
		}
	}

	return line
}

func processLine(i int, line string, coloring bool) {
	if coloring {
		switch {
		case declaredVar.MatchString(line):
			gfx.SetColor(gfx.Violet)
		case declFuncStart.MatchString(line):
			gfx.SetColor(gfx.Fuschia)
		case declFuncEnd.MatchString(line):
			gfx.SetColor(gfx.Fuschia)
		case calledFunc.MatchString(line): // FIXME: hardwired for 2 params each
			gfx.SetColor(gfx.Fuschia)
		case comment.MatchString(line): // allow "//" comments    FIXME to allow this at any later point in the line
			gfx.SetColor(gfx.GrayDark)
		case line == "":
			// just ignore
		default:
			gfx.SetColor(gfx.White)
		}
	} else {
		// scan for high level pieces
		switch {
		case declaredVar.MatchString(line):
			result := declaredVar.FindStringSubmatch(line)

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
		case declFuncStart.MatchString(line):
			result := declFuncStart.FindStringSubmatch(line)

			gfx.Con.Add(fmt.Sprintf("%d: func (%s) declared, with params: %s\n", i, result[1], result[3]))

			if currBlock.Name == "main" {
				currBlock = &CodeBlock{Name: result[1]}
				mainBlock.SubBlocks = append(mainBlock.SubBlocks, currBlock) // FUTURE FIXME: methods in structs shouldn't be on main/root func
			} else {
				gfx.Con.Add("Func'y func-ception! CAN'T PUT A FUNC INSIDE A FUNC!\n")
			}
		case declFuncEnd.MatchString(line):
			gfx.Con.Add(fmt.Sprintf("func close...\n"))
			//printIntsFrom(mainBlock)
			//printIntsFrom(currBlock)

			if currBlock.Name == "main" {
				gfx.Con.Add(fmt.Sprintf("ERROR! Main\\Root level function doesn't need enclosure!\n"))
			} else {
				currBlock = mainBlock
			}
		case calledFunc.MatchString(line): // FIXME: hardwired for 2 params each
			result := calledFunc.FindStringSubmatch(line)

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
		case line == "":
			// just ignore
		case comment.MatchString(line): // allow "//" comments    FIXME to allow this at any later point in the line
		default:
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

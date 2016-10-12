package main

import (
	"fmt"
	"math"
	"regexp"
	"strconv"
	//"regexp/syntax"
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
var mainFunc = &CodeBlock{Name: "main"} // the root/top/alpha/entry level of the program
var currFunc = mainFunc

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
	VarBools    []VarBool
	VarInt32s   []VarInt32
	VarStrings  []VarString
	CodeBlocks  []*CodeBlock
	Expressions []*string
	Parameters  []string // unused atm
}

func initParser() {
	/*
		for _, f := range funcs {
			con.Add(fmt.Sprintf(f))
		}
	*/
	makeHighlyVisibleRuntimeLogHeader(`PARSING`, 5)
	parse()
	makeHighlyVisibleRuntimeLogHeader("RUNNING", 5)
	run(mainFunc)
}

func parse() {
	for i, line := range textRend.Focused.Body {
		switch {
		case declaredVar.MatchString(line):
			result := declaredVar.FindStringSubmatch(line)
			var s = fmt.Sprintf("%d: var (%s) declared", i, result[3])
			//printIntsFrom(currFunc)

			if result[8] == "" {
				currFunc.VarInt32s = append(currFunc.VarInt32s, VarInt32{result[3], 0})
			} else {
				value, err := strconv.Atoi(result[8])
				if err != nil {
					s = fmt.Sprintf("%s... BUT COULDN'T CONVERT ASSIGNMENT (%s) TO A NUMBER!", s, result[8])
				} else {
					currFunc.VarInt32s = append(currFunc.VarInt32s, VarInt32{result[3], int32(value)})
					s = fmt.Sprintf("%s & assigned: %d", s, value)
				}
			}

			con.Add(fmt.Sprintf("%s\n", s))
			//printIntsFrom(currFunc)
		case declFuncStart.MatchString(line):
			result := declFuncStart.FindStringSubmatch(line)
			con.Add(fmt.Sprintf("%d: func (%s) declared, with params: %s\n", i, result[1], result[3]))

			if currFunc.Name == "main" {
				currFunc = &CodeBlock{Name: result[1]}
				mainFunc.CodeBlocks = append(mainFunc.CodeBlocks, currFunc) // FUTURE FIXME: methods in structs shouldn't be on main/root func
			} else {
				con.Add("Func'y func-ception! CAN'T PUT A FUNC INSIDE A FUNC!\n")
			}
		case declFuncEnd.MatchString(line):
			con.Add(fmt.Sprintf("func close...\n"))
			//printIntsFrom(mainFunc)
			//printIntsFrom(currFunc)

			if currFunc.Name == "main" {
				con.Add(fmt.Sprintf("ERROR! Main\\Root level function doesn't need enclosure!\n"))
			} else {
				currFunc = mainFunc
			}
		case calledFunc.MatchString(line): // FIXME: hardwired for 2 params each
			result := calledFunc.FindStringSubmatch(line)
			con.Add(fmt.Sprintf("%d: func call (%s) expressed\n", i, result[2]))
			con.Add(fmt.Sprintf("currFunc: %s\n", currFunc))
			currFunc.Expressions = append(currFunc.Expressions, &line)
			/*
				currFunc.Expressions = append(currFunc.Expressions, result[2])
				currFunc.Parameters = append(currFunc.Parameters, result[3])
				currFunc.Parameters = append(currFunc.Parameters, result[5])
			*/
			//printIntsFrom(currFunc)

			/*
				// prints out all captures
				for i, v := range result {
					con.Add(fmt.Sprintf("%d. %s\n", i, v))
				}
			*/
		case comment.MatchString(line): // allow "//" comments
			fallthrough
		case line == "":
			// just ignore
		default:
			con.Add(fmt.Sprintf("SYNTAX ERROR on line %d: \"%s\"\n", i, line))
		}
	}
}

func run(pb *CodeBlock) { // passed block of code
	con.Add(fmt.Sprintf("running function: '%s'\n", pb.Name))

	for i, line := range pb.Expressions {
		con.Add(fmt.Sprintf("running expression: '%s' in function: '%s'\n", *line, pb.Name))

		switch {
		case calledFunc.MatchString(*line): // FIXME: hardwired for 2 params each
			result := calledFunc.FindStringSubmatch(*line)
			con.Add(fmt.Sprintf("%d: calling func (%s) with params: %s, %s\n", i, result[2], result[3], result[5]))

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
				con.Add(fmt.Sprintf("%d + %d = %d\n", a, b, a+b))
			case "sub32":
				con.Add(fmt.Sprintf("%d - %d = %d\n", a, b, a-b))
			case "mult32":
				con.Add(fmt.Sprintf("%d * %d = %d\n", a, b, a*b))
			case "div32":
				con.Add(fmt.Sprintf("%d / %d = %d\n", a, b, a/b))
			default:
				for _, fun := range pb.CodeBlocks {
					//con.Add((fmt.Sprintf("user func: %s == %s\n", fun.Name, result[2])))

					if fun.Name == result[2] {
						con.Add((fmt.Sprintf("'%s' matched '%s'\n", fun.Name, result[2])))
						run(fun)
					}
				}
			}
		}
	}
}

func getInt32(s string) int32 {
	value, err := strconv.Atoi(s)

	if err != nil {
		for _, v := range currFunc.VarInt32s {
			if s == v.name {
				return v.value
			}
		}

		if currFunc.Name != "main" {
			for _, v := range mainFunc.VarInt32s {
				if s == v.name {
					return v.value
				}
			}
		}

		con.Add(fmt.Sprintf("ERROR!  '%s' IS NOT A VALID VARIABLE/FUNCTION!\n", s))
		return math.MaxInt32
	}

	return int32(value)
}

func printIntsFrom(f *CodeBlock) {
	if len(f.VarInt32s) == 0 {
		con.Add(fmt.Sprintf("%s has no elements!\n", f.Name))
	} else {
		for i, v := range f.VarInt32s {
			con.Add(fmt.Sprintf("%s.VarInt32s[%d]: %s = %d\n", f.Name, i, v.name, v.value))
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

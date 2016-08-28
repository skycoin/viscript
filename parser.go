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
make sure names are unique within a partic scope
*/

var types = []string{"bool", "int32", "string"} // FIXME: should allow [] and [42] prefixes
var rootFunc = &Func{Name: "root"}              // the top/alpha/entry level of the program
var currFunc = rootFunc

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

type Func struct {
	Name      string
	VarBools  []VarBool
	VarInt32s []VarInt32
	VarString []VarString
	Funcs     []Func
	// next is just temporary?
	Parameters []string
}

func initParser() {
	/*
		for _, f := range funcs {
			con.Add(fmt.Sprintf(f))
		}
	*/
	parse()
}

func parse() {
	// Use raw strings to avoid having to quote the backslashes.
	var varInt32 = regexp.MustCompile(`^( +)?var( +)?([a-zA-Z]\w*)( +)?int32(( +)?=( +)?([0-9]+))?$`)
	var declFuncStart = regexp.MustCompile(`^func ([a-zA-Z]\w*)( +)?\((.*)\)( +)?\{$`)
	var declFuncEnd = regexp.MustCompile(`^( +)?\}( +)?$`)
	var funcCall = regexp.MustCompile(`^( +)?(add32|sub32|mult32|div32)\(([0-9]+|[a-zA-Z]\w*),( +)?([0-9]+|[a-zA-Z]\w*)\)$`)

	for i, line := range code.Body {
		switch {
		case varInt32.MatchString(line):
			result := varInt32.FindStringSubmatch(line)
			con.Add(fmt.Sprintf("%d: variable (%s) declaration\n", i, result[3]))
			printIntsFrom(currFunc)

			if result[8] == "" {
				currFunc.VarInt32s = append(currFunc.VarInt32s, VarInt32{result[3], 0})
			} else {
				value, err := strconv.Atoi(result[8])
				if err != nil {
					con.Add(fmt.Sprintf("COULDN'T CONVERT ASSIGNMENT TO NUMBER!"))
				} else {
					currFunc.VarInt32s = append(currFunc.VarInt32s, VarInt32{result[3], int32(value)})
					con.Add(fmt.Sprintf("....assigned value of: %d\n", value))
				}
			}

			printIntsFrom(currFunc)
		case declFuncStart.MatchString(line):
			result := declFuncStart.FindStringSubmatch(line)
			con.Add(fmt.Sprintf("%d: function (%s) declaration, with parameters: %s\n", i, result[1], result[3]))

			if currFunc.Name == "root" {
				currFunc = &Func{Name: result[1]}
				rootFunc.Funcs = append(rootFunc.Funcs, *currFunc) // FIXME: methods in structs shouldn't be on root
			} else {
				con.Add("Func'y func-ception! CAN'T PUT A FUNC INSIDE A FUNC!\n")
			}
		case declFuncEnd.MatchString(line):
			con.Add(fmt.Sprintf("function close...\n"))
			printIntsFrom(rootFunc)
			printIntsFrom(currFunc)

			if currFunc.Name == "root" {
				con.Add(fmt.Sprintf("ERROR! Root level function doesn't need enclosure!\n"))
			} else {
				currFunc = rootFunc
			}
		case funcCall.MatchString(line): // FIXME: hardwired for 4 undefined math functions, with 2 params each
			result := funcCall.FindStringSubmatch(line)
			con.Add(fmt.Sprintf("%d: function call\n", i))
			printIntsFrom(currFunc)

			/*
				// prints out all captures
				for i, v := range result {
					con.Add(fmt.Sprintf("%d. %s\n", i, v))
				}
			*/

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
			}
		case line == "":
			// just ignore
		default:
			con.Add(fmt.Sprintf("SYNTAX ERROR on line %d: %s\n", i, line))
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

		if currFunc.Name != "root" {
			for _, v := range rootFunc.VarInt32s {
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

func printIntsFrom(f *Func) {
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

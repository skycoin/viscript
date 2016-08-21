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
which may come after the opening curly brace, but on the same line

closing curly brace of function only recognized as a 1 character line
*/

// NOTES

/*
depending on whether go allows funcs inside funcs.....
assuming no: what should probably be done is that we have a var that represents what
function we are currently adding to (or none of course).  upon detecting the start
of a function we will start funneling the following expressions into that function.
then pop out when encountering a closing curly brace
*/

var types = []string{"bool", "int32", "string"} // FIXME: should allow [] and [42] prefixes
var varBools = make([]VarBool, 0)
var varInt32s = make([]VarInt32, 0)
var varString = make([]VarString, 0)
var funcs = make([]Func, 0)
var currFunc *Func = nil // expressions go in here, if not nil

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
	Name        string
	Parameters  []string
	Expressions []string
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
	var funcCall = regexp.MustCompile(`^( +)?(add32|sub32|mult32|div32)\(([0-9]+|[a-zA-Z]\w*),( +)?([0-9]+|[a-zA-Z]\w*)\)$`)

	for i, line := range document {
		switch {
		case varInt32.MatchString(line):
			result := varInt32.FindStringSubmatch(line)
			con.Add(fmt.Sprintf("%d: variable (%s) declaration\n", i, result[3]))

			if result[8] == "" {
				varInt32s = append(varInt32s, VarInt32{result[3], 0})
			} else {
				value, err := strconv.Atoi(result[8])
				if err != nil {
					con.Add(fmt.Sprintf("COULDN'T CONVERT ASSIGNMENT TO NUMBER!"))
				}

				con.Add(fmt.Sprintf("....assigned value of: %d\n", value))
				varInt32s = append(varInt32s, VarInt32{result[3], int32(value)})
			}
		case declFuncStart.MatchString(line):
			if currFunc != nil {
				con.Add("CAN'T PUT FUNC INSIDE A FUNC!")
			}

			result := declFuncStart.FindStringSubmatch(line)
			con.Add(fmt.Sprintf("%d: function (%s) declaration, with parameters: %s\n", i, result[1], result[3]))
			funcs = append(funcs, Func{Name: result[1]})
		case funcCall.MatchString(line):
			con.Add(fmt.Sprintf("%d: function call\n", i))
			result := funcCall.FindStringSubmatch(line)

			/*
				// prints out all captures
				for i, v := range result {
					con.Add(fmt.Sprintf("%d. %s\n", i, v))
				}
			*/

			a := getUInt32(result[3])
			if /* not legit num */ a == math.MaxInt32 {
				return
			}
			b := getUInt32(result[5])
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
		default:
			con.Add(fmt.Sprintf("SYNTAX ERROR on line %d: %s\n", i, line))
		}
	}
}

func getUInt32(s string) int32 {
	value, err := strconv.Atoi(s)

	if err != nil {
		for _, v := range varInt32s {
			if s == v.name {
				return v.value
			}
		}

		con.Add(fmt.Sprintf("ERROR!  '%s' IS NOT A VALID VARIABLE/FUNCTION!", s))
		return math.MaxInt32
	}

	return int32(value)
}

/*
The FindAllStringSubmatch-function will, for each match, return an array with the
entire match in the first field and the
content of the groups in the remaining fields.
The arrays for all the matches are then captured in a container array.

the number of fields in the resulting array always matches the number of groups plus one.
*/

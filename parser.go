package main

import (
	"fmt"
	"regexp"
	"strconv"
	//"regexp/syntax"
)

var vars = make([]string, 0)
var funcs = make([]string, 0)

func initParser() {
	funcs = append(funcs, "add32")
	funcs = append(funcs, "sub32")
	funcs = append(funcs, "mult32")
	funcs = append(funcs, "div32")

	for _, f := range funcs {
		fmt.Println(f)
	}

	parse()
}

func parse() {
	// Use raw strings to avoid having to quote the backslashes.
	var varInt32 = regexp.MustCompile(`^( +)?var( +)?([a-zA-Z]\w*)( +)?int32$`)
	var funcCall = regexp.MustCompile(`^( +)?(add32|sub32|mult32|div32)\(([0-9]+),( +)?([0-9]+)\)$`)

	for i, line := range document {
		switch {
		case varInt32.MatchString(line):
			fmt.Printf("%d: variable declaration\n", i)
			vars = append(vars)
		case funcCall.MatchString(line):
			fmt.Printf("%d: function call\n", i)
			result := funcCall.FindStringSubmatch(line)

			/*for i, v := range result {
				fmt.Printf("%d. %s\n", i, v)
			}*/

			a, err := strconv.Atoi(result[3])
			if err != nil {
				fmt.Println("ZOMG!  COULDN'T CONVERT THAT NUMBER!")
			}

			b, err := strconv.Atoi(result[5])
			if err != nil {
				fmt.Println("ZOMG!  COULDN'T CONVERT THAT NUMBER!")
			}

			switch result[2] {
			case "add32":
				fmt.Printf("%d + %d = %d\n", a, b, a+b)
			case "sub32":
				fmt.Printf("%d - %d = %d\n", a, b, a-b)
			case "mult32":
				fmt.Printf("%d * %d = %d\n", a, b, a*b)
			case "div32":
				fmt.Printf("%d / %d = %d\n", a, b, a/b)
			}
		default:
			fmt.Printf("SYNTAX ERROR on line %d: %s\n", i, line)
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

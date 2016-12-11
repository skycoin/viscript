package parser

import (
/*
	"fmt"
	"github.com/corpusc/viscript/app"
	"github.com/corpusc/viscript/gfx"
	"github.com/corpusc/viscript/tree"
	"github.com/corpusc/viscript/ui"
	"math"
	"regexp"
	"strconv"
*/
)

type Token struct {
	LexType int
	value   string
}

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

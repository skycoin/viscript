package main

import (
//"fmt"
)

var code TextPanel = TextPanel{NumCharsX: 80, NumCharsY: 14}
var cons TextPanel = TextPanel{NumCharsX: 80, NumCharsY: 10} // console

type TextPanel struct {
	NumCharsX       int
	NumCharsY       int
	OffsetY         float32
	LenOfOffscreenY float32
}

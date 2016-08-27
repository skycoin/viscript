package main

import (
//"fmt"
)

var code TextPanel = TextPanel{}
var cons TextPanel = TextPanel{} // console

type TextPanel struct {
	OffsetY         float32
	LenOfOffscreenY float32
}

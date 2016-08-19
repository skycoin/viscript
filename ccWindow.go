package main

import (
//"fmt"
)

var view CcViewport = CcViewport{}

type CcViewport struct {
	OffsetY         float32
	LenOfOffscreenY float32
}

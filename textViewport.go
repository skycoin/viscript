package main

import (
//"fmt"
)

var view TextViewport = TextViewport{}

type TextViewport struct {
	OffsetY         float32
	LenOfOffscreenY float32
}

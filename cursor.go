package main

import (
	//"fmt"
	"time"
)

var curs Cursor = Cursor{time.Now(), true, 0, 0}

type Cursor struct {
	NextBlinkChange time.Time
	Visible         bool
	X               int // position at the character level
	Y               int
}

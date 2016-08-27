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

func (c *Cursor) Draw() {
	if c.NextBlinkChange.Before(time.Now()) {
		c.NextBlinkChange = time.Now().Add(time.Millisecond * 170)
		c.Visible = !c.Visible
	}

	if c.Visible == true {
		drawCharAt('_', c.X, c.Y)
	}
}

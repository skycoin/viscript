package gfx

import (
	//"fmt"
	"github.com/corpusc/viscript/app"
	"time"
)

var Curs Cursors = Cursors{NextBlinkChange: time.Now(), Visible: true}

type Cursors struct {
	NextBlinkChange time.Time
	Visible         bool

	// private
	shrinking      bool
	shrinkFraction float32
}

func (c *Cursors) Update() {
	if c.NextBlinkChange.Before(time.Now()) {
		c.NextBlinkChange = time.Now().Add(time.Millisecond * 16) // 170 for simple on/off blinking
		c.Visible = !c.Visible

		if c.shrinking {
			c.shrinkFraction -= PixelSize.Y * 6

			if c.shrinkFraction < 0.2 {
				c.shrinking = false
			}
		} else {
			c.shrinkFraction += PixelSize.Y * 6

			if c.shrinkFraction > 0.8 {
				c.shrinking = true
			}
		}
	}
}

func (c *Cursors) GetAnimationModifiedRect(r app.Rectangle) *app.Rectangle {
	if c.shrinking {
		r.Bottom = r.Top - c.shrinkFraction*r.Height()
		r.Left = r.Right - c.shrinkFraction*r.Width()
	} else {
		r.Top = r.Bottom + c.shrinkFraction*r.Height()
		r.Right = r.Left + c.shrinkFraction*r.Width()
	}

	return &r
}

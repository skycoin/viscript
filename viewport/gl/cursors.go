package gl

import (
	"github.com/corpusc/viscript/app"
	"time"
)

var Curs Cursors = Cursors{NextFrame: time.Now()}

type Cursors struct {
	NextFrame time.Time
	// private
	shrinking      bool
	shrinkFraction float32
}

func (c *Cursors) Update() {
	var speedFactor float32 = 0.06

	if c.NextFrame.Before(time.Now()) {
		c.NextFrame = time.Now().Add(time.Millisecond * 16) // 170 for simple on/off blinking

		if c.shrinking {
			c.shrinkFraction -= speedFactor

			if c.shrinkFraction < 0.2 {
				c.shrinking = false
			}
		} else {
			c.shrinkFraction += speedFactor

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
	} else { // growing
		r.Top = r.Bottom + c.shrinkFraction*r.Height()
		r.Right = r.Left + c.shrinkFraction*r.Width()
	}

	return &r
}

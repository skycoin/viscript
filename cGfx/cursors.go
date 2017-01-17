package cGfx

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

func (c *Cursors) GetAnimationModifiedRect(r app.PicRectangle) *app.PicRectangle {
	if c.shrinking {
		r.Rect.Bottom = r.Rect.Top - c.shrinkFraction*r.Rect.Height()
		r.Rect.Left = r.Rect.Right - c.shrinkFraction*r.Rect.Width()
	} else {
		r.Rect.Top = r.Rect.Bottom + c.shrinkFraction*r.Rect.Height()
		r.Rect.Right = r.Rect.Left + c.shrinkFraction*r.Rect.Width()
	}

	return &r
}

package gfx

import (
	//"fmt"
	"github.com/corpusc/viscript/app"
	"time"
)

var Curs Cursors = Cursors{NextFrame: time.Now()}

var (
	//CanvasExtents   app.Vec2F
	PixelSize app.Vec2F
)

type Cursors struct {
	NextFrame time.Time
	// private
	shrinking      bool
	shrinkFraction float32
}

func (c *Cursors) Update() {
	if c.NextFrame.Before(time.Now()) {
		c.NextFrame = time.Now().Add(time.Millisecond * 16) // 170 for simple on/off blinking

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
		r.Bottom = r.Top - c.shrinkFraction*r.Height()
		r.Left = r.Right - c.shrinkFraction*r.Width()
	} else {
		r.Top = r.Bottom + c.shrinkFraction*r.Height()
		r.Right = r.Left + c.shrinkFraction*r.Width()
	}

	return &r
}

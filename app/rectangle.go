package app

type Rectangle struct {
	Top    float32
	Right  float32
	Bottom float32
	Left   float32
}

func (r *Rectangle) Width() float32 {
	return r.Right - r.Left
}

func (r *Rectangle) Height() float32 {
	return r.Top - r.Bottom
}

func (r *Rectangle) CenterX() float32 {
	return r.Left + r.Width()/2
}

func (r *Rectangle) CenterY() float32 {
	return r.Bottom + r.Height()/2
}

func (r *Rectangle) Contains(x, y float32) bool {
	if x > r.Left && x < r.Right && y > r.Bottom && y < r.Top {
		return true
	}

	return false
}

// ----------------------------------
// rectangles with extra graphic data
// ----------------------------------

// rectangle types
const (
	RectType_9Slice = iota
	RectType_Simple // only requires one quad which is uniformly shrunk or stretched
)

type PicRectangle struct {
	Type     uint8
	AtlasPos Vec2I // x, y position in the atlas
	Rect     *Rectangle
}

package app

type IRect interface {
	/*Top    float32
	Right  float32
	Bottom float32
	Left   float32*/
	Width() float32
	Height() float32
	CenterX() float32
	CenterY() float32
	Contains() bool
}

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

// state
const (
	RectState_Active   = iota
	RectState_Inactive // .../invisible
	RectState_Dead
)

type PicRectangle struct {
	Id    int32
	Type  uint8
	State uint8
	//Color    float32
	AtlasPos Vec2I // x, y position in the atlas
	Rect     *Rectangle
}

type NineSliceRectangle struct { // UNUSED ATM
	X [3]float32 // positions of vertical lines
	Y [3]float32 // positions of horizontal lines
	U [3]float32 // texture coords of vertical lines
	V [3]float32 // texture coords of horizontal lines
}

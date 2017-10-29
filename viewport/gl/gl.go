package gl

import (
	"image"
	"image/color"
	"log"

	"github.com/go-gl/gl/v2.1/gl"
	"github.com/go-gl/glfw/v3.2/glfw"
	"github.com/skycoin/viscript/app"
)

var (
	GlfwWindow *glfw.Window //deprecate eventually
	GlfwCursor *glfw.Cursor
	Texture    uint32

	//private
	img            *image.RGBA
	goldenRatio    = 1.61803398875
	goldenFraction = float32(goldenRatio / (goldenRatio + 1))
)

//only two gfx parameters should be eliminated
//settings in either app or gfx

func InitGlfw() {
	println("<gl>.InitGlfw()")

	if err := glfw.Init(); err != nil {
		log.Fatalln("failed to initialize glfw:", err)
	}

	//defer glfw.Terminate()
	/*
	   Go's defer statement schedules a function call (the deferred function)
	   to be run immediately before the function executing the defer returns.
	   It's an unusual but effective way to deal with situations such as resources
	   that must be released regardless of which path a function takes to return.
	   The canonical examples are unlocking a mutex or closing a file.
	*/

	glfw.WindowHint(glfw.Resizable, glfw.True)
	glfw.WindowHint(glfw.ContextVersionMajor, 2)
	glfw.WindowHint(glfw.ContextVersionMinor, 1)

	var err error
	GlfwWindow, err = glfw.CreateWindow(InitAppWidth, InitAppHeight, app.Name, nil, nil)

	if err != nil {
		panic(err)
	}

	GlfwWindow.MakeContextCurrent()
	if err := gl.Init(); err != nil {
		panic(err)
	}

	//initialize cursor
	CreateAndSetPointer(glfw.ArrowCursor)
}

func CreateAndSetPointer(cursorType glfw.StandardCursor) {
	//FIXME? is this constantly leaking memory because it's
	//not destroying here everytime?  i think its rarely nil
	//(also CreateAndSetCustomPointer() below)
	if GlfwCursor != nil {
		GlfwCursor.Destroy()
	}

	GlfwCursor = glfw.CreateStandardCursor(cursorType)
	GlfwWindow.SetCursor(GlfwCursor)
}

func CreateAndSetCustomPointer() {
	if GlfwCursor != nil {
		GlfwCursor.Destroy()
	}

	max := 16

	if img == nil {
		img = image.NewRGBA(image.Rect(0, 0, max, max))
		//going up & right
		//(top left arrow head)
		plotLine(0, 2, 1, -1, max, color.White)
		plotLine(0, 3, 1, -1, max, color.White)
		plotLine(0, 4, 1, -1, max, color.White)
		plotLine(0, 5, 1, -1, max, color.White)
		plotLine(0, 6, 1, -1, max, color.Black)
		//going down & right
		//(main shaft of arrow)
		plotLine(0, 0, 1, 1, max, color.White)
		plotLine(1, 0, 1, 1, max, color.White)
		plotLine(0, 1, 1, 1, max, color.White)
		plotLine(5, 3, 1, 1, max, color.Black)
		plotLine(3, 5, 1, 1, max, color.Black)
		//going up & right
		//(bottom right arrow head)
		plotLine(max-3, max-1, 1, -1, max, color.White)
		plotLine(max-4, max-1, 1, -1, max, color.White)
		plotLine(max-5, max-1, 1, -1, max, color.White)
		plotLine(max-6, max-1, 1, -1, max, color.White)
		plotLine(max-7, max-1, 1, -1, max, color.Black)
	}

	//img.At(x, y) = color.Black

	GlfwCursor = glfw.CreateCursor(img, max/2, max/2)
	GlfwWindow.SetCursor(GlfwCursor)
}

func SetArrowPointer() {
	CreateAndSetPointer(glfw.ArrowCursor)
}

func SetHResizePointer() {
	CreateAndSetPointer(glfw.HResizeCursor)
}

func SetVResizePointer() {
	CreateAndSetPointer(glfw.VResizeCursor)
}

func SetCornerResizePointer() {
	CreateAndSetCustomPointer()
}

func SetIBeamPointer() {
	CreateAndSetPointer(glfw.IBeamCursor)
}

func SetHandPointer() {
	CreateAndSetPointer(glfw.HandCursor)
}

func LoadTextures() {
	Texture = NewTexture("assets/Bisasam_24x24_Shadowed.png")
}

func InitRenderer() {
	gl.Enable(gl.DEPTH_TEST)
	gl.Enable(gl.LIGHTING)
	//gl.Enable(gl.ALPHA_TEST)

	gl.ClearColor(0.5, 0.5, 0.5, 0.0)
	gl.ClearDepth(1)
	gl.DepthFunc(gl.LEQUAL)

	ambient := []float32{0.5, 0.5, 0.5, 1}
	diffuse := []float32{1, 1, 1, 1}
	lightPosition := []float32{-5, 5, 10, 0}
	gl.Lightfv(gl.LIGHT0, gl.AMBIENT, &ambient[0])
	gl.Lightfv(gl.LIGHT0, gl.DIFFUSE, &diffuse[0])
	gl.Lightfv(gl.LIGHT0, gl.POSITION, &lightPosition[0])
	gl.Enable(gl.LIGHT0)

	gl.MatrixMode(gl.PROJECTION)
	gl.LoadIdentity()
	//gl.Frustum(-1, 1, -1, 1, 1.0, 10.0)
	//gl.Frustum(left, right, bottom, top, zNear, zFar)
	SetOrtho(InitFrustum)
	//SetFrustum(InitFrustum)
	gl.MatrixMode(gl.MODELVIEW)
	gl.LoadIdentity()
}

func SetOrtho(r *app.Rectangle) {
	gl.Ortho( //gl.Frustum(
		float64(r.Left),
		float64(r.Right),
		float64(r.Bottom),
		float64(r.Top),
		-10.0, 10.0) // zNear, zFar
}

func SetFrustum(r *app.Rectangle) {
	gl.Frustum(
		float64(r.Left),
		float64(r.Right),
		float64(r.Bottom),
		float64(r.Top),
		-1.0, 10.0) // zNear, zFar
}

func DrawBegin() {
	//gl.Viewport(0, 0, CurrAppWidth, CurrAppHeight) // OPTIMIZEME?  could set flag upon frame buffer size change event

	//retina displays have larger frame buffers. we can't guess but just
	//get it using window handle.
	//darwin frame buffer width and darwin frame buffer height

	//FIXME?: should this change on framebuffer size change?
	w, h := //width, height
		GlfwWindow.GetFramebufferSize() //println("Frame BUFFER IN DRAW BEGIN: ", w, h)
	gl.Viewport(0, 0, int32(w), int32(h))

	if *PrevFrustum != *CurrFrustum {
		*PrevFrustum = *CurrFrustum
		gl.MatrixMode(gl.PROJECTION)
		gl.LoadIdentity()
		//SetFrustum(CurrFrustum)
		SetOrtho(CurrFrustum)
		println("CHANGE OF FRUSTUM")
	}

	gl.MatrixMode(gl.MODELVIEW) //.PROJECTION) //.MODELVIEW)
	gl.LoadIdentity()
	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
	gl.Translatef(0, 0, -DistanceFromOrigin)

	gl.BindTexture(gl.TEXTURE_2D, Texture)

	gl.Begin(gl.QUADS)
	drawDesktop()
}

func DrawEnd() {
	gl.End()
}

func SwapDrawBuffer() {
	GlfwWindow.SwapBuffers()
}

func ScreenTeardown() {
	glfw.Terminate()
}

//
//
//private
//
//

func plotLine(xStart, yStart, xDelta, yDelta, max int, bulkColor color.Color) {
	//bulk of the color will be as specified, but pixels along the edge will be black

	if xDelta == 0 &&
		yDelta == 0 {

		s := "ERROR!!! Position change deltas can't both be 0 (would cause endless loop)!"
		println(s)
		println(s)
		println(s)
		println(s)
		return
	}

	x := xStart
	y := yStart

	for x >= 0 && x < max &&
		y >= 0 && y < max {

		finalColor := bulkColor

		if x == 0 ||
			y == 0 ||
			x == max-1 ||
			y == max-1 {

			finalColor = color.Black
		}

		img.Set(x, y, finalColor) //.RGBA{0, 0, 0, 55}
		x += xDelta
		y += yDelta
	}
}

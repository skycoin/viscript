package gl

import (
	"fmt"
	"github.com/corpusc/viscript/app"
	"github.com/corpusc/viscript/cGfx"
	"github.com/go-gl/gl/v2.1/gl"
	"github.com/go-gl/glfw/v3.2/glfw"
	"log"
)

var goldenRatio = 1.61803398875
var goldenFraction = float32(goldenRatio / (goldenRatio + 1))

var (
	GlfwWindow  *glfw.Window // deprecate eventually
	CloseWindow chan int     // write to channel to close
	Texture     uint32
)

//gfx in cGfx.CurrAppWidth
//cGfx.InitFrustum

//only two gfx parameters should be eliminated
//gfx import should be eliminated
//settings in either app or gfx

func init() {
	CloseWindow = make(chan int) //write, to attempt to close out
}

func WindowInit() {
	fmt.Printf("Gl: Init glfw \n")

	if err := glfw.Init(); err != nil {
		log.Fatalln("failed to initialize glfw:", err)
	}

	//defer glfw.Terminate()

	fmt.Printf("Gl: set windowhint\n")
	glfw.WindowHint(glfw.Resizable, glfw.True)
	glfw.WindowHint(glfw.ContextVersionMajor, 2)
	glfw.WindowHint(glfw.ContextVersionMinor, 1)

	var err error
	GlfwWindow, err = glfw.CreateWindow(int(cGfx.CurrAppWidth), int(cGfx.CurrAppHeight), app.Name, nil, nil)

	if err != nil {
		panic(err)
	}

	GlfwWindow.MakeContextCurrent()
	if err := gl.Init(); err != nil {
		panic(err)
	}

}

func LoadTextures() {
	fmt.Printf("GL: load texture \n")
	Texture = NewTexture("Bisasam_24x24_Shadowed.png")
}

func InitRenderer() {
	fmt.Println("InitRenderer()")

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
	SetFrustum(cGfx.InitFrustum)
	//gl.Frustum(-1, 1, -1, 1, 1.0, 10.0)
	//gl.Frustum(left, right, bottom, top, zNear, zFar)
	gl.MatrixMode(gl.MODELVIEW)
	gl.LoadIdentity()
}

func SetFrustum(r *app.Rectangle) {
	gl.Frustum(
		float64(r.Left),
		float64(r.Right),
		float64(r.Bottom),
		float64(r.Top), 1.0, 10.0)
}

func DrawScene() {
	gl.Viewport(0, 0, cGfx.CurrAppWidth, cGfx.CurrAppHeight) // OPTIMIZEME?  could set flag upon frame buffer size change event
	if *cGfx.PrevFrustum != *cGfx.CurrFrustum {
		*cGfx.PrevFrustum = *cGfx.CurrFrustum
		gl.MatrixMode(gl.PROJECTION)
		gl.LoadIdentity()
		SetFrustum(cGfx.CurrFrustum)
		fmt.Println("CHANGE OF FRUSTUM")
	}
	gl.MatrixMode(gl.MODELVIEW) //.PROJECTION)                   //.MODELVIEW)
	gl.LoadIdentity()
	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
	gl.Translatef(0, 0, -cGfx.DistanceFromOrigin)

	gl.BindTexture(gl.TEXTURE_2D, Texture)

	gl.Begin(gl.QUADS)
	drawAll()
	gl.End()
}

func SwapDrawBuffer() {
	GlfwWindow.SwapBuffers()
}

func ScreenTeardown() {
	glfw.Terminate()
}

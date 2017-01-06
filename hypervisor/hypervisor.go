package hypervisor

import (
	"fmt"
	"github.com/corpusc/viscript/app"
	"github.com/corpusc/viscript/gfx"
	//"github.com/corpusc/viscript/msg"
	"github.com/corpusc/viscript/script"
	"github.com/go-gl/gl/v2.1/gl"
	"github.com/go-gl/glfw/v3.2/glfw"
	"log"
	"runtime"
)

var GlfwWindow *glfw.Window
var CloseWindow chan int // write to channel to close

func init() {
	CloseWindow = make(chan int)
}

func HypervisorInit() {
	fmt.Println("hypervisor.init()")
	gfx.MakeHighlyVisibleLogHeader(app.Name, 15)
	// GLFW event handling must run on the main OS thread
	// See documentation for functions that are only allowed to be called from the main thread.
	runtime.LockOSThread()
}

func HypervisorScreenTeardown() {
	glfw.Terminate()

}

func HypervisorScreenInit() {
	fmt.Printf("Hypervisor: Init glfw \n")

	if err := glfw.Init(); err != nil {
		log.Fatalln("failed to initialize glfw:", err)
	}

	//defer glfw.Terminate()

	fmt.Printf("Hypervisor: set windowhint\n")
	glfw.WindowHint(glfw.Resizable, glfw.True)
	glfw.WindowHint(glfw.ContextVersionMajor, 2)
	glfw.WindowHint(glfw.ContextVersionMinor, 1)

	var err error
	GlfwWindow, err = glfw.CreateWindow(int(gfx.CurrAppWidth), int(gfx.CurrAppHeight), app.Name, nil, nil)

	if err != nil {
		panic(err)
	}

	GlfwWindow.MakeContextCurrent()
	if err := gl.Init(); err != nil {
		panic(err)
	}
	fmt.Printf("Hypervisor: load texture \n")
	Texture = NewTexture("Bisasam_24x24_Shadowed.png")
	//defer gl.DeleteTextures(1, &Texture)

	fmt.Printf("Hypervisor: init renderer \n")
	InitRenderer()
}

func HypervisorInitInputEvents() {
	fmt.Printf("Hypervisor: init InitInputEvents \n")
	InitInputEvents(GlfwWindow)
}

func PollUiInputEvents() {
	glfw.PollEvents()
}

//could be in messages
func DispatchEvents(ch chan []byte) {
	select {
	case v := <-ch:
		processEvents(v)
	default:
		//fmt.Println("MonitorEvents() default")
	}
}

func UpdateDrawBuffer() {
	DrawScene()
}

func SwapDrawBuffer() {
	GlfwWindow.SwapBuffers()
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
	setFrustum(gfx.InitFrustum)
	//gl.Frustum(-1, 1, -1, 1, 1.0, 10.0)
	//gl.Frustum(left, right, bottom, top, zNear, zFar)
	gl.MatrixMode(gl.MODELVIEW)
	gl.LoadIdentity()

	// future FIXME: finished app would not have a demo program loaded on startup?
	script.Process(false)
}

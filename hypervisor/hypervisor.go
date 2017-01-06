package hypervisor

import (
	"fmt"
	"github.com/corpusc/viscript/app"
	"github.com/corpusc/viscript/gfx"
	//"github.com/corpusc/viscript/msg"
	//"github.com/corpusc/viscript/script"
	//"github.com/go-gl/gl/v2.1/gl"
	"github.com/go-gl/glfw/v3.2/glfw"
	//"log"
	"runtime"

	igl "github.com/corpusc/viscript/gl" //internal gl
)

//var GlfwWindow *glfw.Window //deprecate eventually
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
	fmt.Printf("Hypervisor: screen init\n")
	igl.WindowInit()
	igl.LoadTextures()
	igl.InitRenderer()
}

func HypervisorInitInputEvents() {
	fmt.Printf("Hypervisor: init InitInputEvents \n")
	InitInputEvents(igl.GlfwWindow)
}

func PollUiInputEvents() {
	glfw.PollEvents()
}

//could be in messages
func DispatchInputEvents(ch chan []byte) {
	for len(ch) > 0 { //if channel has elements
		v := <-ch //read from channel
		ProcessInputEvents(v)
	}
}

func UpdateDrawBuffer() {
	igl.DrawScene()
}

func SwapDrawBuffer() {
	igl.SwapDrawBuffer()
}

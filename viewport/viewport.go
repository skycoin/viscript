package viewport

import (
	"fmt"

	"github.com/corpusc/viscript/app"
	//"github.com/corpusc/viscript/script"
	//"log"
	"runtime"

	igl "github.com/corpusc/viscript/viewport/gl" //internal gl
)

//glfw
//glfw.PollEvents()
//only remaining

var (
	CloseWindow bool = false
)

func init() {
	fmt.Println("viewport.init()")
}

func ViewportInit() {
	fmt.Println("viewport.ViewportInit()")
	app.MakeHighlyVisibleLogHeader(app.Name, 15)
	// GLFW event handling must run on the main OS thread
	// See documentation for functions that are only allowed to be called from the main thread.
	runtime.LockOSThread()
}

func ViewportScreenTeardown() {
	igl.ScreenTeardown()
}

func ViewportScreenInit() {
	fmt.Printf("Viewport: screen init\n")
	igl.WindowInit()
	igl.LoadTextures()
	igl.InitRenderer()
}

func ViewportInitInputEvents() {
	fmt.Printf("Viewport: init InitInputEvents \n")
	igl.InitInputEvents(igl.GlfwWindow)
	igl.InitMiscEvents(igl.GlfwWindow)
}

func PollUiInputEvents() {
	igl.PollEvents() //move to gl
}

//could be in messages
func DispatchInputEvents(ch chan []byte) []byte {
	message := []byte{}

	for len(ch) > 0 { //if channel has elements
		v := <-ch //read from channel
		message = ProcessInputEvents(v)

	}

	return message
}

func Update() {
	igl.Update()
}

func UpdateDrawBuffer() {
	igl.DrawScene()
}

func SwapDrawBuffer() {
	igl.SwapDrawBuffer()
}

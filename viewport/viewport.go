package viewport

import (
	"fmt"

	"github.com/corpusc/viscript/app"
	"github.com/corpusc/viscript/viewport/terminal"
	//"github.com/corpusc/viscript/script"
	//"log"
	"runtime"

	igl "github.com/corpusc/viscript/viewport/gl" //internal gl
)

//glfw
//glfw.PollEvents()
//only remaining

var (
	CloseWindow bool                   = false
	Terms       terminal.TerminalStack = terminal.TerminalStack{}
)

func ViewportInit() {
	fmt.Println("viewport.ViewportInit()")
	app.MakeHighlyVisibleLogHeader(app.Name, 15)
	// GLFW event handling must run on the main OS thread
	// See documentation for functions that are only allowed to be called from the main thread.
	runtime.LockOSThread()
}

func ViewportTerminalsInit() {
	fmt.Println("viewport.init()")
	Terms.Init()
	Terms.AddTerminal()
	Terms.AddTerminal()
	Terms.AddTerminal()
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

func InitEvents() {
	fmt.Printf("Viewport: init events \n")
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
		message = UnpackInputEvents(v)
	}

	return message
}

func Tick() {
	igl.Curs.Tick()
	Terms.Tick()
}

func UpdateDrawBuffer() {
	igl.DrawBegin()
	Terms.Draw()
	igl.DrawEnd()
}

func SwapDrawBuffer() {
	igl.SwapDrawBuffer()
}

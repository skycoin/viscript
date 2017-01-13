package hypervisor

import (
	"fmt"

	"github.com/corpusc/viscript/app"
	//"github.com/corpusc/viscript/gfx"
	//"github.com/corpusc/viscript/msg"
	//"github.com/corpusc/viscript/script"
	//"github.com/go-gl/gl/v2.1/gl"
	//"github.com/go-gl/glfw/v3.2/glfw"
	//"log"
	"runtime"

	igl "github.com/corpusc/viscript/gl" //internal gl
)

//glfw
//glfw.PollEvents()
//only remaining

var Focused *Terminal
var Panels []*Terminal

var CloseWindow chan int           // write to channel to close
var runPanelHeiFrac = float32(0.4) // TEMPORARY fraction of vertical strip height which is dedicated to running code

func init() {
	fmt.Println("hypervisor.init()")
	CloseWindow = make(chan int) //write to close windows
	initPanels()
}

func initPanels() {
	Panels = append(Panels, &Terminal{FractionOfStrip: 1 - runPanelHeiFrac, IsEditable: true})
	Panels = append(Panels, &Terminal{FractionOfStrip: runPanelHeiFrac, IsEditable: true}) // console (runtime feedback log)	// FIXME so its not editable once we're done debugging some things
	Focused = Panels[0]

	Panels[0].Init()
	Panels[0].SetupDemoProgram()
	Panels[1].Init()
}

func HypervisorInit() {
	fmt.Println("hypervisor.HypervisorInit()")
	app.MakeHighlyVisibleLogHeader(app.Name, 15)
	// GLFW event handling must run on the main OS thread
	// See documentation for functions that are only allowed to be called from the main thread.
	runtime.LockOSThread()
}

func HypervisorScreenTeardown() {
	igl.ScreenTeardown()
}

func HypervisorScreenInit() {
	fmt.Printf("Hypervisor: screen init\n")
	igl.WindowInit()
	igl.LoadTextures()
	igl.InitRenderer()
}

func HypervisorInitInputEvents() {
	fmt.Printf("Hypervisor: init InitInputEvents \n")
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

func UpdateDrawBuffer() {
	igl.DrawScene()

	for _, pan := range Panels {
		//fmt.Println("drawing a panel", pan.FractionOfStrip)
		pan.Draw()
	}
}

func SwapDrawBuffer() {
	igl.SwapDrawBuffer()
}

// refactoring (possibly termporary) additions
func SetSize() {
	for _, pan := range Panels {
		pan.SetSize()
	}
}

func ScrollPanelThatIsHoveredOver(mousePixelDeltaX, mousePixelDeltaY float32) {
	for _, pan := range Panels {
		pan.ScrollIfMouseOver(mousePixelDeltaX, mousePixelDeltaY)
	}
}

func InsertRuneIntoDocument(s string, message uint32) string {
	f := Focused
	b := f.TextBodies[0]
	resultsDif := f.CursX - len(b[f.CursY])
	fmt.Printf("Rune   [%s: %s]", s, string(message))

	if f.CursX > len(b[f.CursY]) {
		b[f.CursY] = b[f.CursY][:f.CursX-resultsDif] + b[f.CursY][:len(b[f.CursY])] + string(message)
		fmt.Printf("line is %s\n", b[f.CursY])
		f.CursX++
	} else {
		b[f.CursY] = b[f.CursY][:f.CursX] + string(message) + b[f.CursY][f.CursX:len(b[f.CursY])]
		f.CursX++
	}

	return string(message)
}

package main

import (
	"fmt"
	"go/build"
	_ "image/png"
	"log"
	_ "os"
	"runtime"

	"github.com/go-gl/gl/v2.1/gl"
	"github.com/go-gl/glfw/v3.1/glfw"
)

var appName = "V I S C R I P T"

func init() {
	makeHighlyVisibleRuntimeLogHeader(appName, 15)
	// GLFW event handling must run on the main OS thread
	// See documentation for functions that are only allowed to be called from the main thread.
	fmt.Println("init()")
	runtime.LockOSThread()
}

func main() {
	if err := glfw.Init(); err != nil {
		log.Fatalln("failed to initialize glfw:", err)
	}

	defer glfw.Terminate()

	glfw.WindowHint(glfw.Resizable, glfw.False)
	glfw.WindowHint(glfw.ContextVersionMajor, 2)
	glfw.WindowHint(glfw.ContextVersionMinor, 1)
	window, err := glfw.CreateWindow(resX, resY, appName, nil, nil)

	if err != nil {
		panic(err)
	}

	window.MakeContextCurrent()

	if err := gl.Init(); err != nil {
		panic(err)
	}

	texture = newTexture("Bisasam_24x24_Shadowed.png")
	defer gl.DeleteTextures(1, &texture)

	initRenderer()
	initInputEvents(window)
	initParser()

	for !window.ShouldClose() {
		monitorEvents(events)
		pollEventsAndHandleAnInput(window)
		drawScene()
		window.SwapBuffers()
	}
}

// Set the working directory to the root of Go package, so that its assets can be accessed.
func init() {
	/*
		//dir, err := importPathToDir("github.com/go-gl/examples/glfw31-gl21-cube")
		dir, err := importPathToDir("viscript")

		if err != nil {
			log.Fatalln("Unable to find Go package in your GOPATH, it's needed to load assets:", err)
		}

		err = os.Chdir(dir)

		if err != nil {
			log.Panicln("os.Chdir:", err)
		}
	*/
}

// importPathToDir resolves the absolute path from importPath.
// There doesn't need to be a valid Go package inside that import path,
// but the directory must exist.
func importPathToDir(importPath string) (string, error) {
	p, err := build.Import(importPath, "", build.FindOnly)

	if err != nil {
		return "", err
	}

	return p.Dir, nil
}

// numLines: use odd number for an exact middle point
func makeHighlyVisibleRuntimeLogHeader(s string, numLines int) {
	s = " " + s + " "
	fillChar := "#"

	var bar string
	for i := 0; i < 79; i++ {
		bar += fillChar
	}

	var spaces string
	for i := 0; i < len(s); i++ {
		spaces += " "
	}

	var bookend string
	for i := 0; i < (79-len(s))/2; i++ {
		bookend += fillChar
	}

	middle := numLines / 2
	for i := 0; i < numLines; i++ {
		switch {
		case i == middle:
			fmt.Println(bookend + s + bookend)
		case i == middle-1 || i == middle+1:
			fmt.Println(bookend + spaces + bookend)
		default:
			fmt.Println(bar)
		}
	}
}

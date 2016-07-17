// Copyright 2014 The go-gl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Renders a textured spinning cube using GLFW 3.1 and OpenGL 2.1.
package main

import (
	"go/build"
	_ "image/png"
	"log"
	"os"
	"runtime"

	"github.com/go-gl/gl/v2.1/gl"
	"github.com/go-gl/glfw/v3.1/glfw"
)

func init() {
	// GLFW event handling must run on the main OS thread
	// See documentation for functions that are only allowed to be called from the main thread.
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
	window, err := glfw.CreateWindow(resX, resY, "V I S", nil, nil)

	if err != nil {
		panic(err)
	}

	window.MakeContextCurrent()

	if err := gl.Init(); err != nil {
		panic(err)
	}

	texture = newTexture("Bisasam_24x24_Shadowed.png")
	defer gl.DeleteTextures(1, &texture)

	setupScene()
	initDoc()
	initInputEvents(window)

	for !window.ShouldClose() {
		pollEventsAndHandleAnInput(window)
		drawScene()
		window.SwapBuffers()
	}
}

// Set the working directory to the root of Go package, so that its assets can be accessed.
func init() {
	//dir, err := importPathToDir("github.com/go-gl/examples/glfw31-gl21-cube")
	dir, err := importPathToDir("viscript")

	if err != nil {
		log.Fatalln("Unable to find Go package in your GOPATH, it's needed to load assets:", err)
	}

	err = os.Chdir(dir)

	if err != nil {
		log.Panicln("os.Chdir:", err)
	}
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

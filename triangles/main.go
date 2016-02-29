// Modified from Example 1.1: triangles.cpp
// OpenGL Programming Guide (Eighth Edition)
package main

import (
	"log"
	"os"
	"runtime"

	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/glfw/v3.1/glfw"
	"github.com/hurricanerix/gorb/shader"
)

const ( // VAO_IDs
	Triangles = iota
	NumVAOs   = iota
)

const ( // Buffer IDs
	ArrayBuffer = iota
	NumBuffers  = iota
)

const ( // Attrib IDs
	vPosition = 0
)

var (
	VAOs    [NumVAOs]uint32
	Buffers [NumBuffers]uint32
)

const NumVertices = int32(6)

func init() {
	log.Println("init")
	// This is needed to arrange that main() runs on main thread.
	// See documentation for functions that are only allowed to be called from the main thread.
	runtime.LockOSThread()
}

func main() {
	//log.SetOutput(ioutil.Discard)
	log.Println("main")

	// NOTE: Using GLFW instead of GLUT
	if err := glfw.Init(); err != nil {
		panic(err)
	}
	defer glfw.Terminate()

	glfw.WindowHint(glfw.ContextVersionMajor, 4)
	glfw.WindowHint(glfw.ContextVersionMinor, 1) // 4.3 does not work on my Macbook
	glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)
	glfw.WindowHint(glfw.OpenGLForwardCompatible, glfw.True)
	window, err := glfw.CreateWindow(512, 512, os.Args[0], nil, nil)
	if err != nil {
		panic(err)
	}
	window.MakeContextCurrent()

	if err := gl.Init(); err != nil {
		log.Fatalln("unable to initialize Glow ... exiting:", err)
	}

	initGL()

	for !window.ShouldClose() {
		display()

		window.SwapBuffers()
		glfw.PollEvents()
	}
}

func initGL() {
	log.Println("initGL")
	gl.GenVertexArrays(NumVAOs, &VAOs[0])
	gl.BindVertexArray(VAOs[Triangles])

	vertices := [...]float32{
		-0.90, -0.90, // Triangle 1
		0.85, -0.90,
		-0.90, 0.85,
		0.90, -0.85, // Triangle 2
		0.90, 0.90,
		-0.85, 0.90,
	}

	gl.GenBuffers(NumBuffers, &Buffers[0])
	gl.BindBuffer(gl.ARRAY_BUFFER, Buffers[ArrayBuffer])
	gl.BufferData(gl.ARRAY_BUFFER, len(vertices), gl.Ptr(&vertices), gl.STATIC_DRAW)

	shaders := []shader.Info{
		shader.Info{Type: gl.VERTEX_SHADER, Filename: "triangles.vert"},
		shader.Info{Type: gl.FRAGMENT_SHADER, Filename: "triangles.frag"},
	}

	program, err := shader.Load(&shaders)
	if err != nil {
		panic(err)
	}
	gl.UseProgram(program)

	// #define BUFFER_OFFSET(offset) ((void *)(offset))
	// glVertexAttribPointer(vPosition, 2, GL_FLOAT, GL_FALSE, 0, BUFFER_OFFSET(0));
	gl.VertexAttribPointer(vPosition, 2, gl.FLOAT, false, 0, gl.PtrOffset(0))
	gl.EnableVertexAttribArray(vPosition)
}

func display() {
	//log.Println("display")
	gl.Clear(gl.COLOR_BUFFER_BIT)

	gl.BindVertexArray(VAOs[Triangles])
	gl.DrawArrays(gl.TRIANGLES, 0, NumVertices)

	gl.Flush()
}

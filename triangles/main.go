// Modified from Example 1.1: triangles.cpp
// OpenGL Programming Guide (Eighth Edition)
package main

import (
	"fmt"
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
	// NOTE: Using GLFW instead of GLUT
	if err := glfw.Init(); err != nil {
		log.Fatalln("failed to initialize glfw:", err)
	}
	defer glfw.Terminate()

	glfw.WindowHint(glfw.Resizable, glfw.False)
	glfw.WindowHint(glfw.ContextVersionMajor, 4)
	glfw.WindowHint(glfw.ContextVersionMinor, 3)
	glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)
	glfw.WindowHint(glfw.OpenGLForwardCompatible, glfw.True)
	window, err := glfw.CreateWindow(512, 512, os.Args[0], nil, nil)
	if err != nil {
		log.Fatalln("failed to create window:", err)
	}
	window.MakeContextCurrent()

	if err := gl.Init(); err != nil {
		log.Fatalln("unable to initialize Glow ... exiting:", err)
	}

	fmt.Println("OpenGL vendor", gl.GoStr(gl.GetString(gl.VENDOR)))
	fmt.Println("OpenGL renderer", gl.GoStr(gl.GetString(gl.RENDERER)))
	fmt.Println("OpenGL version", gl.GoStr(gl.GetString(gl.VERSION)))
	fmt.Println("GLSL version", gl.GoStr(gl.GetString(gl.SHADING_LANGUAGE_VERSION)))

	initGL()

	for !window.ShouldClose() {
		display()

		window.SwapBuffers()
		glfw.PollEvents()
	}
}

func initGL() {
	shaders := []shader.Info{
		shader.Info{Type: gl.VERTEX_SHADER, Filename: "triangles.vert"},
		shader.Info{Type: gl.FRAGMENT_SHADER, Filename: "triangles.frag"},
	}

	program, err := shader.Load(&shaders)
	if err != nil {
		panic(err)
	}

	gl.UseProgram(program)

	vertices := []float32{
		-0.90, -0.90, // Triangle 1
		0.85, -0.90,
		-0.90, 0.85,
		0.90, -0.85, // Triangle 2
		0.90, 0.90,
		-0.85, 0.90,
	}

	gl.GenVertexArrays(NumVAOs, &VAOs[0])
	gl.BindVertexArray(VAOs[Triangles])

	gl.GenBuffers(NumBuffers, &Buffers[0])
	gl.BindBuffer(gl.ARRAY_BUFFER, Buffers[ArrayBuffer])
	gl.BufferData(gl.ARRAY_BUFFER, len(vertices)*4, gl.Ptr(vertices), gl.STATIC_DRAW)

	gl.VertexAttribPointer(vPosition, 2, gl.FLOAT, false, 0, gl.PtrOffset(0))
	gl.EnableVertexAttribArray(vPosition)

	gl.BindFragDataLocation(program, 0, gl.Str("fColor\x00"))
}

func display() {
	gl.Clear(gl.COLOR_BUFFER_BIT)

	gl.BindVertexArray(VAOs[Triangles])
	gl.DrawArrays(gl.TRIANGLES, 0, NumVertices)

	gl.Flush()
}

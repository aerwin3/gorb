// Modified from Example 3.5: ch03_drawcommands.cpp
// OpenGL Programming Guide (Eighth Edition)
package main

import (
	"go/build"
	"log"
	"os"
	"runtime"

	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/glfw/v3.1/glfw"
	"github.com/go-gl/mathgl/mgl32"
	"github.com/hurricanerix/gorb/shader"
)

const ( // VAO IDs
	Triangles = iota
	NumVAOs   = iota
)

const ( // Buffer IDs
	ArrayBuffer   = iota
	ElementBuffer = iota
	NumBuffers    = iota
)

var ( // Uniform IDs
	renderModelMatrixLoc      int32
	renderProjectionMatrixLoc int32
)

const ( // Attrib IDs
	position = 0
	color    = 1
)

var ( // Program IDs
	RenderProg uint32
)

var (
	VAOs    [NumVAOs]uint32
	Buffers [NumBuffers]uint32
)

var (
	Aspect float32
)

const NumVertices = int32(6)

func init() {
	// GLFW event handling must run on the main OS thread
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
	glfw.WindowHint(glfw.ContextVersionMinor, 1)
	glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)
	glfw.WindowHint(glfw.OpenGLForwardCompatible, glfw.True)
	window, err := glfw.CreateWindow(512, 512, os.Args[0], nil, nil)
	if err != nil {
		log.Fatalln("failed to create window:", err)
	}
	window.MakeContextCurrent()
	Aspect = float32(512) / float32(512)

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
	shaders := []shader.Info{
		shader.Info{Type: gl.VERTEX_SHADER, Filename: "primitive_restart.vert"},
		shader.Info{Type: gl.FRAGMENT_SHADER, Filename: "primitive_restart.frag"},
	}

	RenderProg, err := shader.Load(&shaders)
	if err != nil {
		panic(err)
	}

	gl.UseProgram(RenderProg)

	renderModelMatrixLoc = gl.GetUniformLocation(RenderProg, gl.Str("model_matrix\x00"))
	renderProjectionMatrixLoc = gl.GetUniformLocation(RenderProg, gl.Str("projection_matrix\x00"))

	// A single triangle
	vertexPositions := []float32{
		-1.0, -1.0, 0.0, 1.0,
		1.0, -1.0, 0.0, 1.0,
		-1.0, 1.0, 0.0, 1.0,
		-1.0, -1.0, 0.0, 1.0,
	}

	// Color for each vertex
	vertexColors := []float32{
		1.0, 1.0, 1.0, 1.0,
		1.0, 1.0, 0.0, 1.0,
		1.0, 0.0, 1.0, 1.0,
		0.0, 1.0, 1.0, 1.0,
	}

	// Indices for the triangle strips
	vertexIndices := []int{
		0, 1, 2,
	}

	// Set up the vertex attributes
	gl.GenVertexArrays(NumVAOs, &VAOs[0])
	gl.BindVertexArray(VAOs[Triangles])

	gl.GenBuffers(NumBuffers, &Buffers[0])
	// Set up the element array buffer
	gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, Buffers[ElementBuffer])
	gl.BufferData(gl.ELEMENT_ARRAY_BUFFER, len(vertexIndices)*4, gl.Ptr(vertexIndices), gl.STATIC_DRAW)

	gl.BindBuffer(gl.ARRAY_BUFFER, Buffers[ArrayBuffer])
	gl.BufferData(gl.ARRAY_BUFFER, len(vertexPositions)*4+len(vertexColors)*4, nil, gl.STATIC_DRAW)
	gl.BufferSubData(gl.ARRAY_BUFFER, 0, len(vertexPositions)*4, gl.Ptr(vertexPositions))
	gl.BufferSubData(gl.ARRAY_BUFFER, len(vertexPositions)*4, len(vertexColors)*4, gl.Ptr(vertexColors))

	gl.VertexAttribPointer(0, 4, gl.FLOAT, false, 0, nil)
	//gl.VertexAttribPointer(1, 4, gl.FLOAT, false, 0, gl.Ptr(len(vertexPositions)*4))
	gl.VertexAttribPointer(1, 4, gl.FLOAT, false, 0, gl.Ptr(vertexPositions))

	gl.EnableVertexAttribArray(0)
	gl.EnableVertexAttribArray(1)

	gl.ClearColor(0.0, 0.0, 0.0, 1.0)

	gl.BindFragDataLocation(RenderProg, 0, gl.Str("color\x00"))
}

func display() {
	/*
		//t := float32(GetTickCount() & 0x1FFF) / float(0x1FFF);
		//static float q = 0.0f;
		//static const vmath::vec3 X(1.0f, 0.0f, 0.0f);
		//static const vmath::vec3 Y(0.0f, 1.0f, 0.0f);
		//static const vmath::vec3 Z(0.0f, 0.0f, 1.0f);
	*/

	// Setup
	gl.Enable(gl.CULL_FACE)
	gl.Disable(gl.DEPTH_TEST)

	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

	// Activate simple shading program
	gl.UseProgram(RenderProg)

	// Set up the model and projection matrix
	projectionMatrix := mgl32.Frustum(-1, 1, -Aspect, Aspect, 1, 500)
	gl.UniformMatrix4fv(renderProjectionMatrixLoc, 1, false, &projectionMatrix[0])

	// Set up for a glDrawElements call
	gl.BindVertexArray(VAOs[Triangles])
	gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, Buffers[ElementBuffer])

	var modelMatrix mgl32.Mat4
	// Draw Arrays...
	modelMatrix = mgl32.Translate3D(-3, 0, -5)
	gl.UniformMatrix4fv(renderModelMatrixLoc, 4, false, &modelMatrix[0])
	gl.DrawArrays(gl.TRIANGLES, 0, 3)

	// DrawElements
	modelMatrix = mgl32.Translate3D(-1, 0, -5)
	gl.UniformMatrix4fv(renderModelMatrixLoc, 4, false, &modelMatrix[0])
	gl.DrawElements(gl.TRIANGLES, 3, gl.UNSIGNED_SHORT, nil)

	// DrawElementsBaseVertex
	modelMatrix = mgl32.Translate3D(1, 0, -5)
	gl.UniformMatrix4fv(renderModelMatrixLoc, 4, false, &modelMatrix[0])
	gl.DrawElementsBaseVertex(gl.TRIANGLES, 3, gl.UNSIGNED_SHORT, nil, 1)

	// DrawArraysInstanced
	modelMatrix = mgl32.Translate3D(3, 0, -5)
	gl.UniformMatrix4fv(renderModelMatrixLoc, 4, false, &modelMatrix[0])
	gl.DrawArraysInstanced(gl.TRIANGLES, 0, 3, 1)

	gl.Flush()
}

// Set the working directory to the root of Go package, so that its assets can be accessed.
func init() {

	dir, err := importPathToDir("github.com/hurricanerix/gorb/03/ch03-5_drawcommands")
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

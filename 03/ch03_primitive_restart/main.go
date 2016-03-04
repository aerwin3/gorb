// Example modified from OpenGL Programming Guide (Eighth Edition)
package main

import (
	"go/build"
	"log"
	"os"
	"runtime"
	"unsafe"

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
	Aspect              float32
	UsePrimitiveRestart bool
	ModelMatrix         mgl32.Mat4
	ProjectionMatrix    mgl32.Mat4
)

const NumVertices = int32(6)

func init() {
	// GLFW event handling must run on the main OS thread
	runtime.LockOSThread()
}

func main() {
	// There is a slight modification to set the background color
	// depending on if UsePrimitiveRestart is set.
	// Press the spacebar to enable/disable primitive restart.

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
	window.SetKeyCallback(keyCallback)
	Aspect = float32(512) / float32(512)

	if err := gl.Init(); err != nil {
		log.Fatalln("unable to initialize Glow ... exiting:", err)
	}

	initGL()

	var t float32
	for !window.ShouldClose() {
		// TODO: Make this time based
		t += 0.0001
		update(t)
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

	// 8 corners of a cube, side length 2, centered on the origin
	vertexPositions := []float32{
		-1.0, -1.0, -1.0, 1.0,
		-1.0, -1.0, 1.0, 1.0,
		-1.0, 1.0, -1.0, 1.0,
		-1.0, 1.0, 1.0, 1.0,
		1.0, -1.0, -1.0, 1.0,
		1.0, -1.0, 1.0, 1.0,
		1.0, 1.0, -1.0, 1.0,
		1.0, 1.0, 1.0, 1.0,
	}

	// Color for each vertex
	vertexColors := []float32{
		1.0, 1.0, 1.0, 1.0,
		1.0, 1.0, 0.0, 1.0,
		1.0, 0.0, 1.0, 1.0,
		1.0, 0.0, 0.0, 1.0,
		0.0, 1.0, 1.0, 1.0,
		0.0, 1.0, 0.0, 1.0,
		0.0, 0.0, 1.0, 1.0,
		0.5, 0.5, 0.5, 1.0,
	}

	// Indices for the triangle strips
	vertexIndices := []uint16{
		0, 1, 2, 3, 6, 7, 4, 5, // First strip
		0xFFFF,                 // <<-- This is the restart index
		2, 6, 0, 4, 1, 5, 3, 7, // Second strip
	}

	// Set up the element array buffer
	sizeVertexIndices := len(vertexIndices) * int(unsafe.Sizeof(vertexIndices[0]))

	gl.GenBuffers(1, &Buffers[ElementBuffer])
	gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, Buffers[ElementBuffer])
	gl.BufferData(gl.ELEMENT_ARRAY_BUFFER, sizeVertexIndices, gl.Ptr(vertexIndices), gl.STATIC_DRAW)

	// Set up the vertex attributes
	sizeVertexPositions := len(vertexPositions) * int(unsafe.Sizeof(vertexPositions[0]))
	sizeVertexColors := len(vertexColors) * int(unsafe.Sizeof(vertexColors[0]))

	gl.GenVertexArrays(1, &VAOs[Triangles])
	gl.BindVertexArray(VAOs[Triangles])

	gl.GenBuffers(1, &Buffers[ArrayBuffer])
	gl.BindBuffer(gl.ARRAY_BUFFER, Buffers[ArrayBuffer])
	gl.BufferData(gl.ARRAY_BUFFER, sizeVertexPositions+sizeVertexColors, nil, gl.STATIC_DRAW)
	gl.BufferSubData(gl.ARRAY_BUFFER, 0, sizeVertexPositions, gl.Ptr(vertexPositions))
	gl.BufferSubData(gl.ARRAY_BUFFER, sizeVertexPositions, sizeVertexColors, gl.Ptr(vertexColors))

	gl.VertexAttribPointer(0, 4, gl.FLOAT, false, 0, nil)
	gl.VertexAttribPointer(1, 4, gl.FLOAT, false, 0, gl.PtrOffset(sizeVertexPositions))
	gl.EnableVertexAttribArray(0)
	gl.EnableVertexAttribArray(1)

	UsePrimitiveRestart = true
	gl.ClearColor(0.05, 0.1, 0.05, 1.0)
}

func update(t float32) {
	//static float q = 0.0f;
	//X := mgl32.Vec3{1, 0, 0}
	Y := mgl32.Vec3{0, 1, 0}
	Z := mgl32.Vec3{0, 0, 1}

	// Set up the model and projection matrix
	ModelMatrix = mgl32.Translate3D(0, 0, -5).Mul4(mgl32.HomogRotate3D(t*360, Y)).Mul4(mgl32.HomogRotate3D(t*720, Z))
	ProjectionMatrix = mgl32.Frustum(-1, 1, -Aspect, Aspect, 1, 500)
}

func display() {
	// Setup
	gl.Enable(gl.CULL_FACE)
	gl.Disable(gl.DEPTH_TEST)

	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

	// Activate simple shading program
	// TODO: figure out why enabling this does not work
	//gl.UseProgram(RenderProg)

	gl.UniformMatrix4fv(renderModelMatrixLoc, 1, false, &ModelMatrix[0])
	gl.UniformMatrix4fv(renderProjectionMatrixLoc, 1, false, &ProjectionMatrix[0])

	// Set up for a glDrawElements call
	gl.BindVertexArray(VAOs[Triangles])
	gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, Buffers[ElementBuffer])

	if UsePrimitiveRestart {
		// When primitive restart is on, we can call one draw command
		gl.ClearColor(0.05, 0.1, 0.05, 1.0)
		gl.Enable(gl.PRIMITIVE_RESTART)
		gl.PrimitiveRestartIndex(0xFFFF)
		gl.DrawElements(gl.TRIANGLE_STRIP, 17, gl.UNSIGNED_SHORT, nil)
	} else {
		gl.ClearColor(0.05, 0.05, 0.1, 1.0)
		// Without primitive restart, we need to call two draw commands
		gl.DrawElements(gl.TRIANGLE_STRIP, 8, gl.UNSIGNED_SHORT, nil)
		gl.DrawElements(gl.TRIANGLE_STRIP, 8, gl.UNSIGNED_SHORT, gl.PtrOffset(9*2)) // (const GLvoid *)(9 * sizeof(GLushort))
	}

	gl.Flush()
}

// Set the working directory to the root of Go package, so that its assets can be accessed.
func init() {

	dir, err := importPathToDir("github.com/hurricanerix/gorb/03/ch03_primitive_restart")
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

func keyCallback(w *glfw.Window, key glfw.Key, scancode int, action glfw.Action, mods glfw.ModifierKey) {
	if action == glfw.Release && key == glfw.KeyEscape {
		w.SetShouldClose(true)
	}
	if action == glfw.Release && key == glfw.KeySpace {
		UsePrimitiveRestart = !UsePrimitiveRestart
	}
}

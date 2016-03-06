// Example modified from OpenGL Programming Guide (Eighth Edition)
package main

import (
	"unsafe"

	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/glfw/v3.1/glfw"
	"github.com/go-gl/mathgl/mgl32"
	"github.com/hurricanerix/gorb/utils"
	"github.com/hurricanerix/gorb/utils/app"
	"github.com/hurricanerix/gorb/utils/shader"
)

const ( // Program IDs
	primRestartProgID = iota
	numPrograms       = iota
)

const ( // VAO IDs
	trianglesName = iota
	numVAOs       = iota
)

const ( // Buffer IDs
	arrayBufferName   = iota
	elementBufferName = iota
	numBuffers        = iota
)

const ( // Attrib Locations
	mcVertexLoc = 0 // TODO: rename to mcVertexLoc
	mcColor     = 1 // TODO: rename to mcColorLoc
)

var (
	Aspect              float32
	UsePrimitiveRestart bool
)

const NumVertices = int32(6)

type scene struct {
	Programs    [numPrograms]uint32
	VAOs        [numVAOs]uint32
	NumVertices [numVAOs]int32
	Buffers     [numBuffers]uint32
	// Uniform Locations
	ModelMatrixLoc      int32
	ProjectionMatrixLoc int32
	// App Settings
	ModelMatrix      mgl32.Mat4
	ProjectionMatrix mgl32.Mat4
	Rotation         float32
}

func (s *scene) Setup() error {
	shaders := []shader.Info{
		shader.Info{Type: gl.VERTEX_SHADER, Filename: "primitive_restart.vert"},
		shader.Info{Type: gl.FRAGMENT_SHADER, Filename: "primitive_restart.frag"},
	}

	program, err := shader.Load(&shaders)
	if err != nil {
		panic(err)
	}
	s.Programs[primRestartProgID] = program

	gl.UseProgram(s.Programs[primRestartProgID])

	s.ModelMatrixLoc = gl.GetUniformLocation(s.Programs[primRestartProgID], gl.Str("modelMatrix\x00"))
	s.ProjectionMatrixLoc = gl.GetUniformLocation(s.Programs[primRestartProgID], gl.Str("projectionMatrix\x00"))

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
	s.NumVertices[trianglesName] = int32(len(vertexPositions))

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

	gl.GenBuffers(numBuffers, &s.Buffers[0])

	gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, s.Buffers[elementBufferName])
	gl.BufferData(gl.ELEMENT_ARRAY_BUFFER, sizeVertexIndices, gl.Ptr(vertexIndices), gl.STATIC_DRAW)

	// Set up the vertex attributes
	sizeVertexPositions := len(vertexPositions) * int(unsafe.Sizeof(vertexPositions[0]))
	sizeVertexColors := len(vertexColors) * int(unsafe.Sizeof(vertexColors[0]))

	gl.GenVertexArrays(numVAOs, &s.VAOs[0])
	gl.BindVertexArray(s.VAOs[trianglesName])

	gl.BindBuffer(gl.ARRAY_BUFFER, s.Buffers[arrayBufferName])
	gl.BufferData(gl.ARRAY_BUFFER, sizeVertexPositions+sizeVertexColors, nil, gl.STATIC_DRAW)
	gl.BufferSubData(gl.ARRAY_BUFFER, 0, sizeVertexPositions, gl.Ptr(vertexPositions))
	gl.BufferSubData(gl.ARRAY_BUFFER, sizeVertexPositions, sizeVertexColors, gl.Ptr(vertexColors))

	gl.VertexAttribPointer(0, 4, gl.FLOAT, false, 0, nil)
	gl.VertexAttribPointer(1, 4, gl.FLOAT, false, 0, gl.PtrOffset(sizeVertexPositions))
	gl.EnableVertexAttribArray(0)
	gl.EnableVertexAttribArray(1)

	UsePrimitiveRestart = true
	gl.ClearColor(0.05, 0.1, 0.05, 1.0)

	s.Rotation = 0
	return nil
}

func (s *scene) Update(dt float32) {
	s.Rotation += dt
	//static float q = 0.0f;
	//X := mgl32.Vec3{1, 0, 0}
	Y := mgl32.Vec3{0, 1, 0}
	Z := mgl32.Vec3{0, 0, 1}

	// Set up the model and projection matrix
	s.ModelMatrix = mgl32.Translate3D(0, 0, -5).Mul4(mgl32.HomogRotate3D(s.Rotation*360, Y)).Mul4(mgl32.HomogRotate3D(s.Rotation*720, Z))
	s.ProjectionMatrix = mgl32.Frustum(-1, 1, -Aspect, Aspect, 1, 500)
}

func (s *scene) Display() {
	// Setup
	gl.Enable(gl.CULL_FACE)
	gl.Disable(gl.DEPTH_TEST)

	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

	// Activate simple shading program
	// TODO: figure out why enabling this does not work
	//gl.UseProgram(RenderProg)

	gl.UniformMatrix4fv(s.ModelMatrixLoc, 1, false, &s.ModelMatrix[0])
	gl.UniformMatrix4fv(s.ProjectionMatrixLoc, 1, false, &s.ProjectionMatrix[0])

	// Set up for a glDrawElements call
	gl.BindVertexArray(s.VAOs[trianglesName])
	gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, s.Buffers[elementBufferName])

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

func (s *scene) Cleanup() {
}

// Main methods
func init() {
	if err := utils.SetWorkingDir("github.com/hurricanerix/gorb/03/primitive_restart"); err != nil {
		panic(err)
	}
}

func KeyCallback(w *glfw.Window, key glfw.Key, scancode int, action glfw.Action, mods glfw.ModifierKey) {
	if action == glfw.Release && key == glfw.KeySpace {
		UsePrimitiveRestart = !UsePrimitiveRestart
	}
}

func main() {
	c := app.Config{
		Name:                "Ch3-PrimitiveRestart",
		DefaultScreenWidth:  512,
		DefaultScreenHeight: 512,
		EscapeToQuit:        true,
		SupportedGLVers: []mgl32.Vec2{
			mgl32.Vec2{4, 3},
			mgl32.Vec2{4, 1},
		},
		KeyCallback: KeyCallback,
	}
	// TODO: Get the w/h to calculate the correct aspect ratio
	Aspect = float32(512) / float32(512)

	s := &scene{}

	a := app.New(c, s)

	if err := a.Run(); err != nil {
		panic(err)
	}
}

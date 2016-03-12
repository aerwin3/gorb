// Example modified from OpenGL Programming Guide (Eighth Edition)
package main

import (
	"runtime"
	"unsafe"

	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/mathgl/mgl32"
	"github.com/hurricanerix/go-gl-utils/app"
	"github.com/hurricanerix/go-gl-utils/path"
	"github.com/hurricanerix/go-gl-utils/shader"
)

func init() {
	runtime.LockOSThread()
	if err := path.SetWorkingDir("github.com/hurricanerix/gorb/03/drawcommands"); err != nil {
		panic(err)
	}
}

func main() {
	c := app.Config{
		Name:                "Ch3-DrawCommands",
		DefaultScreenWidth:  512,
		DefaultScreenHeight: 512,
		EscapeToQuit:        true,
		SupportedGLVers: []mgl32.Vec2{
			mgl32.Vec2{4, 3},
			mgl32.Vec2{4, 1},
		},
	}
	// TODO: Get the w/h to calculate the correct aspect ratio
	Aspect = float32(512) / float32(512)

	s := &scene{}

	a := app.New(c, s)

	if err := a.Run(); err != nil {
		panic(err)
	}
}

// Everything below this line is for the Scene implementation.

const ( // Program IDs
	primRestartProgID = iota
	numPrograms       = iota
)

const ( // VAO Names
	trianglesName = iota
	numVAOs       = iota
)

const ( // Buffer Names
	arrayBufferName   = iota
	elementBufferName = iota
	numBuffers        = iota
)

const ( // Attrib Locations
	mcVertex = 0 // TODO: rename to mcVertexLoc
	mcColor  = 1 // TODO: rename to mcColorLoc
)

var (
	Aspect float32
)

type scene struct {
	Programs    [numPrograms]uint32
	VAOs        [numVAOs]uint32
	NumVertices [numVAOs]int32
	Buffers     [numBuffers]uint32
	// Uniform Locations
	ModelMatrixLoc      int32
	ProjectionMatrixLoc int32
}

func (s *scene) Setup() error {
	shaders := []shader.Info{
		shader.Info{Type: gl.VERTEX_SHADER, Filename: "../primitive-restart/primitive_restart.vert"},
		shader.Info{Type: gl.FRAGMENT_SHADER, Filename: "../primitive-restart/primitive_restart.frag"},
	}

	program, err := shader.Load(&shaders)
	if err != nil {
		panic(err)
	}
	s.Programs[primRestartProgID] = program

	gl.UseProgram(s.Programs[primRestartProgID])

	s.ModelMatrixLoc = gl.GetUniformLocation(s.Programs[primRestartProgID], gl.Str("modelMatrix\x00"))
	s.ProjectionMatrixLoc = gl.GetUniformLocation(s.Programs[primRestartProgID], gl.Str("projectionMatrix\x00"))

	// A single triangle
	vertexPositions := []float32{
		-1.0, -1.0, 0.0, 1.0,
		1.0, -1.0, 0.0, 1.0,
		-1.0, 1.0, 0.0, 1.0,
		-1.0, -1.0, 0.0, 1.0,
	}
	s.NumVertices[trianglesName] = int32(len(vertexPositions))

	// Color for each vertex
	vertexColors := []float32{
		1.0, 1.0, 1.0, 1.0,
		1.0, 1.0, 0.0, 1.0,
		1.0, 0.0, 1.0, 1.0,
		0.0, 1.0, 1.0, 1.0,
	}

	// Indices for the triangle strips
	vertexIndices := []uint16{
		0, 1, 2,
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

	gl.ClearColor(0.0, 0.0, 0.0, 1.0)
	return nil
}

func (s *scene) Update(dt float32) {
	// This is where you would put code to update your scene.
	// This scene does not change, so there is nothing here.
}

func (s *scene) Display() {
	var modelMatrix mgl32.Mat4

	// Setup
	gl.Enable(gl.CULL_FACE)
	gl.Disable(gl.DEPTH_TEST)

	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

	// Activate simple shading program
	// TODO: figure out why enabling this does not work
	//gl.UseProgram(RenderProg)

	// Set up the model and projection matrix
	projectionMatrix := mgl32.Frustum(-1, 1, -Aspect, Aspect, 1, 500)
	gl.UniformMatrix4fv(s.ProjectionMatrixLoc, 1, false, &projectionMatrix[0])

	// Set up for a glDrawElements call
	gl.BindVertexArray(s.VAOs[trianglesName])
	gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, s.Buffers[elementBufferName])

	// Draw Arrays...
	modelMatrix = mgl32.Translate3D(-3, 0, -5)
	// TODO: figure out why the c++ version sends 4 instead of 1.
	//       maybe it is due to its matrix being stored as 4 arrays...?
	gl.UniformMatrix4fv(s.ModelMatrixLoc, 1, false, &modelMatrix[0])
	gl.DrawArrays(gl.TRIANGLES, 0, 3)

	// DrawElements
	modelMatrix = mgl32.Translate3D(-1, 0, -5)
	// TODO: figure out why the c++ version sends 4 instead of 1.
	//       maybe it is due to its matrix being stored as 4 arrays...?
	gl.UniformMatrix4fv(s.ModelMatrixLoc, 1, false, &modelMatrix[0])
	gl.DrawElements(gl.TRIANGLES, 3, gl.UNSIGNED_SHORT, nil)

	// DrawElementsBaseVertex
	modelMatrix = mgl32.Translate3D(1, 0, -5)
	// TODO: figure out why the c++ version sends 4 instead of 1.
	//       maybe it is due to its matrix being stored as 4 arrays...?
	gl.UniformMatrix4fv(s.ModelMatrixLoc, 1, false, &modelMatrix[0])
	gl.DrawElementsBaseVertex(gl.TRIANGLES, 3, gl.UNSIGNED_SHORT, nil, 1)

	// DrawArraysInstanced
	modelMatrix = mgl32.Translate3D(3, 0, -5)
	// TODO: figure out why the c++ version sends 4 instead of 1.
	//       maybe it is due to its matrix being stored as 4 arrays...?
	gl.UniformMatrix4fv(s.ModelMatrixLoc, 1, false, &modelMatrix[0])
	gl.DrawArraysInstanced(gl.TRIANGLES, 0, 3, 1)

	gl.Flush()
}

func (s *scene) Cleanup() {
	var id uint32
	for i := 0; i < numPrograms; i++ {
		id = s.Programs[i]
		gl.UseProgram(id)
		gl.DeleteProgram(id)
	}
}

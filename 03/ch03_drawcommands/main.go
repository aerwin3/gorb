// Example modified from OpenGL Programming Guide (Eighth Edition)
package main

import (
	"runtime"
	"unsafe"

	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/mathgl/mgl32"
	"github.com/hurricanerix/gorb/utils"
	"github.com/hurricanerix/gorb/utils/app"
	"github.com/hurricanerix/gorb/utils/shader"
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

const ( // Attrib IDs
	position = 0
	color    = 1
)

var ( // Uniform IDs
	renderModelMatrixLoc      int32
	renderProjectionMatrixLoc int32
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

type scene struct{}

func (s *scene) Setup() error {
	shaders := []shader.Info{
		shader.Info{Type: gl.VERTEX_SHADER, Filename: "../ch03_primitive_restart/primitive_restart.vert"},
		shader.Info{Type: gl.FRAGMENT_SHADER, Filename: "../ch03_primitive_restart/primitive_restart.frag"},
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
	vertexIndices := []uint16{
		0, 1, 2,
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

	gl.ClearColor(0.0, 0.0, 0.0, 1.0)
	return nil
}

func (s *scene) Update(dt float32) {

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
	gl.UniformMatrix4fv(renderProjectionMatrixLoc, 1, false, &projectionMatrix[0])

	// Set up for a glDrawElements call
	gl.BindVertexArray(VAOs[Triangles])
	gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, Buffers[ElementBuffer])

	// Draw Arrays...
	modelMatrix = mgl32.Translate3D(-3, 0, -5)
	// TODO: figure out why the c++ version sends 4 instead of 1.
	//       maybe it is due to its matrix being stored as 4 arrays...?
	gl.UniformMatrix4fv(renderModelMatrixLoc, 1, false, &modelMatrix[0])
	gl.DrawArrays(gl.TRIANGLES, 0, 3)

	// DrawElements
	modelMatrix = mgl32.Translate3D(-1, 0, -5)
	// TODO: figure out why the c++ version sends 4 instead of 1.
	//       maybe it is due to its matrix being stored as 4 arrays...?
	gl.UniformMatrix4fv(renderModelMatrixLoc, 1, false, &modelMatrix[0])
	gl.DrawElements(gl.TRIANGLES, 3, gl.UNSIGNED_SHORT, nil)

	// DrawElementsBaseVertex
	modelMatrix = mgl32.Translate3D(1, 0, -5)
	// TODO: figure out why the c++ version sends 4 instead of 1.
	//       maybe it is due to its matrix being stored as 4 arrays...?
	gl.UniformMatrix4fv(renderModelMatrixLoc, 1, false, &modelMatrix[0])
	gl.DrawElementsBaseVertex(gl.TRIANGLES, 3, gl.UNSIGNED_SHORT, nil, 1)

	// DrawArraysInstanced
	modelMatrix = mgl32.Translate3D(3, 0, -5)
	// TODO: figure out why the c++ version sends 4 instead of 1.
	//       maybe it is due to its matrix being stored as 4 arrays...?
	gl.UniformMatrix4fv(renderModelMatrixLoc, 1, false, &modelMatrix[0])
	gl.DrawArraysInstanced(gl.TRIANGLES, 0, 3, 1)

	gl.Flush()
}

func (s *scene) Cleanup() {
}

// Main methods
func init() {
	runtime.LockOSThread()
	if err := utils.SetWorkingDir("github.com/hurricanerix/gorb/03/ch03_drawcommands"); err != nil {
		panic(err)
	}
}

func main() {
	c := app.Config{
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

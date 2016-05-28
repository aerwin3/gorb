// Example modified from OpenGL Programming Guide (Eighth Edition)
package main

import (
	"runtime"
	"unsafe"

	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/glfw/v3.1/glfw"
	"github.com/go-gl/mathgl/mgl32"
	"github.com/hurricanerix/gorb/util"
)

func init() {
	runtime.LockOSThread()
	if err := util.SetWorkingDir("github.com/hurricanerix/gorb/03/drawcommands"); err != nil {
		panic(err)
	}
}

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
	programs    [numPrograms]uint32
	vaos        [numVAOs]uint32
	numVertices [numVAOs]int32
	buffers     [numBuffers]uint32
)
var ( // Uniform Locations
	modelMatrixLoc      int32
	projectionMatrixLoc int32
)

var (
	aspect float32
)

func main() {
	var err error

	// Get window context
	window, err := util.NewWindow("Ch3-DrawCommands", 512, 512)
	if err != nil {
		panic(err)
	}
	aspect = float32(512) / float32(512)

	// Load the GLSL program
	shaders := []util.ShaderInfo{
		util.ShaderInfo{Type: gl.VERTEX_SHADER, Filename: "../primitive-restart/primitive_restart.vert"},
		util.ShaderInfo{Type: gl.FRAGMENT_SHADER, Filename: "../primitive-restart/primitive_restart.frag"},
	}
	programs[primRestartProgID], err = util.Load(&shaders)
	if err != nil {
		panic(err)
	}
	gl.UseProgram(programs[primRestartProgID])

	// Setup model to be rendered
	modelMatrixLoc = gl.GetUniformLocation(programs[primRestartProgID], gl.Str("modelMatrix\x00"))
	projectionMatrixLoc = gl.GetUniformLocation(programs[primRestartProgID], gl.Str("projectionMatrix\x00"))
	vertexPositions := []float32{
		-1.0, -1.0, 0.0, 1.0,
		1.0, -1.0, 0.0, 1.0,
		-1.0, 1.0, 0.0, 1.0,
		-1.0, -1.0, 0.0, 1.0,
	}
	numVertices[trianglesName] = int32(len(vertexPositions))

	vertexColors := []float32{
		1.0, 1.0, 1.0, 1.0,
		1.0, 1.0, 0.0, 1.0,
		1.0, 0.0, 1.0, 1.0,
		0.0, 1.0, 1.0, 1.0,
	}

	vertexIndices := []uint16{
		0, 1, 2,
	}

	sizeVertexIndices := len(vertexIndices) * int(unsafe.Sizeof(vertexIndices[0]))

	gl.GenBuffers(numBuffers, &buffers[0])
	gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, buffers[elementBufferName])
	gl.BufferData(gl.ELEMENT_ARRAY_BUFFER, sizeVertexIndices, gl.Ptr(vertexIndices), gl.STATIC_DRAW)

	// Set up the vertex attributes
	sizeVertexPositions := len(vertexPositions) * int(unsafe.Sizeof(vertexPositions[0]))
	sizeVertexColors := len(vertexColors) * int(unsafe.Sizeof(vertexColors[0]))

	gl.GenVertexArrays(numVAOs, &vaos[0])
	gl.BindVertexArray(vaos[trianglesName])

	gl.BindBuffer(gl.ARRAY_BUFFER, buffers[arrayBufferName])
	gl.BufferData(gl.ARRAY_BUFFER, sizeVertexPositions+sizeVertexColors, nil, gl.STATIC_DRAW)
	gl.BufferSubData(gl.ARRAY_BUFFER, 0, sizeVertexPositions, gl.Ptr(vertexPositions))
	gl.BufferSubData(gl.ARRAY_BUFFER, sizeVertexPositions, sizeVertexColors, gl.Ptr(vertexColors))

	gl.VertexAttribPointer(0, 4, gl.FLOAT, false, 0, nil)
	gl.VertexAttribPointer(1, 4, gl.FLOAT, false, 0, gl.PtrOffset(sizeVertexPositions))
	gl.EnableVertexAttribArray(0)
	gl.EnableVertexAttribArray(1)

	gl.ClearColor(0.0, 0.0, 0.0, 1.0)

	// Main loop
	for !window.ShouldClose() {
		// Clear buffer
		gl.Enable(gl.CULL_FACE)
		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
		gl.Disable(gl.DEPTH_TEST)
		// TODO: figure out why enabling this does not work
		//gl.UseProgram(RenderProg)

		// Update
		var modelMatrix mgl32.Mat4
		// Set up the model and projection matrix
		projectionMatrix := mgl32.Frustum(-1, 1, -aspect, aspect, 1, 500)
		gl.UniformMatrix4fv(projectionMatrixLoc, 1, false, &projectionMatrix[0])
		// Set up for a glDrawElements call
		gl.BindVertexArray(vaos[trianglesName])
		gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, buffers[elementBufferName])

		// Render
		modelMatrix = mgl32.Translate3D(-3, 0, -5)
		// TODO: figure out why the c++ version sends 4 instead of 1.
		//       maybe it is due to its matrix being stored as 4 arrays...?
		gl.UniformMatrix4fv(modelMatrixLoc, 1, false, &modelMatrix[0])
		gl.DrawArrays(gl.TRIANGLES, 0, 3)

		// DrawElements
		modelMatrix = mgl32.Translate3D(-1, 0, -5)
		// TODO: figure out why the c++ version sends 4 instead of 1.
		//       maybe it is due to its matrix being stored as 4 arrays...?
		gl.UniformMatrix4fv(modelMatrixLoc, 1, false, &modelMatrix[0])
		gl.DrawElements(gl.TRIANGLES, 3, gl.UNSIGNED_SHORT, nil)

		// DrawElementsBaseVertex
		modelMatrix = mgl32.Translate3D(1, 0, -5)
		// TODO: figure out why the c++ version sends 4 instead of 1.
		//       maybe it is due to its matrix being stored as 4 arrays...?
		gl.UniformMatrix4fv(modelMatrixLoc, 1, false, &modelMatrix[0])
		gl.DrawElementsBaseVertex(gl.TRIANGLES, 3, gl.UNSIGNED_SHORT, nil, 1)

		// DrawArraysInstanced
		modelMatrix = mgl32.Translate3D(3, 0, -5)
		// TODO: figure out why the c++ version sends 4 instead of 1.
		//       maybe it is due to its matrix being stored as 4 arrays...?
		gl.UniformMatrix4fv(modelMatrixLoc, 1, false, &modelMatrix[0])
		gl.DrawArraysInstanced(gl.TRIANGLES, 0, 3, 1)

		// Swap Buffers
		gl.Flush()
		window.SwapBuffers()
		glfw.PollEvents()
	}

	// Cleanup
	for _, s := range shaders {
		s.Delete()
	}
	var id uint32
	for i := 0; i < numPrograms; i++ {
		id = programs[i]
		gl.UseProgram(id)
		gl.DeleteProgram(id)
	}
	util.Terminate()
}

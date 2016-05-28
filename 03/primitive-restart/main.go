// Example modified from OpenGL Programming Guide (Eighth Edition)
package main

import (
	"unsafe"

	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/glfw/v3.1/glfw"
	"github.com/go-gl/mathgl/mgl32"
	"github.com/hurricanerix/gorb/util"
)

func init() {
	if err := util.SetWorkingDir("github.com/hurricanerix/gorb/03/primitive-restart"); err != nil {
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
	mcVertexLoc = 0
	mcColor     = 1 // TODO: rename to mcColorLoc
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

var ( // App Settings
	modelMatrix         mgl32.Mat4
	projectionMatrix    mgl32.Mat4
	rotation            float32
	aspect              float32
	usePrimitiveRestart bool
)

func main() {
	var err error

	// Get window context
	window, err := util.NewWindow("Ch3-PrimitiveRestart", 512, 512)
	if err != nil {
		panic(err)
	}
	aspect = float32(512) / float32(512)

	// Load the GLSL program
	shaders := []util.ShaderInfo{
		util.ShaderInfo{Type: gl.VERTEX_SHADER, Filename: "primitive_restart.vert"},
		util.ShaderInfo{Type: gl.FRAGMENT_SHADER, Filename: "primitive_restart.frag"},
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
		-1.0, -1.0, -1.0, 1.0,
		-1.0, -1.0, 1.0, 1.0,
		-1.0, 1.0, -1.0, 1.0,
		-1.0, 1.0, 1.0, 1.0,
		1.0, -1.0, -1.0, 1.0,
		1.0, -1.0, 1.0, 1.0,
		1.0, 1.0, -1.0, 1.0,
		1.0, 1.0, 1.0, 1.0,
	}
	numVertices[trianglesName] = int32(len(vertexPositions))
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
	vertexIndices := []uint16{
		0, 1, 2, 3, 6, 7, 4, 5, // First strip
		0xFFFF,                 // <<-- This is the restart index
		2, 6, 0, 4, 1, 5, 3, 7, // Second strip
	}
	sizeVertexIndices := len(vertexIndices) * int(unsafe.Sizeof(vertexIndices[0]))
	gl.GenBuffers(numBuffers, &buffers[0])
	gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, buffers[elementBufferName])
	gl.BufferData(gl.ELEMENT_ARRAY_BUFFER, sizeVertexIndices, gl.Ptr(vertexIndices), gl.STATIC_DRAW)

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

	usePrimitiveRestart = true
	gl.ClearColor(0.05, 0.1, 0.05, 1.0)
	rotation = 0

	// Main loop
	for !window.ShouldClose() {
		// Clear buffer
		gl.Enable(gl.CULL_FACE)
		gl.Disable(gl.DEPTH_TEST)
		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
		// Activate simple shading program
		// TODO: figure out why enabling this does not work
		//gl.UseProgram(RenderProg)

		// Update
		rotation += .00005 // TODO: fix dt
		//static float q = 0.0f;
		//X := mgl32.Vec3{1, 0, 0}
		Y := mgl32.Vec3{0, 1, 0}
		Z := mgl32.Vec3{0, 0, 1}
		// Set up the model and projection matrix
		modelMatrix = mgl32.Translate3D(0, 0, -5).Mul4(mgl32.HomogRotate3D(rotation*360, Y)).Mul4(mgl32.HomogRotate3D(rotation*720, Z))
		projectionMatrix = mgl32.Frustum(-1, 1, -aspect, aspect, 1, 500)

		// Render
		gl.UniformMatrix4fv(modelMatrixLoc, 1, false, &modelMatrix[0])
		gl.UniformMatrix4fv(projectionMatrixLoc, 1, false, &projectionMatrix[0])

		// Set up for a glDrawElements call
		gl.BindVertexArray(vaos[trianglesName])
		gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, buffers[elementBufferName])

		if usePrimitiveRestart {
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

		// Swap Buffers
		gl.Flush()
		window.SwapBuffers()
		glfw.PollEvents()
	}

	// Cleanup
	for _, s := range shaders {
		s.Delete()
	}
	util.Terminate()
}

// TODO: wire callback in
func keyCallback(w *glfw.Window, key glfw.Key, scancode int, action glfw.Action, mods glfw.ModifierKey) {
	if action == glfw.Release && key == glfw.KeyM {
		usePrimitiveRestart = !usePrimitiveRestart
	}
}

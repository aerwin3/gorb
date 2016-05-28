// Example modified from OpenGL Programming Guide (Eighth Edition)
package main

import (
	"unsafe"

	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/glfw/v3.1/glfw"

	"github.com/hurricanerix/gorb/util"
)

func init() {
	// This example should always look here for it's resources.
	if err := util.SetWorkingDir("github.com/hurricanerix/gorb/01/triangles"); err != nil {
		panic(err)
	}
}

const ( // Program IDs
	trianglesProgID = iota
	numPrograms     = iota
)

const ( // VAO Names
	trianglesName = iota
	numVAOs       = iota
)

const ( // Buffer Names
	arrayBufferName = iota
	numBuffers      = iota
)

const ( // Attrib Locations
	mcVertexLoc = 0
)

var (
	programs    [numPrograms]uint32
	vaos        [numVAOs]uint32
	numVertices [numVAOs]int32
	buffers     [numBuffers]uint32
)

func main() {
	var err error

	// Get window context
	window, err := util.NewWindow("Ch1-Triangles", 512, 512)
	if err != nil {
		panic(err)
	}

	// Load the GLSL program
	shaders := []util.ShaderInfo{
		util.ShaderInfo{Type: gl.VERTEX_SHADER, Filename: "triangles.vert"},
		util.ShaderInfo{Type: gl.FRAGMENT_SHADER, Filename: "triangles.frag"},
	}
	programs[trianglesProgID], err = util.Load(&shaders)
	if err != nil {
		panic(err)
	}
	gl.UseProgram(programs[trianglesProgID])

	// Setup model to be rendered
	vertices := []float32{
		-0.90, -0.90, // Triangle 1
		0.85, -0.90,
		-0.90, 0.85,
		0.90, -0.85, // Triangle 2
		0.90, 0.90,
		-0.85, 0.90,
	}
	numVertices[trianglesName] = int32(len(vertices))

	gl.GenVertexArrays(numVAOs, &vaos[0])
	gl.BindVertexArray(vaos[trianglesName])

	sizeVertices := len(vertices) * int(unsafe.Sizeof(vertices[0]))
	gl.GenBuffers(numBuffers, &buffers[0])
	gl.BindBuffer(gl.ARRAY_BUFFER, buffers[arrayBufferName])
	gl.BufferData(gl.ARRAY_BUFFER, sizeVertices, gl.Ptr(vertices), gl.STATIC_DRAW)

	gl.VertexAttribPointer(mcVertexLoc, 2, gl.FLOAT, false, 0, gl.PtrOffset(0))
	gl.EnableVertexAttribArray(mcVertexLoc)

	// Main loop
	for !window.ShouldClose() {
		// Clear buffer
		gl.Clear(gl.COLOR_BUFFER_BIT)

		// Render
		gl.BindVertexArray(vaos[trianglesName])
		gl.DrawArrays(gl.TRIANGLES, 0, numVertices[trianglesName])

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

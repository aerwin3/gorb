// Example modified from OpenGL Programming Guide (Eighth Edition)
package main

import (
	"unsafe"

	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/glfw/v3.1/glfw"
	"github.com/hurricanerix/gorb/util"
)

func init() {
	if err := util.SetWorkingDir("github.com/hurricanerix/gorb/04/gouraud"); err != nil {
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
	mcColorLoc  = 1
)

type vertexData struct {
	Pos   [2]float32
	Color [4]uint8
}

var (
	programs    [numPrograms]uint32
	vaos        [numVAOs]uint32
	numVertices [numVAOs]int32
	buffers     [numBuffers]uint32
)

var (
	mode uint32
)

func main() {
	var err error

	// Get window context
	window, err := util.NewWindow("Ch4-Gouraud", 512, 512)
	if err != nil {
		panic(err)
	}

	// Load the GLSL program
	shaders := []util.ShaderInfo{
		util.ShaderInfo{Type: gl.VERTEX_SHADER, Filename: "gouraud.vert"},
		util.ShaderInfo{Type: gl.FRAGMENT_SHADER, Filename: "gouraud.frag"},
	}
	programs[trianglesProgID], err = util.Load(&shaders)
	if err != nil {
		panic(err)
	}
	gl.UseProgram(programs[trianglesProgID])

	// Setup model to be rendered
	vertices := []vertexData{
		vertexData{Pos: [2]float32{-0.90, -0.90}, Color: [4]uint8{255, 0, 0, 255}}, // Triangle 1
		vertexData{Pos: [2]float32{0.85, -0.90}, Color: [4]uint8{0, 255, 0, 255}},
		vertexData{Pos: [2]float32{-0.90, 0.85}, Color: [4]uint8{0, 0, 255, 255}},
		vertexData{Pos: [2]float32{0.90, -0.85}, Color: [4]uint8{10, 10, 10, 255}}, // Triangle 2
		vertexData{Pos: [2]float32{0.90, 0.90}, Color: [4]uint8{100, 100, 100, 255}},
		vertexData{Pos: [2]float32{-0.85, 0.90}, Color: [4]uint8{255, 255, 255, 255}},
	}
	numVertices[trianglesName] = int32(len(vertices))

	gl.GenVertexArrays(numVAOs, &vaos[0])
	gl.BindVertexArray(vaos[trianglesName])

	sizeVertexData := int(unsafe.Sizeof(vertices[0]))
	sizePos := int(unsafe.Sizeof(vertices[0].Pos))
	sizeVertices := len(vertices) * sizeVertexData
	gl.GenBuffers(numBuffers, &buffers[0])
	gl.BindBuffer(gl.ARRAY_BUFFER, buffers[arrayBufferName])
	gl.BufferData(gl.ARRAY_BUFFER, sizeVertices, gl.Ptr(vertices), gl.STATIC_DRAW)

	// TODO: Fix color of triangles, one should be red/green/blue, the other should be black/grey/white
	gl.VertexAttribPointer(mcVertexLoc, 2, gl.FLOAT, false, int32(sizeVertexData), gl.PtrOffset(0))
	gl.VertexAttribPointer(mcColorLoc, 4, gl.UNSIGNED_BYTE, true, int32(sizeVertexData), gl.PtrOffset(sizePos))

	gl.EnableVertexAttribArray(mcVertexLoc)
	gl.EnableVertexAttribArray(mcColorLoc)

	mode = gl.FILL

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
	var id uint32
	for i := 0; i < numPrograms; i++ {
		id = programs[i]
		gl.UseProgram(id)
		gl.DeleteProgram(id)
	}
	util.Terminate()
}

// TODO: fix callback
func keyCallback(w *glfw.Window, key glfw.Key, scancode int, action glfw.Action, mods glfw.ModifierKey) {
	if action == glfw.Release && key == glfw.KeyM {
		if mode == gl.FILL {
			mode = gl.LINE
		} else {
			mode = gl.FILL
		}
		gl.PolygonMode(gl.FRONT_AND_BACK, mode)
	}
}

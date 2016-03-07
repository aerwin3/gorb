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

var (
	mode uint32
)

type scene struct {
	Programs    [numPrograms]uint32
	VAOs        [numVAOs]uint32
	NumVertices [numVAOs]int32
	Buffers     [numBuffers]uint32
}

type vertexData struct {
	Pos   [2]float32
	Color [4]uint8
}

func (s *scene) Setup() error {
	shaders := []shader.Info{
		shader.Info{Type: gl.VERTEX_SHADER, Filename: "gouraud.vert"},
		shader.Info{Type: gl.FRAGMENT_SHADER, Filename: "gouraud.frag"},
	}

	program, err := shader.Load(&shaders)
	if err != nil {
		return err
	}
	s.Programs[trianglesProgID] = program

	gl.UseProgram(s.Programs[trianglesProgID])

	vertices := []vertexData{
		vertexData{Pos: [2]float32{-0.90, -0.90}, Color: [4]uint8{255, 0, 0, 255}}, // Triangle 1
		vertexData{Pos: [2]float32{0.85, -0.90}, Color: [4]uint8{0, 255, 0, 255}},
		vertexData{Pos: [2]float32{-0.90, 0.85}, Color: [4]uint8{0, 0, 255, 255}},
		vertexData{Pos: [2]float32{0.90, -0.85}, Color: [4]uint8{10, 10, 10, 255}}, // Triangle 2
		vertexData{Pos: [2]float32{0.90, 0.90}, Color: [4]uint8{100, 100, 100, 255}},
		vertexData{Pos: [2]float32{-0.85, 0.90}, Color: [4]uint8{255, 255, 255, 255}},
	}
	s.NumVertices[trianglesName] = int32(len(vertices))

	gl.GenVertexArrays(numVAOs, &s.VAOs[0])
	gl.BindVertexArray(s.VAOs[trianglesName])

	sizeVertexData := int(unsafe.Sizeof(vertices[0]))
	sizePos := int(unsafe.Sizeof(vertices[0].Pos))
	sizeVertices := len(vertices) * sizeVertexData
	gl.GenBuffers(numBuffers, &s.Buffers[0])
	gl.BindBuffer(gl.ARRAY_BUFFER, s.Buffers[arrayBufferName])
	gl.BufferData(gl.ARRAY_BUFFER, sizeVertices, gl.Ptr(vertices), gl.STATIC_DRAW)

	// TODO: Fix color of triangles, one should be red/green/blue, the other should be black/grey/white
	gl.VertexAttribPointer(mcVertexLoc, 2, gl.FLOAT, false, int32(sizeVertexData), gl.PtrOffset(0))
	gl.VertexAttribPointer(mcColorLoc, 4, gl.UNSIGNED_BYTE, true, int32(sizeVertexData), gl.PtrOffset(sizePos))

	gl.EnableVertexAttribArray(mcVertexLoc)
	gl.EnableVertexAttribArray(mcColorLoc)

	mode = gl.FILL

	return nil
}

func (s *scene) Update(dt float32) {
	// This is where you would put code to update your scene.
	// This scene does not change, so there is nothing here.
}

func (s *scene) Display() {
	gl.Clear(gl.COLOR_BUFFER_BIT)

	gl.BindVertexArray(s.VAOs[trianglesName])
	gl.DrawArrays(gl.TRIANGLES, 0, s.NumVertices[trianglesName])

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

func KeyCallback(w *glfw.Window, key glfw.Key, scancode int, action glfw.Action, mods glfw.ModifierKey) {
	if action == glfw.Release && key == glfw.KeyM {
		if mode == gl.FILL {
			mode = gl.LINE
		} else {
			mode = gl.FILL
		}
		gl.PolygonMode(gl.FRONT_AND_BACK, mode)
	}
}

// Main methods
func init() {
	if err := utils.SetWorkingDir("github.com/hurricanerix/gorb/04/gouraud"); err != nil {
		panic(err)
	}
}

func main() {
	c := app.Config{
		Name:                "Ch4-Gouraud",
		DefaultScreenWidth:  512,
		DefaultScreenHeight: 512,
		EscapeToQuit:        true,
		SupportedGLVers: []mgl32.Vec2{
			mgl32.Vec2{4, 3},
			mgl32.Vec2{4, 1},
		},
		KeyCallback: KeyCallback,
	}

	s := &scene{}

	a := app.New(c, s)
	if err := a.Run(); err != nil {
		panic(err)
	}
}

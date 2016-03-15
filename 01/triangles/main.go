// Example modified from OpenGL Programming Guide (Eighth Edition)
package main

import (
	"unsafe"

	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/mathgl/mgl32"
	"github.com/hurricanerix/go-gl-utils/app"
	"github.com/hurricanerix/go-gl-utils/path"
	"github.com/hurricanerix/go-gl-utils/shader"
)

func init() {
	if err := path.SetWorkingDir("github.com/hurricanerix/gorb/01/triangles"); err != nil {
		panic(err)
	}
}

func main() {
	c := app.Config{
		Name:                "Ch1-Triangles",
		DefaultScreenWidth:  512,
		DefaultScreenHeight: 512,
		EscapeToQuit:        true,
		SupportedGLVers: []mgl32.Vec2{
			mgl32.Vec2{4, 3},
			mgl32.Vec2{4, 1},
		},
	}

	s := &scene{}

	a := app.New(c, s)
	if err := a.Run(); err != nil {
		panic(err)
	}
}

// Everything below this line is for the Scene implementation.

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

type scene struct {
	Programs    [numPrograms]uint32
	VAOs        [numVAOs]uint32
	NumVertices [numVAOs]int32
	Buffers     [numBuffers]uint32
}

func (s *scene) Setup(ctx *app.Context) error {
	shaders := []shader.Info{
		shader.Info{Type: gl.VERTEX_SHADER, Filename: "triangles.vert"},
		shader.Info{Type: gl.FRAGMENT_SHADER, Filename: "triangles.frag"},
	}

	program, err := shader.Load(&shaders)
	if err != nil {
		return err
	}
	s.Programs[trianglesProgID] = program

	gl.UseProgram(s.Programs[trianglesProgID])

	vertices := []float32{
		-0.90, -0.90, // Triangle 1
		0.85, -0.90,
		-0.90, 0.85,
		0.90, -0.85, // Triangle 2
		0.90, 0.90,
		-0.85, 0.90,
	}
	s.NumVertices[trianglesName] = int32(len(vertices))

	gl.GenVertexArrays(numVAOs, &s.VAOs[0])
	gl.BindVertexArray(s.VAOs[trianglesName])

	sizeVertices := len(vertices) * int(unsafe.Sizeof(vertices[0]))
	gl.GenBuffers(numBuffers, &s.Buffers[0])
	gl.BindBuffer(gl.ARRAY_BUFFER, s.Buffers[arrayBufferName])
	gl.BufferData(gl.ARRAY_BUFFER, sizeVertices, gl.Ptr(vertices), gl.STATIC_DRAW)

	gl.VertexAttribPointer(mcVertexLoc, 2, gl.FLOAT, false, 0, gl.PtrOffset(0))
	gl.EnableVertexAttribArray(mcVertexLoc)

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
}

func (s *scene) Cleanup() {
	var id uint32
	for i := 0; i < numPrograms; i++ {
		id = s.Programs[i]
		gl.UseProgram(id)
		gl.DeleteProgram(id)
	}
}

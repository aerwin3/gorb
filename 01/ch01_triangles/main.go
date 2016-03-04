// Example modified from OpenGL Programming Guide (Eighth Edition)
package main

import (
	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/mathgl/mgl32"
	"github.com/hurricanerix/gorb/utils"
	"github.com/hurricanerix/gorb/utils/app"
	"github.com/hurricanerix/gorb/utils/shader"
)

const ( // VAO_IDs
	Triangles = iota
	NumVAOs   = iota
)

const ( // Buffer IDs
	ArrayBuffer = iota
	NumBuffers  = iota
)

const ( // Attrib IDs
	vPosition = 0
)

type scene struct {
	VAOs        [NumVAOs]uint32
	Buffers     [NumBuffers]uint32
	NumVertices int32
}

func (s *scene) Setup() error {
	shaders := []shader.Info{
		shader.Info{Type: gl.VERTEX_SHADER, Filename: "triangles.vert"},
		shader.Info{Type: gl.FRAGMENT_SHADER, Filename: "triangles.frag"},
	}

	program, err := shader.Load(&shaders)
	if err != nil {
		return err
	}

	gl.UseProgram(program)

	vertices := []float32{
		-0.90, -0.90, // Triangle 1
		0.85, -0.90,
		-0.90, 0.85,
		0.90, -0.85, // Triangle 2
		0.90, 0.90,
		-0.85, 0.90,
	}
	s.NumVertices = int32(len(vertices))

	gl.GenVertexArrays(NumVAOs, &s.VAOs[0])
	gl.BindVertexArray(s.VAOs[Triangles])

	gl.GenBuffers(NumBuffers, &s.Buffers[0])
	gl.BindBuffer(gl.ARRAY_BUFFER, s.Buffers[ArrayBuffer])
	gl.BufferData(gl.ARRAY_BUFFER, len(vertices)*4, gl.Ptr(vertices), gl.STATIC_DRAW)

	gl.VertexAttribPointer(vPosition, 2, gl.FLOAT, false, 0, gl.PtrOffset(0))
	gl.EnableVertexAttribArray(vPosition)

	gl.BindFragDataLocation(program, 0, gl.Str("fColor\x00"))
	return nil
}

func (s *scene) Update(dt float32) {
	// This is where you would put code to update your scene.
	// This scene does not change, so there is nothing here.
}

func (s *scene) Display() {
	gl.Clear(gl.COLOR_BUFFER_BIT)

	gl.BindVertexArray(s.VAOs[Triangles])
	gl.DrawArrays(gl.TRIANGLES, 0, s.NumVertices)

	gl.Flush()
}

func (s *scene) Cleanup() {
	// TODO: cleanup
}

// Main methods
func init() {
	if err := utils.SetWorkingDir("github.com/hurricanerix/gorb/01/ch01_triangles"); err != nil {
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

	s := &scene{}

	a := app.New(c, s)
	if err := a.Run(); err != nil {
		panic(err)
	}
}

// Modified from [lib/LoadShaders.cpp](http://www.opengl-redbook.com/Code/oglpg-8th-edition.zip)
// OpenGL Programming Guide (Eighth Edition)

package shader

import (
	"bytes"
	"fmt"
	"os"

	"github.com/go-gl/gl/v4.1-core/gl"
)

type Info struct {
	Type     uint32
	Filename string
	shader   uint32
}

func (i *Info) Compile(program uint32) error {
	i.shader = gl.CreateShader(i.Type)
	source, err := readShader(i.Filename)
	if err != nil {
		return err
	}

	csrc, free := gl.Strs(source)
	gl.ShaderSource(i.shader, 1, csrc, nil)
	free()
	gl.CompileShader(i.shader)

	var compiled int32
	if gl.GetShaderiv(i.shader, gl.COMPILE_STATUS, &compiled); compiled == gl.FALSE {
		return fmt.Errorf("failed to compile %s: %s", i.Filename, getErrorMsg(true, i.shader))
	}
	gl.AttachShader(program, i.shader)
	return nil
}

func readShader(filename string) (string, error) {
	f, err := os.Open(filename)
	defer f.Close()
	if err != nil {
		return "", err
	}

	buf := new(bytes.Buffer)
	buf.ReadFrom(f)
	source := fmt.Sprintf("%s\x00", buf.String())

	return source, nil
}

func (i *Info) Delete() {
	gl.DeleteShader(i.shader)
}

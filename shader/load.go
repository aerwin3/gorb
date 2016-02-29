// Modified from [lib/LoadShaders.cpp](http://www.opengl-redbook.com/Code/oglpg-8th-edition.zip)
// OpenGL Programming Guide (Eighth Edition)
package shader

import (
	"bytes"
	"fmt"
	"os"
	"runtime"
	"strings"

	"github.com/go-gl/gl/v4.1-core/gl"
)

func init() {
	// This is needed to arrange that main() runs on main thread.
	// See documentation for functions that are only allowed to be called from the main thread.
	runtime.LockOSThread()
}

type Info struct {
	Type     uint32
	Filename string
	shader   uint32
}

func Load(shaders *[]Info) (uint32, error) {
	program := gl.CreateProgram()

	for _, s := range *shaders {
		if err := s.compileShader(program); err != nil {
			cleanup(shaders)
			return 0, err
		}
	}

	//#ifdef GL_VERSION_4_1
	//if ( GLEW_VERSION_4_1 ) {
	//    // glProgramParameteri( program, GL_PROGRAM_SEPARABLE, GL_TRUE );
	//}
	//#endif /* GL_VERSION_4_1 */

	gl.LinkProgram(program)

	cleanup(shaders)

	var linked int32
	if gl.GetProgramiv(program, gl.LINK_STATUS, &linked); linked == gl.FALSE {
		return 0, fmt.Errorf("failed to link program: %s", getErrorMsg(false, program))
	}
	return program, nil
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

func (i *Info) compileShader(program uint32) error {
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

func cleanup(shaders *[]Info) {
	for _, s := range *shaders {
		if s.shader != 0 {
			gl.DeleteShader(s.shader)
		}
	}
}

func getErrorMsg(shader bool, id uint32) string {
	var l int32
	if shader {
		gl.GetShaderiv(id, gl.INFO_LOG_LENGTH, &l)
	} else {
		gl.GetProgramiv(id, gl.INFO_LOG_LENGTH, &l)
	}

	msg := strings.Repeat("\x00", int(l+1))

	if shader {
		gl.GetShaderInfoLog(id, l, nil, gl.Str(msg))
	} else {
		gl.GetProgramInfoLog(id, l, nil, gl.Str(msg))
	}

	return msg
}

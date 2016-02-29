// Modified from [lib/LoadShaders.cpp](http://www.opengl-redbook.com/Code/oglpg-8th-edition.zip)
// OpenGL Programming Guide (Eighth Edition)
package shader

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/go-gl/gl/v4.1-core/gl"
)

type Info struct {
	Type     uint32
	Filename string
	Shader   uint32
}

func readShader(filename string) (string, error) {
	log.Println("readShader:", filename)
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
	log.Println("(", i, ")compileShader:", program)
	i.Shader = gl.CreateShader(i.Type)
	source, err := readShader(i.Filename)
	if err != nil {
		return err
	}

	csrc, free := gl.Strs(source)
	gl.ShaderSource(i.Shader, 1, csrc, nil)
	free()
	gl.CompileShader(i.Shader)

	var compiled int32
	if gl.GetShaderiv(i.Shader, gl.COMPILE_STATUS, &compiled); compiled == gl.FALSE {
		// if debug {
		var logLength int32
		gl.GetShaderiv(i.Shader, gl.INFO_LOG_LENGTH, &logLength)
		msg := strings.Repeat("\x00", int(logLength+1))
		gl.GetShaderInfoLog(i.Shader, logLength, nil, gl.Str(msg))
		log.Println("ERROR:", msg)
		// }
		return fmt.Errorf("failed to compile %s", i.Filename)
	}

	gl.AttachShader(program, i.Shader)
	return nil
}

func cleanup(shaders *[]Info) {
	log.Println("cleanup:", shaders)
	for _, s := range *shaders {
		if s.Shader != 0 {
			//glDeleteShader( entry->shader );
			gl.DeleteShader(s.Shader)
		}
	}
}

func Load(shaders *[]Info) (uint32, error) {
	log.Println("Load:", shaders)

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
		// if debug {
		var logLength int32
		gl.GetProgramiv(program, gl.INFO_LOG_LENGTH, &logLength)
		msg := strings.Repeat("\x00", int(logLength+1))
		gl.GetProgramInfoLog(program, logLength, nil, gl.Str(msg))
		log.Println("ERROR:", msg)
		// }
		gl.DeleteProgram(program)
		return 0, fmt.Errorf("failed to link program")
	}

	return program, nil
}

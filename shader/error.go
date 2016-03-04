package shader

import (
	"strings"

	"github.com/go-gl/gl/v4.1-core/gl"
)

// getErrorsMsg helps to return an error message when compiling/linking goes
// wrong.  If shader is true, check for logs relating to a shader failing to
// compile.  In this context id should be the ID of the shader that failed to
// compile.  If shader is false, check for logs relating to linking shaters
// to a program.  in this context, id is the ID of the program that was
// attempting to link the shaders.
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

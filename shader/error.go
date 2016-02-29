// Modified from [lib/LoadShaders.cpp](http://www.opengl-redbook.com/Code/oglpg-8th-edition.zip)
// OpenGL Programming Guide (Eighth Edition)

package shader

import (
	"strings"

	"github.com/go-gl/gl/v4.1-core/gl"
)

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

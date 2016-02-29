// Modified from [lib/LoadShaders.cpp](http://www.opengl-redbook.com/Code/oglpg-8th-edition.zip)
// OpenGL Programming Guide (Eighth Edition)
package shader

import (
	"fmt"
	"runtime"

	"github.com/go-gl/gl/v4.1-core/gl"
)

func init() {
	// This is needed to arrange that main() runs on main thread.
	// See documentation for functions that are only allowed to be called from the main thread.
	runtime.LockOSThread()
}

func Load(shaders *[]Info) (uint32, error) {
	program := gl.CreateProgram()

	for _, s := range *shaders {
		if err := s.Compile(program); err != nil {
			cleanup(shaders)
			gl.DeleteProgram(program)
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
		gl.DeleteProgram(program)
		return 0, fmt.Errorf("failed to link program: %s", getErrorMsg(false, program))
	}

	return program, nil
}

func cleanup(shaders *[]Info) {
	for _, s := range *shaders {
		if s.shader != 0 {
			s.Delete()
		}
	}
}

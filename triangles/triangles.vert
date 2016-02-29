// Modified from Example 1.2: triangles.vert
// OpenGL Programming Guide (Eighth Edition)
/*
#version 430 core

layout (location = 0) in vec4 vPosition;

void
main()
{
	gl_Position = vPosition;
}
*/
#version 330

in vec4 vPosition;

void
main()
{
	gl_Position = vPosition;
}


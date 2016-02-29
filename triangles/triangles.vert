// Modified from Example 1.2: triangles.vert
// OpenGL Programming Guide (Eighth Edition)
#version 330 core
#extension GL_ARB_explicit_uniform_location : enable // Not needed for 430+

layout (location = 0) in vec4 vPosition;

void
main()
{
	gl_Position = vPosition;
}

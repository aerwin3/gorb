// Modified from OpenGL Programming Guide (Eighth Edition)
#version 410

layout( location = 0 ) in vec4 mcVertex;
layout( location = 1 ) in vec4 mcColor;

out vec4  color;

void
main()
{
    color = mcColor;
    gl_Position = mcVertex;
}

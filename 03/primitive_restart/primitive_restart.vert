// Modified from OpenGL Programming Guide (Eighth Edition)
#version 410

uniform mat4 modelMatrix;
uniform mat4 projectionMatrix;

layout (location = 0) in vec4 mcVertex;
layout (location = 1) in vec4 mcColor;

out vec4 vsColor;

void main(void)
{
    vsColor = mcColor;
    gl_Position = projectionMatrix * (modelMatrix * mcVertex);
}

// Modified from OpenGL Programming Guide (Eighth Edition)
#version 410

in vec4 vsColor;

layout (location = 0) out vec4 fragColor;

void main(void)
{
    fragColor = vsColor;
}

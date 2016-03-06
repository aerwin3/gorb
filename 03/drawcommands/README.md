Draw Commands
=============

Demonstrates draw commands by rendering four multi-color triangles in the window.

Notes
-----

#### Draw Commands

* ```gl.DrawArrays(mode uint32, first int32, count int32)```
[details](https://www.opengl.org/sdk/docs/man/html/glDrawArrays.xhtml)
* ```gl.DrawElements(mode uint32, count int32, xtype uint32, indices unsafe.Pointer)```
[details](https://www.opengl.org/sdk/docs/man/html/glDrawElements.xhtml)
* ```gl.DrawElementsBaseVertex(mode uint32, count int32, xtype uint32, indices unsafe.Pointer, basevertex int32)```
[details](https://www.opengl.org/sdk/docs/man/html/glDrawElementsBaseVertex.xhtml)
* ```gl.DrawArraysInstanced(mode uint32, first int32, count int32, instancecount int32)```
[details](https://www.opengl.org/sdk/docs/man/html/glDrawArraysInstanced.xhtml)

#### Other OpenGL funcs of interest

* ```gl.BufferSubData(target uint32, offset int, size int, data unsafe.Pointer)```
[details](https://www.opengl.org/sdk/docs/man/html/glBufferSubData.xhtml)
* ```gl.Disable(cap uint32)```
[details](https://www.opengl.org/sdk/docs/man/docbook4/xhtml/glEnable.xml)
* ```gl.Enable(cap uint32)```
[details](https://www.opengl.org/sdk/docs/man/html/glEnable.xhtml)
* ```gl.GetUniformLocation(program uint32, name *uint8) int32```
[details](https://www.opengl.org/sdk/docs/man/docbook4/xhtml/glGetUniformLocation.xml)
* ```gl.UniformMatrix4fv(location int32, count int32, transpose bool, value *float32)```
[details](https://www.opengl.org/sdk/docs/man/html/glUniform.xhtml)

Screenshot
----------

![Screenshot](screenshot.png)

Original Source
---------------

[OpenGL Programming Guide,  Eighth Edition](http://www.amazon.com/OpenGL-Programming-Guide-Official-Learning/dp/0321773039/)

* Example 3.5 Setting up for the Drawing Command Example, p. 122
* Example 3.6 Drawing Commands Example, p. 123

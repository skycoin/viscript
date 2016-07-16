# viscript

<<<<<<< HEAD

Renders a textured spinning cube using GLFW 3.1 and OpenGL 2.1.
=======
Feature requirement notes from Brandon:
>>>>>>> d43197e79f54ffc91b1e2abd162e19fdb24ec96f


<<<<<<< HEAD
![Screenshot](./Screenshot.png)

Debian
===

sudo apt-get install libgl1-mesa-dev
sudo apt-get install libxrandr-dev
sudo apt-get install libxcursor-dev
sudo apt-get install libxinerama-dev

go get github.com/go-gl/gl/v{3.2,3.3,4.1,4.4,4.5}-{core,compatibility}/gl
go get github.com/go-gl/gl/v3.3-core/gl
=======
"
I think we would start with opengl in golang and creating a window, that has terminal characters and grid of characters. Then create a list of objects (software) and each object has a list of actions on it.
Then have a scripting language. like one action is "execute for one line" or "execute until stop" etc. we would have a list of all software objects and a program being executed is a software object


we are going to make a little scripting language, where you edit the abstract syntax tree of the language directly
 
[6/14/2016 10:33:07 PM] Brandon: then we will send the key presses, mouse clicks, scroll wheel as messages (length prefixed messages) over a golang channel
[6/14/2016 10:33:19 PM] Brandon: then the application on other side, will respond with length prefixed messages, to set the display


I want the screen to respond to length prefixed messages (32 bit length prefix, followed by binary). and I should be able to get the dimensions in characters and be able to put characters on screen.
later, i want pixels and to be able to create subwindow or 2d plane, and then to blit it to screen and do opengl operations, from the scripting language; but we will do that later


[7/10/2016 6:20:37 AM] Brandon: 1> terminal, with terminal program handling opengl and input (mouse click, scroll, left click right click) and sending the messages to another program over a channel as lenght prefixed messages
2> a simple lisp or C like scripting language on the other end of the terminal
3> tools for scripting language, autocomplete etc

4> extend it to have opengl support
3> An audocad or video game like application in the scripting language where you can create shapes and draw them and apply operations to them
[7/10/2016 6:56:40 AM] Brandon: tell me when you are ready to stat on scripting language


[7/13/2016 3:32:51 PM] Brandon: like an Ssh terminal
[7/13/2016 3:33:15 PM] Brandon: and taking mouse scroll up, mouse scroll down, left click, right click and key presses and serializing them and sending them over an event channel


[7/14/2016 11:39:21 AM] Brandon: i need grid of characters, like a Ssh terminal array
[7/14/2016 11:39:40 AM] Brandon: and all events serialized, through a channel, to be handled by another application (the same application)
[7/14/2016 11:39:47 AM] Brandon: its cross platform terminal in opengl
"
>>>>>>> d43197e79f54ffc91b1e2abd162e19fdb24ec96f

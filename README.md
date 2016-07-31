# viscript


Debian
===
=======

sudo apt-get install libgl1-mesa-dev
sudo apt-get install libxrandr-dev
sudo apt-get install libxcursor-dev
sudo apt-get install libxinerama-dev

go get github.com/go-gl/gl/v{3.2,3.3,4.1,4.4,4.5}-{core,compatibility}/gl
go get github.com/go-gl/gl/v3.3-core/gl
=======






Feature requirement snippets from HaltingState:

"
I think we would start with opengl in golang and creating a window, that has terminal characters and grid of characters. Then create a list of objects (software) and each object has a list of actions on it.
Then have a scripting language. like one action is "execute for one line" or "execute until stop" etc. we would have a list of all software objects and a program being executed is a software object


we are going to make a little scripting language, where you edit the abstract syntax tree of the language directly
 
[6/14/2016 10:33:07 PM] HaltingState: then we will send the key presses, mouse clicks, scroll wheel as messages (length prefixed messages) over a golang channel
[6/14/2016 10:33:19 PM] HaltingState: then the application on other side, will respond with length prefixed messages, to set the display


I want the screen to respond to length prefixed messages (32 bit length prefix, followed by binary). and I should be able to get the dimensions in characters and be able to put characters on screen.
later, i want pixels and to be able to create subwindow or 2d plane, and then to blit it to screen and do opengl operations, from the scripting language; but we will do that later


[7/10/2016 6:20:37 AM] HaltingState: 1> terminal, with terminal program handling opengl and input (mouse click, scroll, left click right click) and sending the messages to another program over a channel as length prefixed messages
2> a simple lisp or C like scripting language on the other end of the terminal
3> tools for scripting language, autocomplete etc

4> extend it to have opengl support
3> An audocad or video game like application in the scripting language where you can create shapes and draw them and apply operations to them
[7/10/2016 6:56:40 AM] HaltingState: tell me when you are ready to stat on scripting language


[7/13/2016 3:32:51 PM] HaltingState: like an Ssh terminal
[7/13/2016 3:33:15 PM] HaltingState: and taking mouse scroll up, mouse scroll down, left click, right click and key presses and serializing them and sending them over an event channel


[7/14/2016 11:39:21 AM] HaltingState: i need grid of characters, like a Ssh terminal array
[7/14/2016 11:39:40 AM] HaltingState: and all events serialized, through a channel, to be handled by another application (the same application)
[7/14/2016 11:39:47 AM] HaltingState: its cross platform terminal in opengl



[7/18/2016 4:32:28 AM] HaltingState: create file called input.go and each key or button press will be sent over
[7/18/2016 4:32:50 AM] HaltingState: then the program will send back length prefixed commands, like "set character" or ask for size of display etc
[7/18/2016 4:33:21 AM] HaltingState: and we will have a program, that only takes length prefixed messages (key inputs) and sends back length prefixed messages (setting display, set characters etc)
[7/18/2016 4:34:38 AM] HaltingState: and then we will start on a simple programming language; like C, but will use abstract syntax tree. will just have structs, functions, int32, byte array; very basic

the programming language will take in length prefixed messages, respond to them and then emit length prefixed messages

Then we will add "modules" which are collections of structs and functions and you can import a module into another module
[7/18/2016 4:44:06 AM] HaltingState: a module is a struct, with a list of struct signatures (the structs in that module) and a list of function signatures (the functions in that module)
[7/18/2016 4:44:49 AM] HaltingState: a function is a struct; a struct for signature (type input list, type output list) and an array of expressions

"

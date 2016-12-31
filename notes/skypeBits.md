Feature requirement snippets from HaltingState (through Skype):

"
create a list of objects, and each has a list of actions on it.
Then have a scripting language. like one action is "execute for one line" or "execute until stop" etc. we would have a list of all software objects and a program being executed is a software object


we are going to make a little scripting language, where you edit the abstract syntax tree of the language directly
 
[6/14/2016 10:33:19 PM] HaltingState: send events, then the application on other side, will respond with length prefixed messages, to set the display


I want the screen to respond to messages.  I should be able to get the dimensions in characters and be able to put characters on screen.
later, i want pixels and to be able to create subwindow or 2d plane, and then to blit it to screen and do opengl operations, from the scripting language


[7/10/2016 6:20:37 AM] HaltingState: 1> terminal program handling opengl and input and sending the messages to another program.
2> a simple lisp or C like scripting language on the other end of the terminal
3> tools for scripting language, autocomplete etc

3> An audocad or video game like application in the scripting language where you can create shapes and draw them and apply operations to them


[7/13/2016 3:32:51 PM] HaltingState: like an Ssh terminal


[7/14/2016 11:39:47 AM] HaltingState: its cross platform terminal in opengl



[7/18/2016 4:32:28 AM] HaltingState: each input event will be sent over
[7/18/2016 4:32:50 AM] HaltingState: then the program will send back length prefixed commands, like "set character" or ask for size of display etc
[7/18/2016 4:33:21 AM] HaltingState: and we will have a program, that only takes input eventss and sends back length prefixed messages (setting display, set characters etc)
[7/18/2016 4:34:38 AM] HaltingState: simple programming language; like C, but will use abstract syntax tree. will just have structs, functions, int32, byte array; very basic

the programming language will take in length prefixed messages, respond to them and then emit length prefixed messages

Then we will add "modules" which are collections of structs and functions and you can import a module into another module
[7/18/2016 4:44:06 AM] HaltingState: a module is a struct, with a list of struct signatures (the structs in that module) and a list of function signatures (the functions in that module)
[7/18/2016 4:44:49 AM] HaltingState: a function is a struct; a struct for signature (type input list, type output list) and an array of expressions




func() {} (), is an anonymous function. that is being called

do window resizing and send length prefixed message to other end, when window is resized


We need a few base types, like "int" which is int32, i32, u32, i64, u64 (uint 64), []byte. []byte (all array types are fixed length).

A set of operators (add_i32, sub_i32, mult_i32, div_i32) etc

then we need
- ability to create structs (new types) from existing types
- functions (a type type in, type tuple out, an array of statements or blocks (a block is an array of statements)

The language will be all structs and everything will be structs and abstract syntax trees. the written part, will be auto correction.

Look at XL programming language. 

https://en.wikipedia.org/wiki/XL_(programming_language)

and do S syntax syntax. 
- look at current thing or choices
- then look at what valid choices are for context, then do dropdown to allow selection or autocompletion
- all choices or actions mut be enumerable at each step


[8/2/2016 3:38:18 AM] HaltingState: so bascly, statically typed scheme for now
[8/2/2016 2:47:57 PM] HaltingState: ???
[8/2/2016 2:48:11 PM] HaltingState: scheme should be simple and we can implement a few commands, then implement everything else in terms of that


[8/3/2016 11:09:02 AM] HaltingState: just text for now
[8/3/2016 11:09:19 AM] HaltingState: just statically typed scheme
[8/3/2016 11:09:41 AM] HaltingState: the language turns in on itself, recursively, eventually
[8/3/2016 6:53:35 PM] HaltingState: you know what  scheme is right?
[8/3/2016 6:54:44 PM] HaltingState: I want
1> structs
2> functions
3> static typing
4> a small number of atomic types and operations

So do (u32.add 3 5) (add 3+5), etc, each expression should be typed (inputs and return type) and have it be a tuple in, tuple out
(u32.div 5 15)
[8/3/2016 6:55:19 PM] HaltingState: func F1(x int) (int) {
line 1
line 2
} etc
[8/3/2016 6:56:25 PM] HaltingState: a function is a 
1> name (text)
2> input tuple (name, type) pair list
3> output tuple (list of types/signatures returned)
4> an array of lines/statements
[8/3/2016 6:56:51 PM] HaltingState: a struct is a 
1> name
2> list of (name, type) pairs
3> later list of functions on the struct but ignore this for now
[8/3/2016 7:09:18 PM] HaltingState: ok; i pushed better instructions; towards bottom
[8/3/2016 7:09:51 PM] HaltingState: everything will be S-expressions like lisp, but statically typed. For now

But in IDE and editor, where you are typing, it will look like normal language (like golang) and have autocomplete



we will have an internal language, which is the s-expressions and parens and (func a b) and then have it look like C to user in IDE, but will just start with S expression and autocomplete



[8/7/2016 1:07:21 PM] HaltingState: lets do some simple types like uint32
[8/7/2016 1:07:26 PM] HaltingState: and do add, sub, new variable etc
[8/7/2016 1:07:33 PM] HaltingState: (add32, a, b)
[8/7/2016 1:07:38 PM] HaltingState: then lets create functions
[8/7/2016 1:07:51 PM] HaltingState: and maybe add functions and structs and some new atomics



[8/9/2016 4:48:21 AM] HaltingState: the most important things are all reflection
[8/9/2016 4:48:29 AM] HaltingState: like "What functions are in this module"
[8/9/2016 4:48:48 AM] HaltingState: "What elements/types/names are in this struct"
[8/9/2016 4:49:03 AM] HaltingState: "what variables are in the local scope"
[8/9/2016 4:49:13 AM] HaltingState: "What actions can I do on this type, in this context"
[8/9/2016 5:40:01 AM] HaltingState: scheme is a very simple language really
[8/9/2016 10:45:08 PM] HaltingState: do you know lisp style?
[8/9/2016 10:45:20 PM] HaltingState: (add a b)
[8/9/2016 10:45:32 PM] HaltingState: (sub (add 1 3) 4)
[8/9/2016 10:45:36 PM] HaltingState: (div 1 4)



[8/10/2016 9:38:40 AM] HaltingState: it will be golang like
[8/10/2016 9:38:46 AM] HaltingState: the back end is trees and S expression
[8/10/2016 9:38:56 AM] HaltingState: but to user and editor, it will displayed as golang type
[8/10/2016 9:41:42 AM] HaltingState: create an "object" type for embedding in the document. These are like software objects or code.  Then for (add32, x, 3), have a type that is like a list and have element1, element2, element3; store the type, and "tag" 

like 3 is int32 literal (or literal and then store type and value in the literal type sruct)
add32 is func or operation 
x is int32 variable or variable and store the type in variable

and (add32 5 (add32 4 5))

the (add32 4 5) is an "Expression"
[8/10/2016 9:42:21 AM] HaltingState: and have a value for "unknown" and as you are typing, enumerate in a dropdown all the possibilities for what can go there or what you can do
[8/10/2016 9:42:52 AM] HaltingState: so arrow keys wont select the type directly, but will select the object and go from field to field.
[8/10/2016 9:43:18 AM] HaltingState: also color code each type differently so (add32 x 3), the add32, will be in different color
[8/10/2016 9:44:27 AM] HaltingState: but you can write a parser at start, for a text language; but the source code really should be a software object, that has a representation as text on screen; but you are applying operations to the software object itself
[8/10/2016 9:44:58 AM] HaltingState: (var label type), (var X int32)  is define variable


[8/10/2016 9:25:37 PM] HaltingState: the idea is that there are a finite set of actions at each point, you can apply to the program (finite action selection). Like a video game
[8/10/2016 9:27:10 PM] HaltingState: So you at at a node,with the node selected, then it sees (add32 X 3) and you have X selected and then pop down menu and it knows what can go in x like "expression, literal, variable) and then if you hit variable, then it would look up all the variables in the local scope and only show the int32 and look up the variables in the current module and look up variables that are inputs to the functions (checks the variables in scope and filters them)

etc
[8/10/2016 9:27:30 PM] HaltingState: and so the computer program is actually a data structure, in memory (like a program), made out of structs and objects) and you are applying functions to it, to modify it
[8/10/2016 9:27:41 PM] HaltingState: there is a "line editor" and it looks like terminal or text application
[8/10/2016 9:28:08 PM] HaltingState: You can ignore this for now and just do a text based programming language, get the basic one working. Then we will build the advanced language on top of that
[8/10/2016 9:29:08 PM] HaltingState: so when you have a module selected, you can add function to the module, add new structs or add a variable to the global scope of the module and the module is a struct object
[8/10/2016 9:29:58 PM] HaltingState: and when you have a function object selected , you can edit its name (its tag, it will have a description as text string, but will have a UUID that is used internally; so you can change its name without the UUID chanigng and everywhere it is called, the name will change automatically, because its really referencing a UUID)
[8/10/2016 9:30:24 PM] HaltingState: and when the function object is selected, you can edit its name, add/remove lines, edit lines, change its type signature etc
[8/10/2016 9:30:51 PM] HaltingState: and when a struct object is selected, you can add new fields, remove fields, or add/remove functions to the struct
[8/10/2016 9:32:19 PM] HaltingState: so it is sort of like a "programming language video game", because it would show children at each step, everything they can do to the program, at the current cursor position, or what object is selected. And will tell you what actiosn you can do on the object.

And when you mouse over an object, it will tell you its type and information.
[8/10/2016 9:32:42 PM] HaltingState: this is long term stuff; but good to think about, but probably wont be there for years
[8/10/2016 9:34:10 PM] HaltingState: we can start with a few basic types like int32 and then do some operations and have a text based programming language; that looks like lisp. Have structs and funcs, like golang has.

Then from there, we can build this more advanced language, on top of the structs and funcs, of the base language
[8/10/2016 9:34:22 PM] HaltingState: also, look at iPython notebook
[8/10/2016 9:39:37 PM] HaltingState: there is also a command called (eval) and eval takes in an expression and it executes it. So if you have a function and it is N lines, then when you run the function, a loop will just go line by line, calling eval on each line


http://ergoemacs.org/emacs/i/emacs_lisp_interactive_command_line_interface_REPL_ielm.png


https://en.wikipedia.org/wiki/Lisp_(programming_language)#Conses_and_lists


[8/10/2016 9:43:29 PM] HaltingState: this is how lisp does it, but not sure  we need this, but might help, but not sure if its best way. I like the struct, func, module, lines as basic units. then each line is an "expression"
[8/10/2016 9:47:40 PM] HaltingState: lisp does not have "objects". I want to focus on objects (essentially structs) and focus on the functions of (things that dont modify) and functions on (things that modify) the objects. 

So that the modules are object, the structs in the modules are objects, the functions in the modules are objects and the expressions in the functions are objects


https://upload.wikimedia.org/wikipedia/commons/a/af/IPython-notebook.png

http://2.bp.blogspot.com/-W-mMDGXDLiQ/UM0fIZEJwwI/AAAAAAAAAYc/khS0LdI2UI4/s1600/notebook_shot.png



[8/10/2016 9:53:21 PM] HaltingState: we need 3 modes. first mode is like editor mode and for code and browsing and editing lines etc.

Then one mode is like command like REPL mode which is just boring

The ipython mode is probably best. it gives you a box and you put code in box and under the box it runs it
[8/10/2016 9:53:46 PM] HaltingState: and you can have a notebook page, which has a bunch of boxes and can run code in them and shows code examples


http://ipython.org/_static/screenshots/ipython-notebook.png



[8/10/2016 9:59:12 PM] HaltingState: this language CX, is also for controlling my robot swarms and CNC machines, because I was pissed at how difficult it was and how shitty software was and how much stuff i had to install
[8/10/2016 9:59:45 PM] HaltingState: and will be for video games also and have native support for component systems



http://www.idryman.org/images/ipython_notebook.png

http://www2.cisl.ucar.edu/sites/default/files/users/bjsmith/ipython-notebook_0.png


[8/10/2016 10:20:54 PM] HaltingState: we need an "Action button" or something, so we can insert code boxes like this into the terminal or do actions on an object
[8/10/2016 10:21:09 PM] HaltingState: and then insert code box, then type some stuff and then runi t and see result in box below the code box
[8/10/2016 10:21:23 PM] HaltingState: and scroll through the document, and looking at the code boxes and results
[8/10/2016 10:21:35 PM] HaltingState: and so we have software objects embedded in the text document essentially



[8/11/2016 11:33:07 PM] HaltingState: i like bash promp or python notebook style for prototyping and testing and runnign small snippets (like python)
[8/11/2016 11:33:25 PM] HaltingState: but then functions, structs, modules; are good for libraries and stuff that does not break
[8/11/2016 11:34:08 PM] HaltingState: and meta programming (ability to write program that creates new functions in the module etc) is used for video game programming and "templates" and generating code for other things, using existing code
[8/11/2016 11:34:46 PM] HaltingState: because a template is just a program that takes in some parameters and then outputs code that calls the module and creates new structs and new functions, (such as object list template takes in an object type and outputs the code for that object)
[8/11/2016 11:35:14 PM] HaltingState: I have some more documentation that is important and will add to readme later


[8/12/2016 4:02:32 AM] HaltingState: by editing software objects or embedding software objects, what I mean is. Look at the text document. It can either be an array of lines and characters on each line.

or it can be an array, that can contain structs or types. And the text lines are actually a "text lines" type
[8/12/2016 4:03:03 AM] HaltingState: but ignore that right now, too complex; expecially because it has to be reverse draw back to the grid of characters


[8/22/2016 10:06:54 PM] HaltingState: 'i haven't the foggiest idea how we would do pointers, if that's something you want.' Me either
[8/22/2016 10:07:50 PM] HaltingState: I think we have "functions of" and "functions on" an object, but i have no idea what a reference to an object is; a function of an object cannot change it; so if you copied whole object, it would be same as a pointer because you would get same result, so does not matter
[8/22/2016 10:08:34 PM] HaltingState: and for a "function on" an object it changes the object, so you either perform the function on the object itself, or a copy of it and have to choose
[8/22/2016 10:11:06 PM] HaltingState: I actually think, that each object has a UUID and its sort of like a "pointer", that identifies the object. Internally it is a pointer because all objects sit in memory



[8/23/2016 5:10:27 AM] HaltingState: looks good so far
[8/23/2016 5:11:06 AM] HaltingState: console log is good; can use console to toggle displaying mouse updates on and off etc
[8/23/2016 11:35:48 AM] HaltingState: ok how about this
[8/23/2016 11:38:52 AM] HaltingState: for each object
1> first 32 bytes is object type
2> next 32 bytes is the object "UUID"
[8/23/2016 11:40:12 AM] HaltingState: when you malloc something
1> 4 byte length, then data (as malloc overhead header)
[8/23/2016 11:46:11 AM] HaltingState: pointers and objects are difficult.
[8/23/2016 11:52:52 AM] HaltingState: when an object contains another object as a member and you want a reference or pointer to the object contained then it gets complex
[8/23/2016 11:53:05 AM] HaltingState: I think we can do what C does for now and keep it same and keep more complex "objects" as higher level than the base
[8/23/2016 11:55:07 AM] HaltingState: a pointer is just an "int" that references memory and program state or memory is just array of bytes, laid out one by one. then pointer is an index into that. A pointer needs a type so you can look up the object's members and index.

So if you have struct A {int x, int y}, then you do A.y, then it looks at the object metadata and looks at what offset for y is and then adds the pointer to the offset and casts to int; i guess
[8/23/2016 11:55:27 AM] HaltingState: another important thing or very important thing; I need to be able to see all the objects in the program.
[8/23/2016 11:56:05 AM] HaltingState: 
so all the variables that have been allocated in the program, all the types. I want to be able to browse through them and have an "object browser" and look at the programs state and interograte it
[8/23/2016 12:01:06 PM] HaltingState: using notation (func <NAME> (<TYPE> <NAME>) .... code)

have an "Expression class" which the () thing and have (op in1 in2) and have an op for introducing variable at each scope. op for defining new struct, op for defining a function.

ex.

A func def op, takes NAME, then array of (TYPE, NAME) for inputs and then array of (TYPE) for outputs for functions.

A pointer type is struct POINTER { type TYPE, index uint32)
[8/23/2016 12:05:55 PM] HaltingState: except, i want to be able to do 

(VAR X INT32) //defines variables
(uint32.add X 3) //next line,expression, that uses a variable defined in local scope

Then go line by line, executing but maybe write the lines in infix
[8/23/2016 12:06:22 PM] HaltingState: or just do 

(+ X 3) and that is shorthand
[8/23/2016 12:10:38 PM] HaltingState: So like BASIC (Where it executes line 1, then line 2, line 3 etc and one expression per line

But made the expressions in scheme like infix notation.

define a macro or somethign so you can write it in itself

then add functions, modules etc to make it structurally like golang
[8/23/2016 12:10:51 PM] HaltingState: I think we will figure it out
[8/23/2016 12:11:17 PM] HaltingState: S notation (The parents and doing (OP var1 var2) is allot easier to parser than doing the text and making it look like golang
[8/23/2016 12:11:24 PM] HaltingState: look at scheme or lisp

http://zhehaomao.com/project/2014/03/28/glisp.html

[8/23/2016 12:23:30 PM] HaltingState: except I want
1> functions (a function has a name, a list of inputs that are statically typed and return type, and a list of expressions; one expression per line)
[8/23/2016 12:23:37 PM] HaltingState: 2> structs
[8/23/2016 12:23:42 PM] HaltingState: 3> statically typed
[8/23/2016 12:23:48 PM] HaltingState: look at lisp implementations in golang

https://pkelchte.wordpress.com/2013/12/31/scm-go/

https://gist.github.com/pkelchte/c2bd76b9f8f9cd603b3c

http://norvig.com/lispy.html

https://gist.github.com/pkelchte/c2bd76b9f8f9cd603b3c

like this, but statically typed and with structs and funcs


[8/24/2016 1:24:58 AM] HaltingState: 1> like basic (line one, line two, line thee, one expression per line
2> expressions like C in terms of line introducing variable or assigning variable etc
3> statically typed like C (constraints on expressions that must be satisifed, ability to enumerate the valid set of things you could put into a place), getting set of variables that are in local scope, enumerating set of variables in global scope
4> etc
[8/24/2016 10:53:32 AM] HaltingState: notice how small the lisp interpreter is? like 800 lines or 200 lines
[8/24/2016 10:53:49 AM] HaltingState: once you have the language working. then you can write everything in the language itself as a macro



[8/28/2016 7:09:04 AM] HaltingState: ideally, there should be a runtime process and then it should accept messages and then should emit messages (for rendering etc)
[8/28/2016 7:09:13 AM] HaltingState: and eventually the thing running the console will be in CX itself
[8/28/2016 7:09:33 AM] HaltingState: also look at DOM or "Virtual DOM". a DOM is a  datastructure, and then there is a function that renders it
[8/28/2016 7:09:52 AM] HaltingState: virtual DOM is where you create the DOM in javascript and build it out as a software object
[8/28/2016 7:10:11 AM] HaltingState: for instance, a DOM element might be a code box, where you can insert code and run it


https://github.com/Matt-Esch/virtual-dom



Me: i'm not sure what "emit messages " means



http://stackoverflow.com/questions/21109361/why-is-reacts-concept-of-virtual-dom-said-to-be-more-performant-than-dirty-mode

https://github.com/teropa/angular-virtual-dom

http://stackoverflow.com/questions/21965738/what-is-virtual-dom


[8/28/2016 7:11:58 AM] HaltingState: "virtual DOM" is just a DOM made out of structs; so you can do thing like, goto a DOM node and then do "add code box" and it will insert an element into the document, that you can add code to and run from the box
[8/28/2016 7:12:12 AM] HaltingState: or a DOM for a document, might be a list of paragraphs of text and each paragraph is one DOM element
[8/28/2016 7:12:32 AM] HaltingState: then you have a DOM renderer, that maps the DOM to the characters on screen and window and how it looks
[8/28/2016 7:12:58 AM] HaltingState: and DOM elements can have names or tags and then you have a CSS type language that gives hints on how to display them
[8/28/2016 7:13:12 AM] HaltingState: but ignore that for now; maybe virtual DOM will be in CX eventually and not golang
[8/28/2016 7:13:27 AM] HaltingState: golang should just be base layer to build up a good language for this maybe
[8/28/2016 7:14:24 AM] HaltingState: or one line of text or code can be a DOM element



[8/29/2016 6:30:06 AM] HaltingState: eventually may have multiple "windows" as an abstraction you might do; 

But windows are really recursive. there is a window and then another window in the windows and its windows all the way down; because each box or grid, is contained within the parent
[8/29/2016 7:11:08 AM] HaltingState: A "DOM" is a datastructure. that represents the document. So you can have a "function" DOM element or a "struct" DOM element, or a "Code box" DOM element, or a "text paragraph" DOM element
[8/29/2016 7:11:55 AM] HaltingState: and there needs to be an "active object" that you have selected (what the cursor is on), then you need to have an "action key" and then dropbox box for what actions you can take (add dom element, enter element , run element etc).



https://media.8ch.net/file_store/182f81c50b3263c84c2269056bde7221026adc7971d119b1c04062c26f81c7c9.png



[9/24/2016 9:26:18 AM] HaltingState: you can make the text window bigger and even resizable if you want
[9/24/2016 9:32:25 AM] HaltingState: I like lisp style (add32 x y) or (add32 (add32 3 5) 8) and color coding for (operator, int32 literal, etc) and mouse over for type
[9/24/2016 9:34:28 AM] HaltingState: things do not need a representation as text, but can be data (that has way of being rendered as text and a way for it to be interacted with)

[12:27:36 PM] HaltingState: make each window a program
[12:27:50 PM] HaltingState: eventually each window or pane will be driven by a program, that accepts/receives length prefixed messages
[12:28:04 PM] HaltingState: and if you have 3 panes inside of a bigger pane, they will pass the messages down to lower level
[12:28:22 PM] HaltingState: eventually the programs handling the panes will be in the scripting language you are implementing
[12:30:22 PM] HaltingState: definately look at scheme as a stating point
[12:31:22 PM] HaltingState: also, label line numbers in the editor; each line needs to have line numbers on left hand side, like sublime






[9/27/2016 7:05:12 AM] HaltingState: 
1> All code must be in a function
2> you call function
[9/27/2016 7:07:42 AM] HaltingState: 
3> functions have line 1, line 2, line 3, line 4 etc.
4> Functions have "blocks" which are like "if X, block Y else Block Y" the blocks are embedded recursively. So there is a top block, then subblocks. 
5> each line is an expression. start with basic things like "introduce variable X" and "set x value" and (set x (add32 5 8))   , (new x int32), or (int32.new x) //introduce type int32, with label x






[9/30/2016 3:00:31 AM] HaltingState: yes
[9/30/2016 3:01:00 AM] HaltingState: so (* 3 (+ 5 6)) instead of 3* (5+6) so you do not have to do any work to parse and no order of precidence
[9/30/2016 3:01:13 AM] HaltingState: and it is clear exactly what is being called and no ambiguity
[9/30/2016 3:02:17 AM] HaltingState: however, for coding, we can display at it as 3*(5+6) for the programmer but on backend it is (*3 (+ 5 6)) or (mult32 3 (add32 5 6)), or (int32.mult 3 ( int32.add 5 6))
[9/30/2016 3:02:40 AM] HaltingState: right now, we are just doing the backend and the abstract syntax tree;
[9/30/2016 3:26:53 AM] HaltingState: do
1. a command for introducing a variable into scope
2. a command for setting the value of a variable
3. int32 type and add, mult, div, sub, int32.add, int32.sub, int32.mult etc is ok/good. sort of like assembly
4. function for creating a new module
5. function for creating a new function (Adding a function to a module)
6. function for adding a new type (struct, to a module)
[9/30/2016 3:27:17 AM] HaltingState: the program has to be built up, from calling the commands to build up the program
[9/30/2016 3:27:35 AM] HaltingState: and then being able to call functions
[9/30/2016 3:27:48 AM] HaltingState: and a function is a series of lines
[9/30/2016 3:29:44 AM] HaltingState: you must do the Citrulline Malate, for ~3 months, to get rid of aluminum and then will have more energy
[9/30/2016 3:44:40 AM] HaltingState: so to start
1> define a function (input, output)
2> add lines to the tuple (Expressions) and keep the expressions simple and types minumum
3> be able to run the function

being able to create structs etc and just start simple
[9/30/2016 3:45:21 AM] HaltingState: look up "abstract syntax tree", that is what it is
[9/30/2016 5:26:08 AM] HaltingState: also, hurry, have customers who need it
[9/30/2016 10:52:14 AM] HaltingState: look up "Abstract Syntax Tree" , that is what it is (not even lisp)
[9/30/2016 11:11:33 AM] HaltingState: https://en.wikipedia.org/wiki/Abstract_syntax_tree
[9/30/2016 11:11:39 AM] HaltingState: that what it really is



[7:05:21 AM] HaltingState: there are different types of stylign an rendering, such as horizontal scroll on current line, or only expanding to multiple lines for the selection
[7:05:52 AM] HaltingState: and we should have a DOM or virtual DOM, which is virtual document object model, like react has and then have an object that renders the DOM and have the options in that object

But doesnt matter for now
[7:19:28 AM] HaltingState: A virtual DOM, then an object that renders it into the square buffer





[10/11/2016 5:10:54 AM] HaltingState: go on rise up and create an etherpad
[10/11/2016 5:11:50 AM] HaltingState: https://pad.riseup.net/
[10/11/2016 6:37:45 AM] HaltingState: write down what you think it is
[10/11/2016 6:37:56 AM] HaltingState: and then i will make notes or reply
[10/11/2016 6:40:40 AM] HaltingState: basicly, you are implementing 
1> a very simply "base language"
2> the more complicated things will be in the base language itself, written in terms of itself, instead of golang

you need to have an atom, an operation, a piece of data. Then a rule for when it can be used (A context is fed in) and there will be a slot (ex for a data item in an operand) and there will be a rule for what can occupy that.

So if you add (int32.add X 3), then X int32.add may say "X has to be an int32 or an unknown and then valid X would be a list of (literal, or int32 in local scope, or module.var which is an int 32 etc) and we can enumerate exhaustively, everything valid that can go in that slot. People could select it from a dropdown

The environment is current line, current function, current module; then we can look at vars in the local scope and can look at vars at the module scope and can look at vars passed into the function and see if any of them matched
[10/11/2016 6:43:01 AM] HaltingState: There is an idea of "enumeration" (ability to list all possibilties and it must be finite) (Similar idea to counting).

We build up programs by applying operations to them.

So we do
1> add new line operator
2> Add operator (int32.add X Y)
3> substitute or set X, or select X and see options for thigns that could go in there (a function that returns int32, an int32 literal, a function that casts something else to int, an int32 at local scope etc)

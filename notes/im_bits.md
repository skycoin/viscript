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



[7:05:21 AM] HaltingState: there are different types of styling an rendering, such as horizontal scroll on current line, or only expanding to multiple lines for the selection
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
3> substitute or set X, or select X and see options for things that could go in there (a function that returns int32, an int32 literal, a function that casts something else to int, an int32 at local scope etc)






you can eventually write own log library or have someone do this
then use ~ in upper left to open console
and tab through different types of logs, like input event, parsing etc and etc





maybe want to have a command line flag
for opengl version and disable or enable resizing
and to use advanced gl or use normal gl and fixed size

i have loop and pass in channel from main, then i push an int into the channel, when i want to tell it to stop the loop
look at github.com/skycoin/skycoin in cmd/skycoin/skycoin.go for an example of top level main and some stuff you might want to copy
make sure you have the gosublime plugin installed, for sublime 3
i have single entry point, default settings and then flags





we can avoid writing parser, there is another way
which is representing the computer program as software objects and then modifying the objects through the program/IDE

look at iPython notebook
there are two panels; you put the code in the first pane and the output is in the second pane
create a DOM object, and ability to embed "objects" into the terminal, inside of a "Document" object.
then add ability to do that. like have a button you can press for popdown and "add code"
and then can script in the box






when I put my cursor over an object, I want to see its type.

you can have "4" and 4 can be displayed the same way, but could be "uint32 4" or could be "uint64 4" etc.

There are "nodes" and you put your cursor over the node and want to know what its type is.

There is a type for "line" or "expression", which is one line in a function and that is root node. Then node can be (var x uint32) which is function definition.

There should be a function for saying how the node is displayed on screen. For taking a line and the root node and displaying it. And for displaying a function node.
eventually.

We should be able to navigate between the "nodes". and go left and right through an expression or go up and down between expressions. or up and down between functions/structs/ and higher level stuff.

Then if we are at a node, we should be able to hold a key and get a pop down, for the actions we can take to edit a node.

We also need a symbol for "Anything" or "Wild card", which is a node that has not been filled in
A function F. has name, type sig etc.

Then next lines start the expressions (one per line)

Then when we are at the expressions, we can do "insert expression", which is like new line, to insert a new line of code into the function.

Then we can do (var x int32)  , etc var x int32

and when we hold down key will see "create var" and will add a var node, but the name and type will not be filled in yet. So we put a wild card there.

so we have a (var * *), then we nagivate to the first * and put in a name. Then we navigate to the second * and we give it a type.

(varset x int32 5)

is like create var and set var, or var x int32 = 5

(set x 5) is "assign variable", that already exists, where as var defines new variables and varset, defines new variable, then sets it to expression value
(var x int32) is a "node". and it has 3 entries. first is node type. Then 2nd is "name" which is restricted to a string. Then third node inside of the parent node is "type" and must be a type.

So we have "restrictions" on the "slots". A node has slots, and other nodes can go into those slots. and there are requirements for what can go into the slots or what is permissible














-------------------------------------
this is one way of doing messags as structs, converting the messages to []byte and then converting them back again, and givng each message an ID, and labeling what component the message is to and from


we have to get the UI and system to do the length prefixed message channels and Tick() methods for dispatching them


get the opengl stuff in main.go into the hyper visor folder


we have provide an api through hypevisor.
that will allow a task bar to exist.
and task bar is just another program.
we do not need a task bar, but do need method to get the GUI windows.
and then methods to get the running processes and will eventually have terminal app for that.
I will do a process list, and a method for running multiple processes and message dispatch.

eventually the processes need to be in the scripting language, along with the UI library and handling and processing.
we have one computer process in golang, but inside that, we have "terminals" or discrete processes (can accept messages and can emit messages).

for menu button, make a 3x3 or 4x4 button of characters, no text.
look at uplink, hacker elite for interface.
for tabs, tab character is 2 spaces.

eventually you will have
- a DOM< a document object (document object model in the scripting language)
- a rendering object in the scripting language, which has hints or suggestions for how to render or constraints it will try to satisfy.

how long do you think?
and how do we deal with windows/terminals?

do you want to assign an int32 id to each process?
and assign an int32 id for each terminal?
and have each process specify the id of terminal they are sending message to?
or will hypervisor track an id for attached terminal?

there will be a object that receives length prefixed messages, then it will emit length prefixed messages back to the hyper visor.
and it will have a tick() method to clear its length prefixed messages from the queue.

i think terminals or display objects might need own id.
generally all resources like files or network  connections get int id,
and we can bind the terminal object to the process object.
but a lot of process will be in background and wont have display.
or one display may have ~10 process objects feeding it into for each pane and will multiplex them.

create a channel, for messages going to hypervisor, and create channel for each process.
and have a Tick() method to pop the messages off and return the messages to hypervisor.

in this framework, the UI changes have to only occur on user events,
or in response to a length prefixed message.



i will only give you "set character" and then you have to do things like "set cursor position" and the scrolll bar and advanced rendering stuff.
and write golang functions to wrap the messaging functionality,
and the messaging functionality is all state change based.
so setting cursor position is a state change,  or screen resize etc.

so we will need to seperate out the rendering thing from the process and then have API and it will probably be similar to opengl with the state change functions etc.

also, we have to support korean, russian and japanese and chinese.



i am doing input as length prefixed messages and doing rendering as length prefixed messages and I cannot remember why i am wrapping this or what the point is.



do we want the top level of hypervisor to be special?
or we want top level to be a terminal itself?
that has other terminals inside of it, and turtles all the way down?

do we want to pass messages to the parent?
or pass it to hypervisor, who passes to parent?
and is terminal a process object?



i have a thing called a process and it accepts length prefixed messages and it emits them.
do we want it setting cursor and screen with the length prefixed messages?
and we have a screen or terminal thing, which might be different type of resource than a process, but the processes may drive the screens or terminals.

how do we do this?


there are three things:

1> A blank screen, with start button in button left. ability to click the button to get popup of programs to run
2> A terminal window or monospaced thing pops up (a terminal and a process)
3> we will have multiple of these things (different programs/process types)


we can allow the terminal windows to call out to opengl and set their graphics etc and use your graphics library from opengl

but will restrict input and key presses and mouse through length prefixed message channel I guess



for simple app, it needs to know terminal size, and needs only two functions:
set cursor position and set character.



we need
1> to be able to have multiple windows (how to do this)
2> to be able to run different programs (the thing driving the terminal may be different)
3> to have good library for creating UIs for terminal applications (because we will be writing multiple applications), so good if this can be done quickly



one program might be a bulletin board system.
another might be a configuration tool.





----------------------
we need to create a "terminal" class,then have it receive the input messages and set the gfx stuff
----------------------





Do you want a "region" or "terminal object" so that we can pass input down like a widget?
where the big terminal gets event, then it determines the sub terminal to reply event to based upon focus and region.


the top level hypervisor, we have some "terminal object" and list of these, then the terminal object has "set cursor position" and "Set char" and that is all it has called on it.
then it gets length prefixed messages from the hypervisor for input.



get game Uplink, and study the gui.


we want want on global renderer, that loops through the terminals or terminal renders itself.


we will have
1> 1 type of terminal (it will do set cursor and do will do set char x,y)
2> a "process" interface, and this will determine how the key inputs are handled and will be responsible for rendering
3> there will be a "master process" for hypervisor, to deal with the sub-terminals or windows and moving them and things like that?
4> in theory we can run a master process inside another master process for recursion/containment and for lulz


we have pseudo widgets and a lot of stuff flushed out ,but the module interfaces are not correct or simple enough or contained




https://upload.wikimedia.org/wikipedia/en/c/cd/Fmli_lu.PNG


how do you want to do this?
1> an array of terminal classes in hyper visor
2> on method for rendering them all
3> an associated "process"  interface for handling input messages
4> messages terminal process sends to the terminal or hypervisor for "set character" and for "set cursor"




we should have keyboard messages going into the terminal from hypervisor
and have the process interface, sending messages for "set cursor" and "set character" to the terminal




this is basicly a full operating system.
and handles file system, networking and graphics/rendering, sound and keyboard input.
an the hypervisor gives cross platform, uniform, abstracted interface to all those things.




8888888888888888
inputs are part of terminal objects and should be handled in the object that is receiving the inputs, so even inputs should not be in terminal.
there will be a "process" interface attached to the terminal.
who will handle input messages.
and the process interface implementation will emit messages to hypervisor for setcursor and set character for terminal.
8888888888888888



look at process list.
the process gets message in, through incoming channel.
and process sends messages back to hypervisor to do stuff with terminal if it has an attached terminal.  see api.go.
any state, except for window size, and cursor position and that, can be moved into the process.
in the state variable.
and terminal should only keep cursor position and array of characters.



we will have multiple implementations of the process interface; do one that does one line input like bash and you hit enter and it goes to next line etc.




in viscript.go
do "defer teardown" or whatever.
and make the while 1 loop, exit on write to the close channel.




in the viscript.go
have a method to check if app should close
and if so, then let the loop finish in viscript.go
and it calls the teardown methods.





there are two ways to do this:
have a process that covers whole screen and deal with process stack in hypervisor.
or have just a terminal with one window, and have master process and have the sub terminals communicating back to the master process instead of the hypervisor.
the viscript editor will be one type of process object interface implementation.
the start button window and that will be another.
and then we want a bash like terminal with just "print line" type API functionality.




push as much state as you can into the process object.

and terminal only has

> set cursor
> set char
> creation methods
> its size, position, layer etc





do you see what I did with process list?
processes receive messages
processes communicate back with messages
888888888888888888888888888888888888888888888888888888888
processes control the terminal via message api; processes get the inputs from a terminal object
and almost all of the terminal state, has to be pushed into a process interface implementatoin; and the state is specific to a process
etc its type and what it does
while set character, get screen size, set cursor etc, is generic to the terminal




these variables and all the input handlers will be moved into the process
and is not part of terminal.
the hypervisor just figures out what terminal is in focus, what process is associated with it and forwards along the messages (maybe with offsets for mouse so correct for window position etc)





the "start button" program,
will be another process implementation





what does not work right now
- terminal object does not accept length prefixed message for set char, set cursor, from process object
- terminal does not have attached process object in process list, to relay input events to





the app, gfx, these libraries that call out for rendering will be called insido the process implementation and will be imported from there; and will be boiled down to messages to drive terminal going over the messaging channel to hypervisor

i made small example in process/example/api.go

and those are "atomic" terminal operations, and then the gui elements will be drawn using compount operations (multiple atomic operations)




anything inside of process, will not be allowed to import gl or glfw or any of this; and will only communicate with terminal over the event interface





bash, new line, read command, print some stuff; allow scrolling for back log; that is simpliest thing and we need that




each process has its own state struct
and its own isolate state.
process state can only be modified by inputing length prefixed messages to the process
and the process responds with length prefixed messages.





there is no reason a process needs a terminal, but i do not have any process instances right now that are not using the terminal interface to set gui.
it is impossible to start a process without a terminal, because processes that take in user input and have termial, exist.
becuase you need a process with a terminal, to run a process that does not have a terminal




what is a bash thing:

i mean, lines, and have line 1, line 2, line 3, then something that renders it to terminal and lets you scroll and which wraps the output to the termimal size; and have a max buffer size and maybe scroll.

"what i know of a bash interface...... it has to be a visible box/rectangle on the screen where you can type.  right?"

This is called a terminal instance, attached to a process.

and the terminal has a "set char" and "set cursor" messsage type it uses as API between process and the terminal instance





in our system we will have an example process, attached to a terminal.
and the process will render the interface and we will type something and when we hit enter it will echo it.
and will have scrollable buffer, that goes through the length prefixed message API and event channgel.

it will
- receive input events over the event channel
- it will send setCursor and setCharacter (set character in terminal), messages to hypervisor, which will be forwarded to the attached terminal display.

the process is what will receive the input events and the process is what will sent the commands to the terminal object






you hit key, the hypervisor takes that key and determines which terminal is in focus.
and then relays the key to the "process" for that terminal/window.

then the process gets the key and does something. the process emits mesages to chang the GUI state. the messages go to the hypervisor. the hypervisor forwards the messages to the terminal if they are terminal messages






the the stack of terminals has  method for draw()
that will draw the hypervisor
and hypervisor might keep track of which terminal is in focus
a terminal is a rectangle in the window, drawn by hypervisor





I would be happy if
- terminals draw to screen
- the processes can animate the terminals
- the proecess user input handle is already implemented

the process gets events in on the events in channel, you call tick and it processes them and events go out on the events out channel
and the process is determinstic.
go into /process/example/api.go

that is function that writes terminal messages to the output event channel

issues:
2> the terminals are not associated to processes yet






once you have terminal working and it responds to event channel and you can draw a single character etc
then you can implement bash as a process implemntation instance
right now, you are implementing the stuff you would nee to implement any process or application.






viscript/process/example/api.go

and multiple process implementations can use one library and just import the library; and will have library for atomic, and library for compound or widget like GUI stuff






use ~ to toggle terminal 0?






we have six applications that will be running on top of this
and soon you can just focus on making it nice and bugs and etc
and making it like a video game





the hypervisor will process the commands for display and controlling display
the hyper visor will send input messages to the process/task
the process/task will send draw messages back to the hypervisor
the hypervisor terminal object will receive the messages from hypervisor, and will have a draw method and will import the opengl
the process/task will import a gui library, that wraps the operations through the length prefixed message channel






I have to get contract, then have money from customer, then I have sixty days to get thing done





needs:
5> ability to resize, drag, move terminals around
6> ability to cycle through infocus terminal with the ~ key
once we have bash terminal setup
then we do commands for
list terminals
list processes
start process
and we do reflection/introspection self tests
so system self boot straps





look at uplink
and imagine the windows are bash terminals
and you are clicking the buttons to open bash terminnals





 hypervisor will look at the terminal "InFocus" and forward all the input to the process acsociated with that terminal






 and we need bash done NOW
 i have six applications i have to import into this interface
 bash only have fmt and printf and new line and add character, etc
and drawing them around the screen and opening and closing and minimizing them




https://en.wikipedia.org/wiki/D-Bus
 https://en.wikipedia.org/wiki/File:D-Feet.png





window border is annoying
inclusive window border is probably better than window border on whole thing
make it cyber punk and 1-3 pixels
light accent.
interior of border has 1-3 pixels of special stuff on border interior pixels.
and interior border means that if window is 256 pixels wide, the border stuff is drawn in pixel 0,1,2  and pixels 255,254,253 etc
a gradient could be good, because would indicate visually which side of the line is "inside" and which side is "outside"
interior borders are better than exterior borders, because the border does not fuck with the grid spacing.





> but then how does it respond to resizing the OS window?

That is good question./important question.

We send a message/struct to the process controlling the terminal, that resized happens. Then the process will send a clear screen command and rerender the terminal buffer probably, with the new size






https://en.wikipedia.org/wiki/File:Stdstreams-notitle.svg







https://upload.wikimedia.org/wikipedia/en/3/3b/Norton_Utilities_6.01_UI.png
this is user defined characters
where application sets the character grid itself, that it will use to render
and then can make gui type stuff, while actually being monospaced







there will be terminal 0, which is terminal that is lowest in stack; that takes up whole screen
and this is the "Desktop background" terminal
where we will have the start menu and if i right click on it, will give me dropdown etc
and this terminal always covers the whole screen

we have a special, terminal or terminal 0, that is size of whole viewport and is behind the other terminals (the desktop background viewport)








there is a general app or process type for bash type commands like this
but other types of processes are for other apps and not a general type of terminal







we are supposed to have a resource registry, where all viewports, terminals, processes are given an id and we have ResourceId and ResourceType, but will leave it for now or ignore it and I am not sure its used for anything right now
and it is part of dbus right now, but not sure that is best place
dbus creates channels for communication between resources
a resource has an in and a resource type
so we can see that a particular channe has an id, that it is owned by a process with a particular ID and that it is to a terminal resource object, which has a particular id, etc and we can enumerate all of the resources and their types and the channels connecting them






[1/28/2017 8:50:45 AM] chattanoo: looking it over.  what does "pubsub" mean?
[1/28/2017 8:51:02 AM] HaltingState: it means publisher and subscriber
[1/28/2017 8:51:09 AM] HaltingState: its not a one to one channel
[1/28/2017 8:51:19 AM] HaltingState: but list of people subscribed
[1/28/2017 8:51:34 AM] HaltingState: the terminal subscibes to process and the process subscribes to terminal
[1/28/2017 8:51:42 AM] HaltingState: so technically only one subscriber in each direction for now

each object will have a single channel for "in"
and when publisher writes, the message will get written to all the in channels for subscribers
dbus will eventually need a socket type
where there is a server and you get a "connection" and its one to one, and bidirectional and like a socket; but we dont have or need that yet







-------------------------------WECHAT transition---------------------------------

if the terminal publishes messsage to its pubsub channel (it is out channel),
that message gets written to the subscribers (the process) instantly




wrapper for printf type

there will be a library terminal imports
for scrolling, the terminal object might have to keep a copy of the stuff to scroll
and update the terminal state etc

the process has a document object model, then you scroll through it and the
terminal sends messages about how to render it, etc

process is implemented as an interface because we might have several types of processes
or implementations.

for instance, a bash type printf driven process.  Then have another process
implementation that has gui with widgets etc



we will have a background terminal, that takes up whole screen and is behind all
terminals.  this is our "virtual desktop"

Then we will have a widget library, for drawing buttons

Then we will add a start button, that you can click and have popup menu, of applications
and clicking on the application will launch a terminal window with the app.

this is a package imported only by terminal.  and which will not import anything

by terminal

i don't know if anything in the ui package is used or imported

but this is separate




the process instance is self contained and not allow to import objects like terminals;
it only communicates over the channel where messages come in and messages leave on the channel.

api should be a module imported by process instance.  because multiple process
implementations will use same library

one library will be for text based putChar, etc and eventually will have graphics/widget library

the sprintf function would write a series of commands to the command channel (out channel?)





does terminal write input messages to its out channel?  and does process read those input
messages?  and can process write messages to its out channel via dbus?

it should only communicate through dbus
you write to dbus and dbus writes to channel


terminal must not import process and vice versa





[PublishTo()] also appends prefix of 32 bit, to tell app id of channel.





dbus right now does not store or manage the channels
the InChannel is the channel dbus is writing to, to send messages to terminal




dbus will prefix the channel id, so the program knows which subscribed channel
wrote the message




skycoin/skycoin

the CLI for skycoin.  you will write a lib for task to run a cli like that inside of
the viscript terminal but als in normal terminal with same lib

and look in src mesh and run node for that and cli





red: mesh doesn't appear to be having GUI i think it doesn't compile

it compiles.  it has a node of command line
and it has a CLI interface to the node

its in /src/mesh folder.  look at the readme






do the mesh api





run the meshnet cli and meshnet daemon.  we need to get cli like this.  working in viscript





Terminal might even store its own character of text array or have DOM so that it can do
back scroll through history





[i believe this was a quote from Vyacheslav Zgord...]: you can just use web interface for
command-line client - the manual for it is in skycoin/src/mesh/README.md and it let's you
run web interface on port 9999.






for scrolling the terminal task needs to store internal redundant text backlog etc.






eventually we need master api
to list all terminals and attached processes
or list all dbus objects and subscribers and publishers
and to start a new terminal or process of a specific type from a list.
it might only need one packet type.  for receiving text and another for responding







do we have library for task to get user input and printf to terminal and
maintaining a back buffer and scroll in terminal etc.






what about a meshnet process
have list of task types which is a process interface implementation
then do task implementation for meshnet CLI





how would we share printf library implementation between different task implementations?





task type.  we do not have text labels for task type yet
and probably only have one task type or task interface implementation.  right now.
we should have a task type id and a text label for each task interface implementation






Red:
1) go install && viscript
2) go run rpc/cli.go






HaltingState:
for CLI apps/daemons, getting self inspection terminal to run inside of viscript,
for listing the tasks etc






skycoind
terminald
cxod

the terminal implementation needs to import some library that manages everything.
like println and etc.







i would not worry about changing the back process for now
if we spawn new process, we will have new terminal for now

but i need to get meshnet CLI and daemon or skycoind running in viscript.






we just need to have library that is drop in replacement for log or fmt.  and which can
wrap log or format and work in normal terminal.

can you run the meshnet daemon.  and meshnet CLI and get it working as a process type in
viscript
2 separate process types







we will replace all fmt and log imports with our library








it has 2 options.
just use printf normally
unless it is configured for viscript.  then send to channel







we could create a type "external process" which is process run with exec...
how does gotty work







how quickly can we dump exes into a folder and run them from inside viscript?
like bash does for its bin folder?

Lets ls or have command to give us list of commands in folder.
then we type the command name or "run command" and it will pop open new terminal & run it

and we will throw the daemons and meshnet etc in the bin folder







Red:
so I run: "go install" at first in the skycoin server directory and once it's installed
I use same package as gotty it's called pty and I run 2 routines that wait for the output
from the skycoin server as you can see






HaltingState:
we want to be able to control the directory the exes we are running are in.
you can use path if you want

yes like you have 'go run command.go'
and we could run apps from source if needed and not just exe






HaltingState:
And good api for attachment and detachment etc

But we could do more advanced.

Can we hsve teo commands.

One runs in new terminal window.

Other runs in current window and sets parent which will resume.

Another or third opens process in background








can we get the meshnet daemon and clie packaged in there by default






Red: i was thinking of one api that could provide basics to all the
process types

HaltingState: i like being able to import one thing








i do not have problem with each new process starting in its own window.

but emphasis is on working apps afn viscript we can put on the website and
people can use.


but we could do more advanced.

can we hsve teo commands.

one runs in new terminal window.

other runs in current window and sets parent which will resume.

another or third opens process in background

and good api for attachment and detachment etc








CorpusC: but we need to be using inheritance

HaltingState: we have api standard interface abd library to import.
could also use some type of inheritance but api thing is better

not sure how to use interface for this







[re: CorpusC's command-prompt-scrolling related questions]
You can generally put character for carriage return on end of line and go to next line







there is supposed to be thing for managing open gl window.
and thing for managing input
and hypervisor was supposed to manage tasks and list of running tasks
and then thing for objects and object ids.
eventually the viewport itself will just be an object and you can have multiple
viewports.  on same screen or multiple screens or same computer  or across
multiple computers







eventually we will have to make viewport modular and put task management and object
management, as its own thing that can run headless







viewport handles
- opengl window
- keyboard events
- terminal drawing

hypervisor handles
- object list
- dbus, interprocess communication
- task list etc

the hypervisor could in theory, run without a viewport







"in virtual desktops, the viewport is the visible portion of a 2D area which is
larger than the visualization device"

it means specific thing in virtualization







!!! as soon as i type "help" it needs to print a command list
how do i list task ids etc
and get list of terminals, & tasks

1> how do i get list of tasks from within viscript
2> how do i get list of terminals etc, and what is attached (introspection)
3> how do i use the viscript rpc from within viscript
4> when will "help" print a command list
5> what is rpc and why cant i use it from within viscript
6> where is viscript rpc
7> how do i use the viscript rpc from within viscript
8> hwere is the viscript command for opening a new viscript terminal







go on
skycoin.atlassian.net
create ticket for what you are working on and tag the upwork work log with ticket









several steps need to be taken, before this is a working app that can be put on
website.  we are trying to go into production now
and it means we need tickets for tracking bugs and ticket for tracking what works and
what does not work









Red:
CTRL+L could work for clearing it think, it wouldn't be a command thus won't
interrupt running atached external task.
it's like that in a normal terminal.

go run rpc/cli/cli.go won't work alone as stated in README.md inside the rpc package
https://github.com/skycoin/viscript/tree/master/rpc
you need to have viscript running at first

use "exec go run rpc/cli/cli.go"

Could we have an interface for external process like containing start, tick, shutdown
and process to use it because we have to move the external process outside of terminal
as in tickets, right now commented all the additional functionalities like ctrl+z for
running in background ctrl+c for stopping process and also it doesn't stop automatically
when it sees EOF from running process, cleaned up a bit and I'll work more tomorrow.
Main thing is that it works for now. We might even need to remove that State variable
from external process and have ProcessOut and have Process to watch that for output
sequentially and printing it to viscript terminal.









Red: isn't it easy to run "e srv"
at first and then "e cli" in another terminal?









Red: there's 3 goroutines.  1 that reads, 1 that writes upon receiving input from
viscript, and 1 that waits for those goroutines to end cleans up









HaltingState: if command is in folder we want it to show up in terminal
as native command

look at skycoin rep folder.  Cmd skycoin skycoin.go as example of command line
settings and defaults.

we will have viscript binary.  Then a subfolder called bin.  And exes will go in there.

Then we will have script for loading the defaults into bin.

Bin will be in .gitignore









The external task is not attached to the terminal
But to a list of external tasks








there will be a program to bind the terminal to the external task

or have a parameter for terminal, which if set, means foreground is attached
to external task








the only programs we are running are our programs.  In golang








each exe would be in own folder and may have assets







we should pass in some rpc library thing like port on local host to skycoin exe.

then the exe will connect to hypervisor via tcpip local and register itself and we will
send shutdown commands via that channel


can send length prefixed messages and commands between the hypervisor and the child task

later we might have special log wrapper and will hack our log to send the logs over
this channel instead of operating system terminal channel.

And we will also use this channel to implement sigint and send signals or events
back and forth without relying on OS dependent stuff

Have a library we include in our special apps.  If command line parameter is set,
then app will connect back to hypervisor over tcp/ip on local port.

And we can send length prefixed message like hypervisor users.  For "shutdown" to
ask task to shutdown.

And will use this channel for communication

And eventually commands and terminal and sigint can go over this channel.  in platform
independent way

then we can ping the app every once in a while to see if its still running and responding









Steve: use go vendoring, gvt
for dependencies










Red: this is how skycoin is run inside run.sh
--gui-dir+"${DIR}/src/gui/static/" $@
so my guess is we only need the /src/gui/static/ directory for running skycoin right?

Steve: yes
$@ means copy the arguments given to the script











HaltingState: create a channel.  in dbusm for terminal global commands.
and create messages for deleting terminal by id, listing terminals by id,
and creating a new terminal.
Then if something wants to modify terminal it will write to that dbus channel.

maybe write smalle wrapper library for terminal control over dbus









if you print a line and it runs over to next line, then you must put a block character
at the end of the line.  To visually indicate that.

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


buy "mesogold, collodial gold" and take 1 tea spoon a day (or less if its too strong)


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


"












=== Spec ===

Macros + Reflection
Statically typed scheme

Define types
- i32 (int 32)
- u8 (uint 8)
- []byte (byte array, all byte arrays are fixed length)

Define operators
- i32_add
- i32_sub
- i32_mult
- i32_div

Syntax
- (i32_add 3, 5)

A function
- has tuple of inputs, tuple of outputs
- i32_add has type (i32, i32) (i32); takes in two i32 and returns one i32
- a function has a list of statements

A struct
- a tuple of types
- (i32, i32) is a tuple of two i32

(def_struct Tag,
	tag []byte
	uuid [16] byte
)

A union
- means "has to choose A or B"
- enumerates possibilities
- a choice of types

A function has an array of "statesments" or "blocks" (an array of statements) or a control flow (if/then, for)

(def_func, NAME, (u32 a, u32 b), (u32), 
	(u32_add a, b)
	)

Type types of operators
- functions on (can modify object)
- functions of (cannot modify object)
- functions on and functions of, should be in different colors

The type of an object is a type
- the type of an object is a struct, defining the type

(def_func_on ...)
(def_func_of ...)

def_func_of cannot call functions def_func_on
- can only read program, cannot modify

(def_var_m i32 x)
- defines new entity x

(def_var_imutable i32 x)
- defines new variable, which cannot be changed after creation

An "assert" is something that must be true

An "affordance" is something that CAN be done to an object
- all affordances, must be enumerable
- a user must be able to select each possible action, from a list

A "restriction" is something that CANNOT be done to an object
- restrictions must be checked and the affordance list filtered
- Some types of restrictions, can be applied to a program as an operation, to remove an affordance
- an example, is that structs/functions are defined, then the affordance for modifying them is removed (allowing compilation or simplification to a static binary)

A "context" is current state
- list of functions
- list of structs
- list of defined things
- the context contains the current module, the stack, the current line and function
- current function the program is on
- current statement the program is on
- list of variables in the current scope
- list of functions in the current module
- list of variables defined in a module

reflection
- each type has a list of functions and operators on it
- reflection on an object is a func_of the object meta_type
- modifying or extending an object is a func_on the object meta_type

(is_def x)
- returns if thing is defined

(is_type x, type)
- determines of object is of type x

(type_of x)
- returns type

A "choice" is a place where A or B (where program can choose A or B)
- unions are for objects (structs)
- choices are for code
- a choice has a list of preconditions (things that must be true for choice to be made) and a list of statements (which can only be chosen if the preconditions are met)
- a choice may have a signature (if it returns an object)
- a choice could be modeled as a special type of function, that returns something
- the choice operator, is a special function, that takes in the current context (state of program)

(def_func, NAME, (uint a, uint b), (uint),
	(choice, 
		(true, (u32_add a, b)),
		(true, (u32_add b,a)))
	)
)

UIDs
- struct { text_tag []byte , uuid [16]byte}
- all variables, functions, modules have 128 bit UIDS
- a UUID has a text "text_tag", or keyboard/display name
- The UUID is used to look up the object in a table

Modules
- all code occurs in modules
- a module has a name
- a module contains functions and function definitions
- a module contains structs and struct definitions

(as x y)
(as x z)
- another way of writing choice operator
- see XL programming language
- when conditions occurs, x can be choice Y or Z
- as operator, defines choices, that are not bound to a context or object
- the as operator, defines the conditions, when an affordance is available

- The program itself is an object (a struct)
- The program begins with a default object (the null object)
- The program starts with a series of actions that can be applied to it (affordances)
- The program is built up, by a series of operators applied to it (affordances)

Programs accept and emit only length prefixed messages
- there is a function for checking if there is an incoming length prefixed message
- there is a function for emitting a length prefixed message
- there is a function for receiving the length prefixed message from the queue
- there is a function for halting the program, until a length prefixed message is available

A program can spawn, isolated sub-programs that it can only communicate to over length prefixed messages

An "agent" is a program that can apply affordances
- the agent program reads the curent affordances on objects and applies them
- agents are often restricted to a subset of objects and a subset of affordances

A "behavior" is a set of criteria or goals or states, that the agent attempts to maintain

Reflection
- the fields and functions on a struct can be enumerated
- the signature and body of a function can be enumerated
- the list of structs, variables and objects in a given module can be enumerated
- the list of modules, imported by a given module can be enumerated
- the list of local scope variables, in a given function/context/stack frame can be enumerated
- the types of each variable, can be enumerated

Reflection on dependecy graphs

Define a program object
- then apply its operators on it, to construct the program

---

!!!
Crash logs
- overlay all crashes on the source code
- find and trace crash logs

Interactions between classes/functions
- graph all interactions


====

atomic types (int32, uin32, []byte)
operatiosn on atomic types (uint32.add a b)
structs
functions
modules
type signatures

Atomic types
- ints
- byte arrays
- "Type" objects

====

A function is a 
1> name (text)
2> input tuple (name, type) pair list
3> output tuple (list of types/signatures returned)
4> an array of lines/statements

A struct is a 
1> name
2> list of (name, type) pairs
3> later list of functions on the struct but ignore this for now

A module is
0> The name of the module (string)
1> A list of modules imported by the current module
2> A list of structs defined in the current module
3> A list of function defined in the current module
4> A list of variables at the global scope of the module

===

Note:
- modules
- structs
- functions

Should all have unique ids, to be used as references
- unique IDs can be 64 bit (effectively pointers to def)

A function is

struct Function {
	name []byte
	input []struct{[]byte name, type Type)} //name/type pair array
	output []struct{[]byte name, type Type} //optional name
	lines []Expressions
}

A struct is

struct Struct {
	name []bytes
	fields []struct{[]byte name, type Type)} //name/type pair aray
}

A module is

struct Module {
	name []bytes
	module_imports []*Modules
	module_functions []*function
	module_structs []*structs
}

Each of these is written as S notation

(def_func Name (in...) (out...) (expression_array...) )

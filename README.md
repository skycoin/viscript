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

"

=== Spec ===

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

A "restriction" is something that CANNOT be done to an object
- restrictions must be checked and the affordance list filtered

A "context" is current state
- list of functions
- list of structs
- list of defined things
- the context contains the current module, the stack, the current line and function
- current function the program is on
- current statement the program is on
- list of variables in the current scope
- list of functions in the current module
- list of 

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
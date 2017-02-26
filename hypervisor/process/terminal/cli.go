package process

/*import (
	"github.com/corpusc/viscript/app"
)*/

var (
	log      []string
	commands []string
	currCmd  int    //index
	cursPos  int    //cursor/insert position
	prompt   string = ">"
	//this assumes 64 horizontal characters.  so we dedicate 2 lines for each command
	maxCommandSize = 128 - len(prompt) - 1 //& ending space for cursor at the end of (potentially the 2nd) line
)

func init() {
	println("(process/terminal/cli).init()")
	log = []string{}
	commands = []string{}
	commands = append(commands, "OLDEST command that you typed (not really, just an example of functionality)")
	commands = append(commands, "older line that you typed (nah, not really)")
	commands = append(commands, prompt)
	currCmd = 2
}

func EchoWholeCommand() {
	println("(process/terminal/cli).EchoWholeCommand()")
}

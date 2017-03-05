package process

import (
	"github.com/corpusc/viscript/hypervisor"
	"github.com/corpusc/viscript/hypervisor/dbus"
	"github.com/corpusc/viscript/msg"
)

var (
	log      []string
	commands []string
	currCmd  int    //index
	cursPos  int    //cursor/insert position, local to 1 commands space (2 lines)
	prompt   string = ">"
	//FIXME to work with Terminal's dynamic self.GridSize.X
	//assumes 64 horizontal characters, then dedicates 2 lines for each command.
	maxCommandSize = 128 - 1 //reserve ending space for cursor at the end of last line
)

func init() {
	println("(process/terminal/cli).init()")
	log = []string{}
	commands = []string{}
	commands = append(commands, prompt+"OLDEST command that you typed (not really, just an example of functionality)")
	commands = append(commands, prompt+"older line that you typed (nah, not really)")
	commands = append(commands, prompt)
	cursPos = 1
	currCmd = 2
}

func EchoWholeCommand(outChanId uint32) {
	println("(process/terminal/cli).EchoWholeCommand()")

	//FIXME? actual terminal id really needed?  i just gave it 0 for now
	//message := msg.Serialize(msg.TypePutChar, msg.MessagePutChar{0, m.Char})

	message := msg.Serialize(
		msg.TypeCommandLine, msg.MessageCommandLine{0, commands[currCmd], uint32(cursPos)})

	hypervisor.DbusGlobal.PublishTo(
		dbus.ChannelId(outChanId), message) //EVERY publish action prefixes another chan id
}

func traverseCommands(delta int) {
	if delta > 1 || delta < -1 {
		println("FIXME if we ever want to stride/jump by more than 1")
		return
	}

	currCmd += delta

	if currCmd < 0 {
		currCmd = 0
	}

	if currCmd >= len(commands) {
		currCmd = len(commands) - 1
	}

	cursPos = len(commands[currCmd])
}

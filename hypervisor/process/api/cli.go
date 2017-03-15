package api

import (
	"strings"

	"github.com/corpusc/viscript/hypervisor"
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

	m := msg.Serialize(msg.TypeCommandLine, msg.MessageCommandLine{0, commands[currCmd], uint32(cursPos)})
	hypervisor.DbusGlobal.PublishTo(outChanId, m) //EVERY publish action prefixes another chan id
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

func moveOneStepLeft() bool { //returns whether moved successfully
	cursPos--

	if cursPos < len(prompt) {
		cursPos = len(prompt)
		return false
	}

	return true
}

func moveOneStepRight() bool { //returns whether moved successfully
	cursPos++

	if cursPos > len(commands[currCmd]) {
		cursPos = len(commands[currCmd])
		return false
	} else if cursPos > maxCommandSize { //allows cursor to be one position beyond last char
		cursPos = maxCommandSize
		return false
	}

	return true
}

func moveOrJumpCursorLeft(mod uint8) {
	if msg.ModifierKey(mod) == msg.ModControl {
		numSpaces := 0
		numVisible := 0 //NON-space

		for moveOneStepLeft() == true {
			if commands[currCmd][cursPos] == ' ' {
				numSpaces++

				if numVisible > 0 || numSpaces > 1 {
					moveOneStepRight()
					break
				}
			} else {
				numVisible++
			}
		}
	} else {
		moveOneStepLeft()
	}
}

func moveOrJumpCursorRight(mod uint8) {
	if msg.ModifierKey(mod) == msg.ModControl {
		for moveOneStepRight() == true {
			if cursPos < len(commands[currCmd]) &&
				commands[currCmd][cursPos] == ' ' {
				moveOneStepRight()
				break
			}
		}
	} else {
		moveOneStepRight()
	}
}

func goUpCommandHistory(mod uint8) {
	if msg.ModifierKey(mod) == msg.ModControl {
		currCmd = 0
	} else {
		traverseCommands(-1)
	}
}

func goDownCommandHistory(mod uint8) {
	if msg.ModifierKey(mod) == msg.ModControl {
		currCmd = len(commands) - 1 // this could cause crash if we don't make sure at least 1 command always exists
	} else {
		traverseCommands(+1)
	}
}

func (st *State) actOnEnter(serializedMsg []byte) {
	numLines := 1

	if cursPos >= 64 { //FIXME using Terminal's st.GridSize.X
		numLines++
	}

	for numLines > 0 {
		numLines--
		// hypervisor.DbusGlobal.PublishTo(st.proc.OutChannelId, serializedMsg)
	}

	words := strings.Split(commands[currCmd][len(prompt):], " ")
	st.actOnCommand(words)
	log = append(log, commands[currCmd])
	commands = append(commands, prompt)
	currCmd = len(commands) - 1
	cursPos = len(commands[currCmd])
}

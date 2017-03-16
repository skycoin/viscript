package process

import (
	"strings"

	"github.com/corpusc/viscript/hypervisor"
	"github.com/corpusc/viscript/msg"
)

type Cli struct {
	Log      []string
	Commands []string
	CurrCmd  int //index
	CursPos  int //cursor/insert position, local to 1 commands space (2 lines)
	Prompt   string
	//FIXME to work with Terminal's dynamic self.GridSize.X
	//assumes 64 horizontal characters, then dedicates 2 lines for each command.
	MaxCommandSize int //reserve ending space for cursor at the end of last line
}

func NewCli() *Cli {
	println("(process/terminal/cli).init()")
	var cli Cli
	cli.Log = []string{}
	cli.Commands = []string{}
	cli.Prompt = ">"
	cli.Commands = append(cli.Commands, cli.Prompt+"OLDEST command that you typed (not really, just an example of functionality)")
	cli.Commands = append(cli.Commands, cli.Prompt+"older line that you typed (nah, not really)")
	cli.Commands = append(cli.Commands, cli.Prompt)
	cli.CursPos = 1
	cli.CurrCmd = 2
	cli.MaxCommandSize = 128 - 1

	return &cli
}

func (c *Cli) HasEnoughSpace() bool {
	return len(c.Commands[c.CurrCmd]) < c.MaxCommandSize
}

func (c *Cli) AddCharAndMoveRight(nextChar uint32) {
	c.Commands[c.CurrCmd] =
		c.Commands[c.CurrCmd][:c.CursPos] +
			string(nextChar) +
			c.Commands[c.CurrCmd][c.CursPos:]
	c.moveOneStepRight()
}

func (c *Cli) OnBackSpace() {
	c.Commands[c.CurrCmd] = c.Commands[c.CurrCmd][:c.CursPos] +
		c.Commands[c.CurrCmd][c.CursPos+1:]
}

func (c *Cli) OnDelete() {
	c.Commands[c.CurrCmd] = c.Commands[c.CurrCmd][:c.CursPos] +
		c.Commands[c.CurrCmd][c.CursPos+1:]
}

func (c *Cli) EchoWholeCommand(outChanId uint32) {
	println("(process/terminal/cli).EchoWholeCommand()")

	//FIXME? actual terminal id really needed?  i just gave it 0 for now
	//message := msg.Serialize(msg.TypePutChar, msg.MessagePutChar{0, m.Char})

	m := msg.Serialize(msg.TypeCommandLine,
		msg.MessageCommandLine{0, c.Commands[c.CurrCmd], uint32(c.CursPos)})
	hypervisor.DbusGlobal.PublishTo(outChanId, m) //EVERY publish action prefixes another chan id
}

func (c *Cli) traverseCommands(delta int) {
	if delta > 1 || delta < -1 {
		println("FIXME if we ever want to stride/jump by more than 1")
		return
	}

	c.CurrCmd += delta

	if c.CurrCmd < 0 {
		c.CurrCmd = 0
	}

	if c.CurrCmd >= len(c.Commands) {
		c.CurrCmd = len(c.Commands) - 1
	}

	c.CursPos = len(c.Commands[c.CurrCmd])
}

func (c *Cli) moveOneStepLeft() bool { //returns whether moved successfully
	c.CursPos--

	if c.CursPos < len(c.Prompt) {
		c.CursPos = len(c.Prompt)
		return false
	}

	return true
}

func (c *Cli) moveOneStepRight() bool { //returns whether moved successfully
	c.CursPos++

	if c.CursPos > len(c.Commands[c.CurrCmd]) {
		c.CursPos = len(c.Commands[c.CurrCmd])
		return false
	} else if c.CursPos > c.MaxCommandSize { //allows cursor to be one position beyond last char
		c.CursPos = c.MaxCommandSize
		return false
	}

	return true
}

func (c *Cli) moveOrJumpCursorLeft(mod uint8) {
	if msg.ModifierKey(mod) == msg.ModControl {
		numSpaces := 0
		numVisible := 0 //NON-space

		for c.moveOneStepLeft() == true {
			if c.Commands[c.CurrCmd][c.CursPos] == ' ' {
				numSpaces++

				if numVisible > 0 || numSpaces > 1 {
					c.moveOneStepRight()
					break
				}
			} else {
				numVisible++
			}
		}
	} else {
		c.moveOneStepLeft()
	}
}

func (c *Cli) moveOrJumpCursorRight(mod uint8) {
	if msg.ModifierKey(mod) == msg.ModControl {
		for c.moveOneStepRight() == true {
			if c.CursPos < len(c.Commands[c.CurrCmd]) &&
				c.Commands[c.CurrCmd][c.CursPos] == ' ' {
				c.moveOneStepRight()
				break
			}
		}
	} else {
		c.moveOneStepRight()
	}
}

func (c *Cli) goUpCommandHistory(mod uint8) {
	if msg.ModifierKey(mod) == msg.ModControl {
		c.CurrCmd = 0
	} else {
		c.traverseCommands(-1)
	}
}

func (c *Cli) goDownCommandHistory(mod uint8) {
	if msg.ModifierKey(mod) == msg.ModControl {
		c.CurrCmd = len(c.Commands) - 1 // this could cause crash if we don't make sure at least 1 command always exists
	} else {
		c.traverseCommands(+1)
	}
}

func (c *Cli) GetCommandWithArgs() (string, []string) {
	words := strings.Split(c.Commands[c.CurrCmd][len(c.Prompt):], " ")
	return words[0], words[1:]
}

func (c *Cli) OnEnter(st *State, serializedMsg []byte) {
	numLines := 1

	if c.CursPos >= 64 { //FIXME using Terminal's st.GridSize.X
		numLines++
	}

	for numLines > 0 {
		numLines--
		hypervisor.DbusGlobal.PublishTo(st.proc.OutChannelId, serializedMsg)
	}

	st.actOnCommand()
	c.Log = append(c.Log, c.Commands[c.CurrCmd])
	c.Commands = append(c.Commands, c.Prompt)
	c.CurrCmd = len(c.Commands) - 1
	c.CursPos = len(c.Commands[c.CurrCmd])
}

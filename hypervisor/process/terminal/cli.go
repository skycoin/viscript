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
	NumLinesFromTheEnd int //backscroll offset
	MaxCommandSize     int //reserve ending space for cursor at the end of last line
}

func NewCli() *Cli {
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

func (c *Cli) AdjustBackscrollOffset(delta int) {
	println("BACKSCROLLING delta: ", delta)
	c.NumLinesFromTheEnd += delta

	if c.NumLinesFromTheEnd >= len(c.Log) {
		c.NumLinesFromTheEnd = len(c.Log)
	}

	if c.NumLinesFromTheEnd < 0 {
		c.NumLinesFromTheEnd = 0
	}
}

func (c *Cli) InsertCharIfItFits(char uint32, state *State) {
	if len(c.Commands[c.CurrCmd]) < c.MaxCommandSize {
		c.InsertCharAtCursor(char)
		c.EchoWholeCommand(state.proc.OutChannelId)
	}
}

func (c *Cli) InsertCharAtCursor(char uint32) {
	c.Commands[c.CurrCmd] =
		c.Commands[c.CurrCmd][:c.CursPos] +
			string(char) +
			c.Commands[c.CurrCmd][c.CursPos:]
	c.moveCursorOneStepRight()
}

func (c *Cli) DeleteCharAtCursor() {
	c.Commands[c.CurrCmd] =
		c.Commands[c.CurrCmd][:c.CursPos] +
			c.Commands[c.CurrCmd][c.CursPos+1:]
}

func (c *Cli) EchoWholeCommand(outChanId uint32) {
	termId := uint32(0) //FIXME? correct terminal id really needed?
	//message := msg.Serialize(msg.TypePutChar, msg.MessagePutChar{0, m.Char})

	m := msg.Serialize(msg.TypeCommandLine,
		msg.MessageCommandLine{termId, c.Commands[c.CurrCmd], uint32(c.CursPos)})
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

func (c *Cli) moveCursorOneStepLeft() bool { //returns whether moved successfully
	c.CursPos--

	if c.CursPos < len(c.Prompt) {
		c.CursPos = len(c.Prompt)
		return false
	}

	return true
}

func (c *Cli) moveCursorOneStepRight() bool { //returns whether moved successfully
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

		for c.moveCursorOneStepLeft() == true {
			if c.Commands[c.CurrCmd][c.CursPos] == ' ' {
				numSpaces++

				if numVisible > 0 || numSpaces > 1 {
					c.moveCursorOneStepRight()
					break
				}
			} else {
				numVisible++
			}
		}
	} else {
		c.moveCursorOneStepLeft()
	}
}

func (c *Cli) moveOrJumpCursorRight(mod uint8) {
	if msg.ModifierKey(mod) == msg.ModControl {
		for c.moveCursorOneStepRight() == true {
			if c.CursPos < len(c.Commands[c.CurrCmd]) &&
				c.Commands[c.CurrCmd][c.CursPos] == ' ' {
				c.moveCursorOneStepRight()
				break
			}
		}
	} else {
		c.moveCursorOneStepRight()
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
		c.CurrCmd = len(c.Commands) - 1 //this could cause crash if we don't make sure at least 1 command always exists
	} else {
		c.traverseCommands(+1)
	}
}

func (c *Cli) CurrentCommandLine() string {
	return strings.ToLower(c.Commands[c.CurrCmd][len(c.Prompt):])
}

func (c *Cli) CurrentCommandAndArgs() (string, []string) {
	tokens := strings.Split(c.CurrentCommandLine(), " ")
	return tokens[0], tokens[1:]
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

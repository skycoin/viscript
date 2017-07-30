package task

import (
	"strings"

	"github.com/skycoin/viscript/hypervisor"
	"github.com/skycoin/viscript/msg"
)

type Cli struct {
	Log        []string
	Commands   []string
	VisualRows []string //each log entry can be multiple fragments/rows to fit current columns
	CurrCmd    int      //index
	CursPos    int      //cursor/insert position, local to command space (2 lines dedicated ATM)
	Prompt     string
	//FIXME to work with Terminal's dynamic .GridSize.X
	//assumes 64 horizontal characters, then dedicates 2 lines for each command.
	BackscrollAmount int //number of VISUAL LINES...
	//(each could be merely a SECTION of a larger (than NumColumns) log entry)
	MaxCommandSize int
}

func NewCli() *Cli {
	var cli Cli
	cli.Log = []string{}
	cli.Commands = []string{}
	cli.VisualRows = []string{}
	cli.Prompt = ">"
	cli.Commands = append(cli.Commands, cli.Prompt+"OLDEST command that you typed (not really, just an example of functionality)")
	cli.Commands = append(cli.Commands, cli.Prompt+"older command that you typed (nah, not really)")
	cli.Commands = append(cli.Commands, cli.Prompt)
	cli.CursPos = 1
	cli.CurrCmd = 2
	cli.MaxCommandSize = 128 - 1 //results in 2 lines at current initial default
	//of 64 columns.  -1 reserves space for cursor at the end of last line.

	return &cli
}

func (c *Cli) BuildRowsFromLogEntryFragments(vi msg.MessageVisualInfo) {
	//println("BuildRowsFromLogEntryFragments()   START")
	c.VisualRows = []string{}

	for _, entry := range c.Log { // 'entry' shrinks as we cut out fitting fragments
		for len(entry) > int(vi.NumColumns) {
			lff := "" /* largest fitting fragment */
			lff, entry = c.breakStringIn2(entry, int(vi.NumColumns))
			c.VisualRows = append(c.VisualRows, lff)
		}

		//last fragment is less than .NumColumns
		if /* something remains */ len(entry) > 0 {
			println("what's left of current log entry:", entry)
			c.VisualRows = append(c.VisualRows, entry) //add last fragment
		}
	}

	c.printLogInOsBox(vi)
}

func (c *Cli) printLogInOsBox(vi msg.MessageVisualInfo) {
	println("printLogInOsBox()")

	for _, entry := range c.VisualRows {
		for len(entry) < int(vi.NumColumns) {
			entry += "*"
		}

		println(entry)
	}
}

func (c *Cli) AdjustBackscrollOffset(delta int) {
	c.BackscrollAmount += delta
	println("BACKSCROLLING --- delta:", delta)

	switch {

	case c.BackscrollAmount < 0:
		c.BackscrollAmount = 0

	case c.BackscrollAmount > len(c.VisualRows):
		c.BackscrollAmount = len(c.VisualRows)

	}

	println("BACKSCROLLING --- amount:", c.BackscrollAmount)
}

func (c *Cli) InsertCharIfItFits(char uint32, state *State) {
	if len(c.Commands[c.CurrCmd]) < c.MaxCommandSize {
		c.InsertCharAtCursor(char)
		c.EchoWholeCommand(state.task.OutChannelId)
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

	m := msg.Serialize(msg.TypeCommandPrompt,
		msg.MessageCommandPrompt{termId, c.Commands[c.CurrCmd], uint32(c.CursPos)})
	hypervisor.DbusGlobal.PublishTo(outChanId, m) //EVERY publish action prefixes another chan id
}

func (c *Cli) CurrentCommandLineInLowerCase() string {
	return strings.ToLower(c.Commands[c.CurrCmd][len(c.Prompt):])
}

func (c *Cli) CurrentCommandAndArgsInLowerCase() (string, []string) {
	tokens := strings.Split(c.CurrentCommandLineInLowerCase(), " ")
	return tokens[0], tokens[1:]
}

func (c *Cli) OnEnter(st *State, serializedMsg []byte) {
	//FIXME IF we ever want more than 2 rows dedicated to command prompt.
	numRows := 1 //...to advance
	if c.CursPos >= int(st.VisualInfo.NumColumns) {
		numRows++
	}

	for numRows > 0 { //for each row of command prompt:
		numRows-- //...pass key event (value: Enter) to terminal
		//...(which advances it's y position).
		hypervisor.DbusGlobal.PublishTo(st.task.OutChannelId, serializedMsg)
	}

	//append to log history & make a "blank" new command line (which user modifies when they type)
	c.Log = append(c.Log, c.Commands[c.CurrCmd])
	c.Commands = append(c.Commands, c.Prompt)

	//action
	st.onUserCommand()

	//reset prompt & position
	c.CurrCmd = len(c.Commands) - 1
	c.CursPos = len(c.Commands[c.CurrCmd])
}

//
//
//private

//num == number of columns
//a == leftmost fragment that fits num columns
//b == remaining fragment which still may need breaking
func (c *Cli) breakStringIn2(s string, num int) (a, b string) {
	x := num - 1 //current position in x
	foundSpaceChar := false

	//fb /* feedback */ := "breakStringIn2 () "
	//println(fb+"STARTING   -   x:", x, "   -   num(columns):", num)
	//println(fb+"passed string:", s)

	for /* fragment A is smaller than a row */ len(s[x:num]) < num {
		//scan for line break
		if s[x] == ' ' {
			foundSpaceChar = true
			break
		}

		x--
	}

	//show both frags with no space removed
	//println(x, " \"", s[:x], "\" \"", s[x:], "\"")

	if foundSpaceChar {
		a = s[:x]
		b = s[x+1:] //eliminate space between final 2 pieces
	} else {
		a = s[:num]
		b = s[num:]
	}

	return a, b
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

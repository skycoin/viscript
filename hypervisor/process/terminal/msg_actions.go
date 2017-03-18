package process

import (
	"fmt"
	"strings"

	"github.com/corpusc/viscript/hypervisor"
	"github.com/corpusc/viscript/msg"
)

func (st *State) onChar(m msg.MessageChar) {
	//println("process/terminal/events.onChar()")

	if st.Cli.HasEnoughSpace() {
		// (we have free space to put character into)
		st.Cli.InsertCharAtCursor(m.Char)
		st.Cli.EchoWholeCommand(st.proc.OutChannelId)
	}
}

func (st *State) onKey(m msg.MessageKey, serializedMsg []byte) {
	//println("process/terminal/events.onKey()")

	switch msg.Action(m.Action) {

	case msg.Press:
		fallthrough
	case msg.Repeat:
		switch m.Key {

		case msg.KeyHome:
			st.Cli.CursPos = len(st.Cli.Prompt)
		case msg.KeyEnd:
			st.Cli.CursPos = len(st.Cli.Commands[st.Cli.CurrCmd])

		case msg.KeyUp:
			st.Cli.goUpCommandHistory(m.Mod)
		case msg.KeyDown:
			st.Cli.goDownCommandHistory(m.Mod)

		case msg.KeyLeft:
			st.Cli.moveOrJumpCursorLeft(m.Mod)
		case msg.KeyRight:
			st.Cli.moveOrJumpCursorRight(m.Mod)

		case msg.KeyBackspace:
			if st.Cli.moveCursorOneStepLeft() { //...succeeded
				st.Cli.DeleteCharAtCursor()
			}
		case msg.KeyDelete:
			if st.Cli.CursPos < len(st.Cli.Commands[st.Cli.CurrCmd]) {
				st.Cli.DeleteCharAtCursor()
			}

		case msg.KeyEnter:
			st.Cli.OnEnter(st, serializedMsg)
		}

		st.Cli.EchoWholeCommand(st.proc.OutChannelId)
	case msg.Release:
		// most keys will do nothing upon release
	}
}

func (st *State) onMouseScroll(m msg.MessageMouseScroll, serializedMsg []byte) {
	//println("process/terminal/events.onMouseScroll()")
	hypervisor.DbusGlobal.PublishTo(st.proc.OutChannelId, serializedMsg)
}

func (st *State) actOnCommand() {
	command, _ := st.Cli.GetCommandWithArgs()

	if len(strings.ToLower(command)) > 0 {
		st.CmdOut <- []byte(strings.ToLower(command))
	}

	// switch strings.ToLower(command) {

	// case "?":
	// 	fallthrough
	// case "h":
	// 	fallthrough
	// case "help":
	// 	st.PrintLn("Yes master, help is coming 'very soon'. (TM)")
	// default:
	// 	st.PrintLn("ERROR: \"" + command + "\" is an unknown command.")
	// }
}

func (st *State) publishToOut(message []byte) {
	hypervisor.DbusGlobal.PublishTo(st.proc.OutChannelId, message)
}

func (st *State) NewLine() {
	keyEnter := msg.MessageKey{
		Key:    msg.KeyEnter,
		Scan:   0,
		Action: uint8(msg.Action(msg.Press)),
		Mod:    0}

	st.publishToOut(msg.Serialize(msg.TypeKey, keyEnter))
}

func (st *State) PrintLn(s string) {
	for _, c := range s {
		st.sendChar(uint32(c))
	}

	st.NewLine()
}

func (st *State) Printf(format string, vars ...interface{}) {
	formattedString := fmt.Sprintf(format, vars)
	for _, c := range formattedString {
		st.sendChar(uint32(c))
	}
}

func (st *State) sendChar(c uint32) {
	var s string

	switch c {
	case msg.EscNewLine:
		st.NewLine()
		return
	case msg.EscTab:
		s = "Tab"
	case msg.EscCarriageReturn:
		s = "Carriage Return"
	case msg.EscBackSpace:
		s = "BackSpace"
	case msg.EscBackSlash:
		s = "BackSlash"
	}

	if s != "" {
		println("TASK ENCOUNTERED ESCAPE CHARACTER FOR [" + s + "], & WON'T SEND IT TO TERMINAL!")
		return
	}

	m := msg.Serialize(msg.TypePutChar, msg.MessagePutChar{0, c})
	st.publishToOut(m) // EVERY publish action prefixes another chan id
}

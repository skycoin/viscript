package task

import (
	"fmt"

	"github.com/skycoin/viscript/hypervisor"
	"github.com/skycoin/viscript/msg"
)

func (st *State) NewLine() {
	keyEnter := msg.MessageKey{
		Key:    msg.KeyEnter,
		Scan:   0,
		Action: uint8(msg.Action(msg.Press)),
		Mod:    0}

	st.publishToOut(msg.Serialize(msg.TypeKey, keyEnter))
}

func (st *State) PrintError(s string) {
	s = "**** ERROR! ****    " + s

	//to OS box 1st (more reliable)
	for i := 0; i < 4; i++ {
		println(s)
	}

	//THEN to terminal (our code is more likely to crash)
	st.PrintLn(s)
}

func (st *State) Printf(format string, vars ...interface{}) {
	s := fmt.Sprintf(format, vars...)
	for _, c := range s {
		st.sendChar(uint32(c))
	}
}

func (st *State) SendCommand(command string, args []string) {
	m := msg.Serialize(msg.TypeTokenizedCommand,
		msg.MessageTokenizedCommand{Command: command, Args: args})
	st.publishToOut(m)
}

func (st *State) PrintLn(s string) {
	st.printLnAndMAYBELogIt(s, true)
}

//
//
//private
//
//

func (st *State) printLnAndMAYBELogIt(s string, addToLog bool) {
	if addToLog {
		st.Cli.AddEntry(s, st.VisualInfo.NumColumns)
	}

	for _, c := range s {
		st.sendChar(uint32(c))
	}

	if len(s) != int(st.VisualInfo.NumColumns) {
		st.NewLine()
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
		// case msg.EscBackSlash:
		// 	s = "BackSlash"
	}

	if s != "" {
		println("TASK ENCOUNTERED ESCAPE CHAR FOR [" + s + "], NOT SENDING TO TERMINAL")
		return
	}

	m := msg.Serialize(msg.TypePutChar, msg.MessagePutChar{0, c})
	st.publishToOut(m) // EVERY publish action prefixes another chan id
}

func (st *State) publishToOut(message []byte) {
	hypervisor.DbusGlobal.PublishTo(st.task.OutChannelId, message)
}

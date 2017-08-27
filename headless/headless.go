package headless

import (
	"bufio" ///////
	"os"
	"strings"
	"time"

	"github.com/skycoin/viscript/msg"
	"github.com/skycoin/viscript/viewport/terminal"
)

var scanner *bufio.Scanner
var ch chan string = make(chan string)

func Init() {
	// println("<headless>.Init()")

	//reader := bufio.NewReader(os.Stdin)
	scanner = bufio.NewScanner(os.Stdin)
	terminal.Terms.Init()

	go func(ch chan string) {
		/*
			for {
				s, err := reader.ReadString('\n')
				if err != nil { // Maybe log non io.EOF errors, if you want
					close(ch)
					return
				}
				ch <- s
			}
		*/

		for scanner.Scan() {
			inp := scanner.Text()
			ch <- inp
			tokens := strings.Split(inp, " ")

			args := []string{}
			if len(tokens) > 1 {
				args = tokens[1:]
			}

			tc := msg.MessageTokenizedCommand{tokens[0], args}
			m := msg.Serialize(msg.TypeTokenizedCommand, tc)
			terminal.Terms.Focused.RelayToTask(m)
			terminal.Terms.DrawTextMode()
		}

		close(ch)
	}(ch)
}

func Tick() {
	select {

	case _, ok := <-ch:
		if !ok {
			println("!ok")
			return
		} else {
			terminal.Terms.Tick()
			//println("Read input from stdin:", stdin(_))
		}

	case <-time.After(10 * time.Millisecond): //throttles this tangent to 100 frames per second
		terminal.Terms.Tick()
		//println("nothing from stdin")

	}
}

package process

import (
	"os"
	"os/exec"

	"sync"

	"syscall"

	"github.com/kr/pty"
)

type ExternalProcess struct {
	Command    string
	cmd        *exec.Cmd
	currentPty *os.File
	CmdOut     chan []byte
	writeMutex *sync.Mutex

	State *State
}

func NewExternalProcess(st *State, command string) *ExternalProcess {
	println("(process/terminal/process.go).NewExternalProcess()")

	var p ExternalProcess
	p.State = st
	p.InitCmd(command)

	return &p
}

func (pr *ExternalProcess) TearDown() {
	println("(process/terminal/process.go).DeleteProcess()")
	// TODO: tear the external process down here, no remorse :rage: :D
}

func (pr *ExternalProcess) InitCmd(command string) {
	pr.Command = command
	pr.cmd = exec.Command(pr.Command)
	pr.writeMutex = &sync.Mutex{}

	var err error
	pr.currentPty, err = pty.Start(pr.cmd)
	if err != nil {
		println("Failed to execute command.")
		return
	}

	pr.CmdOut = make(chan []byte, 1024)

	exit := make(chan bool, 2)

	// Run Process Send
	go func() {
		defer func() { exit <- true }()

		pr.processSend()
	}()

	// Run Process Receive
	go func() {
		defer func() { exit <- true }()

		pr.processReceive()
	}()

	go func() {
		// TODO: defer cleanup maybe here
		// what happens if the process gets closed or we send
		// a command that makes the running command exit

		// wait for close

		// io.Copy(os.Stdout, pr.currentPty)
		<-exit
		pr.currentPty.Close()

		pr.cmd.Process.Signal(syscall.Signal(1))
		pr.cmd.Wait()
	}()
}

func (pr *ExternalProcess) processSend() {
	buf := make([]byte, 2048)

	for {
		size, err := pr.currentPty.Read(buf)
		if err != nil {
			println("%s exited.", pr.Command)
			return
		}
		pr.writeToSubscribers(buf[:size])
	}
}

func (pr *ExternalProcess) writeToSubscribers(data []byte) {
	pr.writeMutex.Lock()
	defer pr.writeMutex.Unlock()
	pr.State.PrintLn(string(data))
}

func (pr *ExternalProcess) processReceive() {
	for {
		select {
		case data := <-pr.CmdOut:
			_, err := pr.currentPty.Write(append(data, '\n'))
			if err != nil {
				return
			}
		}
	}
}

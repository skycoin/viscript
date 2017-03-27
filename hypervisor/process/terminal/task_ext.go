package process

import (
	"os/exec"

	"sync"

	"io"

	"runtime"

	"strings"

	"github.com/corpusc/viscript/app"
)

const te = "hypervisor/process/terminal/task_ext" //path

type ExternalProcess struct {
	CommandLine string //not just one command/word
	cmd         *exec.Cmd
	CmdOut      chan []byte
	writeMutex  *sync.Mutex

	stdOutPipe io.ReadCloser
	stdInPipe  io.WriteCloser

	State *State
}

//non-instanced
func MakeNewTaskExternal(st *State, tokens []string) (*ExternalProcess, error) {
	app.At(te, "MakeNewTaskExternal")
	var p ExternalProcess

	err := p.InitCmd(tokens)
	if err != nil {
		return nil, err
	}

	p.State = st

	return &p, nil
}

func (pr *ExternalProcess) TearDown() {
	app.At(te, "TearDown")
	println("TODO: tear the external process down here, no remorse :rage: :D")
}

func (pr *ExternalProcess) InitCmd(tokens []string) error {
	pr.CommandLine = strings.Join(tokens, " ")

	ros := runtime.GOOS
	if ros == "linux" || ros == "darwin" {
		pr.cmd = exec.Command(tokens[0], tokens[1:]...)
	} else if ros == "windows" {
		fullCommand := append([]string{"/C"}, tokens...)
		pr.cmd = exec.Command("cmd", fullCommand...)
	}

	pr.writeMutex = &sync.Mutex{}

	var err error
	// save stdoutpipe
	pr.stdOutPipe, err = pr.cmd.StdoutPipe()
	if err != nil {
		return err
	}

	// save stdinpipe
	pr.stdInPipe, err = pr.cmd.StdinPipe()
	if err != nil {
		return err
	}

	err = pr.cmd.Start()
	if err != nil {
		return err
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
		// TODO: what happens when user closes the application
		// and the subprocess is running still.
		// should be a way around this.
		<-exit
		pr.cmd.Wait()
		pr.State.proc.DetachExtProcess() // TODO: quick way to do it
		// pr.cmd.Process.Kill()
		// pr.cmd.Process.Signal(syscall.SIGINT)
	}()

	return nil
}

func (pr *ExternalProcess) processSend() {
	buf := make([]byte, 2048)

	for {
		size, err := pr.stdOutPipe.Read(buf)
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
			_, err := pr.stdInPipe.Write(append(data, '\n'))
			if err != nil {
				return
			}
		}
	}
}

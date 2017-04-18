package task_ext

import (
	"errors"

	"io"

	"runtime"

	"strings"

	"os/exec"

	"fmt"

	"github.com/corpusc/viscript/app"
	"github.com/corpusc/viscript/msg"
)

const te = "hypervisor/task_ext/task_ext" //path

type ExternalProcess struct {
	Id          msg.ExtProcessId
	CommandLine string

	ProcessIn   chan []byte
	ProcessOut  chan []byte
	ProcessExit chan bool // if process needs to exit without user interruption

	CmdOut chan []byte
	CmdIn  chan []byte

	cmd        *exec.Cmd
	stdOutPipe io.ReadCloser
	stdInPipe  io.WriteCloser

	shutdown chan struct{}

	routinesStarted bool
}

//non-instanced
func MakeNewTaskExternal(tokens []string, detached bool) (*ExternalProcess, error) {
	app.At(te, "MakeNewTaskExternal")
	var p ExternalProcess

	err := p.Init(tokens)
	if err != nil {
		return nil, err
	}

	return &p, nil
}

func (pr *ExternalProcess) GetExtProcessInterface() msg.ExtProcessInterface {
	return msg.ExtProcessInterface(pr)
}

func (pr *ExternalProcess) Init(tokens []string) error {
	app.At(te, "Init")

	var err error

	if pr.cmd, err = pr.createCMDAccordingToOS(tokens); err != nil {
		return err
	}

	if pr.stdOutPipe, err = pr.cmd.StdoutPipe(); err != nil {
		return err
	}

	if pr.stdInPipe, err = pr.cmd.StdinPipe(); err != nil {
		return err
	}

	pr.Id = msg.NextExtProcessId()

	//Creates a new process group for the new process
	//to avoid leaving orphan processes.
	pr.CommandLine = strings.Join(tokens, " ")

	pr.CmdOut = make(chan []byte, 2048)
	pr.CmdIn = make(chan []byte, 2048)
	pr.ProcessIn = make(chan []byte, 2048)
	pr.ProcessOut = make(chan []byte, 2048)
	pr.ProcessExit = make(chan bool)
	pr.shutdown = make(chan struct{})
	pr.routinesStarted = false

	return nil
}

func (pr *ExternalProcess) createCMDAccordingToOS(tokens []string) (*exec.Cmd, error) {
	app.At(te, "createCMDAccordingToOS")

	ros := runtime.GOOS
	if ros == "linux" || ros == "darwin" {
		return exec.Command(tokens[0], tokens[1:]...), nil
	} else if ros == "windows" {
		fullCommand := append([]string{"/C"}, tokens...)
		return exec.Command("cmd", fullCommand...), nil
	}

	return nil, errors.New("Unknown Operating System. Aborting command initilization")
}

func (pr *ExternalProcess) cmdInRoutine() {
	app.At(te, "cmdInRoutine")

	for {
		buf := make([]byte, 2048)
		size, err := pr.stdOutPipe.Read(buf[:])
		if err != nil {
			println("Cmd In Routine error:", err.Error())
			pr.ProcessExit <- true
			close(pr.shutdown)
			return
		}

		select {
		case <-pr.shutdown:
			println("!!! Shutting cmdInRoutine down !!!")
			return
		case pr.CmdIn <- buf[:size]:
			fmt.Printf("-- Received data for sending to CmdIn: %s\n",
				string(buf[:size]))
		}
	}
}

func (pr *ExternalProcess) cmdOutRoutine() {
	app.At(te, "cmdOutRoutine")

	for {
		select {
		case <-pr.shutdown:
			println("!!! Shutting cmdOutRoutine down !!!")
			return
		case data := <-pr.CmdOut:
			fmt.Printf("-- Received input to write to external process: %s\n",
				string(data))
			_, err := pr.stdInPipe.Write(append(data, '\n'))
			if err != nil {
				println("!!! Couldn't Write To the std in pipe of the process !!!")
				pr.ProcessExit <- true
				return
			}
		}
	}
}

func (pr *ExternalProcess) startRoutines() error {

	if pr.stdOutPipe == nil {
		return errors.New("Standard out pipe of process is nil")
	}

	if pr.stdInPipe == nil {
		return errors.New("Standard in pipe of process is nil")
	}

	if !pr.routinesStarted {
		pr.shutdown = make(chan struct{})

		// Run the routine which will read and send the data to CmdIn
		go pr.cmdInRoutine()

		// Run the routine which will read from Cmdout and write to process
		go pr.cmdOutRoutine()

		pr.routinesStarted = true
	}

	return nil
}

func (pr *ExternalProcess) stopRoutines() {
	close(pr.shutdown)
}

func (pr *ExternalProcess) processOutput() {
	select {
	case data := <-pr.ProcessIn:
		println("ProcessOutput() - data := ", string(data))
		pr.CmdOut <- data
	default:
	}
}

func (pr *ExternalProcess) processInput() {
	select {
	case data := <-pr.CmdIn:
		println("ProcessInput() - data := ", string(data))
		pr.ProcessOut <- data
	default:
	}
}

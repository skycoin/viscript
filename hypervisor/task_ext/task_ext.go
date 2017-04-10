package task_ext

import (
	"errors"

	"io"

	"runtime"

	"strings"

	"os/exec"

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

	detached bool
}

//non-instanced
func MakeNewTaskExternal(tokens []string, detached bool) (*ExternalProcess, error) {
	app.At(te, "MakeNewTaskExternal")
	var p ExternalProcess

	err := p.Init(tokens)
	if err != nil {
		return nil, err
	}

	p.detached = detached
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
		println("!!! INSIDE THE CMD IN !!!")
		if pr.detached {
			println("!!! DETACHED IS TRUE IN CMDIN !!!")
			return
		}

		if pr.stdOutPipe == nil {
			println("!!! Standard output pipe is nil. Sending Exit Request !!!")
			pr.ProcessExit <- true
			pr.detached = true
			return
		}

		buf := make([]byte, 2048)

		size, err := pr.stdOutPipe.Read(buf)
		if err != nil {
			// s := fmt.Sprintf("**** ERROR! **** From \"%s\".  Returning. %s", pr.CommandLine, err.Error())
			// for i := 0; i < 5; i++ {
			// 	println(s) //to OS box
			// }
			pr.ProcessExit <- true
			pr.detached = true
			return
		}

		println("--- Received input to write to the terminal:", string(buf[:size]))
		pr.CmdIn <- buf[:size]

	}
}

func (pr *ExternalProcess) cmdOutRoutine() {
	app.At(te, "cmdOutRoutine")

	for {
		println("!!! INSIDE THE CMD OUT !!!")
		if pr.detached {
			println("!!! DETACHED IS TRUE IN CMD OUT !!!")
			return
		}

		select {
		case data := <-pr.CmdOut:

			if pr.stdInPipe == nil {
				println("!!! Standard input pipe is nil. Sending Exit Request !!!")
				pr.ProcessExit <- true
				pr.detached = true
				return
			}

			println("--- Received input to write to external process:", string(data))
			_, err := pr.stdInPipe.Write(append(data, '\n'))
			if err != nil {
				pr.ProcessExit <- true
				pr.detached = true
				return
			}
		}
	}
}

func (pr *ExternalProcess) startRoutines() {
	// Run the routine which will read and send the data to CmdIn
	go pr.cmdInRoutine()
	// Run the routine which will read from Cmdout and write to process
	go pr.cmdOutRoutine()
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

package extprocess

import (
	"errors"
	"fmt"

	"io"

	"runtime"

	"strings"

	"os/exec"

	"syscall"

	"github.com/corpusc/viscript/app"
	"github.com/corpusc/viscript/msg"
)

const te = "hypervisor/process/terminal/task_ext" //path

type ExternalProcess struct {
	CommandLine string

	ProcessIn   chan []byte
	ProcessOut  chan []byte
	ProcessExit chan bool

	CmdOut chan []byte
	CmdIn  chan []byte

	cmd        *exec.Cmd
	stdOutPipe io.ReadCloser
	stdInPipe  io.WriteCloser
}

//non-instanced
func MakeNewTaskExternal(tokens []string) (*ExternalProcess, error) {
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

	// Creates a new process group for the new process
	// to avoid leaving orphan processes.
	pr.cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}
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

	for pr.stdOutPipe != nil {
		buf := make([]byte, 2048)

		size, err := pr.stdOutPipe.Read(buf)
		if err != nil {
			f := err.Error()
			s := fmt.Sprintf("**** ERROR! **** From \"%s\".  Returning. %s", pr.CommandLine, f)
			for i := 0; i < 5; i++ {
				println(s) //OS box print
			}
			// close(pr.CmdIn)
			return
		}

		pr.CmdIn <- buf[:size]
	}
}

func (pr *ExternalProcess) cmdOutRoutine() {
	app.At(te, "cmdOutRoutine")

	for pr.stdInPipe != nil {
		select {
		case data := <-pr.CmdOut:
			println("RECEIVED ____ ", string(data), " IN CMDOUTROUTINE")
			_, err := pr.stdInPipe.Write(append(data, '\n'))
			if err != nil {
				return
			}
		}
	}
}

func (pr *ExternalProcess) ProcessOutput() {
	for len(pr.ProcessIn) > 0 {
		// println("ProcessOutput() - data := ", string(data))
		data := <-pr.ProcessIn
		pr.CmdOut <- data
	}
}

func (pr *ExternalProcess) ProcessInput() {
	for len(pr.CmdIn) > 0 {
		// println("ProcessInput() - data := ", string(data))
		data := <-pr.CmdIn
		pr.ProcessOut <- data
	}
}

package process

import (
	"errors"
	"fmt"

	"io"

	"runtime"

	"strings"

	"os/exec"

	"syscall"

	"github.com/corpusc/viscript/app"
)

const te = "hypervisor/process/terminal/task_ext" //path

type ExternalProcess struct {
	CommandLine string

	ProcessIn chan []byte
	CmdOut    chan []byte
	CmdIn     chan []byte

	cmd        *exec.Cmd
	stdOutPipe io.ReadCloser
	stdInPipe  io.WriteCloser

	State *State
}

//non-instanced
func MakeNewTaskExternal(st *State, tokens []string) (*ExternalProcess, error) {
	app.At(te, "MakeNewTaskExternal")
	var p ExternalProcess

	err := p.Init(tokens)
	if err != nil {
		return nil, err
	}

	p.State = st

	return &p, nil
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

func (pr *ExternalProcess) Start() error {
	app.At(te, "Start")

	err := pr.cmd.Start()
	if err != nil {
		return err
	}

	// Run the routine which will read and send the data to CmdIn
	go pr.cmdInRoutine()
	// Run the routine which will read from Cmdout and write to process
	go pr.cmdOutRoutine()

	return nil
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

func (pr *ExternalProcess) Tick() {
	// var err error
	pr.ProcessInput()
	// if err = pr.ProcessInput(); err != nil {
	// 	// return err
	// 	app.At(te, "ProcessInput() - returned error: "+err.Error()+" ! Deleting Process")
	// 	pr.State.proc.DeleteAttachedExtProcess()
	// }
	pr.ProcessOutput()
	// if err = pr.ProcessOutput(); err != nil {
	// 	// return err
	// 	app.At(te, "ProcessOutput() - returned error: "+err.Error()+" ! Deleting Process")
	// 	pr.State.proc.DeleteAttachedExtProcess()
	// }
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
		pr.State.PrintLn(string(data))
	}
}

func (pr *ExternalProcess) ShutDown() {
	app.At(te, "ShutDown")
	close(pr.CmdOut)
	pr.cmd.Process.Kill()
	pr.cmd = nil
	pr.stdOutPipe = nil
	pr.stdInPipe = nil
	pr.State = nil
}

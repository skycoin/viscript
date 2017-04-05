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
	Id          msg.ExtProcessId
	CommandLine string

	ProcessIn   chan []byte
	ProcessOut  chan []byte
	ProcessExit chan bool // if process needs to exit without user interruption

	ProcessQuit chan bool // when process should end upon user's command

	CmdOut chan []byte
	CmdIn  chan []byte

	cmd        *exec.Cmd
	stdOutPipe io.ReadCloser
	stdInPipe  io.WriteCloser

	runningInBg bool
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

	pr.Id = msg.NextExtProcessId()

	// Creates a new process group for the new process
	// to avoid leaving orphan processes.
	pr.cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}
	pr.CommandLine = strings.Join(tokens, " ")

	pr.CmdOut = make(chan []byte, 2048)
	pr.CmdIn = make(chan []byte, 2048)
	pr.ProcessIn = make(chan []byte, 2048)
	pr.ProcessOut = make(chan []byte, 2048)
	pr.ProcessExit = make(chan bool)
	pr.ProcessQuit = make(chan bool)

	// TODO: maybe passed as an arg on the start?
	pr.runningInBg = true

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

		select {
		case shouldQuit := <-pr.ProcessQuit:
			if shouldQuit == true {
				// In this case process should:
				// 1) Send true to the ProcessShouldEnd channel
				// 2) Call external process's TearDown
				// 3) Remove it from the GlobalList and further cleanup whatever will be needed
				return
			}
		default:
		}

		if !pr.runningInBg && pr.stdOutPipe != nil {
			if pr.stdOutPipe == nil {
				println("!!! Standard output pipe is nil. Sending Exit Request !!!")
				pr.ProcessExit <- true
				return
			}
			buf := make([]byte, 2048)

			size, err := pr.stdOutPipe.Read(buf)
			if err != nil {
				f := err.Error()
				s := fmt.Sprintf("**** ERROR! **** From \"%s\".  Returning. %s", pr.CommandLine, f)
				for i := 0; i < 5; i++ {
					println(s) //OS box print
				}
				pr.ProcessExit <- true
				return
			}
			println("--- Received input to write to the terminal:", string(buf[:size]))
			pr.CmdIn <- buf[:size]
		}
	}
}

func (pr *ExternalProcess) cmdOutRoutine() {
	app.At(te, "cmdOutRoutine")

	for {
		select {
		case shouldEnd := <-pr.ProcessQuit:
			if shouldEnd == true {
				return
			}
		case data := <-pr.CmdOut:
			if !pr.runningInBg {

				if pr.stdInPipe == nil {
					println("!!! Standard input pipe is nil. Sending Exit Request !!!")
					pr.ProcessExit <- true
					return
				}

				println("--- Received input to write to external process:", string(data))
				_, err := pr.stdInPipe.Write(append(data, '\n'))
				if err != nil {
					pr.ProcessExit <- true
					return
				}
			}
		default:
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

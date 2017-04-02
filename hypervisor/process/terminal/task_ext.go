package process

import (
	"fmt"

	"sync"

	"io"

	"runtime"

	"strings"

	"os/exec"

	"syscall"

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

	RunningInBg bool

	State *State

	shouldEnd bool
}

//non-instanced
func MakeNewTaskExternal(st *State, tokens []string) (*ExternalProcess, error) {
	app.At(te, "MakeNewTaskExternal")
	var p ExternalProcess

	err := p.InitCmd(tokens)
	if err != nil {
		return nil, err
	}

	p.shouldEnd = false
	p.RunningInBg = false
	p.State = st

	return &p, nil
}

func (pr *ExternalProcess) TearDown() {
	app.At(te, "TearDown")
	close(pr.CmdOut)
	pr.cmd.Process.Kill()
	pr.cmd = nil
	pr.writeMutex = nil
	pr.stdOutPipe = nil
	pr.stdInPipe = nil
	pr.State = nil
	// syscall.Kill(-pr.cmd.Process.Pid, syscall.SIGKILL)
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

	// Creates a new process group for the new process
	// to avoid leaving orphan processes.
	pr.cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}
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
		// does external process become an orphan process?
		// should be a way around this.
		<-exit
		pr.cmd.Wait()
		_ = pr.State.proc.DeleteAttachedExtProcess()
		// pr.cmd.Process.Kill()
		// pr.cmd.Process.Signal(syscall.SIGINT)
	}()

	return nil
}

func (pr *ExternalProcess) processSend() {
	buf := make([]byte, 2048)

	for !pr.shouldEnd {
		if pr.stdOutPipe == nil {
			return
		}

		size, err := pr.stdOutPipe.Read(buf)
		if err != nil {
			s := fmt.Sprintf("**** ERROR! ****    %s\n", err.Error())
			s += fmt.Sprintf("**** ERROR! ****    (command: \"%s\").  Returning.", pr.CommandLine)

			for i := 0; i < 5; i++ {
				println(s)
			}

			// having an err set to something means the stdOutPipe was closed or process was finished
			// unable to read again. I'll look more into the Read func doc, just to be sure.
			return
		}

		if !pr.RunningInBg {
			pr.writeToSubscribers(buf[:size])
		}
	}
}

func (pr *ExternalProcess) writeToSubscribers(data []byte) {
	pr.writeMutex.Lock()
	defer pr.writeMutex.Unlock()
	pr.State.PrintLn(string(data))
}

func (pr *ExternalProcess) processReceive() {
	for !pr.shouldEnd {
		if pr.stdInPipe == nil {
			return
		}
		select {
		case data := <-pr.CmdOut:
			_, err := pr.stdInPipe.Write(append(data, '\n'))
			if err != nil {
				return
			}
		}
	}
}

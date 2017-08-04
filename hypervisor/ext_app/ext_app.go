package ext_app

import (
	"errors"
	"sync"

	"io"

	"runtime"

	"strings"

	"os/exec"

	"fmt"

	"strconv"

	"github.com/skycoin/viscript/app"
	"github.com/skycoin/viscript/msg"
)

const te = "hypervisor/ext_app/ext_app" //path

type ExternalTask struct {
	Id          msg.ExtAppId
	CommandLine string

	TaskIn   chan []byte
	TaskOut  chan []byte
	TaskExit chan struct{} //this way it's easy to cleanup multiple places

	cmdOut chan []byte
	cmdIn  chan []byte

	cmd        *exec.Cmd
	stdOutPipe io.ReadCloser
	stdInPipe  io.WriteCloser

	shutdown chan struct{}

	routinesStarted bool

	wg sync.WaitGroup
}

//non-instanced
func MakeNewTaskExternal(tokens []string, detached bool) (*ExternalTask, error) {
	app.At(te, "MakeNewTaskExternal")
	var et ExternalTask

	err := et.Init(tokens)
	if err != nil {
		return nil, err
	}

	return &et, nil
}

func (pr *ExternalTask) GetExtTaskInterface() msg.ExtTaskInterface {
	return msg.ExtTaskInterface(pr)
}

func (pr *ExternalTask) Init(tokens []string) error {
	app.At(te, "Init")

	var err error

	pr.Id = msg.NextExtTaskId()

	//append app id before creating command
	tokens = append(tokens, "-signal-client-id")
	tokens = append(tokens, strconv.Itoa(int(pr.Id)))
	tokens = append(tokens, "-signal-client")

	//TODO: think about this here if we have daemon should we attach anything?

	if pr.cmd, err = pr.createCMDAccordingToOS(tokens); err != nil {
		return err
	}

	if pr.stdOutPipe, err = pr.cmd.StdoutPipe(); err != nil {
		return err
	}

	if pr.stdInPipe, err = pr.cmd.StdinPipe(); err != nil {
		return err
	}

	pr.CommandLine = strings.Join(tokens, " ")

	pr.cmdOut = make(chan []byte, 2048)
	pr.cmdIn = make(chan []byte, 2048)

	pr.TaskIn = make(chan []byte, 2048)
	pr.TaskOut = make(chan []byte, 2048)
	pr.TaskExit = make(chan struct{})

	pr.shutdown = make(chan struct{})

	pr.routinesStarted = false

	return nil
}

func (pr *ExternalTask) createCMDAccordingToOS(tokens []string) (*exec.Cmd, error) {
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

func (pr *ExternalTask) cmdInRoutine() {
	app.At(te, "cmdInRoutine")

	for {
		buf := make([]byte, 2048)
		size, err := pr.stdOutPipe.Read(buf[:])
		if err != nil {
			println("Cmd In Routine error:", err.Error())
			close(pr.TaskExit)
			close(pr.shutdown)
			return
		}

		select {
		case <-pr.shutdown:
			println("!!! Shutting cmdInRoutine down !!!")
			return
		case pr.cmdIn <- buf[:size]:
			fmt.Printf("-- Received data for sending to CmdIn: %s\n",
				string(buf[:size]))
		}
	}
}

func (pr *ExternalTask) cmdOutRoutine() {
	app.At(te, "cmdOutRoutine")

	for {
		select {
		case <-pr.shutdown:
			println("!!! Shutting cmdOutRoutine down !!!")
			return
		case data := <-pr.cmdOut:
			fmt.Printf("-- Received input to write to external task: %s\n",
				string(data))
			_, err := pr.stdInPipe.Write(append(data, '\n'))
			if err != nil {
				println("!!! Couldn't Write To the std in pipe of the task!!!")
				close(pr.TaskExit)
				close(pr.shutdown)
				return
			}
		}
	}
}

func (pr *ExternalTask) startRoutines() error {

	if pr.stdOutPipe == nil {
		return errors.New("Standard out pipe of task is nil")
	}

	if pr.stdInPipe == nil {
		return errors.New("Standard in pipe of task is nil")
	}

	if !pr.routinesStarted {

		pr.wg = sync.WaitGroup{}

		pr.TaskExit = make(chan struct{})
		pr.shutdown = make(chan struct{})

		pr.wg.Add(2)

		//Run the routine which will read and send the data to CmdIn
		go pr.cmdInRoutine()

		//Run the routine which will read from Cmdout and write to task
		go pr.cmdOutRoutine()

		pr.routinesStarted = true
	}

	return nil
}

func (pr *ExternalTask) stopRoutines() {
	close(pr.shutdown)
}

func (pr *ExternalTask) taskOutput() {
	select {
	case data := <-pr.TaskIn:
		println("taskOutput() - data := ", string(data))
		pr.cmdOut <- data
	default:
	}
}

func (pr *ExternalTask) taskInput() {
	select {
	case data := <-pr.cmdIn:
		println("taskInput() - data := ", string(data))
		pr.TaskOut <- data
	default:
	}
}
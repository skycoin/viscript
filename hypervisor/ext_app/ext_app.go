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

const path = "hypervisor/ext_app/ext_app"

type ExternalApp struct {
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
func MakeNewTaskExternal(tokens []string, detached bool) (*ExternalApp, error) {
	app.At(path, "MakeNewTaskExternal")
	var ea ExternalApp

	err := ea.Init(tokens)
	if err != nil {
		return nil, err
	}

	return &ea, nil
}

func (ea *ExternalApp) GetExtTaskInterface() msg.ExtTaskInterface {
	return msg.ExtTaskInterface(ea)
}

func (ea *ExternalApp) Init(tokens []string) error {
	app.At(path, "Init")

	var err error

	ea.Id = msg.NextExtTaskId()

	//append app id before creating command
	tokens = append(tokens, "-signal-client-id")
	tokens = append(tokens, strconv.Itoa(int(ea.Id)))
	tokens = append(tokens, "-signal-client")

	//TODO: think about this here if we have daemon should we attach anything?

	if ea.cmd, err = ea.createCMDAccordingToOS(tokens); err != nil {
		return err
	}

	if ea.stdOutPipe, err = ea.cmd.StdoutPipe(); err != nil {
		return err
	}

	if ea.stdInPipe, err = ea.cmd.StdinPipe(); err != nil {
		return err
	}

	ea.CommandLine = strings.Join(tokens, " ")

	ea.cmdOut = make(chan []byte, 2048)
	ea.cmdIn = make(chan []byte, 2048)

	ea.TaskIn = make(chan []byte, 2048)
	ea.TaskOut = make(chan []byte, 2048)
	ea.TaskExit = make(chan struct{})

	ea.shutdown = make(chan struct{})

	ea.routinesStarted = false

	return nil
}

func (ea *ExternalApp) createCMDAccordingToOS(tokens []string) (*exec.Cmd, error) {
	app.At(path, "createCMDAccordingToOS")

	ros := runtime.GOOS
	if ros == "linux" || ros == "darwin" {
		return exec.Command(tokens[0], tokens[1:]...), nil
	} else if ros == "windows" {
		fullCommand := append([]string{"/C"}, tokens...)
		return exec.Command("cmd", fullCommand...), nil
	}

	return nil, errors.New("Unknown Operating System. Aborting command initilization")
}

func (ea *ExternalApp) cmdInRoutine() {
	app.At(path, "cmdInRoutine")

	for {
		buf := make([]byte, 2048)
		size, err := ea.stdOutPipe.Read(buf[:])
		if err != nil {
			println("Cmd In Routine error:", err.Error())
			close(ea.TaskExit)
			close(ea.shutdown)
			return
		}

		select {
		case <-ea.shutdown:
			println("!!! Shutting cmdInRoutine down !!!")
			return
		case ea.cmdIn <- buf[:size]:
			fmt.Printf("-- Received data for sending to CmdIn: %s\n",
				string(buf[:size]))
		}
	}
}

func (ea *ExternalApp) cmdOutRoutine() {
	app.At(path, "cmdOutRoutine")

	for {
		select {
		case <-ea.shutdown:
			println("!!! Shutting cmdOutRoutine down !!!")
			return
		case data := <-ea.cmdOut:
			fmt.Printf("-- Received input to write to external task: %s\n",
				string(data))
			_, err := ea.stdInPipe.Write(append(data, '\n'))
			if err != nil {
				println("!!! Couldn't Write To the std in pipe of the task!!!")
				close(ea.TaskExit)
				close(ea.shutdown)
				return
			}
		}
	}
}

func (ea *ExternalApp) startRoutines() error {
	if ea.stdOutPipe == nil {
		return errors.New("Standard out pipe of task is nil")
	}

	if ea.stdInPipe == nil {
		return errors.New("Standard in pipe of task is nil")
	}

	if !ea.routinesStarted {
		ea.wg = sync.WaitGroup{}
		ea.TaskExit = make(chan struct{})
		ea.shutdown = make(chan struct{})
		ea.wg.Add(2)

		//Run the routine which will read and send the data to CmdIn
		go ea.cmdInRoutine()

		//Run the routine which will read from Cmdout and write to task
		go ea.cmdOutRoutine()

		ea.routinesStarted = true
	}

	return nil
}

func (ea *ExternalApp) stopRoutines() {
	close(ea.shutdown)
}

func (ea *ExternalApp) taskOutput() {
	select {

	case data := <-ea.TaskIn:
		println("taskOutput() - data := ", string(data))
		ea.cmdOut <- data
	default:

	}
}

func (ea *ExternalApp) taskInput() {
	select {

	case data := <-ea.cmdIn:
		println("taskInput() - data := ", string(data))
		ea.TaskOut <- data
	default:

	}
}

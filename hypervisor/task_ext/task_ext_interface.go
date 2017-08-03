package ext_app

import (
	"github.com/skycoin/viscript/app"
	"github.com/skycoin/viscript/msg"
)

//ExtAppInterface implementation

func (et *ExternalTask) Tick() {
	et.taskInput()
	et.taskOutput()
}

func (et *ExternalTask) Start() error {
	app.At(te, "Start")

	err := et.cmd.Start()
	if err != nil {
		return err
	}

	return nil
}

func (et *ExternalTask) TearDown() {
	app.At(te, "TearDown")

	et.cmd.Process.Kill()

	close(et.cmdIn)
	close(et.cmdOut)

	close(et.TaskIn)
	close(et.TaskOut)
	// close(et.TaskExit)

	if et.cmd != nil {
		et.cmd = nil
	}

	if et.stdOutPipe != nil {
		et.stdOutPipe = nil
	}

	if et.stdInPipe != nil {
		et.stdInPipe = nil
	}
}

func (et *ExternalTask) Attach() error {
	app.At(te, "Attach")
	return et.startRoutines()
}

func (et *ExternalTask) Detach() {
	app.At(te, "Detach")
	// TODO: detach using channels maybe
	et.stopRoutines()
}

func (et *ExternalTask) GetId() msg.ExtAppId {
	return et.Id
}

func (et *ExternalTask) GetFullCommandLine() string {
	return et.CommandLine
}

func (et *ExternalTask) GetTaskInChannel() chan []byte {
	return et.TaskIn
}

func (et *ExternalTask) GetTaskOutChannel() chan []byte {
	return et.TaskOut
}

func (et *ExternalTask) GetTaskExitChannel() chan struct{} {
	return et.TaskExit
}

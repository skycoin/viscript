package ext_app

import (
	"github.com/skycoin/viscript/app"
	"github.com/skycoin/viscript/msg"
)

//ExtAppInterface implementation

func (ea *ExternalApp) Tick() {
	et.taskInput()
	et.taskOutput()
}

func (ea *ExternalApp) Start() error {
	app.At(te, "Start")

	err := et.cmd.Start()
	if err != nil {
		return err
	}

	return nil
}

func (ea *ExternalApp) TearDown() {
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

func (ea *ExternalApp) Attach() error {
	app.At(te, "Attach")
	return et.startRoutines()
}

func (ea *ExternalApp) Detach() {
	app.At(te, "Detach")
	// TODO: detach using channels maybe
	et.stopRoutines()
}

func (ea *ExternalApp) GetId() msg.ExtAppId {
	return et.Id
}

func (ea *ExternalApp) GetFullCommandLine() string {
	return et.CommandLine
}

func (ea *ExternalApp) GetTaskInChannel() chan []byte {
	return et.TaskIn
}

func (ea *ExternalApp) GetTaskOutChannel() chan []byte {
	return et.TaskOut
}

func (ea *ExternalApp) GetTaskExitChannel() chan struct{} {
	return et.TaskExit
}

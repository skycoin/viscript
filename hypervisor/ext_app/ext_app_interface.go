package ext_app

import (
	"github.com/skycoin/viscript/app"
	"github.com/skycoin/viscript/msg"
)

//ExtAppInterface implementation

func (ea *ExternalApp) Tick() {
	ea.taskInput()
	ea.taskOutput()
}

func (ea *ExternalApp) Start() error {
	app.At(path, "Start")

	err := ea.cmd.Start()
	if err != nil {
		return err
	}

	return nil
}

func (ea *ExternalApp) TearDown() {
	app.At(path, "TearDown")

	ea.cmd.Process.Kill()

	close(ea.cmdIn)
	close(ea.cmdOut)

	close(ea.TaskIn)
	close(ea.TaskOut)
	// close(ea.TaskExit)

	if ea.cmd != nil {
		ea.cmd = nil
	}

	if ea.stdOutPipe != nil {
		ea.stdOutPipe = nil
	}

	if ea.stdInPipe != nil {
		ea.stdInPipe = nil
	}
}

func (ea *ExternalApp) Attach() error {
	app.At(path, "Attach")
	return ea.startRoutines()
}

func (ea *ExternalApp) Detach() {
	app.At(path, "Detach")
	// TODO: detach using channels maybe
	ea.stopRoutines()
}

func (ea *ExternalApp) GetId() msg.ExternalAppId {
	return ea.Id
}

func (ea *ExternalApp) GetFullCommandLine() string {
	return ea.CommandLine
}

func (ea *ExternalApp) GetTaskInChannel() chan []byte {
	return ea.TaskIn
}

func (ea *ExternalApp) GetTaskOutChannel() chan []byte {
	return ea.TaskOut
}

func (ea *ExternalApp) GetTaskExitChannel() chan struct{} {
	return ea.TaskExit
}

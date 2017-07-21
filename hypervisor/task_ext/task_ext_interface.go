package task_ext

import (
	"github.com/skycoin/viscript/app"
	"github.com/skycoin/viscript/msg"
)

//ExtTaskInterface implementation

func (pr *ExternalTask) Tick() {
	pr.taskInput()
	pr.taskOutput()
}

func (pr *ExternalTask) Start() error {
	app.At(te, "Start")

	err := pr.cmd.Start()
	if err != nil {
		return err
	}

	return nil
}

func (pr *ExternalTask) TearDown() {
	app.At(te, "TearDown")

	pr.cmd.Process.Kill()

	close(pr.cmdIn)
	close(pr.cmdOut)

	close(pr.TaskIn)
	close(pr.TaskOut)
	// close(pr.ProcessExit)

	if pr.cmd != nil {
		pr.cmd = nil
	}

	if pr.stdOutPipe != nil {
		pr.stdOutPipe = nil
	}

	if pr.stdInPipe != nil {
		pr.stdInPipe = nil
	}
}

func (pr *ExternalTask) Attach() error {
	app.At(te, "Attach")
	return pr.startRoutines()
}

func (pr *ExternalTask) Detach() {
	app.At(te, "Detach")
	// TODO: detach using channels maybe
	pr.stopRoutines()
}

func (pr *ExternalTask) GetId() msg.ExtTaskId {
	return pr.Id
}

func (pr *ExternalTask) GetFullCommandLine() string {
	return pr.CommandLine
}

func (pr *ExternalTask) GetTaskInChannel() chan []byte {
	return pr.TaskIn
}

func (pr *ExternalTask) GetTaskOutChannel() chan []byte {
	return pr.TaskOut
}

func (pr *ExternalTask) GetTaskExitChannel() chan struct{} {
	return pr.ProcessExit
}

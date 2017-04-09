package task_ext

import (
	"github.com/corpusc/viscript/app"
	"github.com/corpusc/viscript/msg"
)

//ExtProcessInterface implementation

func (pr *ExternalProcess) Tick() {
	pr.processInput()
	pr.processOutput()
}

func (pr *ExternalProcess) Start() error {
	app.At(te, "Start")

	err := pr.cmd.Start()
	if err != nil {
		return err
	}

<<<<<<< refs/remotes/corpusc/master:hypervisor/task_ext/task_ext_int.go
	//Run the routine which will read and send the data to CmdIn
	go pr.cmdInRoutine()
	//Run the routine which will read from Cmdout and write to process
	go pr.cmdOutRoutine()
=======
	if !pr.detached {
		pr.startRoutines()
	}
>>>>>>> ...:hypervisor/extprocess/task_extimpl.go

	return nil
}

func (pr *ExternalProcess) TearDown() {
	app.At(te, "TearDown")

	pr.cmd.Process.Kill()

	close(pr.CmdIn)
	close(pr.CmdOut)

	close(pr.ProcessIn)
	close(pr.ProcessOut)
	close(pr.ProcessExit)

	pr.cmd = nil
	pr.stdOutPipe = nil
	pr.stdInPipe = nil
}

func (pr *ExternalProcess) Attach() {
	app.At(te, "Attach")
	pr.detached = false
	pr.startRoutines()
}

func (pr *ExternalProcess) Detach() {
	app.At(te, "Detach")
	pr.detached = true
}

func (pr *ExternalProcess) GetId() msg.ExtProcessId {
	return pr.Id
}

func (pr *ExternalProcess) GetFullCommandLine() string {
	return pr.CommandLine
}

func (pr *ExternalProcess) GetProcessInChannel() chan []byte {
	return pr.ProcessIn
}

func (pr *ExternalProcess) GetProcessOutChannel() chan []byte {
	return pr.ProcessOut
}

func (pr *ExternalProcess) GetProcessExitChannel() chan bool {
	return pr.ProcessExit
}

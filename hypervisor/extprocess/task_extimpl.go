package extprocess

import (
	"github.com/corpusc/viscript/app"
	"github.com/corpusc/viscript/msg"
)

// ExtProcessInterface implementation

func (pr *ExternalProcess) Tick() {
	pr.ProcessInput()
	pr.ProcessOutput()
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

func (pr *ExternalProcess) TearDown() {
	app.At(te, "ShutDown")

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

func (pr *ExternalProcess) GetId() msg.ExtProcessId {
	return pr.Id
}

func (pr *ExternalProcess) GetRunningInBg() bool {
	return pr.runningInBg
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

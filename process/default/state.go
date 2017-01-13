package process

import ()

//put all your process state here
type State struct {
	proc *DefaultProcess
}

func (self *State) InitState(proc *DefaultProcess) {
	self.proc = proc
}

func HandleMessages() {

}

//p.MessageIn = make(chan []byte)
//p.MessageOut = make(chan []byte)

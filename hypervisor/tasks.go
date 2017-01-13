package hypervisor

import ()

var Tasks *[]Task

type Task struct {
	In       [][]byte
	Out      [][]byte
	Terminal *Terminal // if nil, no visual needed
}

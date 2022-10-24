package core

import "sync"

type CancelCmds struct {
	tasks map[int32]bool
}

func (cmd *CancelCmds) Insert(id int32) {
	var m sync.Mutex
	m.Lock()
	cmd.tasks[id] = false
	m.Unlock()
}

func (cmd *CancelCmds) Cancel(id int32) {
	var m sync.Mutex
	m.Lock()
	cmd.tasks[id] = true
	m.Unlock()
}

func (cmd *CancelCmds) GetCancelStatus(id int32) bool {
	return cmd.tasks[id]
}

func (cmd *CancelCmds) remove(id int32) {
	var m sync.Mutex
	m.Lock()
	delete(cmd.tasks, id)
	m.Unlock()
}

func NewCancelCmd() *CancelCmds {
	return &CancelCmds{
		tasks: make(map[int32]bool),
	}
}

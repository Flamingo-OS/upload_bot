package core

import "sync"

type CancelCmds struct {
	tasks map[uint64]bool
}

func (cmd *CancelCmds) Insert(id uint64) {
	var m sync.Mutex
	m.Lock()
	cmd.tasks[id] = false
	m.Unlock()
}

func (cmd *CancelCmds) Cancel(id uint64) {
	var m sync.Mutex
	m.Lock()
	cmd.tasks[id] = true
	m.Unlock()
}

func (cmd *CancelCmds) GetCancelStatus(id uint64) bool {
	return cmd.tasks[id]
}

func (cmd *CancelCmds) Remove(id uint64) {
	var m sync.Mutex
	m.Lock()
	delete(cmd.tasks, id)
	m.Unlock()
}

func NewCancelCmd() *CancelCmds {
	return &CancelCmds{
		tasks: make(map[uint64]bool),
	}
}

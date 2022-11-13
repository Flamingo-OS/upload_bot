package core

type CancelCmds struct {
	tasks map[uint64]bool
}

func (cmd *CancelCmds) Insert(id uint64) {
	Mut.Lock()
	cmd.tasks[id] = false
	Mut.Unlock()
}

func (cmd *CancelCmds) Cancel(id uint64) {
	Mut.Lock()
	cmd.tasks[id] = true
	Mut.Unlock()
}

func (cmd *CancelCmds) GetCancelStatus(id uint64) bool {
	return cmd.tasks[id]
}

func (cmd *CancelCmds) Remove(id uint64) {
	Mut.Lock()
	delete(cmd.tasks, id)
	Mut.Unlock()
}

func NewCancelCmd() *CancelCmds {
	return &CancelCmds{
		tasks: make(map[uint64]bool),
	}
}

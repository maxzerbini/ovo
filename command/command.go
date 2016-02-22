package command

import (
	"github.com/maxzerbini/ovo/storage"
)

type Command struct {
	OpCode string
	Obj    *storage.MetaDataUpdObj
}

func (cmd Command) RpcCommand() *RpcCommand {
	return &RpcCommand{OpCode: cmd.OpCode, Obj: cmd.Obj}
}

type RpcCommand struct {
	Source string
	OpCode string
	Obj    *storage.MetaDataUpdObj
}

func (rpccmd RpcCommand) Command() *Command {
	return &Command{OpCode: rpccmd.OpCode, Obj: rpccmd.Obj}
}

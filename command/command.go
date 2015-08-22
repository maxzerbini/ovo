package command

import(
	"github.com/maxzerbini/ovo/storage"
)

type Command struct {
	OpCode string
	Obj *storage.MetaDataUpdObj
}

type RpcCommand struct {
	OpCode string
	Obj storage.MetaDataUpdObj
}

func (rpccmd RpcCommand) Command() (*Command){
	return &Command{OpCode:rpccmd.OpCode, Obj:&(rpccmd.Obj)}
}

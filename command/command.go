package command

import(
	"github.com/maxzerbini/ovo/storage"
)

type Command struct {
	OpCode string
	Obj *storage.MetaDataUpdObj
}



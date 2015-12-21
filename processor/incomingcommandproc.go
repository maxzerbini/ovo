package processor

import (
	"github.com/maxzerbini/ovo/command"
	"github.com/maxzerbini/ovo/storage"
)

const commands_buffer_size = 10000

type InCommandQueue struct {
	commands   chan *command.Command
	keystorage storage.OvoStorage
}

func NewCommandQueue(ks storage.OvoStorage) *InCommandQueue {
	cq := new(InCommandQueue)
	cq.commands = make(chan *command.Command, commands_buffer_size)
	cq.keystorage = ks
	go cq.backend()
	return cq
}

func (cq *InCommandQueue) KeyStore() storage.OvoStorage {
	return cq.keystorage
}

func (cq *InCommandQueue) Enqueu(cmd *command.Command) {
	cq.commands <- cmd
}

func (cq *InCommandQueue) backend() {
	for cmd := range cq.commands {
		if cmd != nil {
			switch cmd.OpCode {
			case "put":
				cq.put(cmd.Obj)
			case "delete":
				cq.delete(cmd.Obj)
			case "touch":
				cq.touch(cmd.Obj)
			case "updatevalue":
				cq.updatevalue(cmd.Obj)
			case "updatekey":
				cq.updatekey(cmd.Obj)
			case "updatekeyvalue":
				cq.updatekeyvalue(cmd.Obj)
			case "setcounter":
				cq.setcounter(cmd.Obj)
			default:
				println("usupported command: " + cmd.OpCode)
			}
		}
	}
}

func (cq *InCommandQueue) put(obj *storage.MetaDataUpdObj) {
	cq.keystorage.Put(obj.MetaDataObj())
}

func (cq *InCommandQueue) delete(obj *storage.MetaDataUpdObj) {
	cq.keystorage.Delete(obj.Key)
}

func (cq *InCommandQueue) touch(obj *storage.MetaDataUpdObj) {
	cq.keystorage.Touch(obj.Key)
}

func (cq *InCommandQueue) updatevalue(obj *storage.MetaDataUpdObj) {
	cq.keystorage.UpdateValueIfEqual(obj)
}

func (cq *InCommandQueue) updatekey(obj *storage.MetaDataUpdObj) {
	cq.keystorage.UpdateKey(obj)
}

func (cq *InCommandQueue) updatekeyvalue(obj *storage.MetaDataUpdObj) {
	cq.keystorage.UpdateKeyAndValueIfEqual(obj)
}

func (cq *InCommandQueue) setcounter(obj *storage.MetaDataUpdObj) {
	cq.keystorage.SetCounter(obj.MetaDataCounter())
}

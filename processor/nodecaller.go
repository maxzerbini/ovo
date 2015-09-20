package processor
import(
	"github.com/maxzerbini/ovo/storage"
	"github.com/maxzerbini/ovo/cluster"
)

type NodeCaller struct {
	// TODO
}

// Ask remote server to add the obj 
func (nc *NodeCaller) Add(obj *storage.MetaDataObj, destination *cluster.OvoNode) {
	// TODO
}

// Ask remote server to delete the obj
func (nc *NodeCaller) Delete(obj *storage.MetaDataObj, destination *cluster.OvoNode) {
	// TODO
}

// Ask remote server to touch the obj
func (nc *NodeCaller) Touch(obj *storage.MetaDataObj, destination *cluster.OvoNode) {
	// TODO
}

// Ask remote server to change the obj
func (nc *NodeCaller) Change(obj *storage.MetaDataUpdObj, destination *cluster.OvoNode) {
	// TODO
}
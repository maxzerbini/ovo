package cluster

import(
	"github.com/maxzerbini/ovo/command"
)

type OvoNode struct {
	Name string
	HashRange []int
	Host string
	CommandPort int
	APIHost string
	APIPort int
}

type ClusterTopologyNode struct {
	Name string
	Twins []OvoNode
}

type ClusterTopology struct {
	Nodes []ClusterTopologyNode
}

type OvoNodeCommander interface {
	SendCommandToTwins(cmd *command.Command)
	ProcessIncomingCommand(cmd *command.Command)
}


type OvoCluster interface {
	GetTwins()([]OvoNode)
}

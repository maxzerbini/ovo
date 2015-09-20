package processor

import (
	"github.com/maxzerbini/ovo/command"
)

type OutCommandQueue struct {
	commands chan *command.Command
	caller *NodeCaller
}
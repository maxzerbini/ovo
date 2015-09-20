package server

import (
	"github.com/maxzerbini/ovo/cluster"
	"testing"
)

func TestConfigurationWrite(t *testing.T) {
	t.Log("TestConfigurationWrite started")
	var conf ServerConf
	var node = &cluster.OvoNode{Name:"testnode",Host:"localhost"}
	conf.Node = *node
	WriteConfiguration("../conf/severconf.json", conf)
}

func TestConfigurationLoad(t *testing.T) {
	t.Log("TestConfigurationLoad started")
	var conf = LoadConfiguration("../conf/severconf.json")
	t.Logf("conf = %s", conf)
}
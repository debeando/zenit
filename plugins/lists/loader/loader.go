package loader

import (
	"fmt"

	"github.com/swapbyt3s/zenit/common/log"
)

type Plugin interface {
	Collect()
}

type Creator func() Plugin

var Plugins = map[string]Creator{}

func Add(name string, creator Creator) {
	log.Info(fmt.Sprintf("Plugin - %s", name))

	Plugins[name] = creator
}

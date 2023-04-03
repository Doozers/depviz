package dvstore

import (
	"github.com/Doozers/depviz/internal/dvmodel"
	"github.com/cayleygraph/cayley/schema"
)

func Schema() *schema.Config {
	config := schema.NewConfig()
	// temporarily forced to register it globally :(
	schema.RegisterType("dv:Owner", dvmodel.Owner{})
	schema.RegisterType("dv:Task", dvmodel.Task{})
	schema.RegisterType("dv:Topic", dvmodel.Topic{})
	return config
}

package plugins

import (
	"fmt"

	"github.com/mitchellh/mapstructure"
	"github.com/nextrevision/runfile/plugins/shell"
	"github.com/nextrevision/runfile/plugins/template"
)

type Plugin interface {
	Validate() error
	Run() error
	Type() string
}

type genericPlugin struct{}

func (g genericPlugin) Validate() error { return nil }
func (g genericPlugin) Type() string    { return "generic" }
func (g genericPlugin) Run() error      { return nil }

type Type struct {
	Type string `yaml:"type"`
}

// NewPlugin ...
func NewPlugin(config interface{}) (Plugin, error) {
	t := Type{}
	if err := mapstructure.Decode(config, &t); err != nil {
		return genericPlugin{}, err
	}

	switch t.Type {
	case "shell":
		return shell.New(config)
	case "template":
		return template.New(config)
	default:
		return genericPlugin{}, fmt.Errorf("Invalid type: %s", t.Type)
	}
}

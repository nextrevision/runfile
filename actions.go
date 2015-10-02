package rf

import (
	"github.com/nextrevision/runfile/action"
)

type RfFile struct {
	Check action.CheckAction `yaml:"check"`
	Build action.BuildAction `yaml:"build"`
}

type Actions []Action

type Action struct {
	Name     string
	Commands Commands
}

type Commands interface {
	Validate() error
	Run() error
	Nil() bool
}

func GetActions(r RfFile) Actions {
	return Actions{
		{Name: "check", Commands: &r.Check},
		{Name: "build", Commands: &r.Build},
	}
}

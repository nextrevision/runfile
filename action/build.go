package action

import (
	"fmt"
	"os/exec"
)

type BuildAction struct {
	Command string `yaml:"command"`
}

func (b BuildAction) Nil() bool {
	if b == (BuildAction{}) {
		return true
	}
	return false
}

func (b *BuildAction) Validate() error {
	println("[Build] Validating:", b.Command)
	return nil
}

func (b *BuildAction) Run() error {
	println("[Build] Running:", b.Command)
	out, err := exec.Command(b.Command).CombinedOutput()
	fmt.Println(string(out[:]))
	return err
}

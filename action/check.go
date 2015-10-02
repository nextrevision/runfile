package action

import (
	"fmt"
	"os/exec"
)

type CheckAction struct {
	Command string `yaml:"command"`
}

func (c CheckAction) Nil() bool {
	if c == (CheckAction{}) {
		return true
	}
	return false
}

func (c *CheckAction) Validate() error {
	println("[Check] Validating:", c.Command)
	return nil
}

func (c *CheckAction) Run() error {
	println("[Check] Running:", c.Command)
	out, err := exec.Command(c.Command).CombinedOutput()
	fmt.Println(string(out[:]))
	return err
}

package shell

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"

	"github.com/mitchellh/mapstructure"
)

// Config ...
type Config struct {
	Script  string `yaml:"script"`
	Content string `yaml:"content"`
	Quiet   bool   `yaml:"quiet"`
}

func New(config interface{}) (Config, error) {
	c := Config{}
	if err := mapstructure.Decode(config, &c); err != nil {
		return c, err
	}
	return c, nil
}

func (c Config) Validate() error {
	if c.Script == "" && c.Content == "" {
		return fmt.Errorf("Content cannot be empty")
	}

	if c.Script != "" && c.Content != "" {
		return fmt.Errorf("Both script and content cannot be set")
	}
	return nil
}

func (c Config) Type() string {
	return "shell"
}

// Run ...
func (c Config) Run() error {
	var script string

	if c.Content != "" {
		tmpfile, err := ioutil.TempFile("", "runshell")
		if err != nil {
			return err
		}

		defer os.Remove(tmpfile.Name()) // clean up

		if err = tmpfile.Chmod(0755); err != nil {
			return err
		}

		if bytes.Index([]byte(c.Content), []byte("#!")) != 0 {
			if _, err = tmpfile.Write([]byte("#!/bin/sh\n")); err != nil {
				return err
			}
		}

		if _, err = tmpfile.Write([]byte(c.Content)); err != nil {
			return err
		}
		if err = tmpfile.Close(); err != nil {
			return err
		}
		script = tmpfile.Name()
	} else {
		script = c.Script
	}

	out, err := exec.Command(script).CombinedOutput()
	output := string(out[:])
	if !c.Quiet && output != "" {
		fmt.Printf("%s", output)
	}
	return err
}

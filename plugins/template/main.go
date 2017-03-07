package template

import (
	"bytes"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/mitchellh/mapstructure"
)

// Config ...
type Config struct {
	Source      string                 `yaml:"source"`
	Content     string                 `yaml:"content"`
	Destination string                 `yaml:"destination"`
	Vars        map[string]interface{} `yaml:"vars"`
}

func New(config interface{}) (Config, error) {
	c := Config{}
	if err := mapstructure.Decode(config, &c); err != nil {
		return c, err
	}
	return c, nil
}

func (c Config) Validate() error {
	if c.Source == "" && c.Content == "" {
		return fmt.Errorf("Source or content must be set")
	}
	return nil
}

func (c Config) Type() string {
	return "template"
}

// Run ...
func (c Config) Run() error {
	content := c.Content

	if c.Source != "" {
		if strings.Split(c.Source, "s:")[0] == "http" {
			resp, err := http.Get(c.Source)
			if err != nil {
				return err
			}
			defer resp.Body.Close()
			body, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				return err
			}
			content = string(body)
		} else {
			body, err := ioutil.ReadFile(c.Source)
			if err != nil {
				return err
			}
			content = string(body)
		}
	}

	if content == "" {
		return fmt.Errorf("Template cannot be empty")
	}

	t := template.Must(template.New("runfile").Parse(content))

	if c.Destination == "" {
		tmpfile, err := ioutil.TempFile("", "runfiletemplate")
		if err != nil {
			return err
		}
		c.Destination = tmpfile.Name()
	}

	b := bytes.NewBuffer(nil)
	if err := t.Execute(b, c.Vars); err != nil {
		return err
	}

	if err := ioutil.WriteFile(c.Destination, b.Bytes(), 0644); err != nil {
		return err
	}

	log.Printf("Template written to %s", c.Destination)
	return nil
}

func fileExists(filepath string) bool {
	_, err := os.Stat(filepath)
	return !os.IsNotExist(err)
}

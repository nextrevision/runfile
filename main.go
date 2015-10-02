package rf

import (
	"io/ioutil"
	"log"
	"os"

	"gopkg.in/yaml.v2"
)

func Run(args []string) {
	if err := checkRf(); err != nil {
		log.Fatalln("rf.yml not found in current directory")
	}

	actions := loadRf()

	if len(args) == 0 {
		println("Available actions:")
		for _, action := range *actions {
			if action.Commands.Nil() {
				continue
			}
			println(action.Name)
		}
	} else {
		for _, arg := range args {
			for _, action := range *actions {
				if action.Name == arg {
					if action.Commands.Nil() {
						println(action.Name, "is not declared in rf.yml")
						continue
					}
					if err := action.Commands.Validate(); err != nil {
						log.Fatalln(err)
					}
					if err := action.Commands.Run(); err != nil {
						log.Fatalln(err)
					}
				}
			}
		}
	}
}

func checkRf() error {
	if _, err := os.Stat("rf.yml"); err != nil {
		return err
	}
	return nil
}

func loadRf() *Actions {
	contents, err := ioutil.ReadFile("rf.yml")
	if err != nil {
		log.Fatalln("rf.yml not found in current directory")
	}

	var rfFile RfFile
	if err = yaml.Unmarshal(contents, &rfFile); err != nil {
		log.Fatalln(err)
	}

	actions := GetActions(rfFile)
	return &actions
}

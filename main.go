package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"time"

	"github.com/ghodss/yaml"
	"github.com/nextrevision/runfile/plugins"
)

var (
	cfgFile      string
	printTasks   bool
	printVersion bool
)

func init() {
	flag.StringVar(&cfgFile, "c", "run.yml", "Run config file")
	flag.BoolVar(&printTasks, "l", false, "Print tasks and exit")
	flag.BoolVar(&printVersion, "v", false, "Print version and exit")
}

func main() {
	flag.Parse()
	log.SetFlags(0)
	log.SetPrefix("[runfile] ")

	if printVersion {
		fmt.Println(Version)
		return
	}

	runConfig, err := loadRunFile(cfgFile)
	if err != nil {
		log.Fatal(err)
	}

	pluginMap, err := pluginsFromConfig(runConfig)
	if err != nil {
		log.Fatal(err)
	}

	if printTasks {
		for task, plugin := range pluginMap {
			fmt.Printf("%s (%s)\n", task, plugin.Type())
		}
		return
	}

	log.Println("Validating tasks...")
	for task, plugin := range pluginMap {
		if err := plugin.Validate(); err != nil {
			log.Fatalf("Invalid config in '%s' task: %s", task, err)
		}
	}

	for _, task := range flag.Args() {
		plugin, ok := pluginMap[task]
		if !ok {
			log.Fatalf("No such task: %s", task)
		}

		log.SetPrefix(fmt.Sprintf("[%s] ", task))
		start := time.Now()
		log.Printf("Executing ...")
		if err := plugin.Run(); err != nil {
			log.Fatalf("ERROR: %s", err)
		}
		log.Printf("Completed in %s", time.Since(start))
		log.SetPrefix("")
	}
}

func loadRunFile(filename string) (map[string]interface{}, error) {
	config := make(map[string]interface{})

	contents, err := ioutil.ReadFile(filename)
	if err != nil {
		return config, err
	}

	if err = yaml.Unmarshal(contents, &config); err != nil {
		return config, err
	}
	return config, nil
}

func pluginsFromConfig(config map[string]interface{}) (map[string]plugins.Plugin, error) {
	pluginMap := make(map[string]plugins.Plugin)
	for name, config := range config {
		plugin, err := plugins.NewPlugin(config)
		if err != nil {
			return pluginMap, err
		}
		pluginMap[name] = plugin
	}
	return pluginMap, nil
}

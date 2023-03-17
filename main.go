package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	flag "github.com/spf13/pflag"

	"github.com/keilerkonzept/dockerfile-json/pkg/dockerfile"
)

var config struct {
	Quiet       bool
	Expand      bool
	BuildArgs   AssignmentsMap
	NonzeroExit bool
}

var name = "dockerfile-json"
var version = "dev"
var jsonOut = json.NewEncoder(os.Stdout)

func init() {
	log.SetOutput(os.Stderr)
	log.SetFlags(0)
	log.SetPrefix(fmt.Sprintf("[%s %s] ", filepath.Base(name), version))

	config.Expand = true
	flag.BoolVar(&config.Quiet, "quiet", config.Quiet, "suppress log output (stderr)")
	flag.BoolVar(&config.Expand, "expand-build-args", config.Expand, "expand build args")
	flag.Var(&config.BuildArgs, "build-arg", config.BuildArgs.Help())
	flag.Parse()

	if config.Quiet {
		log.SetOutput(ioutil.Discard)
	}

	if flag.NArg() == 0 {
		flag.Usage()
	}
}

func buildArgEnvExpander() dockerfile.SingleWordExpander {
	env := make(map[string]string, len(config.BuildArgs.Value))
	for key, value := range config.BuildArgs.Value {
		if value != nil {
			env[key] = *value
			continue
		}
		if value, ok := os.LookupEnv(key); ok {
			env[key] = value
		}
	}
	return func(word string) (string, error) {
		if value, ok := env[word]; ok {
			return value, nil
		}
		return "", fmt.Errorf("not defined: $%s", word)
	}
}

func main() {
	var dockerfiles []*dockerfile.Dockerfile
	for _, path := range flag.Args() {
		dockerfile, err := dockerfile.Parse(path)
		if err != nil {
			log.Printf("error: parse %q: %v", path, err)
			config.NonzeroExit = true
			continue
		}
		dockerfiles = append(dockerfiles, dockerfile)
	}
	if config.Expand {
		env := buildArgEnvExpander()
		for _, dockerfile := range dockerfiles {
			dockerfile.Expand(env)
		}
	}

	for _, dockerfile := range dockerfiles {
		jsonOut.Encode(dockerfile)
	}

	if config.NonzeroExit {
		os.Exit(1)
	}
}

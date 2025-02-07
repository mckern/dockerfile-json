package dockerfile

import (
	"fmt"
	"os"

	"github.com/moby/buildkit/frontend/dockerfile/instructions"
)

type SingleWordExpander instructions.SingleWordExpander

func (d *Dockerfile) Expand(env SingleWordExpander) {
	d.expand(instructions.SingleWordExpander(env))
	d.analyzeStages()
}

func (d *Dockerfile) expand(env instructions.SingleWordExpander) {
	metaArgsEnvExpander := d.metaArgsEnvExpander(env)
	for i, stage := range d.Stages {
		d.Stages[i].BaseName = os.Expand(stage.BaseName, func(key string) string {
			value, err := metaArgsEnvExpander(key)
			if err != nil {
				return ""
			}
			return value
		})

		for i := range stage.Commands {
			cmdExpander, ok := stage.Commands[i].Command.(instructions.SupportsSingleWordExpansion)
			if ok {
				cmdExpander.Expand(metaArgsEnvExpander)
			}
		}
	}
}

func (d *Dockerfile) metaArgsEnvExpander(env instructions.SingleWordExpander) instructions.SingleWordExpander {
	metaArgsEnv := make(map[string]string, len(d.MetaArgs))
	for _, arg := range d.MetaArgs {
		if arg.DefaultValue != nil {
			metaArgsEnv[arg.Key] = *arg.DefaultValue
		}

		if value, err := env(arg.Key); err == nil {
			arg.ProvidedValue = &value
			metaArgsEnv[arg.Key] = value
			arg.Value = &value
		}

		err := arg.Expand(env)
		if err != nil {
			continue
		}

		// process any key-value pairs in a given arg
		for _, kv := range arg.ArgCommand.Args {
			if kv.Value != nil {
				metaArgsEnv[arg.Key] = kv.ValueString()
			}
		}
	}
	return func(key string) (string, error) {
		if value, ok := metaArgsEnv[key]; ok {
			return value, nil
		}
		return "", fmt.Errorf("not defined: $%s", key)
	}
}

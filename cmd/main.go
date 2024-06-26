package main

import (
	"errors"
	"flag"
	"fmt"
	"octopus/internal/requests"
	"octopus/internal/templates"
	"strings"
)

type VarArg map[string]string

func (v *VarArg) String() string {
	return fmt.Sprintf("%s", *v)
}

func (v *VarArg) Set(value string) error {
	if name, value, found := strings.Cut(value, ":"); found {
		(*v)[name] = value
	} else {
		return errors.New("wrong variable format")
	}
	return nil
}

type Args struct {
	path        *string
	variables   VarArg
	parallelism *int
}

func parseArgs() *Args {
	args := Args{
		variables: make(VarArg),
	}
	flag.Var(&args.variables, "v", "Define a variable. Example -v=\"key:value\"")
	args.path = flag.String("f", "", "Templates file path")
	args.parallelism = flag.Int("p", 1, "Parallelism level")
	flag.Parse()
	return &args
}

func main() {
	args := parseArgs()
	repo := templates.NewTemplatesRepositoryBuilder().
		LoadFile(*args.path).
		Inject(args.variables).
		Build()
	manager := requests.NewHttpRequestsManagerBuilder().
		Templates(repo).
		Parallelism(*args.parallelism).
		Build()
	manager.Execute()
}

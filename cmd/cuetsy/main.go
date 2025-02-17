package main

import (
	"fmt"
	"os"

	"cuelang.org/go/cue"
	"cuelang.org/go/cue/errors"
	"cuelang.org/go/cue/load"
	"github.com/grafana/cuetsy"
)

func main() {
	wd, err := os.Getwd()
	if err != nil {
		os.Exit(1)
	}

	// For now at least, be higly restrictive and only allow passing one file.
	if len(os.Args) != 2 {
		fmt.Fprint(os.Stderr, "must provide path to exactly one .cue file\n")
		os.Exit(1)
	}

	fi, err := os.Stat(os.Args[1])
	if err != nil {
		fmt.Fprint(os.Stderr, err)
		os.Exit(1)
	}
	if fi.IsDir() || !fi.Mode().IsRegular() {
		fmt.Fprint(os.Stderr, "must provide path to exactly one .cue file\n")
		os.Exit(1)
	}

	loadedInstances := load.Instances([]string{os.Args[1]}, nil)
	instances := cue.Build(loadedInstances)
	// Given the above input constraints, there _should_ only ever be a
	// single element in this slice.
	for _, inst := range instances {
		b, err := cuetsy.Generate(inst.Value(), cuetsy.Config{})
		if err != nil {
			errors.Print(os.Stderr, err, &errors.Config{
				Cwd: wd,
			})
			os.Exit(1)
		}
		// For now, write results to a file adjacent to the input cue file.
		fd, err := os.Create(os.Args[1][:len(os.Args[1])-3] + "ts")
		if err != nil {
			fmt.Fprint(os.Stderr, err)
		}
		fmt.Fprint(fd, string(b))
	}
}

package main

import (
	"fmt"
	"os"

	"github.com/fission/fission-workflows/pkg/parse/yaml"

	"github.com/fission/fission-workflows/pkg/types/validate"
	"github.com/urfave/cli"
)

// TODO also validate with backend (optional)
var cmdValidate = cli.Command{
	Name:        "validate",
	Usage:       "Validate <file>",
	Description: "Validate a Workflow",
	Action: commandContext(func(ctx Context) error {
		// Get path from args
		if ctx.NArg() == 0 {
			fail("No file provided.")
		}

		var failed bool
		for _, path := range ctx.Args() {

			printErr := func(msg string) {
				fmt.Fprintf(os.Stderr, "%s: %s\n", path, msg)
				failed = true
			}

			// Get file
			file, err := os.Open(path)
			if err != nil {
				printErr(fmt.Sprintf("Failed to read file: %v", err))
				continue
			}

			// Read file into workflowSpec (assume yaml for now)
			spec, err := yaml.Parse(file)
			if err != nil {
				printErr(fmt.Sprintf("Failed to parse yaml definition: %v", err))
				continue
			}

			// Validate workflowSpec
			err = validate.WorkflowSpec(spec)
			if err != nil {
				invalid, ok := err.(validate.Error)
				if ok {
					printErr(validate.Format(invalid))
				} else {
					printErr(fmt.Sprintf("Unknown error: %v", err))
				}
				continue
			}
		}

		if failed {
			os.Exit(1)
		}

		return nil
	}),
}

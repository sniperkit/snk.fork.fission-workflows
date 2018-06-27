package main

import (
	"fmt"
	"os"
	"sort"

	"github.com/fission/fission-workflows/pkg/parse/yaml"
	"github.com/urfave/cli"
)

var cmdWorkflow = cli.Command{
	Name:    "workflow",
	Aliases: []string{"wf", "workflows"},
	Usage:   "Workflow-related commands",
	Subcommands: []cli.Command{
		{
			Name:  "get",
			Usage: "get <Workflow-id> <task-id>",
			Action: commandContext(func(ctx Context) error {
				client := getClient(ctx)

				switch ctx.NArg() {
				case 0:
					// List workflows
					resp, err := client.Workflow.List(ctx)
					if err != nil {
						panic(err)
					}
					wfs := resp.Workflows
					sort.Strings(wfs)
					var rows [][]string
					for _, wfID := range wfs {
						wf, err := client.Workflow.Get(ctx, wfID)
						if err != nil {
							panic(err)
						}
						updated := wf.Status.UpdatedAt.String()
						created := wf.Metadata.CreatedAt.String()

						rows = append(rows, []string{wfID, wf.Spec.Name, string(wf.Status.Status),
							created, updated})
					}
					table(os.Stdout, []string{"id", "NAME", "STATUS", "CREATED", "UPDATED"}, rows)
				case 1:
					// Get Workflow
					wfID := ctx.Args().Get(0)
					wf, err := client.Workflow.Get(ctx, wfID)
					if err != nil {
						panic(err)
					}
					b, err := yaml.Marshal(wf)
					if err != nil {
						panic(err)
					}
					fmt.Printf("%v\n", string(b))
				case 2:
					// Get Workflow task
					fallthrough
				default:
					wfID := ctx.Args().Get(0)
					taskID := ctx.Args().Get(1)
					wf, err := client.Workflow.Get(ctx, wfID)
					if err != nil {
						panic(err)
					}
					task, ok := wf.Spec.Tasks[taskID]
					if !ok {
						fmt.Println("Task not found.")
						return nil
					}
					b, err := yaml.Marshal(task)
					if err != nil {
						panic(err)
					}
					fmt.Printf("%v\n", string(b))
				}

				return nil
			}),
		},
	},
}

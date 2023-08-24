// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package command

import (
	"fmt"
	"os"
	"strings"

	"github.com/hashicorp/consul/snapshot"
	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/raft"
	"github.com/mitchellh/cli"
	"github.com/posener/complete"
)

var (
	_ cli.Command             = (*OperatorRaftSnapshotInspectCommand)(nil)
	_ cli.CommandAutocomplete = (*OperatorRaftSnapshotInspectCommand)(nil)
)

type OperatorRaftSnapshotInspectCommand struct {
	*BaseCommand
}

func (c *OperatorRaftSnapshotInspectCommand) Synopsis() string {
	return "Inspects raft snapshot"
}

func (c *OperatorRaftSnapshotInspectCommand) Help() string {
	helpText := `
Usage: vault operator raft snapshot inspect <snapshot_file>

  Inspects a snapshot file.

	  $ vault operator raft snapshot inspect raft.snap

` + c.Flags().Help()

	return strings.TrimSpace(helpText)
}

func (c *OperatorRaftSnapshotInspectCommand) Flags() *FlagSets {
	set := c.flagSet(FlagSetHTTP | FlagSetOutputFormat)

	return set
}

func (c *OperatorRaftSnapshotInspectCommand) AutocompleteArgs() complete.Predictor {
	return complete.PredictAnything
}

func (c *OperatorRaftSnapshotInspectCommand) AutocompleteFlags() complete.Flags {
	return c.Flags().Completions()
}

func (c *OperatorRaftSnapshotInspectCommand) Run(args []string) int {
	fmt.Println("I AM INSPECTING A RAFT SNAPSHOT")

	// TODO: how to add other flags???
	flags := c.Flags()

	if err := flags.Parse(args); err != nil {
		c.UI.Error(err.Error())
		return 1
	}

	fmt.Printf("flag set after parse %+v\n", *flags)

	var file string
	args = c.flags.Args()
	fmt.Printf("args %+v\n", args)

	switch len(args) {
	case 0:
		c.UI.Error("Missing FILE argument")
		return 1
	case 1:
		file = args[0]
	default:
		c.UI.Error(fmt.Sprintf("Too many arguments (expected 1, got %d)", len(args)))
		return 1
	}

	fmt.Println("file", file)

	// Open the file.
	f, err := os.Open(file)
	if err != nil {
		c.UI.Error(fmt.Sprintf("Error opening snapshot file: %s", err))
		return 1
	}
	defer f.Close()

	// TODO: skipping state.bin logic for now
	var readFile *os.File
	var meta *raft.SnapshotMeta

	// TODO: should we use consul's snapshot reader or raft-snapshot repo?
	// Decided to use this because it was more complete that read from raft-snapshot
	readFile, meta, err = snapshot.Read(hclog.New(nil), f)
	if err != nil {
		c.UI.Error(fmt.Sprintf("Error reading snapshot: %s", err))
		return 1
	}
	fmt.Println("readFile", readFile.Name())
	defer func() {
		if err := readFile.Close(); err != nil {
			c.UI.Error(fmt.Sprintf("Failed to close temp snapshot: %v", err))
		}
		if err := os.Remove(readFile.Name()); err != nil {
			c.UI.Error(fmt.Sprintf("Failed to clean up temp snapshot: %v", err))
		}
	}()

	fmt.Printf("META %+v\n", *meta)

	// OLD for snapshot save
	// f := c.Flags()

	// if err := f.Parse(args); err != nil {
	// 	c.UI.Error(err.Error())
	// 	return 1
	// }

	// path := ""

	// args = f.Args()
	// switch len(args) {
	// case 1:
	// 	path = strings.TrimSpace(args[0])
	// default:
	// 	c.UI.Error(fmt.Sprintf("Incorrect arguments (expected 1, got %d)", len(args)))
	// 	return 1
	// }

	// if len(path) == 0 {
	// 	c.UI.Error("Output file name is required")
	// 	return 1
	// }

	// w := &lazyOpenWriter{
	// 	openFunc: func() (io.WriteCloser, error) {
	// 		return os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0o600)
	// 	},
	// }

	// client, err := c.Client()
	// if err != nil {
	// 	w.Close()
	// 	c.UI.Error(err.Error())
	// 	return 2
	// }

	// err = client.Sys().RaftSnapshot(w)
	// if err != nil {
	// 	w.Close()
	// 	c.UI.Error(fmt.Sprintf("Error taking the snapshot: %s", err))
	// 	return 2
	// }

	// err = w.Close()
	// if err != nil {
	// 	c.UI.Error(fmt.Sprintf("Error taking the snapshot: %s", err))
	// 	return 2
	// }
	return 0
}

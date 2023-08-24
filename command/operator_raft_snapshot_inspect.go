// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package command

import (
	"fmt"
	"io"
	"math"
	"os"
	"strings"

	protoio "github.com/gogo/protobuf/io"
	"github.com/hashicorp/go-hclog"
	iradix "github.com/hashicorp/go-immutable-radix"
	"github.com/hashicorp/raft"
	snapshot "github.com/hashicorp/raft-snapshot"
	"github.com/hashicorp/vault/sdk/plugin/pb"
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
	// var readFile *os.File
	// var meta *raft.SnapshotMeta

	// testing...
	parseSnapshotNew(hclog.New(nil), f)

	// TODO: should we use consul's snapshot reader or raft-snapshot repo?
	// Decided to use this because it was more complete that read from raft-snapshot
	// readFile, meta, err = snapshot.Read(hclog.New(nil), f)
	// if err != nil {
	// 	c.UI.Error(fmt.Sprintf("Error reading snapshot: %s", err))
	// 	return 1
	// }
	// fmt.Println("readFile", readFile.Name())
	// defer func() {
	// 	if err := readFile.Close(); err != nil {
	// 		c.UI.Error(fmt.Sprintf("Failed to close temp snapshot: %v", err))
	// 	}
	// 	if err := os.Remove(readFile.Name()); err != nil {
	// 		c.UI.Error(fmt.Sprintf("Failed to clean up temp snapshot: %v", err))
	// 	}
	// }()

	// fmt.Printf("META %+v\n", *meta)

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

func parseSnapshotNew(logger hclog.Logger, in io.Reader) (*raft.SnapshotMeta, *iradix.Tree, error) {
	// file, err := os.Open(filename)
	// if err != nil {
	// 	return nil, nil, err
	// }
	// defer file.Close()

	reader, writer := io.Pipe()

	protoReader := protoio.NewDelimitedReader(reader, math.MaxInt32)
	defer protoReader.Close()

	errCh := make(chan error, 2)

	var meta *raft.SnapshotMeta
	go func() {
		var err error
		meta, err = snapshot.Parse(file, writer)
		writer.Close()
		errCh <- err
	}()

	// txn := iradix.New().Txn()

	// go func() {
	// 	for {
	// 		s := new(pb.StorageEntry)
	// 		err := protoReader.ReadMsg(s)
	// 		if err != nil {
	// 			if err == io.EOF {
	// 				errCh <- nil
	// 				return
	// 			}
	// 			errCh <- err
	// 			return
	// 		}

	// 		var value interface{} = struct{}{}
	// 		if loadValues {
	// 			value = s.Value
	// 		}

	// 		txn.Insert([]byte(s.Key), value)
	// 	}
	// }()

	// err = <-errCh
	// if err != nil && err != io.EOF {
	// 	return nil, nil, err
	// }

	// err = <-errCh
	// if err != nil && err != io.EOF {
	// 	return nil, nil, err
	// }

	// data := txn.Commit()

	// return meta, data, nil
}

func parseSnapshot(filename string, loadValues bool) (*raft.SnapshotMeta, *iradix.Tree, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, nil, err
	}
	defer file.Close()

	reader, writer := io.Pipe()

	protoReader := protoio.NewDelimitedReader(reader, math.MaxInt32)
	defer protoReader.Close()

	errCh := make(chan error, 2)

	var meta *raft.SnapshotMeta
	go func() {
		var err error
		meta, err = snapshot.Parse(file, writer)
		writer.Close()
		errCh <- err
	}()

	txn := iradix.New().Txn()

	go func() {
		for {
			s := new(pb.StorageEntry)
			err := protoReader.ReadMsg(s)
			if err != nil {
				if err == io.EOF {
					errCh <- nil
					return
				}
				errCh <- err
				return
			}

			var value interface{} = struct{}{}
			if loadValues {
				value = s.Value
			}

			txn.Insert([]byte(s.Key), value)
		}
	}()

	err = <-errCh
	if err != nil && err != io.EOF {
		return nil, nil, err
	}

	err = <-errCh
	if err != nil && err != io.EOF {
		return nil, nil, err
	}

	data := txn.Commit()

	return meta, data, nil
}

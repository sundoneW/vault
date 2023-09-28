// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package command

import (
	"fmt"
	"testing"

	"github.com/mitchellh/cli"
)

func testOperatorRaftSnapshotInspectCommand(tb testing.TB) (*cli.MockUi, *OperatorRaftSnapshotInspectCommand) {
	tb.Helper()

	ui := cli.NewMockUi()
	return ui, &OperatorRaftSnapshotInspectCommand{
		BaseCommand: &BaseCommand{
			UI: ui,
		},
	}
}

func TestOperatorRaftSnapshotInspectCommand_Run(t *testing.T) {
	t.Parallel()

	// cases := []struct {
	// 	name string
	// 	args []string
	// 	out  string
	// 	code int
	// }{
	// 	{
	// 		"too_many_args",
	// 		[]string{"foo"},
	// 		"Too many arguments",
	// 		1,
	// 	},
	// 	{
	// 		"default",
	// 		nil,
	// 		"Success! Stepped down: ",
	// 		0,
	// 	},
	// }

	t.Run("validations", func(t *testing.T) {
		t.Parallel()

		// tc := tc

		t.Run("passing test", func(t *testing.T) {
			t.Parallel()

			client, closer := testVaultServer(t)
			defer closer()

			ui, cmd := testOperatorRaftSnapshotInspectCommand(t)
			cmd.client = client

			code := cmd.Run([]string{"complete.snap"})
			if code != 0 {
				t.Errorf("expected %d to be %d", code, 0)
			}

			combined := ui.OutputWriter.String() + ui.ErrorWriter.String()
			fmt.Println("combined", combined)
			// if !strings.Contains(combined, tc.out) {
			// 	t.Errorf("expected %q to contain %q", combined, tc.out)
			// }
		})
	})
}

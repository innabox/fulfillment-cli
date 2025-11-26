/*
Copyright (c) 2025 Red Hat Inc.

Licensed under the Apache License, Version 2.0 (the "License"); you may not use this file except in compliance with the
License. You may obtain a copy of the License at

  http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software distributed under the License is distributed on an
"AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the License for the specific
language governing permissions and limitations under the License.
*/

package create_test

import (
	"testing"

	"github.com/innabox/fulfillment-cli/internal/cmd/create"
	"github.com/innabox/fulfillment-cli/internal/cmd/create/cluster"
	"github.com/innabox/fulfillment-cli/internal/cmd/create/hostpool"
	"github.com/innabox/fulfillment-cli/internal/cmd/create/hub"
	"github.com/innabox/fulfillment-cli/internal/cmd/create/virtualmachine"
	"github.com/spf13/cobra"
)

func TestCreateSubcommandAliases(t *testing.T) {
	tests := []struct {
		name           string
		cmdFunc        func() *cobra.Command
		expectedAlias  string
		subcommandName string
	}{
		{
			name:           "cluster subcommand has fully-qualified alias",
			cmdFunc:        cluster.Cmd,
			expectedAlias:  "fulfillment.v1.Cluster",
			subcommandName: "cluster",
		},
		{
			name:           "hostpool subcommand has fully-qualified alias",
			cmdFunc:        hostpool.Cmd,
			expectedAlias:  "fulfillment.v1.HostPool",
			subcommandName: "hostpool",
		},
		{
			name:           "hub subcommand has fully-qualified alias",
			cmdFunc:        hub.Cmd,
			expectedAlias:  "private.v1.Hub",
			subcommandName: "hub",
		},
		{
			name:           "virtualmachine subcommand has fully-qualified alias",
			cmdFunc:        virtualmachine.Cmd,
			expectedAlias:  "fulfillment.v1.VirtualMachine",
			subcommandName: "virtualmachine",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := tt.cmdFunc()
			aliases := cmd.Aliases
			found := false
			for _, alias := range aliases {
				if alias == tt.expectedAlias {
					found = true
					break
				}
			}
			if !found {
				t.Errorf("Expected alias %q not found in aliases %v for %s subcommand",
					tt.expectedAlias, aliases, tt.subcommandName)
			}
		})
	}
}

func TestCreateCommandHasSubcommands(t *testing.T) {
	cmd := create.Cmd()
	subcommands := cmd.Commands()

	expectedSubcommands := []string{"cluster", "hostpool", "hub", "virtualmachine"}
	for _, expected := range expectedSubcommands {
		found := false
		for _, subcmd := range subcommands {
			if subcmd.Name() == expected {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("Expected subcommand %q not found in create command", expected)
		}
	}
}

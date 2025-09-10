/*
Copyright (c) 2025 Red Hat Inc.

Licensed under the Apache License, Version 2.0 (the "License"); you may not use this file except in compliance with the
License. You may obtain a copy of the License at

  http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software distributed under the License is distributed on an
"AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the License for the specific
language governing permissions and limitations under the License.
*/

package hub

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/innabox/fulfillment-cli/internal/config"
	"github.com/innabox/fulfillment-cli/internal/terminal"
	privatev1 "github.com/innabox/fulfillment-common/api/private/v1"
)

func Cmd() *cobra.Command {
	runner := &runnerContext{}
	result := &cobra.Command{
		Use:   "hub",
		Short: "Create a hub",
		RunE:  runner.run,
	}
	flags := result.Flags()
	flags.StringVar(
		&runner.id,
		"id",
		"",
		"Unique identifier of the hub.",
	)
	flags.StringVar(
		&runner.kubeconfig,
		"kubeconfig",
		"",
		"Kubeconfig file containing the details to connect to the Kubernetes API.",
	)
	flags.StringVar(
		&runner.namespace,
		"namespace",
		"",
		"Namespace where cluster orders will be created.",
	)
	return result
}

type runnerContext struct {
	console    *terminal.Console
	id         string
	kubeconfig string
	namespace  string
}

func (c *runnerContext) run(cmd *cobra.Command, args []string) error {
	// Get the context:
	ctx := cmd.Context()

	// Get the console:
	c.console = terminal.ConsoleFromContext(ctx)

	// Get the configuration:
	cfg, err := config.Load()
	if err != nil {
		return err
	}
	if cfg.Address == "" {
		return fmt.Errorf("there is no configuration, run the 'login' command")
	}

	// Check the parameters:
	if c.id == "" {
		return fmt.Errorf("identifier is required")
	}
	if c.namespace == "" {
		return fmt.Errorf("namespace name is required")
	}
	if c.kubeconfig == "" {
		return fmt.Errorf("kubeconfig file is required")
	}
	if c.namespace == "" {
		return fmt.Errorf("namespace name is required")
	}

	// Create the gRPC connection from the configuration:
	conn, err := cfg.Connect(ctx, cmd.Flags())
	if err != nil {
		return fmt.Errorf("failed to create gRPC connection: %w", err)
	}

	// Create the client:
	client := privatev1.NewHubsClient(conn)

	// Read the kubeconfig file:
	kubeconfig, err := os.ReadFile(c.kubeconfig)
	if err != nil {
		return fmt.Errorf("failed to read kubeconfig file '%s': %w", c.kubeconfig, err)
	}

	// Prepare the cluster:
	hub := privatev1.Hub_builder{
		Id:         c.id,
		Kubeconfig: kubeconfig,
		Namespace:  c.namespace,
	}.Build()

	// Create the hub:
	response, err := client.Create(ctx, privatev1.HubsCreateRequest_builder{
		Object: hub,
	}.Build())
	if err != nil {
		return fmt.Errorf("failed to create hub: %w", err)
	}

	// Display the result:
	hub = response.Object
	c.console.Printf(ctx, "Created hub `%s`.\n", hub.GetId())

	return nil
}

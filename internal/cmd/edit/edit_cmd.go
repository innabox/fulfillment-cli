/*
Copyright (c) 2025 Red Hat Inc.

Licensed under the Apache License, Version 2.0 (the "License"); you may not use this file except in compliance with the
License. You may obtain a copy of the License at

  http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software distributed under the License is distributed on an
"AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the License for the specific
language governing permissions and limitations under the License.
*/

package edit

import (
	"bytes"
	"context"
	"embed"
	"encoding/json"
	"fmt"
	"log/slog"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/innabox/fulfillment-common/logging"
	"github.com/innabox/fulfillment-common/templating"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
	"gopkg.in/yaml.v3"

	"github.com/innabox/fulfillment-cli/internal/config"
	"github.com/innabox/fulfillment-cli/internal/reflection"
	"github.com/innabox/fulfillment-cli/internal/terminal"
	ffv1 "github.com/innabox/fulfillment-common/api/fulfillment/v1"
)

//go:embed templates
var templatesFS embed.FS

// Possible output formats:
const (
	outputFormatJson = "json"
	outputFormatYaml = "yaml"
)

func Cmd() *cobra.Command {
	runner := &runnerContext{
		marshalOptions: protojson.MarshalOptions{
			UseProtoNames: true,
		},
	}
	result := &cobra.Command{
		Use:   "edit OBJECT ID",
		Short: "Edit objects",
		RunE:  runner.run,
	}
	flags := result.Flags()
	flags.StringVarP(
		&runner.format,
		"output",
		"o",
		outputFormatYaml,
		fmt.Sprintf(
			"Output format, one of '%s' or '%s'.",
			outputFormatJson, outputFormatYaml,
		),
	)
	flags.BoolVarP(
		&runner.watch,
		"watch",
		"w",
		false,
		"Watch for changes after update (clusters only)",
	)
	return result
}

type runnerContext struct {
	logger         *slog.Logger
	engine         *templating.Engine
	console        *terminal.Console
	format         string
	watch          bool
	conn           *grpc.ClientConn
	marshalOptions protojson.MarshalOptions
	helper         *reflection.ObjectHelper
}

func (c *runnerContext) run(cmd *cobra.Command, args []string) error {
	var err error

	// Get the context:
	ctx := cmd.Context()

	// Get the logger and the console:
	c.logger = logging.LoggerFromContext(ctx)
	c.console = terminal.ConsoleFromContext(ctx)

	// Get the configuration:
	cfg, err := config.Load(ctx)
	if err != nil {
		return err
	}
	if cfg == nil {
		return fmt.Errorf("there is no configuration, run the 'login' command")
	}

	// Create the gRPC connection from the configuration:
	c.conn, err = cfg.Connect(ctx, cmd.Flags())
	if err != nil {
		return fmt.Errorf("failed to create gRPC connection: %w", err)
	}
	defer c.conn.Close()

	// Create the reflection helper:
	helper, err := reflection.NewHelper().
		SetLogger(c.logger).
		SetConnection(c.conn).
		AddPackages(cfg.Packages()...).
		Build()
	if err != nil {
		return fmt.Errorf("failed to create reflection tool: %w", err)
	}

	// Create the templating engine:
	c.engine, err = templating.NewEngine().
		SetLogger(c.logger).
		SetFS(templatesFS).
		SetDir("templates").
		Build()
	if err != nil {
		return fmt.Errorf("failed to create templating engine: %w", err)
	}

	// Check that the object type has been specified:
	if len(args) == 0 {
		c.console.Render(ctx, c.engine, "no_object.txt", map[string]any{
			"Helper": helper,
			"Binary": os.Args[0],
		})
		return nil
	}

	// Get the information about the object type:
	c.helper = helper.Lookup(args[0])
	if c.helper == nil {
		c.console.Render(ctx, c.engine, "wrong_object.txt", map[string]any{
			"Helper": helper,
			"Binary": os.Args[0],
			"Object": args[0],
		})
		return nil
	}

	// Check the flags:
	if c.format != outputFormatJson && c.format != outputFormatYaml {
		return fmt.Errorf(
			"unknown output format '%s', should be '%s' or '%s'",
			c.format, outputFormatJson, outputFormatYaml,
		)
	}

	// Check that the object identifier has been specified:
	if len(args) < 2 {
		c.console.Render(ctx, c.engine, "no_id.txt", map[string]any{
			"Binary": os.Args[0],
		})
		return nil
	}
	objectId := args[1]

	// Create the gRPC connection from the configuration:
	c.conn, err = cfg.Connect(ctx, cmd.Flags())
	if err != nil {
		return fmt.Errorf("failed to create gRPC connection: %w", err)
	}
	defer c.conn.Close()

	// Get the current representation of the object:
	object, err := c.get(ctx, objectId)
	if err != nil {
		return fmt.Errorf("failed to get object of type '%s' with identifier '%s': %w", c.helper, objectId, err)
	}

	// Render the object:
	var render func(proto.Message) ([]byte, error)
	switch c.format {
	case outputFormatJson:
		render = c.renderJson
	default:
		render = c.renderYaml
	}
	data, err := render(object)
	if err != nil {
		return err
	}

	// Write the rendered object to a temporary file:
	tmpDir, err := os.MkdirTemp("", "")
	if err != nil {
		return err
	}
	defer func() {
		err := os.RemoveAll(tmpDir)
		if err != nil {
			c.logger.ErrorContext(
				ctx,
				"Failed to remove temporary directory",
				slog.String("dir", tmpDir),
				slog.Any("error", err),
			)
		}
	}()
	tmpFile := filepath.Join(tmpDir, fmt.Sprintf("%s-%s.%s", c.helper, objectId, c.format))
	err = os.WriteFile(tmpFile, data, 0600)
	if err != nil {
		return fmt.Errorf("failed to create temporary file '%s': %w", tmpFile, err)
	}

	// Run the editor:
	editorName := c.findEditor(ctx)
	editorPath, err := exec.LookPath(editorName)
	if err != nil {
		return fmt.Errorf("failed to find editor command '%s': %w", editorName, err)
	}
	editorCmd := &exec.Cmd{
		Path: editorPath,
		Args: []string{
			editorName,
			tmpFile,
		},
		Stdin:  os.Stdin,
		Stdout: os.Stdout,
		Stderr: os.Stderr,
	}
	err = editorCmd.Run()
	if err != nil {
		return fmt.Errorf("failed to edit: %w", err)
	}

	// Load the potentiall modified file:
	data, err = os.ReadFile(tmpFile)
	if err != nil {
		return fmt.Errorf("failed to read back temporary file '%s': %w", tmpFile, err)
	}

	// Parse the result:
	var parse func([]byte) (proto.Message, error)
	switch c.format {
	case outputFormatJson:
		parse = c.parseJson
	default:
		parse = c.parseYaml
	}
	object, err = parse(data)
	if err != nil {
		return fmt.Errorf("failed to parse modified object: %w", err)
	}

	// Save the result:
	updated, err := c.update(ctx, object)
	if err != nil {
		return err
	}

	// Show feedback about the update
	return c.showUpdateFeedback(ctx, updated, objectId)
}

// findEditor tries to find the name of the editor command. It will first try with the content of the `EDITOR` and
// `VISUAL` environment variables, and if those are empty it defaults to `vi`.
func (c *runnerContext) findEditor(ctx context.Context) string {
	for _, editorEnvVar := range editorEnvVars {
		value, ok := os.LookupEnv(editorEnvVar)
		if ok && value != "" {
			c.logger.DebugContext(
				ctx,
				"Found editor using environment variable",
				slog.String("var", editorEnvVar),
				slog.String("value", value),
			)
			return value
		}
	}
	c.logger.InfoContext(
		ctx,
		"Didn't find a editor in the environment, will use the default",
		slog.Any("vars", editorEnvVars),
		slog.String("default", defaultEditor),
	)
	return defaultEditor
}

func (c *runnerContext) get(ctx context.Context, id string) (result proto.Message, err error) {
	result, err = c.helper.Get(ctx, id)
	return
}

func (c *runnerContext) update(ctx context.Context, object proto.Message) (result proto.Message, err error) {
	result, err = c.helper.Update(ctx, object)
	return
}

func (c *runnerContext) renderJson(object proto.Message) (result []byte, err error) {
	result, err = c.marshalOptions.Marshal(object)
	return
}

func (c *runnerContext) renderYaml(object proto.Message) (result []byte, err error) {
	data, err := c.renderJson(object)
	if err != nil {
		return
	}
	var value any
	err = json.Unmarshal(data, &value)
	if err != nil {
		return
	}
	buffer := &bytes.Buffer{}
	encoder := yaml.NewEncoder(buffer)
	encoder.SetIndent(2)
	err = encoder.Encode(value)
	if err != nil {
		return
	}
	result = buffer.Bytes()
	return
}

func (c *runnerContext) parseJson(data []byte) (result proto.Message, err error) {
	object := c.helper.Instance()
	err = protojson.Unmarshal(data, object)
	if err != nil {
		return
	}
	result = object
	return
}

func (c *runnerContext) parseYaml(data []byte) (result proto.Message, err error) {
	var value any
	err = yaml.Unmarshal(data, &value)
	if err != nil {
		return
	}
	data, err = json.Marshal(value)
	if err != nil {
		return
	}
	result, err = c.parseJson(data)
	return
}

// editorEnvVars is the list of environment variables that will be used to obtain the name of the editor command.
var editorEnvVars = []string{
	"EDITOR",
	"VISUAL",
}

// defualtEditor is the editor used when the environment variables don't indicate any other editor.
const defaultEditor = "vi"

// watchPollInterval is the interval between polls when watching for changes.
const watchPollInterval = 5 * time.Second

// showUpdateFeedback displays feedback about the update operation, and optionally watches for changes.
func (c *runnerContext) showUpdateFeedback(ctx context.Context, updated proto.Message, objectId string) error {
	// Always show confirmation
	c.console.Printf(ctx, "Updated %s '%s'.\n", c.helper.Singular(), objectId)

	// Check if this is a cluster (only clusters support watch currently)
	cluster, isCluster := updated.(*ffv1.Cluster)
	if !isCluster {
		// For non-cluster objects, just confirm the update
		return nil
	}

	// Check if change is in progress
	if !c.isClusterChanging(cluster) {
		// No change in progress, we're done
		return nil
	}

	// Show that change is in progress
	c.console.Printf(ctx, "\nChange in progress...\n")
	c.showClusterConditions(ctx, cluster)

	if c.watch {
		// Watch mode: poll for status updates
		return c.watchClusterProgress(ctx, objectId)
	}

	// Non-watch mode: tell user how to check progress
	c.console.Printf(ctx, "\nUse '%s get %s %s' to check progress.\n",
		os.Args[0], c.helper.Singular(), objectId)

	return nil
}

// isClusterChanging checks if a cluster has changes in progress by comparing desired spec vs actual status.
// Note: The cluster state remains READY during scaling operations, so we cannot rely on checking the state.
// Instead, we compare spec.node_sets vs status.node_sets to detect ongoing changes.
func (c *runnerContext) isClusterChanging(cluster *ffv1.Cluster) bool {
	for name, desired := range cluster.Spec.NodeSets {
		actual, exists := cluster.Status.NodeSets[name]
		if !exists || desired.Size != actual.Size {
			return true
		}
	}

	return false
}

// showClusterConditions displays cluster conditions with their messages.
func (c *runnerContext) showClusterConditions(ctx context.Context, cluster *ffv1.Cluster) {
	if cluster.Status == nil || cluster.Status.Conditions == nil {
		return
	}

	for _, cond := range cluster.Status.Conditions {
		// Show any condition with a message, regardless of status
		// TRUE status means "this is happening", FALSE means "this is NOT happening (blocked)"
		if cond.Message != nil && *cond.Message != "" {
			c.console.Printf(ctx, "  %s\n", *cond.Message)
		}
	}
}

// watchClusterProgress polls the cluster status and displays progress until it reaches a terminal state.
func (c *runnerContext) watchClusterProgress(ctx context.Context, objectId string) error {
	c.console.Printf(ctx, "\nWatching for changes (Ctrl+C to stop)...\n\n")

	lastState := ""
	shownMessages := make(map[string]bool)

	// Do an immediate poll before starting the ticker
	done, err := c.pollClusterStatus(ctx, objectId, &lastState, shownMessages)
	if err != nil {
		return err
	}
	if done {
		return nil
	}

	ticker := time.NewTicker(watchPollInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-ticker.C:
			done, err := c.pollClusterStatus(ctx, objectId, &lastState, shownMessages)
			if err != nil {
				return err
			}
			if done {
				return nil
			}
		}
	}
}

// pollClusterStatus polls the cluster once and displays any status updates.
// Returns (true, nil) if polling should stop, (false, nil) to continue, or (false, err) on error.
func (c *runnerContext) pollClusterStatus(ctx context.Context, objectId string, lastState *string, shownMessages map[string]bool) (bool, error) {
	// Poll for current status
	current, err := c.helper.Get(ctx, objectId)
	if err != nil {
		c.logger.WarnContext(
			ctx,
			"Failed to get status",
			slog.String("id", objectId),
			slog.Any("error", err),
		)
		return false, nil
	}

	cluster, ok := current.(*ffv1.Cluster)
	if !ok {
		c.logger.ErrorContext(ctx, "Unexpected object type")
		return false, fmt.Errorf("unexpected object type")
	}

	// Show state changes
	stateStr := cluster.Status.State.String()
	if stateStr != *lastState {
		timestamp := time.Now().Format(time.TimeOnly)
		displayState := strings.TrimPrefix(stateStr, "CLUSTER_STATE_")
		c.console.Printf(ctx, "[%s] State: %s\n", timestamp, displayState)
		*lastState = stateStr
	}

	// Show new condition messages
	if cluster.Status.Conditions != nil {
		for _, cond := range cluster.Status.Conditions {
			// Show any condition with a message, regardless of status
			if cond.Message != nil && *cond.Message != "" {
				// Only show each unique message once
				if !shownMessages[*cond.Message] {
					timestamp := time.Now().Format(time.TimeOnly)
					c.console.Printf(ctx, "[%s] %s\n", timestamp, *cond.Message)
					shownMessages[*cond.Message] = true
				}
			}
		}
	}

	// Check if changes are complete (spec matches status)
	if !c.isClusterChanging(cluster) {
		c.console.Printf(ctx, "\nCluster update complete.\n")
		return true, nil
	}

	return false, nil
}

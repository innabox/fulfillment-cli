/*
Copyright (c) 2025 Red Hat Inc.

Licensed under the Apache License, Version 2.0 (the "License"); you may not use this file except in compliance with the
License. You may obtain a copy of the License at

  http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software distributed under the License is distributed on an
"AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the License for the specific
language governing permissions and limitations under the License.
*/

package virtualmachine

import (
	"context"
	"embed"
	"errors"
	"fmt"
	"log/slog"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/spf13/cobra"
	grpccodes "google.golang.org/grpc/codes"
	grpcstatus "google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/anypb"
	"google.golang.org/protobuf/types/known/durationpb"
	"google.golang.org/protobuf/types/known/timestamppb"
	"google.golang.org/protobuf/types/known/wrapperspb"

	ffv1 "github.com/innabox/fulfillment-cli/internal/api/fulfillment/v1"
	"github.com/innabox/fulfillment-cli/internal/config"
	"github.com/innabox/fulfillment-cli/internal/exit"
	"github.com/innabox/fulfillment-cli/internal/logging"
	"github.com/innabox/fulfillment-cli/internal/templating"
	"github.com/innabox/fulfillment-cli/internal/terminal"
)

//go:embed templates
var templatesFS embed.FS

func Cmd() *cobra.Command {
	runner := &runnerContext{}
	result := &cobra.Command{
		Use:   "virtualmachine [flags]",
		Short: "Create a virtual machine",
		RunE:  runner.run,
	}
	flags := result.Flags()
	flags.StringVarP(
		&runner.template,
		"template",
		"t",
		"",
		"Template identifier",
	)
	flags.StringSliceVarP(
		&runner.templateParameterValues,
		"template-parameter",
		"p",
		[]string{},
		"Template parameter in the format 'name=value'",
	)
	flags.StringSliceVarP(
		&runner.templateParameterFiles,
		"template-parameter-file",
		"f",
		[]string{},
		"Template parameter from file in the format 'name=filename'",
	)
	return result
}

type runnerContext struct {
	template                string
	templateParameterValues []string
	templateParameterFiles  []string
	logger                  *slog.Logger
	console                 *terminal.Console
	engine                  *templating.Engine
	templatesClient         ffv1.VirtualMachineTemplatesClient
	virtualMachinesClient   ffv1.VirtualMachinesClient
}

func (c *runnerContext) run(cmd *cobra.Command, args []string) error {
	var err error

	// Get the context:
	ctx := cmd.Context()

	// Get the logger and console:
	c.logger = logging.LoggerFromContext(ctx)
	c.console = terminal.ConsoleFromContext(ctx)

	// Create the templating engine:
	c.engine, err = templating.NewEngine().
		SetLogger(c.logger).
		SetFS(templatesFS).
		SetDir("templates").
		Build()
	if err != nil {
		return fmt.Errorf("failed to create templating engine: %w", err)
	}

	// Get the configuration:
	cfg, err := config.Load()
	if err != nil {
		return err
	}
	if cfg.Address == "" {
		return fmt.Errorf("there is no configuration, run the 'login' command")
	}

	// Check that we have a template:
	if c.template == "" {
		return fmt.Errorf("template identifier is required")
	}

	// Create the gRPC connection from the configuration:
	conn, err := cfg.Connect(ctx, cmd.Flags())
	if err != nil {
		return fmt.Errorf("failed to create gRPC connection: %w", err)
	}
	defer conn.Close()

	// Create the gRPC clients:
	c.templatesClient = ffv1.NewVirtualMachineTemplatesClient(conn)
	c.virtualMachinesClient = ffv1.NewVirtualMachinesClient(conn)

	// Fetch the virtual machine template:
	templateResponse, err := c.templatesClient.Get(ctx, ffv1.VirtualMachineTemplatesGetRequest_builder{
		Id: c.template,
	}.Build())
	if err != nil {
		status, ok := grpcstatus.FromError(err)
		if ok {
			if status.Code() == grpccodes.NotFound {
				templatesResponse, err := c.templatesClient.List(ctx, ffv1.VirtualMachineTemplatesListRequest_builder{
					Limit: proto.Int32(50),
				}.Build())
				if err != nil {
					return fmt.Errorf("failed to list templates: %w", err)
				}
				templates := templatesResponse.GetItems()
				sort.Slice(templates, func(i, j int) bool {
					return templates[i].GetId() < templates[j].GetId()
				})
				c.console.Render(ctx, c.engine, "template_not_found.txt", map[string]any{
					"Binary":    os.Args[0],
					"Template":  c.template,
					"Templates": templates,
				})
				return exit.Error(1)
			}
			return fmt.Errorf("failed to get template '%s': %w", c.template, err)
		}
		return fmt.Errorf("failed to get template '%s': %w", c.template, err)
	}
	template := templateResponse.Object
	if template == nil {
		return exit.Error(1)
	}

	// Parse the template parameters:
	templateParameterValues, templateParameterIssues := c.parseTemplateParameters(ctx, template)
	if len(templateParameterIssues) > 0 {
		validTemplateParameters := c.validTemplateParameters(template)
		c.console.Render(ctx, c.engine, "template_parameter_issues.txt", map[string]any{
			"Binary":     os.Args[0],
			"Template":   c.template,
			"Parameters": validTemplateParameters,
			"Issues":     templateParameterIssues,
		})
		return exit.Error(1)
	}

	// Prepare the virtual machine:
	virtualMachine := ffv1.VirtualMachine_builder{
		Spec: ffv1.VirtualMachineSpec_builder{
			Template:           c.template,
			TemplateParameters: templateParameterValues,
		}.Build(),
	}.Build()

	// Create the virtual machine:
	response, err := c.virtualMachinesClient.Create(ctx, ffv1.VirtualMachinesCreateRequest_builder{
		Object: virtualMachine,
	}.Build())
	if err != nil {
		return fmt.Errorf("failed to create virtual machine: %w", err)
	}

	// Display the result:
	virtualMachine = response.Object
	c.console.Printf(ctx, "Created virtual machine '%s'.\n", virtualMachine.Id)

	return nil
}

// parseTemplateParameters parses the '--template-parameter' and '--template-parameter-file' flags into a map of
// parameter name to value, and a list of issues found. The issues are intended for display to the user.
func (c *runnerContext) parseTemplateParameters(ctx context.Context,
	template *ffv1.VirtualMachineTemplate) (result map[string]*anypb.Any, issues []string) {
	// Prepare empty results and issues:
	result = map[string]*anypb.Any{}

	// Make a map of parameter definitions indexed by name for quick lookup:
	definitions := map[string]*ffv1.VirtualMachineTemplateParameterDefinition{}
	for _, definition := range template.GetParameters() {
		definitions[definition.GetName()] = definition
	}

	// Parse '--template-parameter' flags:
	for _, flag := range c.templateParameterValues {
		parts := strings.SplitN(flag, "=", 2)
		if len(parts) != 2 {
			name := strings.TrimSpace(flag)
			definition := definitions[name]
			if definition == nil {
				issues = append(
					issues,
					fmt.Sprintf(
						"In '%s' parameter '%s' doesn't exist, and if it existed the value "+
							"would be missing",
						flag, name,
					),
				)
			} else {
				issues = append(
					issues,
					fmt.Sprintf(
						"In '%s' parameter value is missing",
						flag,
					),
				)
			}
			continue
		}
		name := strings.TrimSpace(parts[0])
		if name == "" {
			issues = append(
				issues,
				fmt.Sprintf(
					"In '%s' parameter name is missing",
					flag,
				),
			)
			continue
		}
		definition := definitions[name]
		if definition == nil {
			issues = append(
				issues,
				fmt.Sprintf(
					"In '%s' parameter '%s' doesn't exist",
					flag, name,
				),
			)
			continue
		}
		text := strings.TrimSpace(parts[1])
		value, issue := c.convertTextToTemplateParameterValue(ctx, text, definition.GetType())
		if issue != "" {
			issues = append(issues, fmt.Sprintf("In '%s' %s", flag, issue))
			continue
		}
		result[name] = value
	}

	// Parse '--template-parameter-file' flags:
	for _, flag := range c.templateParameterFiles {
		parts := strings.SplitN(flag, "=", 2)
		if len(parts) != 2 {
			name := strings.TrimSpace(flag)
			definition := definitions[name]
			if definition == nil {
				issues = append(issues, fmt.Sprintf(
					"In '%s' parameter '%s' doesn't exist, and if existed the file would be "+
						"missing",
					flag, name,
				))
			} else {
				issues = append(
					issues,
					fmt.Sprintf(
						"In '%s' file is missing",
						flag,
					))
			}
			continue
		}
		name := strings.TrimSpace(parts[0])
		if name == "" {
			issues = append(
				issues,
				fmt.Sprintf(
					"In '%s' parameter name is missing",
					flag,
				),
			)
			continue
		}
		definition := definitions[name]
		if definition == nil {
			issues = append(
				issues,
				fmt.Sprintf(
					"In '%s' parameter '%s' doesn't exist",
					flag, name,
				),
			)
			continue
		}
		file := strings.TrimSpace(parts[1])
		if file == "" {
			issues = append(
				issues,
				fmt.Sprintf(
					"In '%s' file is missing",
					flag,
				),
			)
			continue
		}
		data, err := os.ReadFile(file)
		if errors.Is(err, os.ErrNotExist) {
			issues = append(
				issues, fmt.Sprintf(
					"In '%s' file '%s' doesn't exist",
					flag, file,
				),
			)
			continue
		}
		if err != nil {
			issues = append(
				issues,
				fmt.Sprintf(
					"In '%s' failed to read file '%s': %w",
					file, err,
				),
			)
			continue
		}
		text := string(data)
		value, issue := c.convertTextToTemplateParameterValue(ctx, text, definition.GetType())
		if issue != "" {
			issues = append(
				issues,
				fmt.Sprintf("In '%s' %s'", flag, issue),
			)
			continue
		}
		result[name] = value
	}

	// Add issues for missing required parameters, at the end of the list and sorted by parameter name:
	var missing []*ffv1.VirtualMachineTemplateParameterDefinition
	for _, definition := range template.GetParameters() {
		if definition.GetRequired() && result[definition.GetName()] == nil {
			missing = append(missing, definition)
		}
	}
	sort.Slice(missing, func(i, j int) bool {
		return missing[i].GetName() < missing[j].GetName()
	})
	for _, definition := range missing {
		issues = append(
			issues,
			fmt.Sprintf("Parameter '%s' is required", definition.GetName()),
		)
	}

	return
}

// convertTextToTemplateParameterValue converts a string value to the appropriate protobuf type based on the kind. It
// returns the value and a string descibing the issue if the conversion fails.
func (c *runnerContext) convertTextToTemplateParameterValue(ctx context.Context, text,
	kind string) (result *anypb.Any, issue string) {
	var wrapper proto.Message
	switch kind {
	case "type.googleapis.com/google.protobuf.StringValue":
		wrapper = &wrapperspb.StringValue{Value: text}
	case "type.googleapis.com/google.protobuf.BoolValue":
		text = strings.TrimSpace(text)
		value, err := strconv.ParseBool(text)
		if err != nil {
			c.logger.DebugContext(
				ctx,
				"Failed to parse boolean",
				slog.String("text", text),
				slog.Any("error", err),
			)
			issue = fmt.Sprintf(
				"value '%s' isn't a valid boolean, valid values are 'true' and 'false'",
				text,
			)
			return
		}
		wrapper = &wrapperspb.BoolValue{Value: value}
	case "type.googleapis.com/google.protobuf.Int32Value":
		text = strings.TrimSpace(text)
		var value int64
		value, err := strconv.ParseInt(text, 10, 32)
		if err != nil {
			c.logger.DebugContext(
				ctx,
				"Failed to parse 32-bit integer number",
				slog.String("text", text),
				slog.Any("error", err),
			)
			issue = fmt.Sprintf("value '%s' isn't a valid 32-bit integer", text)
			return
		}
		wrapper = &wrapperspb.Int32Value{Value: int32(value)}
	case "type.googleapis.com/google.protobuf.Int64Value":
		text = strings.TrimSpace(text)
		var value int64
		value, err := strconv.ParseInt(text, 10, 64)
		if err != nil {
			c.logger.DebugContext(
				ctx,
				"Failed to parse 64-bit integer number",
				slog.String("text", text),
				slog.Any("error", err),
			)
			issue = fmt.Sprintf("value '%s' isn't a valid 64-bit integer", text)
			return
		}
		wrapper = &wrapperspb.Int64Value{Value: value}
	case "type.googleapis.com/google.protobuf.FloatValue":
		text = strings.TrimSpace(text)
		var value float64
		value, err := strconv.ParseFloat(text, 32)
		if err != nil {
			c.logger.DebugContext(
				ctx,
				"Failed to parse 32-bit floating point number",
				slog.String("text", text),
				slog.Any("error", err),
			)
			issue = fmt.Sprintf("value '%s' isn't a valid 32-bit floating point number", text)
			return
		}
		wrapper = &wrapperspb.FloatValue{Value: float32(value)}
	case "type.googleapis.com/google.protobuf.DoubleValue":
		text = strings.TrimSpace(text)
		var value float64
		value, err := strconv.ParseFloat(text, 64)
		if err != nil {
			c.logger.DebugContext(
				ctx,
				"Failed to parse 64-bit floating point number",
				slog.String("text", text),
				slog.Any("error", err),
			)
			issue = fmt.Sprintf("value '%s' isn't a valid 64-bit floating point numberw", text)
			return
		}
		wrapper = &wrapperspb.DoubleValue{Value: value}
	case "type.googleapis.com/google.protobuf.BytesValue":
		wrapper = &wrapperspb.BytesValue{Value: []byte(text)}
	case "type.googleapis.com/google.protobuf.Timestamp":
		text = strings.TrimSpace(text)
		var value time.Time
		value, err := time.Parse(time.RFC3339, text)
		if err != nil {
			c.logger.DebugContext(
				ctx,
				"Failed to parse RFC3339 timestamp",
				slog.String("text", text),
				slog.Any("error", err),
			)
			issue = fmt.Sprintf("value '%s' isn't a valid RFC3339 timestamp", text)
			return
		}
		wrapper = timestamppb.New(value)
	case "type.googleapis.com/google.protobuf.Duration":
		var value time.Duration
		value, err := time.ParseDuration(text)
		if err != nil {
			c.logger.DebugContext(
				ctx,
				"Failed to parse duration",
				slog.String("text", text),
				slog.Any("error", err),
			)
			issue = fmt.Sprintf("value '%s' isn't a valid duration", text)
			return
		}
		wrapper = durationpb.New(value)
	default:
		issue = fmt.Sprintf("flag has is of an unsupported type '%s'", kind)
		return
	}
	if issue != "" {
		return
	}
	result, err := anypb.New(wrapper)
	if err != nil {
		c.logger.DebugContext(
			ctx,
			"Failed to create protobuf value for template parameter",
			slog.String("text", text),
			slog.String("kind", kind),
			slog.Any("error", err),
		)
		issue = fmt.Sprintf("Failed to create protobuf value for template parameter: %w", err)
		return
	}
	return
}

// validTemplateParameter contains the information about a valid template parameter, for use in the error messages that
// display them.
type validTemplateParameter struct {
	// Name is the name of the parameter.
	Name string

	// Type is the type of the parameter.
	Type string

	// Title is the title of the parameter.
	Title string
}

// validTemplateParameters returns the list of valid template parameters for the given template.
func (c *runnerContext) validTemplateParameters(template *ffv1.VirtualMachineTemplate) []validTemplateParameter {
	// Prepare the results:
	results := []validTemplateParameter{}
	for _, parameter := range template.GetParameters() {
		result := validTemplateParameter{
			Name:  parameter.GetName(),
			Title: parameter.GetTitle(),
		}
		switch parameter.GetType() {
		case "type.googleapis.com/google.protobuf.StringValue":
			result.Type = "string"
		case "type.googleapis.com/google.protobuf.BoolValue":
			result.Type = "boolean"
		case "type.googleapis.com/google.protobuf.Int32Value":
			result.Type = "int32"
		case "type.googleapis.com/google.protobuf.Int64Value":
			result.Type = "int64"
		case "type.googleapis.com/google.protobuf.FloatValue":
			result.Type = "float"
		case "type.googleapis.com/google.protobuf.DoubleValue":
			result.Type = "double"
		case "type.googleapis.com/google.protobuf.BytesValue":
			result.Type = "bytes"
		case "type.googleapis.com/google.protobuf.Timestamp":
			result.Type = "timestamp"
		case "type.googleapis.com/google.protobuf.Duration":
			result.Type = "duration"
		default:
			result.Type = "unknown"
		}
		results = append(results, result)
	}

	// Sort the result by name so that the output will be predictable:
	sort.Slice(results, func(i, j int) bool {
		return results[i].Name < results[j].Name
	})

	return results
}

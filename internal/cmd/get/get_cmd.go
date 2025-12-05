/*
Copyright (c) 2025 Red Hat Inc.

Licensed under the Apache License, Version 2.0 (the "License"); you may not use this file except in compliance with the
License. You may obtain a copy of the License at

  http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software distributed under the License is distributed on an
"AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the License for the specific
language governing permissions and limitations under the License.
*/

package get

import (
	"context"
	"embed"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"os"
	"path"
	"slices"
	"strconv"
	"strings"
	"text/tabwriter"

	"github.com/google/cel-go/cel"
	"github.com/google/cel-go/common/types"
	"github.com/google/cel-go/common/types/ref"
	"github.com/google/cel-go/ext"
	"github.com/innabox/fulfillment-common/logging"
	"github.com/innabox/fulfillment-common/templating"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/reflect/protoregistry"
	"google.golang.org/protobuf/types/dynamicpb"
	"google.golang.org/protobuf/types/known/anypb"
	"gopkg.in/yaml.v3"

	"github.com/innabox/fulfillment-cli/internal/cmd/get/kubeconfig"
	"github.com/innabox/fulfillment-cli/internal/cmd/get/password"
	"github.com/innabox/fulfillment-cli/internal/cmd/get/token"
	"github.com/innabox/fulfillment-cli/internal/config"
	"github.com/innabox/fulfillment-cli/internal/reflection"
	"github.com/innabox/fulfillment-cli/internal/terminal"
)

//go:embed templates
var templatesFS embed.FS

//go:embed tables
var tablesFS embed.FS

// Possible output formats:
const (
	outputFormatTable = "table"
	outputFormatJson  = "json"
	outputFormatYaml  = "yaml"
)

func Cmd() *cobra.Command {
	runner := &runnerContext{
		marshalOptions: protojson.MarshalOptions{
			UseProtoNames: true,
		},
	}
	result := &cobra.Command{
		Use:   "get OBJECT [OPTION]... [ID|NAME]...",
		Short: "Get objects",
		RunE:  runner.run,
	}
	result.AddCommand(kubeconfig.Cmd())
	result.AddCommand(password.Cmd())
	result.AddCommand(token.Cmd())
	flags := result.Flags()
	flags.StringVarP(
		&runner.args.format,
		"output",
		"o",
		outputFormatTable,
		fmt.Sprintf(
			"Output format, one of '%s', '%s' or '%s'.",
			outputFormatTable, outputFormatJson, outputFormatYaml,
		),
	)
	flags.StringVar(
		&runner.args.filter,
		"filter",
		"",
		"CEL expression used for filtering results.",
	)
	flags.BoolVar(
		&runner.args.includeDeleted,
		"include-deleted",
		false,
		"Include deleted objects.",
	)
	flags.BoolVarP(
		&runner.args.watch,
		"watch",
		"w",
		false,
		"Watch for changes to objects",
	)
	return result
}

type runnerContext struct {
	args struct {
		format         string
		filter         string
		includeDeleted bool
		watch          bool
	}
	ctx            context.Context
	logger         *slog.Logger
	engine         *templating.Engine
	console        *terminal.Console
	conn           *grpc.ClientConn
	marshalOptions protojson.MarshalOptions
	globalHelper   *reflection.Helper
	objectHelper   *reflection.ObjectHelper
	lookupCache    map[protoreflect.FullName]map[string]string
}

func (c *runnerContext) run(cmd *cobra.Command, args []string) error {
	var err error

	// Get the context:
	ctx := cmd.Context()

	// Save the context. This is needed because some of the CEL functions that we create need the context, but
	// there is no way to pass it directly. Refrain from using tis for other purposes.
	c.ctx = ctx

	// Get the logger and console:
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
	c.globalHelper, err = reflection.NewHelper().
		SetLogger(c.logger).
		SetConnection(c.conn).
		AddPackages(cfg.Packages()).
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
			"Helper": c.globalHelper,
			"Binary": os.Args[0],
		})
		return nil
	}

	// Get the object helper:
	c.objectHelper = c.globalHelper.Lookup(args[0])
	if c.objectHelper == nil {
		c.console.Render(ctx, c.engine, "wrong_object.txt", map[string]any{
			"Helper": c.globalHelper,
			"Binary": os.Args[0],
			"Object": args[0],
		})
		return nil
	}

	// Check the flags:
	if c.args.format != outputFormatTable && c.args.format != outputFormatJson && c.args.format != outputFormatYaml {
		return fmt.Errorf(
			"unknown output format '%s', should be '%s', '%s' or '%s'",
			c.args.format, outputFormatTable, outputFormatJson, outputFormatYaml,
		)
	}

	// If watch mode is enabled, watch for events instead of listing
	if c.args.watch {
		return c.watch(ctx, args[1:])
	}

	// Get the objects using the list method, which will handle filtering by identifiers or names if provided.
	objects, err := c.list(ctx, args[1:])
	if err != nil {
		return err
	}

	// Render the items:
	var render func(context.Context, []proto.Message) error
	switch c.args.format {
	case outputFormatJson:
		render = c.renderJson
	case outputFormatYaml:
		render = c.renderYaml
	default:
		render = c.renderTable
	}
	return render(ctx, objects)
}

func (c *runnerContext) list(ctx context.Context, keys []string) (results []proto.Message, err error) {
	var options reflection.ListOptions

	// If keys (identifiers or names) were provided, build a CEL filter to match them.
	if len(keys) > 0 {
		var values []string
		for _, key := range keys {
			values = append(values, strconv.Quote(key))
		}
		list := strings.Join(values, ", ")
		options.Filter = fmt.Sprintf(
			`this.id in [%[1]s] || this.metadata.name in [%[1]s]`,
			list,
		)
	}

	// Apply the user-provided filter if specified.
	if c.args.filter != "" {
		if options.Filter != "" {
			options.Filter = fmt.Sprintf("(%s) && (%s)", options.Filter, c.args.filter)
		} else {
			options.Filter = c.args.filter
		}
	}

	// Exclude deleted objects unless explicitly requested.
	if !c.args.includeDeleted {
		const notDeletedFilter = "!has(this.metadata.deletion_timestamp)"
		if options.Filter != "" {
			options.Filter = fmt.Sprintf("%s && (%s)", notDeletedFilter, options.Filter)
		} else {
			options.Filter = notDeletedFilter
		}
	}

	listResult, err := c.objectHelper.List(ctx, options)
	if err != nil {
		return
	}
	results = listResult.Items
	return
}

func (c *runnerContext) loadTable() (result *Table, err error) {
	file := fmt.Sprintf("%s.yaml", c.objectHelper.FullName())
	data, err := tablesFS.ReadFile(path.Join("tables", file))
	if err != nil {
		err = fmt.Errorf(
			"failed to read table definition file '%s': %w",
			file, err,
		)
		return
	}
	var table Table
	err = yaml.Unmarshal(data, &table)
	if err != nil {
		err = fmt.Errorf(
			"failed to unmarshal table definition file '%s': %w",
			file, err,
		)
		return
	}
	result = &table
	return
}

func (c *runnerContext) defaultTable() *Table {
	return &Table{
		Columns: []*Column{
			{
				Header: "ID",
				Value:  "this.id",
			},
			{
				Header: "NAME",
				Value:  "has(this.metadata.name)? this.metadata.name: '-'",
			},
		},
	}
}

func (c *runnerContext) renderTable(ctx context.Context, objects []proto.Message) error {
	// Check if there are results:
	if len(objects) == 0 {
		c.console.Render(ctx, c.engine, "no_matching_objects.txt", nil)
		return nil
	}

	// Try to load the table that matches the object type:
	table, err := c.loadTable()
	if err != nil {
		return err
	}
	if table == nil {
		table = c.defaultTable()
	}

	// If the user has asked to include deleted objects then add the deletion timestamp column:
	if c.args.includeDeleted {
		deletedCol := &Column{
			Header: "DELETED",
			Value:  "has(this.metadata.deletion_timestamp)? string(this.metadata.deletion_timestamp): '-'",
		}
		table.Columns = slices.Insert(table.Columns, 1, deletedCol)
	}

	// Initialize the lookup cache:
	c.lookupCache = map[protoreflect.FullName]map[string]string{}

	// Add all file descriptors from the current object's package:
	thisDesc := c.objectHelper.Descriptor()

	// Build CEL environment:
	celEnv, err := cel.NewEnv(
		cel.Types(dynamicpb.NewMessage(thisDesc)),
		cel.Variable("this", cel.ObjectType(string(thisDesc.FullName()))),
		ext.Strings(),
	)
	if err != nil {
		return fmt.Errorf("failed to create CEL environment: %w", err)
	}

	// Compile the CEL expressions for the columns:
	prgs := make([]cel.Program, len(table.Columns))
	for i, col := range table.Columns {
		ast, issues := celEnv.Compile(col.Value)
		err = issues.Err()
		if err != nil {
			return fmt.Errorf(
				"failed to compile CEL expression '%s' for column '%s' of type '%s': %w",
				col.Value, col.Header, c.objectHelper, err,
			)
		}
		prg, err := celEnv.Program(ast)
		if err != nil {
			return fmt.Errorf(
				"failed to create CEL program from expression '%s' for column '%s' of type '%s': %w",
				col.Value, col.Header, c.objectHelper, err,
			)
		}
		prgs[i] = prg
	}

	// Render the table:
	writer := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	c.renderTableHeader(writer, table.Columns)
	for _, object := range objects {
		err := c.renderTableRow(writer, table.Columns, prgs, object)
		if err != nil {
			return err
		}
	}
	writer.Flush()

	return nil
}

func (c *runnerContext) renderTableHeader(writer io.Writer, cols []*Column) error {
	for i, col := range cols {
		if i > 0 {
			fmt.Fprint(writer, "\t")
		}
		fmt.Fprintf(writer, "%s", col.Header)
	}
	fmt.Fprintf(writer, "\n")
	return nil
}

func (c *runnerContext) renderTableRow(writer io.Writer, cols []*Column, prgs []cel.Program,
	object proto.Message) error {
	// Wrap the object in a top-level "this" field to avoid conflicts with reserved words
	in := map[string]any{
		"this": object,
	}
	celVars, err := cel.PartialVars(in)
	if err != nil {
		return fmt.Errorf(
			"failed to set variables for CEL expression for type '%s': %w",
			c.objectHelper, err,
		)
	}
	for i := range len(cols) {
		if i > 0 {
			fmt.Fprintf(writer, "\t")
		}
		if err != nil {
			return err
		}
		col := cols[i]
		prg := prgs[i]
		var out ref.Val
		out, _, err = prg.Eval(celVars)
		if err != nil {
			return fmt.Errorf(
				"failed to evaluate CEL expression '%s' for column '%s' of type '%s': %w",
				col.Value, col.Header, c.objectHelper, err,
			)
		}
		err = c.renderTableCell(writer, col, out)
		if err != nil {
			return fmt.Errorf(
				"failed to render value '%s' for column '%s' of type '%s': %w",
				out, col.Header, c.objectHelper, err,
			)
		}
	}
	fmt.Fprintf(writer, "\n")
	return nil
}

func (c *runnerContext) renderTableCell(writer io.Writer, col *Column, val ref.Val) error {
	switch val := val.(type) {
	case types.Int:
		if col.Type != "" {
			enumType, _ := protoregistry.GlobalTypes.FindEnumByName(col.Type)
			if enumType != nil {
				return c.renderTableCellEnumType(writer, val, enumType.Descriptor())
			}
			c.logger.Error(
				"Failed to find enum type",
				slog.String("type", string(col.Type)),
			)
		}
	case types.String:
		if col.Lookup && col.Type != "" {
			messageType, _ := protoregistry.GlobalTypes.FindMessageByName(col.Type)
			if messageType != nil {
				return c.renderTableCellLookup(writer, val, messageType.Descriptor())
			}
		}
	}
	return c.renderTableCellAnyType(writer, val)
}

func (c *runnerContext) renderTableCellEnumType(writer io.Writer, val types.Int,
	enumDesc protoreflect.EnumDescriptor) error {
	// Get the text of the name of the enum value:
	valueDescs := enumDesc.Values()
	valueDesc := valueDescs.ByNumber(protoreflect.EnumNumber(val))
	if valueDesc == nil {
		_, err := fmt.Fprintf(writer, "UNKNOWN:%d", val)
		if err != nil {
			return err
		}
	}
	valueTxt := string(valueDesc.Name())

	// If the enum has been created according to our style guide then all the values should have a prefix with the
	// name of the type, for example `CLUSTER_ORDER_STATUS_STATE`. That prefix is not useful for humans, so we try
	// to remove it. To do so we find the value with number zero, which should end with `_UNSPECIFIED`, extract the
	// prefix from that and remove it from the representation of the value.
	unspecifiedDesc := valueDescs.ByNumber(protoreflect.EnumNumber(0))
	unspecifiedText := string(unspecifiedDesc.Name())
	prefixIndex := strings.LastIndex(unspecifiedText, "_")
	if prefixIndex != -1 {
		prefixTxt := unspecifiedText[0:prefixIndex]
		if strings.HasPrefix(valueTxt, prefixTxt) {
			valueTxt = valueTxt[prefixIndex+1:]
		}
	}

	_, err := fmt.Fprintf(writer, "%s", valueTxt)
	return err
}

func (c *runnerContext) renderTableCellLookup(writer io.Writer, val types.String,
	messageDesc protoreflect.MessageDescriptor) error {
	key := string(val)
	var text string
	if key != "" {
		text = c.lookupName(c.ctx, messageDesc.FullName(), key)
	} else {
		text = "-"
	}
	_, err := fmt.Fprintf(writer, "%s", text)
	return err
}

func (c *runnerContext) lookupName(ctx context.Context, messageFullName protoreflect.FullName,
	key string) (result string) {
	// Check if the result is already in the cache and return it immediately if so, otherwise
	// remember to update the cache when done.:
	lookupCache, ok := c.lookupCache[messageFullName]
	if !ok {
		lookupCache = map[string]string{}
		c.lookupCache[messageFullName] = lookupCache
	}
	result, ok = lookupCache[key]
	if ok {
		return result
	}
	defer func() {
		lookupCache[key] = result
	}()

	// Find the object helper:
	objectHelper := c.globalHelper.Lookup(string(messageFullName))
	if objectHelper == nil {
		c.logger.ErrorContext(
			ctx,
			"Failed to find object helper for type",
			slog.String("type", string(messageFullName)),
		)
		result = key
		return
	}

	// Find the objects whose identifier or name matches the key:
	filter := fmt.Sprintf(
		`this.id == %[1]q || this.metadata.name == %[1]q`,
		key,
	)
	listResult, err := objectHelper.List(ctx, reflection.ListOptions{
		Filter: filter,
	})
	if err != nil {
		c.logger.ErrorContext(
			ctx,
			"Failed to list objects for lookup",
			slog.String("type", string(messageFullName)),
			slog.String("key", key),
			slog.Any("error", err),
		)
		result = key
		return
	}

	// If there is no match, or multiple matches, return the original key:
	if len(listResult.Items) == 0 {
		result = key
		return
	}

	// Return the name of the first object:
	object := listResult.Items[0]
	metadata := objectHelper.GetMetadata(object)
	if metadata == nil {
		c.logger.ErrorContext(
			ctx,
			"Failed to get metadata for object",
			slog.String("type", string(messageFullName)),
			slog.String("key", key),
		)
		result = key
		return
	}
	result = metadata.GetName()
	return
}

func (c *runnerContext) renderTableCellAnyType(writer io.Writer, val ref.Val) error {
	_, err := fmt.Fprintf(writer, "%s", val)
	return err
}

func (c *runnerContext) renderJson(ctx context.Context, objects []proto.Message) error {
	values, err := c.encodeObjects(objects)
	if err != nil {
		return err
	}
	if len(values) == 1 {
		c.console.RenderJson(ctx, values[0])
	} else {
		c.console.RenderJson(ctx, values)
	}
	return nil
}

func (c *runnerContext) renderYaml(ctx context.Context, objects []proto.Message) error {
	values, err := c.encodeObjects(objects)
	if err != nil {
		return err
	}
	if len(values) == 1 {
		c.console.RenderYaml(ctx, values[0])
	} else {
		c.console.RenderYaml(ctx, values)
	}
	return nil
}

func (c *runnerContext) encodeObjects(objects []proto.Message) (result []any, err error) {
	values := make([]any, len(objects))
	for i, object := range objects {
		values[i], err = c.encodeObject(object)
		if err != nil {
			return
		}
	}
	result = values
	return
}

func (c *runnerContext) encodeObject(object proto.Message) (result any, err error) {
	wrapper, err := anypb.New(object)
	if err != nil {
		return
	}
	var data []byte
	data, err = c.marshalOptions.Marshal(wrapper)
	if err != nil {
		return
	}
	err = json.Unmarshal(data, &result)
	return
}

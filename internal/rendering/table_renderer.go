/*
Copyright (c) 2025 Red Hat Inc.

Licensed under the Apache License, Version 2.0 (the "License"); you may not use this file except in compliance with the
License. You may obtain a copy of the License at

  http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software distributed under the License is distributed on an
"AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the License for the specific
language governing permissions and limitations under the License.
*/

package rendering

import (
	"context"
	"embed"
	"fmt"
	"io"
	"io/fs"
	"log/slog"
	"os"
	"path"
	"slices"
	"strings"
	"text/tabwriter"

	"github.com/google/cel-go/cel"
	"github.com/google/cel-go/common/types"
	"github.com/google/cel-go/common/types/ref"
	"github.com/google/cel-go/ext"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/reflect/protoregistry"
	"google.golang.org/protobuf/types/dynamicpb"
	"gopkg.in/yaml.v3"

	"github.com/innabox/fulfillment-cli/internal/reflection"
)

//go:embed tables
var tablesFS embed.FS

// tableLayout describes how to render protocol buffers messages in tabular form.
type tableLayout struct {
	// Columns describes how fields of the message are mapped to columns.
	Columns []*columnLayout `yaml:"columns,omitempty"`
}

// columnLayout describes how to render a field of a protocol buffers message as a column in a table.
type columnLayout struct {
	// Header is the text of the header for the colum. The default is to use the name of the field in upper case
	// and replacing underscores with spaces.
	Header string `yaml:"header,omitempty"`

	// Value is a CEL expression that will be used to calculate the rendered value. The expression can access
	// the message via the `this` built-in variable.
	Value string `yaml:"value,omitempty"`

	// Type is the name of the type of the result of the expression. This is only needed when the result of the
	// expression is an enum value or an identifier that needs to be translated into a type.
	//
	// When the result is a enum value, then the 'type' field should contain the name of the enum type, and it
	// will be used to translate the integer value into the name of the enum value shortened to eliminate the
	// prefix common to all the enum values of that type.
	//
	// When the result is an identifier the 'type' field should be the name of the type, and it will be used to
	// find the name of the object.
	Type protoreflect.FullName `yaml:"type,omitempty"`

	// Lookup indicates if the result of the expression is an identifier that needs to be translated into a name.
	// When this is set to true the 'type' field also needs to be specified, and should contain the name of the
	// type to use for the lookup. For example, if the result of the expression is a cluster, then the 'type'
	// should be 'fulfillment.v1.Cluster'.
	Lookup bool `yaml:"lookup,omitempty"`
}

// TableRendererBuilder is used to create table renderers. Don't create instances of this type directly, use the
// NewTableRenderer function instead.
type TableRendererBuilder struct {
	logger         *slog.Logger
	helper         *reflection.Helper
	includeDeleted bool
}

// TableRenderer is responsible for rendering protocol buffer messages as tables. Don't create instances of this type
// directly, use the NewTableRenderer function instead.
type TableRenderer struct {
	logger         *slog.Logger
	helper         *reflection.Helper
	includeDeleted bool
}

// NewTableRenderer creates a new builder for table renderers.
func NewTableRenderer() *TableRendererBuilder {
	return &TableRendererBuilder{}
}

// SetLogger sets the logger that the renderer will use to write messages to the log. This is mandatory.
func (b *TableRendererBuilder) SetLogger(value *slog.Logger) *TableRendererBuilder {
	b.logger = value
	return b
}

// SetHelper sets the reflection helper that will be used to introspect objects. This is mandatory.
func (b *TableRendererBuilder) SetHelper(value *reflection.Helper) *TableRendererBuilder {
	b.helper = value
	return b
}

// SetIncludeDeleted sets whether to include the DELETED column in the output.
func (b *TableRendererBuilder) SetIncludeDeleted(value bool) *TableRendererBuilder {
	b.includeDeleted = value
	return b
}

// Build uses the data stored in the builder to create a new table renderer.
func (b *TableRendererBuilder) Build() (result *TableRenderer, err error) {
	// Check parameters:
	if b.logger == nil {
		err = fmt.Errorf("logger is mandatory")
		return
	}
	if b.helper == nil {
		err = fmt.Errorf("helper is mandatory")
		return
	}

	// Create and populate the object:
	result = &TableRenderer{
		logger:         b.logger,
		helper:         b.helper,
		includeDeleted: b.includeDeleted,
	}
	return
}

// Render renders the given objects as a table to stdout.
func (r *TableRenderer) Render(ctx context.Context, objects []proto.Message) error {
	// Check if there are any objects:
	if len(objects) == 0 {
		return nil
	}

	// Get the object helper from the first object:
	firstObject := objects[0]
	descriptor := firstObject.ProtoReflect().Descriptor()
	objectHelper := r.helper.Lookup(string(descriptor.FullName()))
	if objectHelper == nil {
		return fmt.Errorf("failed to find object helper for type %q", descriptor.FullName())
	}

	// Try to load the table definition for this object type:
	table, err := r.loadTable(objectHelper)
	if err != nil {
		return err
	}
	if table == nil {
		table = r.defaultTable()
	}

	// If the user has asked to include deleted objects then add the deletion timestamp column:
	if r.includeDeleted {
		deletedCol := &columnLayout{
			Header: "DELETED",
			Value:  "has(this.metadata.deletion_timestamp)? string(this.metadata.deletion_timestamp): '-'",
		}
		table.Columns = slices.Insert(table.Columns, 1, deletedCol)
	}

	// Initialize the lookup cache:
	lookupCache := map[protoreflect.FullName]map[string]string{}

	// Get the descriptor for the object type:
	thisDesc := objectHelper.Descriptor()

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
				"failed to compile CEL expression %q for column %q of type %q: %w",
				col.Value, col.Header, objectHelper, err,
			)
		}
		prg, err := celEnv.Program(ast)
		if err != nil {
			return fmt.Errorf(
				"failed to create CEL program from expression %q for column %q of type %q: %w",
				col.Value, col.Header, objectHelper, err,
			)
		}
		prgs[i] = prg
	}

	// Create a tab writer for proper column alignment to stdout:
	writer := tabwriter.NewWriter(io.Writer(os.Stdout), 0, 0, 2, ' ', 0)
	defer writer.Flush()

	// Render the header:
	err = r.renderTableHeader(writer, table.Columns)
	if err != nil {
		return err
	}

	// Render each row:
	for _, message := range objects {
		err := r.renderTableRow(ctx, writer, table.Columns, prgs, message, objectHelper, lookupCache)
		if err != nil {
			return err
		}
	}

	return nil
}

// loadTable loads the table definition for the given object type from the embedded filesystem.
func (r *TableRenderer) loadTable(objectHelper *reflection.ObjectHelper) (result *tableLayout, err error) {
	// Try to read the table definition file:
	file := fmt.Sprintf("%s.yaml", objectHelper.FullName())
	data, err := fs.ReadFile(tablesFS, path.Join("tables", file))
	if err != nil {
		// If the file doesn't exist, that's okay - we'll use the default table.
		return
	}

	// Unmarshal the table definition:
	var table tableLayout
	err = yaml.Unmarshal(data, &table)
	if err != nil {
		err = fmt.Errorf(
			"failed to unmarshal table definition file %q: %w",
			file, err,
		)
		return
	}
	result = &table
	return
}

// defaultTable returns a default table definition with ID and NAME columns.
func (r *TableRenderer) defaultTable() *tableLayout {
	return &tableLayout{
		Columns: []*columnLayout{
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

// renderTableHeader renders the table header with column names.
func (r *TableRenderer) renderTableHeader(writer io.Writer, cols []*columnLayout) error {
	for i, col := range cols {
		if i > 0 {
			fmt.Fprint(writer, "\t")
		}
		fmt.Fprintf(writer, "%s", col.Header)
	}
	fmt.Fprintf(writer, "\n")
	return nil
}

// renderTableRow renders a single row of the table.
func (r *TableRenderer) renderTableRow(ctx context.Context, writer io.Writer, cols []*columnLayout, prgs []cel.Program,
	object proto.Message, objectHelper *reflection.ObjectHelper, lookupCache map[protoreflect.FullName]map[string]string) error {
	// Wrap the object in a top-level "this" field to avoid conflicts with reserved words:
	in := map[string]any{
		"this": object,
	}
	celVars, err := cel.PartialVars(in)
	if err != nil {
		return fmt.Errorf(
			"failed to set variables for CEL expression for type %q: %w",
			objectHelper, err,
		)
	}

	// Render each column:
	for i := range len(cols) {
		if i > 0 {
			fmt.Fprintf(writer, "\t")
		}
		col := cols[i]
		prg := prgs[i]

		// Evaluate the CEL expression:
		var out ref.Val
		out, _, err = prg.Eval(celVars)
		if err != nil {
			return fmt.Errorf(
				"failed to evaluate CEL expression %q for column %q of type %q: %w",
				col.Value, col.Header, objectHelper, err,
			)
		}

		// Render the cell value:
		err = r.renderTableCell(ctx, writer, col, out, lookupCache)
		if err != nil {
			return fmt.Errorf(
				"failed to render value %q for column %q of type %q: %w",
				out, col.Header, objectHelper, err,
			)
		}
	}
	fmt.Fprintf(writer, "\n")
	return nil
}

// renderTableCell renders a single cell in the table.
func (r *TableRenderer) renderTableCell(ctx context.Context, writer io.Writer, col *columnLayout, val ref.Val,
	lookupCache map[protoreflect.FullName]map[string]string) error {
	switch val := val.(type) {
	case types.Int:
		if col.Type != "" {
			enumType, _ := protoregistry.GlobalTypes.FindEnumByName(col.Type)
			if enumType != nil {
				return r.renderTableCellEnumType(writer, val, enumType.Descriptor())
			}
			r.logger.Error(
				"Failed to find enum type",
				slog.String("type", string(col.Type)),
			)
		}
	case types.String:
		if col.Lookup && col.Type != "" {
			messageType, _ := protoregistry.GlobalTypes.FindMessageByName(col.Type)
			if messageType != nil {
				return r.renderTableCellLookup(ctx, writer, val, messageType.Descriptor(), lookupCache)
			}
		}
	}
	return r.renderTableCellAnyType(writer, val)
}

// renderTableCellEnumType renders an enum value as a string.
func (r *TableRenderer) renderTableCellEnumType(writer io.Writer, val types.Int,
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
	// name of the type, for example `CLUSTER_STATE_PENDING`. That prefix is not useful for humans, so we try
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

// renderTableCellLookup renders a lookup value (identifier to name translation).
func (r *TableRenderer) renderTableCellLookup(ctx context.Context, writer io.Writer, val types.String,
	messageDesc protoreflect.MessageDescriptor, lookupCache map[protoreflect.FullName]map[string]string) error {
	key := string(val)
	var text string
	if key != "" {
		text = r.lookupName(ctx, messageDesc.FullName(), key, lookupCache)
	} else {
		text = "-"
	}
	_, err := fmt.Fprintf(writer, "%s", text)
	return err
}

// lookupName looks up a name from an identifier.
func (r *TableRenderer) lookupName(ctx context.Context, messageFullName protoreflect.FullName,
	key string, lookupCache map[protoreflect.FullName]map[string]string) (result string) {
	// Check if the result is already in the cache and return it immediately if so, otherwise
	// remember to update the cache when done:
	cache, ok := lookupCache[messageFullName]
	if !ok {
		cache = map[string]string{}
		lookupCache[messageFullName] = cache
	}
	result, ok = cache[key]
	if ok {
		return result
	}
	defer func() {
		cache[key] = result
	}()

	// Find the object helper:
	objectHelper := r.helper.Lookup(string(messageFullName))
	if objectHelper == nil {
		r.logger.ErrorContext(
			ctx,
			"Failed to find object helper for type",
			slog.String("type", string(messageFullName)),
		)
		result = key
		return
	}

	// Find the objects whose identifier or name matches the key:
	filter := fmt.Sprintf(
		"this.id == %q || this.metadata.name == %q",
		key, key,
	)
	listResult, err := objectHelper.List(ctx, reflection.ListOptions{
		Filter: filter,
	})
	if err != nil {
		r.logger.ErrorContext(
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
	result = metadata.GetName()
	return
}

// renderTableCellAnyType renders any value type as a string.
func (r *TableRenderer) renderTableCellAnyType(writer io.Writer, val ref.Val) error {
	_, err := fmt.Fprintf(writer, "%s", val)
	return err
}

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
	"google.golang.org/protobuf/reflect/protoreflect"
)

// Table describes how to render protocol buffers messages in tabular form.
type Table struct {
	// Columns describes how fields of the message are mapped to columns.
	Columns []*Column `yaml:"columns,omitempty"`
}

// Columns describes how to render a field of a protocol buffers message as a column in a table.
type Column struct {
	// Header is the text of the header for the colum. The default is to use the name of the field in upper case
	// and replacing underscores with spaces.
	Header string `yaml:"header,omitempty"`

	// Value is a CEL expression that will be used to calculate the rendered value. The expression can access
	// the message via the `this` built-in variable.
	Value string `yaml:"value,omitempty"`

	// Type is the name of the type of the result of the expression. Thi this only needed when the result of the
	// expression is an enum value or an identifier that needs to be translated into a type.
	//
	// When the result is a enum value, then the 'type' field should contain the name of the enum type, and it
	// will be used to translate the integer value into th ename of the enum value shortened to elimiate the
	// prefix common to all the enum values of tat type.
	//
	// When the result is an identifier the 'type' field should be the name of the type, and it will be used to
	// find the name of the object.
	Type protoreflect.FullName `yaml:"type,omitempty"`

	// Lookup indicates if the result of the expression is an identifier that needs to be translated into a name.
	// When this is set to true the 'type' field also needs to be specified, and should contain the name of the
	// type to use for the looup. For example, is the result of the expression is a cluster, then the 'type'
	// should be 'fulfillment.v1.Cluster'.
	Lookup bool `yaml:"lookup,omitempty"`
}

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
	// Domains is a list of domains that are available in CEL expressions.
	Domains []*Domain `yaml:"domains,omitempty"`

	// Columns describes how fields of the message are mapped to columns.
	Columns []*Column `yaml:"columns,omitempty"`
}

// Domain describes a domain that is available in CEL expressions. The name is the name of the domain that will be
// used as the first argument to the 'lookup' function. The type is the type of objects that are membmers of the
// domain. For example, if the domain is 'templates' and the type is 'private.v1.ClusterTemplate' then the 'lookup'
// function will be called with the first argument 'templates' and the second argument will be the name or identifier
// of the template and it will return the identifier.
type Domain struct {
	// Name is the name of the domain.
	Name string `yaml:"name,omitempty"`

	// Type is the type of objects that are members of the domain. This is used to decide which service to query
	// to get the results.
	Type protoreflect.FullName `yaml:"type,omitempty"`
}

// Columns describes how to render a field of a protocol buffers message as a column in a table.
type Column struct {
	// Header is the text of the header for the colum. The default is to use the name of the field in upper case
	// and replacing underscores with spaces.
	Header string `yaml:"header,omitempty"`

	// Value is a CEL expression that will be used to calculate the rendered value. The expression can access
	// the message via the `this` built-in variable.
	Value string `yaml:"value,omitempty"`

	// Type is the name of a enum type that is the result of evaluationg the expression. This is needed because
	// CEL doesn't have a notion of enum types: they are all translated to integers. When this is specified the
	// result of the CEL expression will be then translated into the name of the enum value.
	Type protoreflect.FullName `yaml:"type,omitempty"`
}

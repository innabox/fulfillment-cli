/*
Copyright (c) 2025 Red Hat Inc.

Licensed under the Apache License, Version 2.0 (the "License"); you may not use this file except in compliance with the
License. You may obtain a copy of the License at

  http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software distributed under the License is distributed on an
"AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the License for the specific
language governing permissions and limitations under the License.
*/

package terminal

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	iofs "io/fs"
	"log/slog"
	"os"
	"strings"

	"github.com/alecthomas/chroma/v2/formatters"
	"github.com/alecthomas/chroma/v2/lexers"
	"github.com/alecthomas/chroma/v2/styles"
	"github.com/mattn/go-colorable"
	"github.com/mattn/go-isatty"
	"gopkg.in/yaml.v3"

	"github.com/innabox/fulfillment-common/templating"
)

// ConsoleBuilder contains the data and logic needed to create a console. Don't create objects of this type directly,
// use the NewConsole function instead.
type ConsoleBuilder struct {
	logger *slog.Logger
	file   *os.File
}

// Console is helps writing messages to the console. Don't create objects of this type directly, use the NewConsole
// function instead.
type Console struct {
	logger *slog.Logger
	file   *os.File
	engine *templating.Engine
}

// NewConsole creates a builder that can the be used to create a template engine.
func NewConsole() *ConsoleBuilder {
	return &ConsoleBuilder{}
}

// SetLogger sets the logger that the console will use to write messages to the log. This is mandatory.
func (b *ConsoleBuilder) SetLogger(value *slog.Logger) *ConsoleBuilder {
	b.logger = value
	return b
}

// SetFile sets the file that the console will use to write messages to the console. This is optional, the default
// is to use os.Stdout and there is usually no need to change it; it is intended for unit tests.
func (b *ConsoleBuilder) SetFile(value *os.File) *ConsoleBuilder {
	b.file = value
	return b
}

// Build uses the configuration stored in the builder to create a new console.
func (b *ConsoleBuilder) Build() (result *Console, err error) {
	// Check parameters:
	if b.logger == nil {
		err = errors.New("logger is mandatory")
		return
	}

	// Create the template engine:
	engine, err := templating.NewEngine().
		SetLogger(b.logger).
		Build()
	if err != nil {
		err = fmt.Errorf("failed to create template engine: %w", err)
		return
	}

	// Set the default writer if needed:
	file := b.file
	if file == nil {
		file = os.Stdout
	}

	// Create and populate the object:
	result = &Console{
		logger: b.logger,
		file:   file,
		engine: engine,
	}
	return
}

// AddTemplates adds one temlate file system containing templates, including only the templates that are in the given
// directory.
func (c *Console) AddTemplates(fs iofs.FS, dir string) error {
	sub, err := iofs.Sub(fs, dir)
	if err != nil {
		return fmt.Errorf("failed to get templates sub directory '%s': %w", dir, err)
	}
	return c.engine.AddFS(sub)
}

func (c *Console) Printf(ctx context.Context, format string, args ...any) {
	text := fmt.Sprintf(format, args...)
	c.logger.DebugContext(
		ctx,
		"Console printf",
		slog.String("format", format),
		slog.Any("args", args),
		slog.Any("text", text),
	)
	_, err := c.file.WriteString(text)
	if err != nil {
		c.logger.ErrorContext(
			ctx,
			"Failed to write text",
			slog.String("text", text),
			slog.Any("error", err),
		)
	}
}

// Render renders the given template with the given data to stdout. The template should be a template file name that
// was added via AddTemplatesFS. If no template file systems have been added, this method will log an error.
func (c *Console) Render(ctx context.Context, template string, data any) {
	buffer := &bytes.Buffer{}
	err := c.engine.Execute(buffer, template, data)
	if err != nil {
		c.logger.ErrorContext(
			ctx,
			"Failed to execute template",
			slog.String("template", template),
			slog.Any("error", err),
		)
		return
	}
	text := buffer.String()
	lines := strings.Split(text, "\n")
	previousEmpty := true
	for _, line := range lines {
		currentEmpty := len(line) == 0
		if currentEmpty {
			if !previousEmpty {
				fmt.Fprintf(os.Stdout, "\n")
				previousEmpty = true
			}
		} else {
			fmt.Fprintf(os.Stdout, "%s\n", line)
			previousEmpty = false
		}
	}
}

// RenderJson renders the given data as JSON to stdout. If the terminal supports color, the output will be colorized
// using the chroma syntax highlighter.
func (c *Console) RenderJson(ctx context.Context, data any) {
	bytes, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		c.logger.ErrorContext(
			ctx,
			"Failed to encode JSON",
			slog.Any("error", err),
		)
		return
	}
	text := string(bytes) + "\n"
	c.renderColored(ctx, text, "json")
}

// RenderYaml renders the given data as YAML to stdout. If the terminal supports color, the output will be colorized
// using the chroma syntax highlighter.
func (c *Console) RenderYaml(ctx context.Context, data any) {
	buffer := &bytes.Buffer{}
	encoder := yaml.NewEncoder(buffer)
	encoder.SetIndent(2)
	err := encoder.Encode(data)
	if err != nil {
		c.logger.ErrorContext(
			ctx,
			"Failed to encode YAML",
			slog.Any("error", err),
		)
		return
	}
	encoder.Close()
	c.renderColored(ctx, buffer.String(), "yaml")
}

// renderColored renders the given text to stdout with syntax highlighting using the specified lexer. If the terminal
// doesn't support color or an error occurs, it falls back to plain text output.
func (c *Console) renderColored(ctx context.Context, text string, format string) error {
	if isatty.IsTerminal(c.file.Fd()) {
		lexer := lexers.Get(format)
		if lexer == nil {
			lexer = lexers.Fallback
		}
		style := styles.Get("friendly")
		if style == nil {
			style = styles.Fallback
		}
		formatter := formatters.Get("terminal256")
		if formatter == nil {
			formatter = formatters.Fallback
		}
		iterator, err := lexer.Tokenise(nil, text)
		if err != nil {
			c.logger.ErrorContext(
				ctx,
				"Failed to tokenize text",
				slog.String("format", format),
				slog.Any("error", err),
			)
			fmt.Fprint(c.file, text)
			return nil
		}
		err = formatter.Format(colorable.NewColorable(c.file), style, iterator)
		if err != nil {
			c.logger.ErrorContext(
				ctx,
				"Failed to format text",
				slog.String("format", format),
				slog.Any("error", err),
			)
			fmt.Fprint(c.file, text)
			return nil
		}
		return nil
	}
	fmt.Fprint(c.file, text)
	return nil
}

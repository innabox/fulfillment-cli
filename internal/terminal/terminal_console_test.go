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
	. "github.com/onsi/ginkgo/v2/dsl/core"
	. "github.com/onsi/gomega"
)

var _ = Describe("Console", func() {
	Describe("Creation", func() {
		It("Can be created with all the default parameters", func() {
			console, err := NewConsole().
				SetLogger(logger).
				Build()
			Expect(err).ToNot(HaveOccurred())
			Expect(console).ToNot(BeNil())
		})

		It("Can't be created without a logger", func() {
			console, err := NewConsole().
				Build()
			Expect(err).To(MatchError("logger is mandatory"))
			Expect(console).To(BeNil())
		})
	})

	Describe("RenderYaml", func() {
		It("Can render a simple map as YAML", func() {
			console, err := NewConsole().
				SetLogger(logger).
				Build()
			Expect(err).ToNot(HaveOccurred())

			data := map[string]any{
				"name":  "test",
				"value": 123,
			}
			console.RenderYaml(ctx, data)
		})

		It("Can render a slice as YAML", func() {
			console, err := NewConsole().
				SetLogger(logger).
				Build()
			Expect(err).ToNot(HaveOccurred())

			data := []map[string]any{
				{"id": "1", "name": "first"},
				{"id": "2", "name": "second"},
			}
			console.RenderYaml(ctx, data)
		})
	})

	Describe("RenderJson", func() {
		It("Can render a simple map as JSON", func() {
			console, err := NewConsole().
				SetLogger(logger).
				Build()
			Expect(err).ToNot(HaveOccurred())

			data := map[string]any{
				"name":  "test",
				"value": 123,
			}
			console.RenderJson(ctx, data)
		})

		It("Can render a slice as JSON", func() {
			console, err := NewConsole().
				SetLogger(logger).
				Build()
			Expect(err).ToNot(HaveOccurred())

			data := []map[string]any{
				{"id": "1", "name": "first"},
				{"id": "2", "name": "second"},
			}
			console.RenderJson(ctx, data)
		})
	})
})

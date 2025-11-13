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

	"github.com/google/cel-go/cel"
	"github.com/google/cel-go/common/types"
	. "github.com/onsi/ginkgo/v2/dsl/core"
	. "github.com/onsi/gomega"
)

var _ = Describe("Lookup function", func() {
	var runner *runnerContext

	BeforeEach(func() {
		runner = &runnerContext{}
		runner.ctx = context.Background()
		runner.logger = logger
		runner.lookupDomains = make(map[string]*Domain)
		runner.lookupCache = make(map[string]map[string]string)
	})

	It("should compile lookup expression successfully", func() {
		// Create a CEL environment with the lookup function:
		env, err := cel.NewEnv(
			runner.createLookupFunc(),
		)
		Expect(err).ToNot(HaveOccurred())

		// Verify that lookup expressions compile:
		ast, issues := env.Compile(`lookup("templates", "my-template")`)
		Expect(issues.Err()).ToNot(HaveOccurred())
		Expect(ast).ToNot(BeNil())
	})

	It("should return key for unknown domain", func() {
		// Create a CEL environment with the lookup function:
		env, err := cel.NewEnv(
			runner.createLookupFunc(),
		)
		Expect(err).ToNot(HaveOccurred())

		// Try to lookup in a non-existent domain:
		ast, issues := env.Compile(`lookup("unknown", "1")`)
		Expect(issues.Err()).ToNot(HaveOccurred())

		prg, err := env.Program(ast)
		Expect(err).ToNot(HaveOccurred())

		result, _, err := prg.Eval(map[string]any{})
		Expect(err).ToNot(HaveOccurred())

		// The result should be the key itself:
		strVal, ok := result.(types.String)
		Expect(ok).To(BeTrue())
		Expect(string(strVal)).To(Equal("1"))
	})

	It("should support conditional access on lookup results", func() {
		// Create a CEL environment with the lookup function:
		env, err := cel.NewEnv(
			runner.createLookupFunc(),
		)
		Expect(err).ToNot(HaveOccurred())

		// Verify that conditional expressions with lookup compile:
		// When lookup doesn't find anything, it returns the key itself, so we can check if result != key
		ast, issues := env.Compile(`lookup("templates", "id") != "id" ? "found" : "not found"`)
		Expect(issues.Err()).ToNot(HaveOccurred())
		Expect(ast).ToNot(BeNil())

		// Execute the expression - should return "not found" since there's no actual domain:
		prg, err := env.Program(ast)
		Expect(err).ToNot(HaveOccurred())

		result, _, err := prg.Eval(map[string]any{})
		Expect(err).ToNot(HaveOccurred())

		strVal, ok := result.(types.String)
		Expect(ok).To(BeTrue())
		Expect(string(strVal)).To(Equal("not found"))
	})

	It("should cache lookup results", func() {
		// Create a CEL environment with the lookup function:
		env, err := cel.NewEnv(
			runner.createLookupFunc(),
		)
		Expect(err).ToNot(HaveOccurred())

		// Compile the lookup expression:
		ast, issues := env.Compile(`lookup("unknown", "test-key")`)
		Expect(issues.Err()).ToNot(HaveOccurred())

		prg, err := env.Program(ast)
		Expect(err).ToNot(HaveOccurred())

		// First call - should populate the cache:
		result1, _, err := prg.Eval(map[string]any{})
		Expect(err).ToNot(HaveOccurred())
		strVal1, ok := result1.(types.String)
		Expect(ok).To(BeTrue())
		Expect(string(strVal1)).To(Equal("test-key"))

		// Verify the cache was populated:
		domainCache, foundDomain := runner.lookupCache["unknown"]
		Expect(foundDomain).To(BeTrue())
		cachedValue, foundKey := domainCache["test-key"]
		Expect(foundKey).To(BeTrue())
		Expect(cachedValue).To(Equal("test-key"))

		// Second call - should use the cached value:
		result2, _, err := prg.Eval(map[string]any{})
		Expect(err).ToNot(HaveOccurred())
		strVal2, ok := result2.(types.String)
		Expect(ok).To(BeTrue())
		Expect(string(strVal2)).To(Equal("test-key"))

		// Verify the cache still has the same value:
		domainCache2, foundDomain2 := runner.lookupCache["unknown"]
		Expect(foundDomain2).To(BeTrue())
		cachedValue2, foundKey2 := domainCache2["test-key"]
		Expect(foundKey2).To(BeTrue())
		Expect(cachedValue2).To(Equal("test-key"))
	})
})

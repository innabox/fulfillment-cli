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
	"context"

	. "github.com/onsi/ginkgo/v2/dsl/core"
	. "github.com/onsi/gomega"
)

var _ = Describe("Context", func() {
	It("Extracts console from the context if previously added", func() {
		console, err := NewConsole().
			SetLogger(logger).
			Build()
		Expect(err).ToNot(HaveOccurred())
		ctx := ConsoleIntoContext(context.Background(), console)
		extracted := ConsoleFromContext(ctx)
		Expect(extracted).To(BeIdenticalTo(console))
	})

	It("Panics if console wasn't added to the context", func() {
		ctx := context.Background()
		Expect(func() {
			ConsoleFromContext(ctx)
		}).To(Panic())
	})
})

/*
Copyright (c) 2025 Red Hat Inc.

Licensed under the Apache License, Version 2.0 (the "License"); you may not use this file except in compliance with the
License. You may obtain a copy of the License at

  http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software distributed under the License is distributed on an
"AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the License for the specific
language governing permissions and limitations under the License.
*/

package auth

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/innabox/fulfillment-cli/internal/testing"
	. "github.com/onsi/ginkgo/v2/dsl/core"
	. "github.com/onsi/gomega"
)

var _ = Describe("Script token source", func() {
	nopTokenLoadFunc := func() (string, error) {
		return "", nil
	}

	nopTokenSaveFunc := func(string) error {
		return nil
	}

	Describe("Creation", func() {
		It("Can be created with all the mandatory parameters", func() {
			source, err := NewScriptTokenSource().
				SetLogger(logger).
				SetScript("echo mytoken").
				SetTokenLoadFunc(nopTokenLoadFunc).
				SetTokenSaveFunc(nopTokenSaveFunc).
				Build()
			Expect(err).ToNot(HaveOccurred())
			Expect(source).ToNot(BeNil())
		})

		It("Can't be created without a logger", func() {
			source, err := NewScriptTokenSource().
				SetScript("echo mytoken").
				SetTokenLoadFunc(nopTokenLoadFunc).
				SetTokenSaveFunc(nopTokenSaveFunc).
				Build()
			Expect(err).To(MatchError("logger is mandatory"))
			Expect(source).To(BeNil())
		})

		It("Can't be created without a script", func() {
			source, err := NewScriptTokenSource().
				SetLogger(logger).
				SetTokenLoadFunc(nopTokenLoadFunc).
				SetTokenSaveFunc(nopTokenSaveFunc).
				Build()
			Expect(err).To(MatchError("token generation script is mandatory"))
			Expect(source).To(BeNil())
		})

		It("Can't be created without a token load function", func() {
			source, err := NewScriptTokenSource().
				SetLogger(logger).
				SetScript("echo mytoken").
				SetTokenSaveFunc(nopTokenSaveFunc).
				Build()
			Expect(err).To(MatchError("token load function is mandatory"))
			Expect(source).To(BeNil())
		})

		It("Can't be created without a token save function", func() {
			source, err := NewScriptTokenSource().
				SetLogger(logger).
				SetScript("echo mytoken").
				SetTokenLoadFunc(nopTokenLoadFunc).
				Build()
			Expect(err).To(MatchError("token save function is mandatory"))
			Expect(source).To(BeNil())
		})
	})

	Describe("Behaviour", func() {
		It("Returns the token generated by the script", func() {
			source, err := NewScriptTokenSource().
				SetLogger(logger).
				SetScript("echo mytoken").
				SetTokenLoadFunc(nopTokenLoadFunc).
				SetTokenSaveFunc(nopTokenSaveFunc).
				Build()
			Expect(err).ToNot(HaveOccurred())
			token, err := source.Token()
			Expect(err).ToNot(HaveOccurred())
			Expect(token).ToNot(BeNil())
			Expect(token.AccessToken).To(Equal("mytoken"))
		})

		It("Generates a new token for each call if the token isn't a JWT", func() {
			// For this test we will use a script that reads a token from a file, so that we can change the
			// token between calls, so we need to create that file.
			tmpDir, err := os.MkdirTemp("", "*.test")
			Expect(err).ToNot(HaveOccurred())
			defer func() {
				err := os.RemoveAll(tmpDir)
				Expect(err).ToNot(HaveOccurred())
			}()
			tmpFile := filepath.Join(tmpDir, "token.txt")

			// Create the source with the script that reads from a file in the temporary directory:
			script := fmt.Sprintf("cat %s", tmpFile)
			source, err := NewScriptTokenSource().
				SetLogger(logger).
				SetScript(script).
				SetTokenLoadFunc(nopTokenLoadFunc).
				SetTokenSaveFunc(nopTokenSaveFunc).
				Build()
			Expect(err).ToNot(HaveOccurred())

			// Write the first token and verify that it is returned:
			err = os.WriteFile(tmpFile, []byte("first"), 0600)
			Expect(err).ToNot(HaveOccurred())
			token, err := source.Token()
			Expect(err).ToNot(HaveOccurred())
			Expect(token).ToNot(BeNil())
			Expect(token.AccessToken).To(Equal("first"))

			// Write the second token and verify that it is returned:
			err = os.WriteFile(tmpFile, []byte("second"), 0600)
			Expect(err).ToNot(HaveOccurred())
			token, err = source.Token()
			Expect(err).ToNot(HaveOccurred())
			Expect(token).ToNot(BeNil())
			Expect(token.AccessToken).To(Equal("second"))
		})

		It("Doesn't use the loaded token if it isn't a JWT", func() {
			source, err := NewScriptTokenSource().
				SetLogger(logger).
				SetScript("echo mytoken").
				SetTokenLoadFunc(func() (token string, err error) {
					token = "myloaded"
					return
				}).
				SetTokenSaveFunc(nopTokenSaveFunc).
				Build()
			Expect(err).ToNot(HaveOccurred())
			token, err := source.Token()
			Expect(err).ToNot(HaveOccurred())
			Expect(token).ToNot(BeNil())
			Expect(token.AccessToken).To(Equal("mytoken"))
		})

		It("Doesn't save the token if it isn't a JWT", func() {
			var saved string
			source, err := NewScriptTokenSource().
				SetLogger(logger).
				SetScript("echo mytoken").
				SetTokenLoadFunc(nopTokenLoadFunc).
				SetTokenSaveFunc(func(token string) error {
					saved = token
					return nil
				}).
				Build()
			Expect(err).ToNot(HaveOccurred())
			token, err := source.Token()
			Expect(err).ToNot(HaveOccurred())
			Expect(token).ToNot(BeNil())
			Expect(token.AccessToken).To(Equal("mytoken"))
			Expect(saved).To(BeEmpty())
		})

		It("Returns the loaded token if it is a JWT that hasn't expired", func() {
			loaded := testing.MakeTokenString("Bearer", 5*time.Minute)
			source, err := NewScriptTokenSource().
				SetLogger(logger).
				SetScript("echo mytoken").
				SetTokenLoadFunc(func() (token string, err error) {
					token = loaded
					return
				}).
				SetTokenSaveFunc(nopTokenSaveFunc).
				Build()
			Expect(err).ToNot(HaveOccurred())
			token, err := source.Token()
			Expect(err).ToNot(HaveOccurred())
			Expect(token).ToNot(BeNil())
			Expect(token.AccessToken).To(Equal(loaded))
		})

		It("Saves the generated token if it is a JWT that hasn't expired", func() {
			generated := testing.MakeTokenString("Bearer", 5*time.Minute)
			script := fmt.Sprintf("echo '%s'", generated)
			var saved string
			source, err := NewScriptTokenSource().
				SetLogger(logger).
				SetScript(script).
				SetTokenLoadFunc(nopTokenLoadFunc).
				SetTokenSaveFunc(func(token string) error {
					saved = token
					return nil
				}).
				Build()
			Expect(err).ToNot(HaveOccurred())
			token, err := source.Token()
			Expect(err).ToNot(HaveOccurred())
			Expect(token).ToNot(BeNil())
			Expect(token.AccessToken).To(Equal(generated))
			Expect(saved).To(Equal(generated))
		})

		It("Generates a new token if the current one is a JWT that has expired", func() {
			loaded := testing.MakeTokenString("Bearer", -5*time.Minute)
			generated := testing.MakeTokenString("Bearer", 5*time.Minute)
			script := fmt.Sprintf("echo '%s'", generated)
			source, err := NewScriptTokenSource().
				SetLogger(logger).
				SetScript(script).
				SetTokenLoadFunc(func() (token string, err error) {
					token = loaded
					return
				}).
				SetTokenSaveFunc(nopTokenSaveFunc).
				Build()
			Expect(err).ToNot(HaveOccurred())
			token, err := source.Token()
			Expect(err).ToNot(HaveOccurred())
			Expect(token).ToNot(BeNil())
			Expect(token.AccessToken).To(Equal(generated))
		})

		It("Returns an error if the token load function fails", func() {
			source, err := NewScriptTokenSource().
				SetLogger(logger).
				SetScript("echo mytoken").
				SetTokenLoadFunc(func() (token string, err error) {
					err = errors.New("my load error")
					return
				}).
				SetTokenSaveFunc(nopTokenSaveFunc).
				Build()
			Expect(err).ToNot(HaveOccurred())
			token, err := source.Token()
			Expect(err).To(MatchError("my load error"))
			Expect(token).To(BeNil())
		})

		It("Returns an error if the token save function fails", func() {
			generated := testing.MakeTokenString("Bearer", 5*time.Minute)
			script := fmt.Sprintf("echo '%s'", generated)
			source, err := NewScriptTokenSource().
				SetLogger(logger).
				SetScript(script).
				SetTokenLoadFunc(nopTokenLoadFunc).
				SetTokenSaveFunc(func(s string) error {
					return errors.New("my save error")
				}).
				Build()
			Expect(err).ToNot(HaveOccurred())
			token, err := source.Token()
			Expect(err).To(MatchError("my save error"))
			Expect(token).To(BeNil())
		})

		It("Returns an error if the script exists with non zero code", func() {
			source, err := NewScriptTokenSource().
				SetLogger(logger).
				SetScript("exit 1").
				SetTokenLoadFunc(nopTokenLoadFunc).
				SetTokenSaveFunc(nopTokenSaveFunc).
				Build()
			Expect(err).ToNot(HaveOccurred())
			token, err := source.Token()
			Expect(err).To(MatchError("failed to execute token generation script 'exit 1': exit status 1"))
			Expect(token).To(BeNil())
		})
	})
})

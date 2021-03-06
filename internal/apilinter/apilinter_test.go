// Copyright 2021 Jon Bennett
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package apilinter

import (
	"io/ioutil"
	"sort"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/descriptorpb"
	"google.golang.org/protobuf/types/pluginpb"
)

func TestNewLinter(t *testing.T) {
	tests := map[string]struct {
		opts LinterOptions
		err  string
	}{
		"Should raise missing file error": {
			opts: LinterOptions{
				ConfigPath: "config.json",
			},
			err: "no such file or directory",
		},
		"Should raise invalid config format error": {
			opts: LinterOptions{
				ConfigPath: "config.text",
			},
			err: "unsupported format",
		},
		"Should load valid config": {
			opts: LinterOptions{
				ConfigPath: "./testdata/config.valid.yaml",
			},
			err: "",
		},
	}
	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			_, err := NewLinter(test.opts)
			switch {
			case err != nil && test.err != "" && strings.Contains(err.Error(), test.err):
				return
			case err != nil:
				t.Fatal(err)
			}
		})
	}
}

func TestFileLinter_LintFiles(t *testing.T) {
	tests := map[string]struct {
		protoFile                   string
		enabledRules, disabledRules []string
		out                         map[string][]string
	}{
		"lint with defaults": {
			out: map[string][]string{
				"service.proto": {
					"core::0131::request-unknown-fields",
					"core::0131::request-name-required",
				},
				"service_ok.proto": {
					"core::0131::request-name-behavior",
					"core::0131::request-name-reference",
					"core::0192::has-comments",
					"core::0192::has-comments",
				},
			},
		},
		"lint with some rules disabled": {
			disabledRules: []string{
				"core::0131::request-name-required",
				"core::0131::request-name-behavior",
				"core::0131::request-name-reference",
			},
			out: map[string][]string{
				"service.proto": {
					"core::0131::request-unknown-fields",
				},
				"service_ok.proto": {
					"core::0192::has-comments",
					"core::0192::has-comments",
				},
			},
		},
		"lint with all expected rules disabled": {
			disabledRules: []string{
				"core::0131::request-unknown-fields",
				"core::0131::request-name-required",
				"core::0131::request-name-behavior",
				"core::0131::request-name-reference",
				"core::0192::has-comments",
			},
			out: map[string][]string{
				"service.proto":    {},
				"service_ok.proto": {},
			},
		},
	}
	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			// read our compiled test proto
			const protoBinFile = "./testdata/protos.bin"
			b, err := ioutil.ReadFile(protoBinFile)
			if err != nil {
				t.Fatalf("proto image (bin) file not found: %s, err: %s", protoBinFile, err)
			}
			// unmarshal to the FDS
			var f descriptorpb.FileDescriptorSet
			if err := proto.Unmarshal(b, &f); err != nil {
				t.Fatalf("proto.Unmarshal: %s", err)
			}

			gen, err := protogen.Options{}.New(&pluginpb.CodeGeneratorRequest{
				ProtoFile: f.GetFile(),
			})
			if err != nil {
				t.Fatalf("protogen.New: %s", err)
			}

			opts := LinterOptions{
				EnabledRules:  test.enabledRules,
				DisabledRules: test.disabledRules,
			}
			fl, err := NewLinter(opts)
			if err != nil {
				t.Fatal(err)
			}

			out, err := fl.LintFiles(gen.Files)
			if err != nil {
				t.Fatalf("fl.LintFile, test.protoFile=%s %s", test.protoFile, err)
			}

			// check out response problems per file
			for _, resp := range out {
				// get our test rules
				tRules, exists := test.out[resp.FilePath]
				if !exists {
					t.Fatalf("test cases missing for file=%s", resp.FilePath)
				}
				sort.Strings(tRules)

				rules := make([]string, 0)
				for _, problem := range resp.Problems {
					rules = append(rules, string(problem.RuleID))
				}
				sort.Strings(rules)

				if diff := cmp.Diff(tRules, rules); diff != "" {
					t.Errorf("LintFiles - expected problems mismatch (-want +have):\n%s", diff)
				}

			}
		})
	}
}

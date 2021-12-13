// Copyright 2020-2021 Jon Bennett
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
	"fmt"

	"github.com/golang/protobuf/protoc-gen-go/descriptor"
	"github.com/googleapis/api-linter/lint"
	"github.com/googleapis/api-linter/rules"
	"github.com/jhump/protoreflect/desc"
	"google.golang.org/protobuf/compiler/protogen"
)

type LinterOptions struct {
	EnabledRules  []string
	DisabledRules []string
}

type Linter struct {
	linter *lint.Linter
}

func NewLinter(opts LinterOptions) (*Linter, error) {
	// set up default rules
	lintRules := lint.NewRuleRegistry()
	if err := rules.Add(lintRules); err != nil {
		return nil, fmt.Errorf("rules.Add: %w", err)
	}
	// configure linter
	lintConfigs := lint.Configs{}
	if len(opts.EnabledRules) > 0 {
		lintConfigs = append(lintConfigs, lint.Config{
			EnabledRules: opts.EnabledRules,
		})
	}
	if len(opts.DisabledRules) > 0 {
		lintConfigs = append(lintConfigs, lint.Config{
			DisabledRules: opts.DisabledRules,
		})
	}

	linter := lint.New(lintRules, lintConfigs)
	return &Linter{linter: linter}, nil
}

func (fl *Linter) LintFiles(files []*protogen.File) ([]lint.Response, error) {
	// convert protogen.File's to desc.FileDescriptor's
	var protos []*descriptor.FileDescriptorProto
	for _, f := range files {
		protos = append(protos, f.Proto)
	}
	fileDescByName, err := desc.CreateFileDescriptors(protos)
	if err != nil {
		return nil, fmt.Errorf("desc.CreateFileDescriptors: %w", err)
	}

	// convert map into files
	var fdList []*desc.FileDescriptor
	for _, fd := range fileDescByName {
		fdList = append(fdList, fd)
	}

	// Lint the proto file
	resp, err := fl.linter.LintProtos(fdList...)
	if err != nil {
		return nil, fmt.Errorf("l.LintProtos: %w", err)
	}

	return resp, nil
}

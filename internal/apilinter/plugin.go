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
	"encoding/json"
	"fmt"

	"github.com/jbendotnet/protoc-gen-api-linter/internal/jsonutil"
	"google.golang.org/protobuf/compiler/protogen"
)

type PluginOptions struct {
	Linter            LinterOptions
	ReportFilename    string
	ReportPrettyPrint bool
	ExitOnError       bool
}

type Plugin struct {
	linter *Linter
	opts   PluginOptions
}

func NewPlugin(opts PluginOptions) (*Plugin, error) {
	fl, err := NewLinter(opts.Linter)
	if err != nil {
		return nil, fmt.Errorf("plugin.Run: %w", err)
	}
	return &Plugin{linter: fl, opts: opts}, nil
}

// Run lints protos and generates reports
// returns a boolean to indicate if problems were found
// only returns an error if it did indeed error
func (p *Plugin) Run(gen *protogen.Plugin) (bool, error) {
	res, err := p.linter.LintFiles(gen.Files)
	if err != nil {
		return false, fmt.Errorf("plugin.Run: %w", err)
	}
	if len(res) == 0 {
		return true, nil
	}

	reportJSON, err := getMarshaller(p.opts.ReportPrettyPrint)(res)
	if err != nil {
		return false, fmt.Errorf("plugin.Run: marshalJSON: %w", err)
	}

	g := gen.NewGeneratedFile(p.opts.ReportFilename, "")
	if _, err := g.Write(reportJSON); err != nil {
		return false, fmt.Errorf("plugin.Run: g.Write, file=%s: %w", p.opts.ReportFilename, err)
	}

	// if we had problems, we need to indicate that to callers
	lintFoundIssues := res == nil
	return lintFoundIssues, nil
}

func getMarshaller(pretty bool) func(v interface{}) ([]byte, error) {
	if pretty {
		return jsonutil.MarshalPretty
	}
	return json.Marshal
}

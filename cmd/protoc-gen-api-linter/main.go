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

package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"

	"github.com/jbendotnet/protoc-gen-api-linter/internal/apilinter"
	"github.com/sgreben/flagvar"
	"google.golang.org/protobuf/compiler/protogen"
)

var (
	versionFlag     bool
	ruleEnableFlag  flagvar.StringSet
	ruleDisableFlag flagvar.StringSet
)

var (
	version = "dev"
	commit  = "unknown"
	date    = "unknown"
)

func main() {

	flag.BoolVar(&versionFlag, "version", false, "Print version and exit.")
	flag.Var(&ruleEnableFlag, "enable-rule", "Enable a rule with the given name.\nMay be specified multiple times.")
	flag.Var(&ruleDisableFlag, "disable-rule", "Disable a rule with the given name.\nMay be specified multiple times.")
	flag.Parse()

	if versionFlag {
		fmt.Printf("Version %v, commit %v, built at %v\n", version, commit, date)
		os.Exit(0)
	}

	protogen.Options{}.Run(runPlugin)
}

func runPlugin(gen *protogen.Plugin) error {
	fl, err := apilinter.NewFileLinter()
	if err != nil {
		return fmt.Errorf("protogen.Plugin.Run: %w", err)
	}

	res, err := fl.LintFiles(gen.Files)
	if err != nil {
		return fmt.Errorf("protogen.Plugin.Run: %w", err)
	}

	reportJSON, err := json.MarshalIndent(res, "", "    ")
	if err != nil {
		return fmt.Errorf("protogen.Plugin.Run: json.Encode: %w", err)
	}

	const filename = "api_linter_report.json"
	g := gen.NewGeneratedFile(filename, "")
	if _, err := g.Write(reportJSON); err != nil {
		return fmt.Errorf("protogen.Plugin.Run: g.Write: %w", err)
	}

	return nil
}

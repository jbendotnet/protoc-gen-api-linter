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
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/jbendotnet/protoc-gen-api-linter/internal/apilinter"

	"github.com/sgreben/flagvar"
	"google.golang.org/protobuf/compiler/protogen"
)

const (
	appName               = "protoc-gen-api-linter"
	defaultReportFilename = "api_linter.json"
)

// Goreleaser vars
var (
	version = "dev"
	commit  = "unknown"
	date    = "unknown"
)

// cli config
type config struct {
	PrintVersion         bool
	LinterEnabledRules   flagvar.StringSet
	LinterDisabledRules  flagvar.StringSet
	PluginReportFilename string
	ReportPrettyPrint    bool
	ExitOnError          bool
}

var cfg config

func main() {
	// set logger to pint to stderr
	log.SetOutput(os.Stderr)
	// parse config options via CLI flags
	fs := flag.NewFlagSet(appName, flag.ExitOnError)
	fs.BoolVar(&cfg.PrintVersion, "version", false, "Print version and exit.")
	fs.Var(&cfg.LinterEnabledRules, "enable_rule", "Enable a rule with the given name.\nMay be specified multiple times.")
	fs.Var(&cfg.LinterDisabledRules, "disable_rule", "Disable a rule with the given name.\nMay be specified multiple times.")
	fs.StringVar(&cfg.PluginReportFilename, "report_filename", defaultReportFilename, "Set the filename of the JSON report")
	fs.BoolVar(&cfg.ReportPrettyPrint, "report_pretty_print", false, "Pretty print JSON reports")
	fs.BoolVar(&cfg.ExitOnError, "exit_on_error", true, "Exit on first error")
	if err := fs.Parse(os.Args[1:]); err != nil {
		fs.Usage()
		fmt.Printf("invalid args, err: %s", err)
		os.Exit(1)
	}

	// Show help
	const help = "help"
	if len(os.Args) == 2 && os.Args[1] == help {
		fs.Usage()
		os.Exit(0)
	}

	if cfg.PrintVersion {
		fmt.Printf("Version %v, commit %v, built at %v\n", version, commit, date)
		os.Exit(0)
	}

	protogen.Options{ParamFunc: fs.Set}.Run(runPlugin)
}

// runPlugin configures our plugin and runs it
func runPlugin(gen *protogen.Plugin) error {
	opts := apilinter.PluginOptions{
		Linter: apilinter.LinterOptions{
			EnabledRules:  cfg.LinterEnabledRules.Values(),
			DisabledRules: cfg.LinterDisabledRules.Values(),
		},
		ReportFilename:    cfg.PluginReportFilename,
		ReportPrettyPrint: cfg.ReportPrettyPrint,
	}

	plg, err := apilinter.NewPlugin(opts)
	if err != nil {
		return fmt.Errorf("runPlugin: %w", err)
	}

	ok, err := plg.Run(gen)
	if err != nil {
		return fmt.Errorf("runPlugin: %w", err)
	}
	if !ok && cfg.ExitOnError {
		return fmt.Errorf("linting problems found, check report: %s", cfg.PluginReportFilename)
	}

	return nil
}

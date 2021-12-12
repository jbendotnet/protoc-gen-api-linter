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

var (
	version = "dev"
	commit  = "unknown"
	date    = "unknown"
)

type config struct {
	PrintVersion         bool
	LinterEnableRule     flagvar.StringSet
	LinterDisableRule    flagvar.StringSet
	PluginReportFilename string
	ReportPrettyPrint    bool
}

var cfg config

func main() {
	// set logger to pint to stderr
	log.SetOutput(os.Stderr)
	fs := flag.NewFlagSet(appName, flag.ExitOnError)
	fs.BoolVar(&cfg.PrintVersion, "version", false, "Print version and exit.")
	fs.Var(&cfg.LinterEnableRule, "enable_rule", "Enable a rule with the given name.\nMay be specified multiple times.")
	fs.Var(&cfg.LinterDisableRule, "disable_rule", "Disable a rule with the given name.\nMay be specified multiple times.")
	fs.StringVar(&cfg.PluginReportFilename, "report_filename", defaultReportFilename, "Disable a rule with the given name.\nMay be specified multiple times.")
	fs.BoolVar(&cfg.ReportPrettyPrint, "report_pretty_print", false, "Pretty print JSON reports")
	if cfg.PrintVersion {
		fmt.Printf("Version %v, commit %v, built at %v\n", version, commit, date)
		os.Exit(0)
	}

	protogen.Options{ParamFunc: fs.Set}.Run(runPlugin)
}

// runPlugin configures our plugin instance and returns a
// well configured Run func
func runPlugin(gen *protogen.Plugin) error {
	opts := apilinter.PluginOptions{
		Linter: apilinter.LinterOptions{
			EnabledRules:  cfg.LinterEnableRule.Values(),
			DisabledRules: cfg.LinterDisableRule.Values(),
		},
		ReportFilename:    cfg.PluginReportFilename,
		ReportPrettyPrint: cfg.ReportPrettyPrint,
	}

	plg, err := apilinter.NewPlugin(opts)
	if err != nil {
		return fmt.Errorf("runPlugin: %w", err)
	}
	return plg.Run(gen)
}

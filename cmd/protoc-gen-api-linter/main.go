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
	versionFlag bool
	ruleEnableFlag flagvar.StringSet
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

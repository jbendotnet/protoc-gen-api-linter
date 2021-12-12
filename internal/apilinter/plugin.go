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

func (p *Plugin) Run(gen *protogen.Plugin) error {
	res, err := p.linter.LintFiles(gen.Files)
	if err != nil {
		return fmt.Errorf("plugin.Run: %w", err)
	}

	reportJSON, err := getMarshaller(p.opts.ReportPrettyPrint)(res)
	if err != nil {
		return fmt.Errorf("plugin.Run: marshalJSON: %w", err)
	}

	g := gen.NewGeneratedFile(p.opts.ReportFilename, "")
	if _, err := g.Write(reportJSON); err != nil {
		return fmt.Errorf("plugin.Run: g.Write, file=%s: %w", p.opts.ReportFilename, err)
	}

	return nil
}

func getMarshaller(pretty bool) func(v interface{}) ([]byte, error) {
	if pretty {
		return jsonutil.MarshalPretty
	}
	return json.Marshal
}

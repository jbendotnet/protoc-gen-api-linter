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
	switch {
	case err != nil:
		return nil, fmt.Errorf("l.LintProtos: %w", err)
	case resp == nil:
		return nil, nil
	default:
		return resp, nil
	}
}

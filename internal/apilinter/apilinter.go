package apilinter

import (
	"fmt"

	"github.com/golang/protobuf/protoc-gen-go/descriptor"
	"github.com/googleapis/api-linter/lint"
	"github.com/googleapis/api-linter/rules"
	"github.com/jhump/protoreflect/desc"
	"google.golang.org/protobuf/compiler/protogen"
)

type FileLinter struct {
	linter *lint.Linter
}

func NewFileLinter() (*FileLinter, error) {
	reg := lint.NewRuleRegistry()
	if err := rules.Add(reg); err != nil {
		return nil, fmt.Errorf("rules.Add: %w", err)
	}

	linter := lint.New(reg, lint.Configs{})
	return &FileLinter{linter: linter}, nil
}

func (fl *FileLinter) LintFiles(files []*protogen.File) ([]lint.Response, error) {
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

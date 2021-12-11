package apilinter

import (
	"io/ioutil"
	"sort"
	"testing"

	"github.com/google/go-cmp/cmp"
	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/descriptorpb"
	"google.golang.org/protobuf/types/pluginpb"
)

func TestFileLinter_LintFiles(t *testing.T) {
	tests := map[string]struct {
		protoFile string
		out       map[string][]string
	}{
		"lint two services": {
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

			fl, err := NewFileLinter()
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

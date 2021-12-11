package apilinter

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"sort"
	"testing"

	"github.com/google/go-cmp/cmp"
	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/descriptorpb"
	"google.golang.org/protobuf/types/pluginpb"
)

func TestFileLinter_LintFile(t *testing.T) {
	tests := map[string]struct {
		protoFile string
		out       []string
	}{
		"invalid file, 2 errors": {
			protoFile: "service.proto",
			out: []string{
				"core::0131::request-unknown-fields",
				"core::0131::request-name-required",
			},
		},
		"valid file, 0 errors": {
			protoFile: "service_ok.proto",
			out: []string{
				"core::0131::request-name-behavior",
				"core::0131::request-name-reference",
				"core::0192::has-comments",
				"core::0192::has-comments",
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

			file, exists := gen.FilesByPath[test.protoFile]
			if !exists {
				t.Fatalf("proto file not found: %s", test.protoFile)
			}

			fl, err := NewFileLinter()
			if err != nil {
				t.Fatal(err)
			}

			out, err := fl.LintFile(gen, file)
			if err != nil {
				t.Fatalf("fl.LintFile, test.protoFile=%s %s", test.protoFile, err)
			}

			// we only care about the rules that were triggered, extract them
			// here to make testing simpler
			rules := make([]string, 0)
			for _, r := range out {
				rules = append(rules, string(r.RuleID))
			}

			// order may not be deterministic from the library, so sort both
			// slices
			sort.Strings(test.out)
			sort.Strings(rules)

			// Ensure we have what we expect
			if diff := cmp.Diff(test.out, rules); diff != "" {
				t.Errorf("LintFileDescriptor mismatch (-want +have):\n%s", diff)
			}
		})
	}
}

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

			s, _ := PrettyString(out)
			t.Logf("res: %s", s)
		})
	}
}

func PrettyPrint(in interface{}) {
	s, err := PrettyString(in)
	if err != nil {
		panic(err)
	}
	fmt.Printf("%s\n", s)
}

func PrettyString(in interface{}) (string, error) {
	pp, err := json.MarshalIndent(in, "", "    ")
	if err != nil {
		return "", err
	}
	return string(pp), nil
}

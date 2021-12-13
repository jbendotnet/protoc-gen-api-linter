# Protoc Gen API Linter

**Currently early ALPHA, APIs can and will change**

Wraps https://github.com/googleapis/api-linter as a protoc plugin.

## Protosets for unit tests

If you add or change the protos in `internal/apilinter/testdata` do this to regenerate:

```
buf build internal/apilinter/testdata/ --output internal/apilinter/testdata/protos.bin
```

## Using the API Linter plugin locally

First, you need to generate the binary:

```
$ go install ./cmd/protoc-gen-api-linter
```

```
$ protoc-gen-api-linter help
Usage of protoc-gen-api-linter:
  -disable_rule value
    	Disable a rule with the given name.
    	May be specified multiple times.
  -enable_rule value
    	Enable a rule with the given name.
    	May be specified multiple times.
  -exit_on_error
    	Exit on first error (default true)
  -report_filename string
    	Set the filename of the JSON report (default "api_linter.json")
  -report_pretty_print
    	Pretty print JSON reports
  -version
    	Print version and exit.
```

> `buf.gen.yaml` is configured in this repo to use the binary installed in `$GOPATH/bin`

Then run `buf generate --debug -vv` and you should see the following:

```
$ tree gen
gen
└── reports
    └── apilinter-report.json

1 directory, 1 file
```

## Roadmap

- [x] Expose base `api-linter` functionality
- [x] Add support for enabling and disabling one or more rules
- [x] Add optional non-zero exit behaviour to halt code-generation process
- [ ] Add verbose output option
- [ ] Add support for using a local config file
- [ ] Add to Buf BSR for remote code generation
- [ ] Add Goreleaser

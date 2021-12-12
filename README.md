# Protoc Gen API Linter

**Currently early ALPHA, APIs can and will change**

Wraps https://github.com/googleapis/api-linter as a protoc plugin.

## Template unit tests

If you add or change the protos in `internal/apilinter/testdata` do this to regenerate:

```
buf build internal/apilinter/testdata/ --output internal/apilinter/testdata/protos.bin
```

## Using the API Linter plugin locally

First, you need to generate the binary:

```
go install ./cmd/protoc-gen-api-linter
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
- [ ] Add support for using a local config file

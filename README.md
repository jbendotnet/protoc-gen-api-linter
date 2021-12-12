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
$ tree reports/
reports/
└── api_linter_report.json

0 directories, 1 file
```



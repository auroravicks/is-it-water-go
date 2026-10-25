# Contributing

Thanks for helping improve Is It Water? Keep changes focused, documented, and
easy to verify.

## Local setup

1. Install the Go version declared in `go.mod`.
2. Clone the repository and run `go mod download`.
3. Install the required ASTWBD archives from
   [`is-on-water-dataset`](https://github.com/auroravicks/is-on-water-dataset).

## Before opening a pull request

```bash
gofmt -w .
go vet ./...
go build ./...
```

For API integration checks, start `go run .` in one terminal and run
`go test ./...` in another. Describe any dataset-dependent behavior you tested
in the pull request.

Please avoid committing generated artifacts, local environments, or unrelated
dataset changes.

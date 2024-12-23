## Setting up the environment

First, make sure you have Go installed (version 1.18 or higher). You can download it from [golang.org](https://golang.org/dl/).

## Dependencies Management

We use Go modules for dependency management. Initialize the module (if not already done):

```shell
go mod init
```

Install dependencies:

```shell
go mod tidy
```

## Running Tests

Run all tests:

```shell
go test ./...
```

Run tests with coverage:

```shell
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out  # View coverage report in browser
```
# SandID

[![Test](https://github.com/aofei/sandid/actions/workflows/test.yaml/badge.svg)](https://github.com/aofei/sandid/actions/workflows/test.yaml)
[![codecov](https://codecov.io/gh/aofei/sandid/branch/master/graph/badge.svg)](https://codecov.io/gh/aofei/sandid)
[![Go Report Card](https://goreportcard.com/badge/github.com/aofei/sandid)](https://goreportcard.com/report/github.com/aofei/sandid)
[![Go Reference](https://pkg.go.dev/badge/github.com/aofei/sandid.svg)](https://pkg.go.dev/github.com/aofei/sandid)

Every grain of sand on Earth has its own ID.

The algorithm used to generate the [`sandid.SandID`](https://pkg.go.dev/github.com/aofei/sandid#SandID) mainly come
from the [UUID](https://tools.ietf.org/html/rfc4122) version 1. Some
[adjustments](https://www.percona.com/blog/2014/12/19/store-uuid-optimized-way/) were made to enhance the efficiency of
database insertion.

## Features

- Extremely easy to use
- Fixed length
  - 16 bytes
  - 22 characters
  - 128-bit
- Huge capacity
  - Up to 2e128
- URL safe
  - `^[A-Za-z0-9-_]{22}$`
- Encoding friendly
  - Implemented [`encoding.TextMarshaler`](https://pkg.go.dev/encoding#TextMarshaler) and [`encoding.TextUnmarshaler`](https://pkg.go.dev/encoding#TextUnmarshaler)
  - Implemented [`encoding.BinaryMarshaler`](https://pkg.go.dev/encoding#BinaryMarshaler) and [`encoding.BinaryUnmarshaler`](https://pkg.go.dev/encoding#BinaryUnmarshaler)
  - Implemented [`json.Marshaler`](https://pkg.go.dev/encoding/json#Marshaler) and [`json.Unmarshaler`](https://pkg.go.dev/encoding/json#Unmarshaler)
- SQL friendly
  - [`sandid.NullSandID`](https://pkg.go.dev/github.com/aofei/sandid#NullSandID) support
  - Implemented [`sql.Scanner`](https://pkg.go.dev/database/sql#Scanner) and [`driver.Valuer`](https://pkg.go.dev/database/sql/driver#Valuer)
- Zero third-party dependencies

## Installation

To use this project programmatically, `go get` it:

```bash
go get github.com/aofei/sandid
```

## Community

If you have any questions or ideas about this project, feel free to discuss them
[here](https://github.com/aofei/sandid/discussions).

## Contributing

If you would like to contribute to this project, please submit issues [here](https://github.com/aofei/sandid/issues)
or pull requests [here](https://github.com/aofei/sandid/pulls).

When submitting a pull request, please make sure its commit messages adhere to
[Conventional Commits 1.0.0](https://www.conventionalcommits.org/en/v1.0.0/).

## License

This project is licensed under the [MIT License](LICENSE).

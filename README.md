# SandID

[![GitHub Actions](https://github.com/aofei/sandid/workflows/Test/badge.svg)](https://github.com/aofei/sandid)
[![codecov](https://codecov.io/gh/aofei/sandid/branch/master/graph/badge.svg)](https://codecov.io/gh/aofei/sandid)
[![Go Report Card](https://goreportcard.com/badge/github.com/aofei/sandid)](https://goreportcard.com/report/github.com/aofei/sandid)
[![PkgGoDev](https://pkg.go.dev/badge/github.com/aofei/sandid)](https://pkg.go.dev/github.com/aofei/sandid)

Every grain of sand on Earth has its own ID.

**Note that the algorithm used to generate the
[`sandid.SandID`](https://pkg.go.dev/github.com/aofei/sandid#SandID) mainly come
from the [UUID](https://tools.ietf.org/html/rfc4122) version 1. Some adjustments
were made to enhance the efficiency of database insertion (see
[this](https://www.percona.com/blog/2014/12/19/store-uuid-optimized-way/)).**

## Features

* Extremely easy to use
* Fixed length
	* 16 bytes
	* 22 characters
	* 128-bit
* Huge capacity
	* Up to 2e128
* URL safe
	* `^[A-Za-z0-9-_]{22}$`
* Encoding friendly
	* Implemented [`encoding.TextMarshaler`](https://pkg.go.dev/encoding#TextMarshaler) and [`encoding.TextUnmarshaler`](https://pkg.go.dev/encoding#TextUnmarshaler)
	* Implemented [`encoding.BinaryMarshaler`](https://pkg.go.dev/encoding#BinaryMarshaler) and [`encoding.BinaryUnmarshaler`](https://pkg.go.dev/encoding#BinaryUnmarshaler)
	* Implemented [`json.Marshaler`](https://pkg.go.dev/encoding/json#Marshaler) and [`json.Unmarshaler`](https://pkg.go.dev/encoding/json#Unmarshaler)
* SQL friendly
	* [`sandid.NullSandID`](https://pkg.go.dev/github.com/aofei/sandid#NullSandID) support
	* Implemented [`sql.Scanner`](https://pkg.go.dev/database/sql#Scanner) and [`driver.Valuer`](https://pkg.go.dev/database/sql/driver#Valuer)
* Zero third-party dependencies

## Installation

Open your terminal and execute

```bash
$ go get github.com/aofei/sandid
```

done.

> The only requirement is the [Go](https://golang.org), at least v1.13.

## Community

If you want to discuss SandID, or ask questions about it, simply post questions
or ideas [here](https://github.com/aofei/sandid/issues).

## Contributing

If you want to help build SandID, simply follow
[this](https://github.com/aofei/sandid/wiki/Contributing) to send pull requests
[here](https://github.com/aofei/sandid/pulls).

## License

This project is licensed under the MIT License.

License can be found [here](LICENSE).

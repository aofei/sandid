# SandID

[![GitHub Actions](https://github.com/aofei/sandid/workflows/Main/badge.svg)](https://github.com/aofei/sandid)
[![codecov](https://codecov.io/gh/aofei/sandid/branch/master/graph/badge.svg)](https://codecov.io/gh/aofei/sandid)
[![Go Report Card](https://goreportcard.com/badge/github.com/aofei/sandid)](https://goreportcard.com/report/github.com/aofei/sandid)
[![GoDoc](https://godoc.org/github.com/aofei/sandid?status.svg)](https://godoc.org/github.com/aofei/sandid)

Every grain of sand on Earth has its own ID.

**Note that the algorithm used to generate the
[`sandid.SandID`](https://godoc.org/github.com/aofei/sandid#SandID) mainly come
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
	* Implemented [`encoding.TextMarshaler`](https://godoc.org/encoding#TextMarshaler) and [`encoding.TextUnmarshaler`](https://godoc.org/encoding#TextUnmarshaler)
	* Implemented [`encoding.BinaryMarshaler`](https://godoc.org/encoding#BinaryMarshaler) and [`encoding.BinaryUnmarshaler`](https://godoc.org/encoding#BinaryUnmarshaler)
	* Implemented [`json.Marshaler`](https://godoc.org/encoding/json#Marshaler) and [`json.Unmarshaler`](https://godoc.org/encoding/json#Unmarshaler)
* SQL friendly
	* [`sandid.NullSandID`](https://godoc.org/github.com/aofei/sandid#NullSandID) support
	* Implemented [`sql.Scanner`](https://godoc.org/database/sql#Scanner) and [`driver.Valuer`](https://godoc.org/database/sql/driver#Valuer)

## Installation

Open your terminal and execute

```bash
$ go get github.com/aofei/sandid
```

done.

> The only requirement is the [Go](https://golang.org), at least v1.5.

## Community

If you want to discuss SandID, or ask questions about it, simply post questions
or ideas [here](https://github.com/aofei/sandid/issues).

## Contributing

If you want to help build SandID, simply follow
[this](https://github.com/aofei/sandid/wiki/Contributing) to send pull requests
[here](https://github.com/aofei/sandid/pulls).

## License

This project is licensed under the Unlicense.

License can be found [here](LICENSE).

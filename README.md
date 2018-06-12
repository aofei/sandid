# SandID

[![Build Status](https://travis-ci.org/sheng/sandid.svg?branch=master)](https://travis-ci.org/sheng/sandid)
[![Coverage Status](https://coveralls.io/repos/github/sheng/sandid/badge.svg?branch=master)](https://coveralls.io/github/sheng/sandid?branch=master)
[![Go Report Card](https://goreportcard.com/badge/github.com/sheng/sandid)](https://goreportcard.com/report/github.com/sheng/sandid)
[![GoDoc](https://godoc.org/github.com/sheng/sandid?status.svg)](https://godoc.org/github.com/sheng/sandid)

Every grain of sand on earth has its own ID.

**It should be noted that the algorithm used to generate the `SandID` in this
project mainly come from UUID version 1. Some adjustments were made to optimize
the efficiency of database insertion (see
[this](https://www.percona.com/blog/2014/12/19/store-uuid-optimized-way/)).**

## Features

* Extremely easy to use
* Fixed length
	* 16 bytes
	* 32 hex characters
	* 128-bit
* Huge capacity
	* Up to 2e128
* Case insensitivity
* URL friendly
* Encoding friendly
	* Implemented `encoding.TextMarshaler` and `encoding.TextUnmarshaler`
	* Implemented `encoding.BinaryMarshaler` and `encoding.BinaryUnmarshaler`
	* Implemented `json.Marshaler` and `json.Unmarshaler`
* SQL friendly
	* `NullSandID` support
	* Implemented `sql.Scanner` and `driver.Valuer`

## Installation

Open your terminal and execute

```bash
$ go get github.com/sheng/sandid
```

done.

> The only requirement is the [Go](https://golang.org), at least v1.2.

## Community

If you want to discuss this project, or ask questions about it, simply post
questions or ideas [here](https://github.com/sheng/sandid/issues).

## Contributing

If you want to help build this project, simply follow
[this](https://github.com/sheng/sandid/wiki/Contributing) to send pull requests
[here](https://github.com/sheng/sandid/pulls).

## License

This project is licensed under the Unlicense.

License can be found [here](LICENSE).

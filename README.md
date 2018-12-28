# SandID

[![Build Status](https://travis-ci.org/aofei/sandid.svg?branch=master)](https://travis-ci.org/aofei/sandid)
[![Coverage Status](https://coveralls.io/repos/github/aofei/sandid/badge.svg?branch=master)](https://coveralls.io/github/aofei/sandid?branch=master)
[![Go Report Card](https://goreportcard.com/badge/github.com/aofei/sandid)](https://goreportcard.com/report/github.com/aofei/sandid)
[![GoDoc](https://godoc.org/github.com/aofei/sandid?status.svg)](https://godoc.org/github.com/aofei/sandid)

Every grain of sand on earth has its own ID.

**It should be noted that the algorithm used to generate the `SandID` in this
project mainly come from the UUID version 1. Some adjustments were made to
enhance the efficiency of database insertion (see
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
	* Implemented `encoding.TextMarshaler` and `encoding.TextUnmarshaler`
	* Implemented `encoding.BinaryMarshaler` and `encoding.BinaryUnmarshaler`
	* Implemented `json.Marshaler` and `json.Unmarshaler`
* SQL friendly
	* `NullSandID` support
	* Implemented `sql.Scanner` and `driver.Valuer`

## Installation

Open your terminal and execute

```bash
$ go get github.com/aofei/sandid
```

done.

> The only requirement is the [Go](https://golang.org), at least v1.2.

## Community

If you want to discuss this project, or ask questions about it, simply post
questions or ideas [here](https://github.com/aofei/sandid/issues).

## Contributing

If you want to help build this project, simply follow
[this](https://github.com/aofei/sandid/wiki/Contributing) to send pull requests
[here](https://github.com/aofei/sandid/pulls).

## License

This project is licensed under the Unlicense.

License can be found [here](LICENSE).

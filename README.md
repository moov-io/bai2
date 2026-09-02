[![Moov Banner Logo](https://user-images.githubusercontent.com/20115216/104214617-885b3c80-53ec-11eb-8ce0-9fc745fb5bfc.png)](https://github.com/moov-io)

<p align="center">
  <a href="https://pkg.go.dev/github.com/moov-io/bai2">GoDoc</a>
  ·
  <a href="https://moov.io/blog/craft/bai2-api-guide/">API Guide</a>
  ·
  <a href="https://github.com/moov-io/bai2/releases">Releases</a>
  ·
  <a href="https://slack.moov.io/">Community</a>
  ·
  <a href="https://moov.io/blog/">Blog</a>
  <br>
  <br>
</p>

[![GoDoc](https://pkg.go.dev/badge/github.com/moov-io/bai2.svg)](https://pkg.go.dev/github.com/moov-io/bai2)
[![Build Status](https://github.com/moov-io/bai2/workflows/Go/badge.svg)](https://github.com/moov-io/bai2/actions)
[![Coverage Status](https://codecov.io/gh/moov-io/bai2/branch/master/graph/badge.svg)](https://codecov.io/gh/moov-io/bai2)
[![Go Report Card](https://goreportcard.com/badge/github.com/moov-io/bai2)](https://goreportcard.com/report/github.com/moov-io/bai2)
[![Repo Size](https://img.shields.io/github/languages/code-size/moov-io/bai2?label=project%20size)](https://github.com/moov-io/bai2)
[![Apache 2 License](https://img.shields.io/badge/license-Apache2-blue.svg)](https://raw.githubusercontent.com/moov-io/bai2/master/LICENSE)
[![Slack Channel](https://slack.moov.io/badge.svg?bg=e01563&fgColor=fffff)](https://slack.moov.io/)
[![Docker Pulls](https://img.shields.io/docker/pulls/moov/bai2)](https://hub.docker.com/r/moov/bai2)
[![GitHub Stars](https://img.shields.io/github/stars/moov-io/bai2)](https://github.com/moov-io/bai2)
[![Twitter](https://img.shields.io/twitter/follow/moov?style=social)](https://twitter.com/moov?lang=en)


# moov-io/bai2

Moov's mission is to give developers an easy way to create and integrate bank processing into their own software products. Our open source projects are each focused on solving a single responsibility in financial services and designed around performance, scalability, and ease-of-use.

This project is a reader, writer, and validator for [BAI2](https://en.wikipedia.org/wiki/BAI_(file_format)) (Cash Management Balance Reporting Specifications Version 2) and X9.121 BTR3 (BAI3) files. It ships as a Go module (`github.com/moov-io/bai2`), a CLI, and an HTTP server in a [Docker image](#docker).

v1.0.0 is a stable release. If you imported `pkg/lib`, see [Migrating](docs/MIGRATING.md).

## Table of contents

- [Project status](#project-status)
- [BAI2 and BAI3](#bai2-and-bai3)
- [Usage](#usage)
  - [As an API](#docker)
  - [As a Go module](#go-library)
  - [As a command line tool](#command-line)
- [Learn about BAI](#learn-about-bai)
- [Getting help](#getting-help)
- [Supported and tested platforms](#supported-and-tested-platforms)
- [Contributing](#contributing)
- [Migrating from pkg/lib](docs/MIGRATING.md)
- [Related projects](#related-projects)

## Project status

Moov Bai2 is used in production. v1.0.0 follows [semver](https://semver.org/): BAI2 lives in `pkg/bai2`, BAI3 / BTR3 in `pkg/bai3`, and `github.com/moov-io/bai2.Read` picks the parser from the file header. Please star the project if you are interested in its progress. Issues and pull requests are welcome.

## BAI2 and BAI3

Both formats use the same 01/02/03/16/88/49/98/99 envelope. The 01 Version Number chooses the grammar.

| Version | Package | Notes |
|---|---|---|
| 2 | [`pkg/bai2`](https://pkg.go.dev/github.com/moov-io/bai2/pkg/bai2) | Group header on record 02. Omitted group currency defaults to USD. |
| 3 | [`pkg/bai3`](https://pkg.go.dev/github.com/moov-io/bai2/pkg/bai3) | Bank header on record 02. Account currency is required. |

`ReadOptions.IgnoreVersion` (CLI `--ignoreVersion`, HTTP `?ignoreVersion=true`) still means **parse as BAI2**. Use that for BAI2 files banks stamp with version `3`. It does not enable the BAI3 grammar.

## Usage

### Docker

We publish a [public Docker image `moov/bai2`](https://hub.docker.com/r/moov/bai2/) on Docker Hub with each tagged release. No configuration is required to serve on `:8208`.

Pull & start the Docker image:
```
docker pull moov/bai2:latest
docker run -p 8208:8208 moov/bai2:latest web
```

Upload a file and parse it:
```
curl -X POST --form "input=@./test/testdata/sample1.txt" http://localhost:8208/parse
```
```
{"status":"valid file"}
```

Print a file after parse:
```
curl -X POST --form "input=@./test/testdata/sample1.txt" http://localhost:8208/print
```
```
01,0004,12345,060321,0829,001,80,1,2/
02,12345,0004,1,060317,,CAD,/
03,10200123456,CAD,040,+000000000000,,,045,+000000000000,,/
88,100,000000000208500,00003,V,060316,,400,000000000208500,00008,V,060316,/
16,409,000000000002500,V,060316,,,,RETURNED CHEQUE     /
...
99,+00000000001280000,000000001,000000027/
```

Format to JSON after parse:
```
curl -X POST --form "input=@./test/testdata/sample1.txt" http://localhost:8208/format | jq .
```
<details>
<summary>JSON Response</summary>

```json
{
  "sender": "0004",
  "receiver": "12345",
  "fileCreatedDate": "060321",
  "fileCreatedTime": "0829",
  "fileIdNumber": "001",
  "physicalRecordLength": 80,
  "blockSize": 1,
  "versionNumber": 2,
  "fileControlTotal": "+00000000001280000",
  "numberOfGroups": 1,
  "numberOfRecords": 27,
  "Groups": [
    {
      "receiver": "12345",
      "originator": "0004",
      "groupStatus": 1,
      "asOfDate": "060317",
      "currencyCode": "CAD",
      "groupControlTotal": "+00000000001280000",
      "numberOfAccounts": 2,
      "numberOfRecords": 25,
      "Accounts": [
        {
          "accountNumber": "10200123456",
          "currencyCode": "CAD",
          "summaries": [
            {
              "TypeCode": "040",
              "Amount": "+000000000000",
              "ItemCount": 0,
              "FundsType": {}
            },
            {
              "TypeCode": "045",
              "Amount": "+000000000000",
              "ItemCount": 0,
              "FundsType": {}
            },
            {
              "TypeCode": "100",
              "Amount": "000000000208500",
              "ItemCount": 3,
              "FundsType": {
                "type_code": "V",
                "date": "060316"
              }
            },
            {
              "TypeCode": "400",
              "Amount": "000000000208500",
              "ItemCount": 8,
              "FundsType": {
                "type_code": "V",
                "date": "060316"
              }
            }
          ],
          "accountControlTotal": "+00000000000834000",
          "numberRecords": 14,
          "Details": [
            {
              "TypeCode": "409",
              "Amount": "000000000002500",
              "FundsType": {
                "type_code": "V",
                "date": "060316"
              },
              "BankReferenceNumber": "",
              "CustomerReferenceNumber": "",
              "Text": "RETURNED CHEQUE     "
            }
          ]
        }
      ]
    }
  ]
}
```
</details>

Query flags on `/parse`, `/print`, and `/format`:

- `?ignoreVersion=true` — parse as BAI2 regardless of the 01 version field
- `?strictControlTotals=true` — reject files whose 49/98/99 totals do not match the body

BAI3 `/format` JSON uses `Banks` (record 02 bank headers) instead of `Groups`.

#### Data persistence
By design, Bai2 **does not persist** any data about the files or entry details created. The only storage occurs in memory of the process and upon restart Bai2 will have no files or data saved. Also, no in-memory encryption of the data is performed.

### Go library

This project uses [Go Modules](https://go.dev/blog/using-go-modules) and Go 1.25 or newer. See [Golang's install instructions](https://golang.org/doc/install) for help setting up Go. Use a [tagged release](https://github.com/moov-io/bai2/releases/latest) in production.

```
go get github.com/moov-io/bai2@v1.0.0
```

`Read` looks at the 01 Version Number and returns a BAI2 or BAI3 document:

```go
import (
    "fmt"
    "os"

    moovbai "github.com/moov-io/bai2"
    "github.com/moov-io/bai2/pkg/bai2"
    "github.com/moov-io/bai2/pkg/bai3"
)

f, err := moovbai.Read(os.Stdin)
if err != nil {
    return err
}
if err := f.Validate(); err != nil {
    return err
}

switch doc := f.(type) {
case *bai2.Bai2:
    fmt.Printf("BAI2 sender %s, %d groups\n", doc.Sender, len(doc.Groups))
case *bai3.File:
    fmt.Printf("BAI3 sender %s, %d banks\n", doc.Sender, len(doc.Banks))
}
```

To parse only BAI2:

```go
scan := bai2.NewBai2Scanner(r)
file := bai2.NewBai2()
if err := file.Read(&scan); err != nil {
    return err
}
```

Appendix A type codes are on parsed summaries and details (`LookupTypeCode`, `TypeInfo()`). Call `Create()` to fill 49/98/99 control totals. The HTTP server loads `configs/config.default.yml` via `go:embed`; set `APP_CONFIG` to override.

```
go doc github.com/moov-io/bai2
go doc github.com/moov-io/bai2/pkg/bai2
go doc github.com/moov-io/bai2/pkg/bai3
```

### Command line

Bai2 has a command line interface to parse files and launch the web service.

```
$ bai2 --help
```
```
Usage:
   [command]

Available Commands:
  completion  Generate the autocompletion script for the specified shell
  format      Format bai2 report
  help        Help about any command
  parse       parse bai2 report
  print       Print bai2 report
  web         Launches web server

Flags:
  -h, --help                  help for this command
      --ignoreVersion          ignore the version number in the file header
      --input string           bai2 report file
      --strictControlTotals    reject files whose trailer totals do not match the body

Use " [command] --help" for more information about a command.
```

`parse`, `print`, and `format` auto-detect BAI2 vs BAI3 from the 01 version field.

## Learn about BAI

- [BAI file format (Wikipedia)](https://en.wikipedia.org/wiki/BAI_(file_format))
- [Cash Management Balance Reporting Specifications Version 2](docs/specifications/Cash%20Management%20Balance%20Reporting%20Specifications%20Version%202.pdf) (BAI2)
- [X9.121 BTRS Version 3 Format Guide](https://x9.org/wp-content/uploads/2018/07/X9.121-2016-BTRS-Version-3.0.pdf) (public samples; the full standard is not vendored)
- [BAI2 API guide](https://moov.io/blog/craft/bai2-api-guide/)

## Getting help

 channel | info
 ------- | -------
Twitter [@moov](https://twitter.com/moov)	| You can follow Moov.io's Twitter feed to get updates on our project(s). You can also tweet us questions or just share blogs or stories.
[GitHub Issue](https://github.com/moov-io/bai2/issues) | If you are able to reproduce a problem please open a GitHub Issue under the specific project that caused the error.
[moov-io slack](https://slack.moov.io/) | Join our (`#bai2`) slack channel to have an interactive discussion about the development of the project.

## Supported and tested platforms

- 64-bit Linux (Ubuntu, Debian), macOS, and Windows

## Contributing

Yes please! Open an issue or pull request. Join (`#bai2`) on [moov-io slack](https://slack.moov.io/) if you want to talk through a change.

This project uses [Go Modules](https://go.dev/blog/using-go-modules) and Go 1.25 or newer. Use a [tagged release](https://github.com/moov-io/bai2/releases/latest) in production.

### Releasing

Release notes live on [GitHub Releases](https://github.com/moov-io/bai2/releases). To cut a release, push a semver tag (`git push origin v1.0.1`); CI tests, publishes binaries, and pushes the Docker image.

### Testing

We maintain a comprehensive suite of unit tests and recommend table-driven testing when a particular function warrants several very similar test cases. To run all test files in the current directory, use `go test`. Current overall coverage can be found on [Codecov](https://app.codecov.io/gh/moov-io/bai2/).

### Fuzzing

We currently run fuzzing over BAI2 and BAI3 in the form of a [GitHub Action](https://github.com/moov-io/bai2/actions/workflows/fuzz.yml). Please report crash examples to [`oss@moov.io`](mailto:oss@moov.io). Thanks!

## Related projects
As part of Moov's initiative to offer open source fintech infrastructure, we have a large collection of active projects you may find useful:

- [Moov Watchman](https://github.com/moov-io/watchman) offers search functions over numerous trade sanction lists from the United States and European Union.

- [Moov Fed](https://github.com/moov-io/fed) implements utility services for searching the United States Federal Reserve System such as ABA routing numbers, financial institution name lookup, and FedACH and Fedwire routing information.

- [Moov Wire](https://github.com/moov-io/wire) implements an interface to write files for the Fedwire Funds Service, a real-time gross settlement funds transfer system operated by the United States Federal Reserve Banks.

- [Moov ACH](https://github.com/moov-io/ach) provides ACH file generation and parsing, supporting all Standard Entry Codes for the primary method of money movement throughout the United States.

- [Moov Image Cash Letter](https://github.com/moov-io/imagecashletter) implements Image Cash Letter (ICL) files used for Check21, X.9 or check truncation files for exchange and remote deposit in the U.S.

## License

Apache License 2.0 - See [LICENSE](LICENSE) for details.

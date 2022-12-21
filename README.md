[![Moov Banner Logo](https://user-images.githubusercontent.com/20115216/104214617-885b3c80-53ec-11eb-8ce0-9fc745fb5bfc.png)](https://github.com/moov-io)

<p align="center">
  <a href="https://moov-io.github.io/bai2/">Project Documentation</a>
  ·
  <a href="https://moov-io.github.io/bai2/api/#overview">API Endpoints</a>
  ·
  <a href="https://moov.io/blog/education/bai2-api-guide/">API Guide</a>
  ·
  <a href="https://slack.moov.io/">Community</a>
  ·
  <a href="https://moov.io/blog/">Blog</a>
  <br>
  <br>
</p>

[![GoDoc](https://godoc.org/github.com/moov-io/bai2?status.svg)](https://godoc.org/github.com/moov-io/bai2)
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

Bai2 implements a reader, writer, and validator for the [Cash Management Balance Reporting Specifications Version 2](https://en.wikipedia.org/wiki/BAI_(file_format)) developed by Bank Administration Institute. This project offers a HTTP server in a [Docker image](#docker) and a Go package `github.com/moov-io/bai2`.

## Table of contents

- [Project status](#project-status)
- [Usage](#usage)
  - As an API
    - [Docker](#docker)
    - [Google Cloud](#google-cloud-run)
    - [Data persistence](#data-persistence)
  - [As a Go module](#go-library)
  - [As a command line tool](#command-line)
- [Learn about Bai 2](#learn-about-bai-2)
- [Getting help](#getting-help)
- [Supported and tested platforms](#supported-and-tested-platforms)
- [Contributing](#contributing)
- [Related projects](#related-projects)

## Project status

Moov Bai2 is actively used in multiple production environments. Please star the project if you are interested in its progress. If you have layers above Bai2 to simplify tasks, perform business operations, or found bugs we would appreciate an issue or pull request. Thanks!

## Usage

The Bai2 project implements an HTTP server and [Go library](https://pkg.go.dev/github.com/moov-io/bai2) for creating and modifying files in bai 2 format, which developed a generic format and widely accepted by most of the Banks in USA.

### Docker

We publish a [public Docker image `moov/bai2`](https://hub.docker.com/r/moov/bai2/) on Docker Hub with each tagged release of Bai2. No configuration is required to serve on `:8080`. We also have Docker images for [OpenShift](https://quay.io/repository/moov/bai2?tab=tags) published as `quay.io/moov/bai2`.

Pull & start the Docker image:
```
docker pull moov/bai2:latest
docker run -p 8080:8080 moov/bai2:latest
```

Upload a file and pase it:
```
curl -X POST --form "input=@sample.txt" http://localhost:8080/parse
```
```
{"status":"valid file"}
```

Print a file after parse:
```
curl -X POST --form "input=@sample.txt" http://localhost:8080/print
```
```
01,0004,12345,060321,0829,001,80,1,2/
02,12345,0004,1,060317,,CAD,/
03,10200123456,CAD,040,+000000000000,,,045,+000000000000,,/
88,100,000000000208500,00003,V,060316,,400,000000000208500,00008,V,060316,/
16,409,000000000002500,V,060316,,,,RETURNED CHEQUE     /
16,409,000000000090000,V,060316,,,,RTN-UNKNOWN         /
16,409,000000000000500,V,060316,,,,RTD CHQ SERVICE CHRG/
16,108,000000000203500,V,060316,,,,TFR 1020 0345678    /
16,108,000000000002500,V,060316,,,,MACLEOD MALL        /
16,108,000000000002500,V,060316,,,,MASCOUCHE QUE       /
16,409,000000000020000,V,060316,,,,1000 ISLANDS MALL   /
16,409,000000000090000,V,060316,,,,PENHORA MALL        /
16,409,000000000002000,V,060316,,,,CAPILANO MALL       /
16,409,000000000002500,V,060316,,,,GALERIES LA CAPITALE/
16,409,000000000001000,V,060316,,,,PLAZA ROCK FOREST   /
49,+00000000000834000,000000014/
03,10200123456,CAD,040,+000000000000,,,045,+000000000000,,/
88,100,000000000111500,00002,V,060317,,400,000000000111500,00004,V,060317,/
16,108,000000000011500,V,060317,,,,TFR 1020 0345678    /
16,108,000000000100000,V,060317,,,,MONTREAL            /
16,409,000000000100000,V,060317,,,,GRANDFALL NB        /
16,409,000000000009000,V,060317,,,,HAMILTON ON         /
16,409,000000000002000,V,060317,,,,WOODSTOCK NB        /
16,409,000000000000500,V,060317,,,,GALERIES RICHELIEU  /
49,+00000000000446000,000000009/
98,+00000000001280000,000000002,000000025/
99,+00000000001280000,000000001,000000027/
...
```

### Google Cloud Run

To get started in a hosted environment you can deploy this project to the Google Cloud Platform.

From your [Google Cloud dashboard](https://console.cloud.google.com/home/dashboard) create a new project and call it:
```
moov-bai2-demo
```

Enable the [Container Registry](https://cloud.google.com/container-registry) API for your project and associate a [billing account](https://cloud.google.com/billing/docs/how-to/manage-billing-account) if needed. Then, open the Cloud Shell terminal and run the following Docker commands, substituting your unique project ID:

```
docker pull moov/bai2
docker tag moov/bai2 gcr.io/<PROJECT-ID>/bai2
docker push gcr.io/<PROJECT-ID>/bai2
```

Deploy the container to Cloud Run:
```
gcloud run deploy --image gcr.io/<PROJECT-ID>/bai2 --port 8080
```

Select your target platform to `1`, service name to `bai2`, and region to the one closest to you (enable Google API service if a prompt appears). Upon a successful build you will be given a URL where the API has been deployed:

```
https://YOUR-BAI2-APP-URL.a.run.app
```

Now you can complete a health check:
```
curl https://YOUR-BAI2-APP-URL.a.run.app/health
```
You should get this response:
```
{"health":true}
```

### Data persistence
By design, Bai2  **does not persist** (save) any data about the files or entry details created. The only storage occurs in memory of the process and upon restart Bai2 will have no files or data saved. Also, no in-memory encryption of the data is performed.

### Go library

This project uses [Go Modules](https://go.dev/blog/using-go-modules) and Go v1.18 or newer. See [Golang's install instructions](https://golang.org/doc/install) for help setting up Go. You can download the source code and we offer [tagged and released versions](https://github.com/moov-io/bai2/releases/latest) as well. We highly recommend you use a tagged release for production.

```
$ git@github.com:moov-io/bai2.git

$ go get -u github.com/moov-io/bai2

$ go doc github.com/moov-io/bai2
```

### Command line

Bai2 has a command line interface to manage Bai 2 files and launch a web service.

```
bai2 --help

Usage:
   [command]

Available Commands:
  help        Help about any command
  print       Print bai2 file
  parse       Parse bai2 file
  web         Launches web server

Flags:
  -h, --help           help for this command
      --input string   bai2 report file

Use " [command] --help" for more information about a command.
```

Each interaction that the library supports is exposed in a command-line option:

 Command | Info
 ------- | -------
`print`  | Print bai2 file
`parse`  | Parse bai2 file
`web`    | Launches web server

## Learn about Bai 2
- [Bai 2](https://www.tdcommercialbanking.com/document/PDF/bai.pdf)
- [Cash Management](https://www.bai.org/docs/default-source/libraries/site-general-downloads/cash_management_2005.pdf)

## Getting help

 channel | info
 ------- | -------
Twitter [@moov](https://twitter.com/moov)	| You can follow Moov.io's Twitter feed to get updates on our project(s). You can also tweet us questions or just share blogs or stories.
[GitHub Issue](https://github.com/moov-io) | If you are able to reproduce a problem please open a GitHub Issue under the specific project that caused the error.
[moov-io slack](https://slack.moov.io/) | Join our slack channel to have an interactive discussion about the development of the project.

## Supported and tested platforms

- 64-bit Linux (Ubuntu, Debian), macOS, and Windows

## Contributing

Yes please! Please review our [Contributing guide](CONTRIBUTING.md) and [Code of Conduct](CODE_OF_CONDUCT.md) to get started!

This project uses [Go Modules](https://go.dev/blog/using-go-modules) and Go v1.18 or newer. See [Golang's install instructions](https://golang.org/doc/install) for help setting up Go. You can download the source code and we offer [tagged and released versions](https://github.com/moov-io/bai2/releases/latest) as well. We highly recommend you use a tagged release for production.

### Releasing

To make a release of bai2 simply open a pull request with `CHANGELOG.md` and `version.go` updated with the next version number and details. You'll also need to push the tag (i.e. `git push origin v1.0.0`) to origin in order for CI to make the release.

### Testing

We maintain a comprehensive suite of unit tests and recommend table-driven testing when a particular function warrants several very similar test cases. To run all test files in the current directory, use `go test`. Current overall coverage can be found on [Codecov](https://app.codecov.io/gh/moov-io/bai2/).

### Fuzzing

We currently run fuzzing over ImageCashLetter in the form of a [`moov/bai2`](https://hub.docker.com/r/moov/bai2fuzz) Docker image. You can [read more](./test/fuzz-reader/README.md) or run the image and report crasher examples to [`security@moov.io`](mailto:security@moov.io). Thanks!

## Related projects
As part of Moov's initiative to offer open source fintech infrastructure, we have a large collection of active projects you may find useful:

- [Moov Watchman](https://github.com/moov-io/watchman) offers search functions over numerous trade sanction lists from the United States and European Union.

- [Moov Fed](https://github.com/moov-io/fed) implements utility services for searching the United States Federal Reserve System such as ABA routing numbers, financial institution name lookup, and FedACH and Fedwire routing information.

- [Moov Wire](https://github.com/moov-io/wire) implements an interface to write files for the Fedwire Funds Service, a real-time gross settlement funds transfer system operated by the United States Federal Reserve Banks.

- [Moov ACH](https://github.com/moov-io/ach) provides ACH file generation and parsing, supporting all Standard Entry Codes for the primary method of money movement throughout the United States.

- [Moov Image Cash Letter](https://github.com/moov-io/imagecashletter) implements Image Cash Letter (ICL) files used for Check21, X.9 or check truncation files for exchange and remote deposit in the U.S.

## License

Apache License 2.0 - See [LICENSE](LICENSE) for details.

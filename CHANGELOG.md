## v0.4.0 (Released 2024-06-17)

IMPROVEMENTS

- feat: allow to ignore version check with `Options` struct

BUILD

- fix(deps): update module github.com/moov-io/base to v0.49.4
- fix(deps): update module github.com/spf13/cobra to v1.8.1

## v0.3.1 (Released 2024-05-15)

IMPROVEMENTS

- feat: remove unused continuation struct
- fix: allow some records to include `/` character (#113)

BUILD

- fix(deps): update module github.com/moov-io/base to v0.49.3

## v0.3.0 (Released 2024-04-16)

IMPROVEMENTS

- feat: Implement aggregate functions to support setting trailer record fields programatically
- feat: return all Details for an Account
- fix: normalize file paths for windows machines to fix failing test
- fix: read BAI2 rune by rune
- fix: separate fuzzing of valid BAI2 files from error files
- fix: validate BAI2 file after parsing
- fuzz: setup runner and scheduled job
- schema: OpenAPI models for Files, Groups, Accounts, and related objects

BUILD

- chore(deps): update golang docker tag to v1.22
- fix(deps): update module github.com/gorilla/mux to v1.8.1
- fix(deps): update module github.com/moov-io/base to v0.48.5
- fix(deps): update module github.com/spf13/cobra to v1.8.0
- fix(deps): update module github.com/stretchr/testify to v1.9.0
- fix(deps): update module golang.org/x/oauth2 to v0.18.0

## v0.2.0 (Released 2023-02-03)

IMPROVEMENTS

- feat: support varialble length records
- feat: Implemented file structure that included File, Group, Account
- feat: Handle contiuation records by returning merged records for callers
- fix: Updated returned types to return better results

BUILD

- chore(deps): update golang docker tag to v1.20
- fix(deps): update module golang.org/x/oauth2 to v0.4.0
- fix(deps): update module github.com/moov-io/base to v0.39.0

## v0.1.0 (Released 2022-12-21)

This is the initial releae of moov-io/bai2. Please join our (`#bai2`) [slack channel](https://slack.moov.io/) for updates and discussions.

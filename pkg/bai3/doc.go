// Copyright 2022 The Moov Authors
// Use of this source code is governed by an Apache License
// license that can be found in the LICENSE file.

// Package bai3 implements a reader, writer, and validator for X9.121 BTR3
// (BAI3) balance reporting files (01 Version Number 3).
//
// Record 02 is a bank header: currency is positional-null (not defaulted to
// USD) and group status is Update (1). Record 03 requires a currency code.
//
// For files that may be BAI2 or BTR3, use github.com/moov-io/bai2.Read.
// Golden files under test/testdata/bai3/ are assembled from public X9
// format-guide samples; the copyrighted standard is not vendored.
package bai3

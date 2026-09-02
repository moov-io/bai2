// Copyright 2022 The Moov Authors
// Use of this source code is governed by an Apache License
// license that can be found in the LICENSE file.

package bai3

// VersionBTR3 is the file header Version Number for X9.121 BTR3 / BAI3.
const VersionBTR3 int64 = 3

// Options controls parser and validator behavior.
type Options struct {
	// IgnoreVersion accepts a file whose 01 Version Number is not 3.
	IgnoreVersion bool
	// StrictControlTotals compares 49/98/99 aggregates to computed sums.
	StrictControlTotals bool
	// StrictBankHeader requires currency to be empty and group status to be 1
	// (BTR3 retired those 02 fields). Off by default so files that still send
	// positional values parse.
	StrictBankHeader bool
}

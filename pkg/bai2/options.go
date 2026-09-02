// Copyright 2022 The Moov Authors
// Use of this source code is governed by an Apache License
// license that can be found in the LICENSE file.

package bai2

// BAI / BTRS version numbers from the file header (01) Version Number field.
const (
	VersionBAI2 int64 = 2
	VersionBTR3 int64 = 3
)

// Options controls parser and validator behavior.
//
// IgnoreVersion accepts a file whose 01 Version Number is not 2. Use this for
// BAI2 files that banks stamp with 3 (or another value) without speaking BTR3.
// It does not enable BTR3 record layouts; those need a BTR3 reader.
//
// StrictControlTotals compares 49/98/99 control totals and record counts to
// values computed from the file body. Off by default because many production
// files have incorrect trailers.
type Options struct {
	IgnoreVersion       bool
	StrictControlTotals bool
}

// DefaultCurrency is used when a group or account omits Currency Code.
// Spec: group default is USD; account default is the group currency.
const DefaultCurrency = "USD"

// DefaultGroupStatus is used when a group omits Group Status.
// Spec: most transmissions are Update (1).
const DefaultGroupStatus int64 = 1

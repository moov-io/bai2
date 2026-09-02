// Copyright 2022 The Moov Authors
// Use of this source code is governed by an Apache License
// license that can be found in the LICENSE file.

package bai2

import (
	"io"

	"github.com/moov-io/bai2/pkg/util"
)

// Bai2Scanner is the envelope scanner shared with the BAI3 parser.
type Bai2Scanner = util.Scanner

// NewBai2Scanner returns a scanner that splits BAI records on '/' and newlines.
func NewBai2Scanner(fd io.Reader) Bai2Scanner {
	return util.NewScanner(fd)
}

// Copyright 2022 The Moov Authors
// Use of this source code is governed by an Apache License
// license that can be found in the LICENSE file.

package bai3

import (
	"errors"
	"fmt"
	"strings"

	"github.com/moov-io/bai2/pkg/util"
)

// NewDetail returns an empty BTR3 transaction detail (record 16).
func NewDetail() *Detail {
	return &Detail{}
}

// Detail Format
type Detail transactionDetail

func (r *Detail) Validate() error {
	if r == nil {
		return nil
	}
	return (*transactionDetail)(r).validate()
}

func (r *Detail) String(opts ...int64) string {
	if r == nil {
		return ""
	}

	return (*transactionDetail)(r).string(opts...)
}

func (r *Detail) Read(scan *util.Scanner, useCurrentLine bool) error {
	if scan == nil {
		return errors.New("invalid bai3 scanner")
	}

	var rawData string
	find := false
	isBreak := false

	for line := scan.ScanLine(useCurrentLine); line != ""; line = scan.ScanLine(useCurrentLine) {
		useCurrentLine = false

		// find record code
		if len(line) < 3 {
			continue
		}

		switch line[:2] {
		case util.TransactionDetailCode:

			if find {
				isBreak = true
				break
			}

			rawData = line
			find = true

		case util.ContinuationCode:
			if rawData == "" {
				return fmt.Errorf("ERROR parsing transaction detail on line %d (continuation without a preceding record)", scan.GetLineIndex())
			}
			if strings.HasSuffix(rawData, "/") {
				rawData = rawData[:len(rawData)-1] + "," + line[3:]
			} else {
				rawData = rawData + "," + line[3:]
			}

		default:
			isBreak = true
		}

		if isBreak {
			break
		}
	}

	_, err := (*transactionDetail)(r).parse(rawData)
	return err
}

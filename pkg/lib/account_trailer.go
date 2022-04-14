// Copyright 2022 The Moov Authors
// Use of this source code is governed by an Apache License
// license that can be found in the LICENSE file.

package lib

import (
	"bytes"
	"fmt"
	"unicode/utf8"

	"github.com/moov-io/bai2/pkg/util"
)

/*

Account Trailer

The account trailer record contains information on the account control total and the number of
records. The account control total is the sum of the amount values in the preceding type 03, 16 and 88
records associated with the account. The number of records is the total of all records in the account,
including the account trailer (type 49) record.

*/

const (
	accountTrailerLength = 32
)

// Creating Account Trailer
func NewAccountTrailer() *AccountTrailer {
	return &AccountTrailer{
		RecordCode: "49",
	}

}

// Account Trailer
type AccountTrailer struct {
	RecordCode          string
	AccountControlTotal string
	NumberRecords       int64
}

func (h *AccountTrailer) Validate() error {
	if h.RecordCode != "49" {
		return fmt.Errorf("AccountTrailer: invalid record code")
	}

	return nil
}

func (h *AccountTrailer) Parse(line string) error {
	if n := utf8.RuneCountInString(line); n < accountTrailerLength {
		return fmt.Errorf("AccountTrailer: length %d is too short", n)
	}

	h.RecordCode, _ = util.EntryParser(line[0:3], ",")
	h.AccountControlTotal, _ = util.EntryParser(line[3:22], ",")
	h.NumberRecords, _ = util.EntryParserToInt(line[22:32], "/")

	return nil
}

func (h *AccountTrailer) String() string {
	var buf bytes.Buffer

	buf.WriteString(fmt.Sprintf("%-2.2v,", h.RecordCode))
	buf.WriteString(fmt.Sprintf("%-18.18v,", h.AccountControlTotal))
	buf.WriteString(fmt.Sprintf("%09.9v/", h.NumberRecords))

	return buf.String()
}

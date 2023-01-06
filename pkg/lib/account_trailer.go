// Copyright 2022 The Moov Authors
// Use of this source code is governed by an Apache License
// license that can be found in the LICENSE file.

package lib

import (
	"bytes"
	"fmt"

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
	atParseErrorFmt    = "AccountTrailer: unable to parse %s"
	atValidateErrorFmt = "AccountTrailer: invalid %s"
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
		return fmt.Errorf(fmt.Sprintf(atValidateErrorFmt, "RecordCode"))
	}
	if h.AccountControlTotal != "" && !util.ValidateAmount(h.AccountControlTotal) {
		return fmt.Errorf(fmt.Sprintf(atValidateErrorFmt, "Amount"))
	}

	return nil
}

func (h *AccountTrailer) Parse(data string) (int, error) {

	var line string
	var err error
	var size, read int

	length := util.GetSize(data)
	if length < 2 {
		return 0, fmt.Errorf(fmt.Sprintf(atParseErrorFmt, "record"))
	} else {
		line = data[:length]
	}

	// RecordCode
	if h.RecordCode, size, err = util.ReadField(line[read:]); err != nil {
		return 0, fmt.Errorf(fmt.Sprintf(atParseErrorFmt, "RecordCode"))
	} else {
		read += size
	}

	// AccountControlTotal
	if h.AccountControlTotal, size, err = util.ReadField(line[read:]); err != nil {
		return 0, fmt.Errorf(fmt.Sprintf(atParseErrorFmt, "AccountControlTotal"))
	} else {
		read += size
	}

	// NumberRecords
	if h.NumberRecords, size, err = util.ReadFieldAsInt(line[read:]); err != nil {
		return 0, fmt.Errorf(fmt.Sprintf(atParseErrorFmt, "NumberRecords"))
	} else {
		read += size
	}

	if err = h.Validate(); err != nil {
		return 0, err
	}

	return read, nil
}

func (h *AccountTrailer) String() string {
	var buf bytes.Buffer

	buf.WriteString(fmt.Sprintf("%s,", h.RecordCode))
	buf.WriteString(fmt.Sprintf("%s,", h.AccountControlTotal))
	buf.WriteString(fmt.Sprintf("%d/", h.NumberRecords))

	return buf.String()
}

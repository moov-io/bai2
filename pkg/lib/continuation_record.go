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

 Continuation of Account Summary Record - Current (CDA) and Personal (PDA) Accounts

This record is a continuation of the account identifier and summary record (see Table 3) and contains
value date and summary information (i.e. total credits and debits as well as total credit and debit dollar
amounts). It always has a record code of 88.

*/

const (
	crParseErrorFmt    = "ContinuationRecord: unable to parse %s"
	crValidateErrorFmt = "ContinuationRecord: invalid %s"
)

// Creating Continuation Record
func NewContinuationRecord() *ContinuationRecord {
	return &ContinuationRecord{
		RecordCode: "88",
	}

}

// Continuation of Account Summary Record
type ContinuationRecord struct {
	RecordCode string
	Composite  []string `json:",omitempty"`
}

func (h *ContinuationRecord) Validate() error {
	if h.RecordCode != "88" {
		return fmt.Errorf(fmt.Sprintf(crValidateErrorFmt, "RecordCode"))
	}

	return nil
}

func (h *ContinuationRecord) Parse(data string) (int, error) {
	var line string
	var err error
	var size, read int

	length := util.GetSize(data)
	if length < 2 {
		return 0, fmt.Errorf(fmt.Sprintf(crParseErrorFmt, "record"))
	} else {
		line = data[:length]
	}

	// RecordCode
	if h.RecordCode, size, err = util.ReadField(util.GetField(line, read)); err != nil {
		return 0, fmt.Errorf(fmt.Sprintf(crParseErrorFmt, "RecordCode"))
	} else {
		read += size
	}

	for int64(read) < length {
		var composite string
		if composite, size, err = util.ReadField(util.GetField(line, read)); err != nil {
			return 0, fmt.Errorf(fmt.Sprintf(crParseErrorFmt, "ExtAmount1"))
		} else {
			read += size
		}
		h.Composite = append(h.Composite, composite)
	}

	if err = h.Validate(); err != nil {
		return 0, err
	}

	return read, nil
}

func (h *ContinuationRecord) String() string {
	var buf bytes.Buffer

	buf.WriteString(h.RecordCode)

	for _, composite := range h.Composite {
		buf.WriteString(fmt.Sprintf(",%s", composite))
	}

	buf.WriteString("/")
	return buf.String()
}

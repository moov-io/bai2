// Copyright 2022 The Moov Authors
// Use of this source code is governed by an Apache License
// license that can be found in the LICENSE file.

package bai3

import (
	"bytes"
	"fmt"

	"github.com/moov-io/bai2/pkg/util"
)

const (
	btParseErrorFmt    = "BankTrailer: unable to parse %s"
	btValidateErrorFmt = "BankTrailer: invalid %s"
)

type bankTrailer struct {
	BankControlTotal string
	NumberOfAccounts int64
	NumberOfRecords  int64
}

func (h *bankTrailer) validate() error {
	if h.BankControlTotal != "" && !util.ValidateAmount(h.BankControlTotal) {
		return fmt.Errorf(btValidateErrorFmt, "BankControlTotal")
	}

	return nil
}

func (h *bankTrailer) parse(data string) (int, error) {

	var line string
	var err error
	var size, read int

	if length := util.GetSize(data); length < 3 {
		return 0, fmt.Errorf(btParseErrorFmt, "record")
	} else {
		line = data[:length]
	}

	// RecordCode
	if util.GroupTrailerCode != data[:2] {
		return 0, fmt.Errorf(btParseErrorFmt, "RecordCode")
	}
	read += 3

	// BankControlTotal
	if h.BankControlTotal, size, err = util.ReadField(line, read); err != nil {
		return 0, fmt.Errorf(btParseErrorFmt, "BankControlTotal")
	} else {
		read += size
	}

	// NumberOfAccounts
	if h.NumberOfAccounts, size, err = util.ReadFieldAsInt(line, read); err != nil {
		return 0, fmt.Errorf(btParseErrorFmt, "NumberOfAccounts")
	} else {
		read += size
	}

	// NumberOfRecords
	if h.NumberOfRecords, size, err = util.ReadFieldAsInt(line, read); err != nil {
		return 0, fmt.Errorf(btParseErrorFmt, "NumberOfRecords")
	} else {
		read += size
	}

	if err = h.validate(); err != nil {
		return 0, err
	}

	return read, nil
}

func (h *bankTrailer) string() string {
	var buf bytes.Buffer

	buf.WriteString(fmt.Sprintf("%s,", util.GroupTrailerCode))
	buf.WriteString(fmt.Sprintf("%s,", h.BankControlTotal))
	buf.WriteString(fmt.Sprintf("%d,", h.NumberOfAccounts))
	buf.WriteString(fmt.Sprintf("%d/", h.NumberOfRecords))

	return buf.String()
}

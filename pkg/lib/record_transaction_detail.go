// Copyright 2022 The Moov Authors
// Use of this source code is governed by an Apache License
// license that can be found in the LICENSE file.

package lib

import (
	"bytes"
	"fmt"

	"github.com/moov-io/bai2/pkg/util"
)

const (
	tdParseErrorFmt    = "AccountTransaction: unable to parse %s"
	tdValidateErrorFmt = "AccountTransaction: invalid %s"
)

// Creating new transaction detail
func NewTransactionDetail() *TransactionDetail {
	return &TransactionDetail{}
}

type TransactionDetail struct {
	TypeCode  string   `json:",omitempty"`
	Amount    string   `json:",omitempty"`
	FundsType string   `json:",omitempty"`
	Composite []string `json:",omitempty"`
}

func (h *TransactionDetail) validate() error {
	if h.TypeCode != "" && !util.ValidateTypeCode(h.TypeCode) {
		return fmt.Errorf(fmt.Sprintf(tdValidateErrorFmt, "TypeCode"))
	}
	if h.Amount != "" && !util.ValidateAmount(h.Amount) {
		return fmt.Errorf(fmt.Sprintf(tdValidateErrorFmt, "Amount"))
	}
	if h.FundsType != "" && !util.ValidateFundsType(h.FundsType) {
		return fmt.Errorf(fmt.Sprintf(tdValidateErrorFmt, "FundsType"))
	}

	return nil
}

func (h *TransactionDetail) parse(data string) (int, error) {

	var line string
	var err error
	var size, read int

	length := util.GetSize(data)
	if length < 3 {
		return 0, fmt.Errorf(fmt.Sprintf(tdParseErrorFmt, "record"))
	} else {
		line = data[:length]
	}

	// RecordCode
	if util.TransactionDetailCode != data[:2] {
		return 0, fmt.Errorf(fmt.Sprintf(fhParseErrorFmt, "RecordCode"))
	}
	read += 3

	// TypeCode
	if h.TypeCode, size, err = util.ReadField(line, read); err != nil {
		return 0, fmt.Errorf(fmt.Sprintf(tdParseErrorFmt, "TypeCode"))
	} else {
		read += size
	}

	// Amount
	if h.Amount, size, err = util.ReadField(line, read); err != nil {
		return 0, fmt.Errorf(fmt.Sprintf(tdParseErrorFmt, "Amount"))
	} else {
		read += size
	}

	// FundsType
	if h.FundsType, size, err = util.ReadField(line, read); err != nil {
		return 0, fmt.Errorf(fmt.Sprintf(tdParseErrorFmt, "FundsType"))
	} else {
		read += size
	}

	for int64(read) < length {
		var composite string
		if composite, size, err = util.ReadField(line, read); err != nil {
			return 0, fmt.Errorf(fmt.Sprintf(tdParseErrorFmt, "Composite"))
		} else {
			read += size
		}
		h.Composite = append(h.Composite, composite)
	}

	if err = h.validate(); err != nil {
		return 0, err
	}

	return read, nil
}

func (h *TransactionDetail) string() string {

	var buf bytes.Buffer

	buf.WriteString(fmt.Sprintf("%s,", util.TransactionDetailCode))
	buf.WriteString(fmt.Sprintf("%s,", h.TypeCode))
	buf.WriteString(fmt.Sprintf("%s,", h.Amount))
	buf.WriteString(h.FundsType)

	for _, composite := range h.Composite {
		buf.WriteString(fmt.Sprintf(",%s", composite))
	}

	buf.WriteString("/")

	return buf.String()
}

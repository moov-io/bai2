// Copyright 2022 The Moov Authors
// Use of this source code is governed by an Apache License
// license that can be found in the LICENSE file.

package lib

import (
	"bytes"
	"fmt"
	"github.com/moov-io/bai2/pkg/util"
	"strings"
)

const (
	tdParseErrorFmt    = "TransactionDetail: unable to parse %s"
	tdValidateErrorFmt = "TransactionDetail: invalid %s"
)

type transactionDetail struct {
	TypeCode                string
	Amount                  string
	FundsType               FundsType
	BankReferenceNumber     string
	CustomerReferenceNumber string
	Text                    string
}

func (r *transactionDetail) validate() error {
	if r.TypeCode != "" && !util.ValidateTypeCode(r.TypeCode) {
		return fmt.Errorf(fmt.Sprintf(tdValidateErrorFmt, "TypeCode"))
	}
	if r.Amount != "" && !util.ValidateAmount(r.Amount) {
		return fmt.Errorf(fmt.Sprintf(tdValidateErrorFmt, "Amount"))
	}
	if r.FundsType.Validate() != nil {
		return fmt.Errorf(fmt.Sprintf(tdValidateErrorFmt, "FundsType"))
	}

	return nil
}

func (r *transactionDetail) parse(data string) (int, error) {

	var line string
	var err error
	var size, read int

	allow_slash_as_character := true
	length := util.GetSize(data, allow_slash_as_character)
	if length < 3 {
		return 0, fmt.Errorf(fmt.Sprintf(tdParseErrorFmt, "record"))
	} else {
		line = data[:length]
	}

	// RecordCode
	if util.TransactionDetailCode != data[:2] {
		return 0, fmt.Errorf(fmt.Sprintf(tdParseErrorFmt, "RecordCode"))
	}
	read += 3

	// TypeCode
	if r.TypeCode, size, err = util.ReadField(line, read); err != nil {
		return 0, fmt.Errorf(fmt.Sprintf(tdParseErrorFmt, "TypeCode"))
	} else {
		read += size
	}

	// Amount
	if r.Amount, size, err = util.ReadField(line, read); err != nil {
		return 0, fmt.Errorf(fmt.Sprintf(tdParseErrorFmt, "Amount"))
	} else {
		read += size
	}

	// FundsType
	if len(line) < read {
		return 0, fmt.Errorf(fmt.Sprintf(tdParseErrorFmt, "FundsType") + " too short")
	}
	if size, err = r.FundsType.parse(line[read:]); err != nil {
		return 0, fmt.Errorf(fmt.Sprintf(tdParseErrorFmt, "FundsType"))
	} else {
		read += size
	}

	// BankReferenceNumber
	if r.BankReferenceNumber, size, err = util.ReadField(line, read, allow_slash_as_character); err != nil {
		return 0, fmt.Errorf(fmt.Sprintf(tdParseErrorFmt, "BankReferenceNumber"))
	} else {
		read += size
	}

	// CustomerReferenceNumber
	if r.CustomerReferenceNumber, size, err = util.ReadField(line, read, allow_slash_as_character); err != nil {
		return 0, fmt.Errorf(fmt.Sprintf(tdParseErrorFmt, "CustomerReferenceNumber"))
	} else {
		read += size
	}

	// Text
	read_remainder_of_line := true
	if r.Text, size, err = util.ReadField(line, read, allow_slash_as_character, read_remainder_of_line); err != nil {
		return 0, fmt.Errorf(fmt.Sprintf(tdParseErrorFmt, "Text"))
	} else {
		read += size
	}

	if err = r.validate(); err != nil {
		return 0, err
	}

	return read, nil
}

func (r *transactionDetail) string(opts ...int64) string {

	var maxLen int64
	if len(opts) > 0 {
		maxLen = opts[0]
	}

	var total, buf bytes.Buffer

	buf.WriteString(fmt.Sprintf("%s,", util.TransactionDetailCode))
	buf.WriteString(fmt.Sprintf("%s,", r.TypeCode))
	buf.WriteString(fmt.Sprintf("%s,", r.Amount))

	util.WriteBuffer(&total, &buf, r.FundsType.String(), maxLen)
	buf.WriteString(",")

	util.WriteBuffer(&total, &buf, r.BankReferenceNumber, maxLen)
	buf.WriteString(",")

	util.WriteBuffer(&total, &buf, r.CustomerReferenceNumber, maxLen)
	buf.WriteString(",")

	util.WriteBuffer(&total, &buf, r.Text, maxLen)
	if !strings.HasSuffix(r.Text, "/") {
		buf.WriteString("/")
	}

	total.WriteString(buf.String())

	return total.String()
}

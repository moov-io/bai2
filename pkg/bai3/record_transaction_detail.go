// Copyright 2022 The Moov Authors
// Use of this source code is governed by an Apache License
// license that can be found in the LICENSE file.

package bai3

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/moov-io/bai2/pkg/util"
)

const (
	tdParseErrorFmt    = "TransactionDetail: unable to parse %s"
	tdValidateErrorFmt = "TransactionDetail: invalid %s"
)

type transactionDetail struct {
	TypeCode                string          `json:"TypeCode"`
	Amount                  string          `json:"Amount"`
	FundsType               FundsType       `json:"FundsType"`
	BankReferenceNumber     string          `json:"BankReferenceNumber"`
	CustomerReferenceNumber string          `json:"CustomerReferenceNumber"`
	Text                    string          `json:"Text"`
	Transaction             TransactionKind `json:"Transaction,omitempty"`
	Level                   TypeLevel       `json:"Level,omitempty"`
	Description             string          `json:"Description,omitempty"`
}

// TypeInfo returns the Appendix A catalog entry for this detail's type code.
func (r *Detail) TypeInfo() (TypeCode, bool) {
	if r == nil {
		return TypeCode{}, false
	}
	return LookupTypeCode(r.TypeCode)
}

func (r *transactionDetail) validate() error {
	if r.TypeCode != "" && !util.ValidateTypeCode(r.TypeCode) {
		return fmt.Errorf(tdValidateErrorFmt, "TypeCode")
	}
	if r.Amount != "" && !util.ValidateAmount(r.Amount) {
		return fmt.Errorf(tdValidateErrorFmt, "Amount")
	}
	if r.FundsType.Validate() != nil {
		return fmt.Errorf(tdValidateErrorFmt, "FundsType")
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
		return 0, fmt.Errorf(tdParseErrorFmt, "record")
	} else {
		line = data[:length]
	}

	// RecordCode
	if util.TransactionDetailCode != data[:2] {
		return 0, fmt.Errorf(tdParseErrorFmt, "RecordCode")
	}
	read += 3

	// TypeCode
	if r.TypeCode, size, err = util.ReadField(line, read); err != nil {
		return 0, fmt.Errorf(tdParseErrorFmt, "TypeCode")
	} else {
		read += size
	}

	// Amount
	if r.Amount, size, err = util.ReadField(line, read); err != nil {
		return 0, fmt.Errorf(tdParseErrorFmt, "Amount")
	} else {
		read += size
	}

	// FundsType
	if len(line) < read {
		return 0, fmt.Errorf(tdParseErrorFmt+" too short", "FundsType")
	}
	if size, err = r.FundsType.parse(line[read:]); err != nil {
		return 0, fmt.Errorf(tdParseErrorFmt, "FundsType")
	} else {
		read += size
	}

	// BankReferenceNumber
	if r.BankReferenceNumber, size, err = util.ReadField(line, read, allow_slash_as_character); err != nil {
		return 0, fmt.Errorf(tdParseErrorFmt, "BankReferenceNumber")
	} else {
		read += size
	}

	// CustomerReferenceNumber
	if r.CustomerReferenceNumber, size, err = util.ReadField(line, read, allow_slash_as_character); err != nil {
		return 0, fmt.Errorf(tdParseErrorFmt, "CustomerReferenceNumber")
	} else {
		read += size
	}

	// Text
	read_remainder_of_line := true
	if r.Text, size, err = util.ReadField(line, read, allow_slash_as_character, read_remainder_of_line); err != nil {
		return 0, fmt.Errorf(tdParseErrorFmt, "Text")
	} else {
		read += size
	}

	// Banks often terminate type 16 records with `/` even though the spec says
	// text records are delimited by the next non-88 record. Strip a single trailing
	// record delimiter; slashes inside the text are kept.
	r.Text = strings.TrimSuffix(r.Text, "/")
	r.Transaction, r.Level, r.Description = typeMeta(r.TypeCode)

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
	if r.Text == "" {
		// Defaulted text field: adjacent delimiters `,/`
		buf.WriteString("/")
	}

	total.WriteString(buf.String())

	return total.String()
}

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

const (
	accountIdentifierContinuationLength = 75
)

func NewAccountIdentifierContinuation() *AccountIdentifierContinuation {
	return &AccountIdentifierContinuation{
		RecordCode: "88",
		TypeCode1:  "100",
		TypeCode2:  "400",
	}

}

type AccountIdentifierContinuation struct {
	RecordCode         string
	TypeCode1          string
	TotalCreditAmount1 string
	TotalOfCredits1    int64
	FundsType1         string
	ValueDate1         string
	TypeCode2          string
	TotalCreditAmount2 string
	TotalOfCredits2    int64
	FundsType2         string
	ValueDate2         string
}

func (h *AccountIdentifierContinuation) Validate() error {
	if h.RecordCode != "88" {
		return fmt.Errorf("AccountIdentifierContinuation: invalid record code")
	}
	if h.TypeCode1 != "100" {
		return fmt.Errorf("AccountIdentifierContinuation: invalid type code")
	}
	if h.TypeCode2 != "400" {
		return fmt.Errorf("AccountIdentifierContinuation: invalid type code")
	}

	return nil
}

func (h *AccountIdentifierContinuation) Parse(line string) error {
	if n := utf8.RuneCountInString(line); n < accountIdentifierContinuationLength {
		return fmt.Errorf("AccountIdentifierContinuation: length %d is too short", n)
	}

	h.RecordCode, _ = util.EntryParser(line[0:3], ",")
	h.TypeCode1, _ = util.EntryParser(line[3:7], ",")
	h.TotalCreditAmount1, _ = util.EntryParser(line[7:23], ",")
	h.TotalOfCredits1, _ = util.EntryParserToInt(line[23:29], ",")
	h.FundsType1, _ = util.EntryParser(line[29:31], ",")
	h.ValueDate1, _ = util.EntryParser(line[31:38], ",")
	h.TypeCode2, _ = util.EntryParser(line[39:43], ",")
	h.TotalCreditAmount2, _ = util.EntryParser(line[43:59], ",")
	h.TotalOfCredits2, _ = util.EntryParserToInt(line[59:65], ",")
	h.FundsType2, _ = util.EntryParser(line[65:67], ",")
	h.ValueDate2, _ = util.EntryParser(line[67:74], ",")

	return nil
}

func (h *AccountIdentifierContinuation) String() string {
	var buf bytes.Buffer

	buf.WriteString(fmt.Sprintf("%-2.2v,", h.RecordCode))
	buf.WriteString(fmt.Sprintf("%-3.3v,", h.TypeCode1))
	buf.WriteString(fmt.Sprintf("%-15.15v,", h.TotalCreditAmount1))
	buf.WriteString(fmt.Sprintf("%05.5v,", h.TotalOfCredits1))
	buf.WriteString(fmt.Sprintf("%-1.1v,", h.FundsType1))
	buf.WriteString(fmt.Sprintf("%-6.6v,", h.ValueDate1))
	buf.WriteString(",")
	buf.WriteString(fmt.Sprintf("%-3.3v,", h.TypeCode2))
	buf.WriteString(fmt.Sprintf("%-15.15v,", h.TotalCreditAmount2))
	buf.WriteString(fmt.Sprintf("%05.5v,", h.TotalOfCredits2))
	buf.WriteString(fmt.Sprintf("%-1.1v,", h.FundsType2))
	buf.WriteString(fmt.Sprintf("%-6.6v,", h.ValueDate2))
	buf.WriteString("/")

	return buf.String()
}

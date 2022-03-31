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
	accountIdentifierLoanLength = 79
)

func NewAccountIdentifierLoan() *AccountIdentifierLoan {
	return &AccountIdentifierLoan{
		RecordCode:   "03",
		CurrencyCode: "USD",
		TypeCode1:    "056",
		TypeCode2:    "056",
	}

}

type AccountIdentifierLoan struct {
	RecordCode     string
	AccountNumber  string
	CurrencyCode   string
	TypeCode1      string
	OpeningBalance string
	FundsType1     string
	ValueDate1     string
	TypeCode2      string
	ClosingBalance string
	FundsType2     string
	ValueDate2     string
}

func (h *AccountIdentifierLoan) Validate() error {
	if h.RecordCode != "03" {
		return fmt.Errorf("AccountIdentifierLoan: invalid record code")
	}
	if h.CurrencyCode != "USD" && h.CurrencyCode != "CAD" {
		return fmt.Errorf("AccountIdentifierLoan: invalid currency code")
	}
	if h.TypeCode1 != "056" {
		return fmt.Errorf("AccountIdentifierLoan: invalid type code")
	}
	if h.TypeCode2 != "056" {
		return fmt.Errorf("AccountIdentifierLoan: invalid type code")
	}

	return nil
}

func (h *AccountIdentifierLoan) Parse(line string) error {
	if n := utf8.RuneCountInString(line); n < accountIdentifierLoanLength {
		return fmt.Errorf("AccountIdentifierLoan: length %d is too short", n)
	}

	h.RecordCode, _ = util.EntryParser(line[0:3], ",")
	h.AccountNumber, _ = util.EntryParser(line[3:17], ",")
	h.CurrencyCode, _ = util.EntryParser(line[17:21], ",")
	h.TypeCode1, _ = util.EntryParser(line[21:25], ",")
	h.OpeningBalance, _ = util.EntryParser(line[25:39], ",")
	h.FundsType1, _ = util.EntryParser(line[40:42], ",")
	h.ValueDate1, _ = util.EntryParser(line[42:49], ",")
	h.TypeCode2, _ = util.EntryParser(line[50:54], ",")
	h.ClosingBalance, _ = util.EntryParser(line[54:68], ",")
	h.FundsType2, _ = util.EntryParser(line[69:71], ",")
	h.ValueDate2, _ = util.EntryParser(line[71:78], ",")

	return nil
}

func (h *AccountIdentifierLoan) String() string {
	var buf bytes.Buffer

	buf.WriteString(fmt.Sprintf("%2.2v,", h.RecordCode))
	buf.WriteString(fmt.Sprintf("%13.13v,", h.AccountNumber))
	buf.WriteString(fmt.Sprintf("%3.3v,", h.CurrencyCode))
	buf.WriteString(fmt.Sprintf("%3.3v,", h.TypeCode1))
	buf.WriteString(fmt.Sprintf("%13.13v,", h.OpeningBalance))
	buf.WriteString(fmt.Sprintf(","))
	buf.WriteString(fmt.Sprintf("%1.1v,", h.FundsType1))
	buf.WriteString(fmt.Sprintf("%6.6v,", h.ValueDate1))
	buf.WriteString(fmt.Sprintf(","))
	buf.WriteString(fmt.Sprintf("%3.3v,", h.TypeCode2))
	buf.WriteString(fmt.Sprintf("%13.13v,", h.ClosingBalance))
	buf.WriteString(fmt.Sprintf(","))
	buf.WriteString(fmt.Sprintf("%1.1v,", h.FundsType2))
	buf.WriteString(fmt.Sprintf("%6.6v,", h.ValueDate2))
	buf.WriteString(fmt.Sprintf("/"))

	return buf.String()
}

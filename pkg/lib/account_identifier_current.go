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
	accountIdentifierCurrentLength = 59
)

/*

Account Identifier and Summary/Status for Current (CDA), Personal (PDA), and Loan accounts

CDA and PDA accounts
This record contains information on opening and closing balances for CDA and PDA accounts. It always
has a record code of 03.

*/

// Creating CDA and PDA accounts
func NewAccountIdentifierCurrent() *AccountIdentifierCurrent {
	return &AccountIdentifierCurrent{
		RecordCode:   "03",
		CurrencyCode: "USD",
		TypeCode1:    "040",
		TypeCode2:    "045",
	}

}

// CDA and PDA accounts
type AccountIdentifierCurrent struct {
	RecordCode     string
	AccountNumber  string
	CurrencyCode   string
	TypeCode1      string
	OpeningBalance string
	TypeCode2      string
	ClosingBalance string
}

func (h *AccountIdentifierCurrent) Validate() error {
	if h.RecordCode != "03" {
		return fmt.Errorf("AccountIdentifierCurrent: invalid record code")
	}
	if h.CurrencyCode != "USD" && h.CurrencyCode != "CAD" {
		return fmt.Errorf("AccountIdentifierCurrent: invalid currency code")
	}
	if h.TypeCode1 != "040" {
		return fmt.Errorf("AccountIdentifierCurrent: invalid type code")
	}
	if h.TypeCode2 != "045" {
		return fmt.Errorf("AccountIdentifierCurrent: invalid type code")
	}

	return nil
}

func (h *AccountIdentifierCurrent) Parse(line string) error {
	if n := utf8.RuneCountInString(line); n < accountIdentifierCurrentLength {
		return fmt.Errorf("AccountIdentifierCurrent: length %d is too short", n)
	}

	h.RecordCode, _ = util.EntryParser(line[0:3], ",")
	h.AccountNumber, _ = util.EntryParser(line[3:15], ",")
	h.CurrencyCode, _ = util.EntryParser(line[15:19], ",")
	h.TypeCode1, _ = util.EntryParser(line[19:23], ",")
	h.OpeningBalance, _ = util.EntryParser(line[23:37], ",")
	h.TypeCode2, _ = util.EntryParser(line[39:43], ",")
	h.ClosingBalance, _ = util.EntryParser(line[43:57], ",")

	return nil
}

func (h *AccountIdentifierCurrent) String() string {
	var buf bytes.Buffer

	buf.WriteString(fmt.Sprintf("%-2.2v,", h.RecordCode))
	buf.WriteString(fmt.Sprintf("%-11.11v,", h.AccountNumber))
	buf.WriteString(fmt.Sprintf("%-3.3v,", h.CurrencyCode))
	buf.WriteString(fmt.Sprintf("%-3.3v,", h.TypeCode1))
	buf.WriteString(fmt.Sprintf("%-13.13v,", h.OpeningBalance))
	buf.WriteString(",")
	buf.WriteString(",")
	buf.WriteString(fmt.Sprintf("%-3.3v,", h.TypeCode2))
	buf.WriteString(fmt.Sprintf("%-13.13v,", h.ClosingBalance))
	buf.WriteString(",")
	buf.WriteString("/")

	return buf.String()
}

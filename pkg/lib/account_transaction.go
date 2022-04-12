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
	accountTransactionLength = 56
)

func NewAccountTransaction() *AccountTransaction {
	return &AccountTransaction{
		RecordCode: "16",
		TypeCode:   "108",
	}

}

type AccountTransaction struct {
	RecordCode  string
	TypeCode    string
	Amount      string
	FundsType   string
	ValueDate   string
	Description string
}

func (h *AccountTransaction) Validate() error {
	if h.RecordCode != "16" {
		return fmt.Errorf("AccountTransaction: invalid record code")
	}
	if h.TypeCode != "108" && h.TypeCode != "409" {
		return fmt.Errorf("AccountTransaction: invalid type code")
	}

	return nil
}

func (h *AccountTransaction) Parse(line string) error {
	if n := utf8.RuneCountInString(line); n < accountTransactionLength {
		return fmt.Errorf("AccountTransaction: length %d is too short", n)
	}

	h.RecordCode, _ = util.EntryParser(line[0:3], ",")
	h.TypeCode, _ = util.EntryParser(line[3:7], ",")
	h.Amount, _ = util.EntryParser(line[7:23], ",")
	h.FundsType, _ = util.EntryParser(line[23:25], ",")
	h.ValueDate, _ = util.EntryParser(line[25:32], ",")
	h.Description, _ = util.EntryParser(line[35:56], "/")

	return nil
}

func (h *AccountTransaction) String() string {
	var buf bytes.Buffer

	buf.WriteString(fmt.Sprintf("%-2.2v,", h.RecordCode))
	buf.WriteString(fmt.Sprintf("%-3.3v,", h.TypeCode))
	buf.WriteString(fmt.Sprintf("%-15.15v,", h.Amount))
	buf.WriteString(fmt.Sprintf("%-1.1v,", h.FundsType))
	buf.WriteString(fmt.Sprintf("%-6.6v,", h.ValueDate))
	buf.WriteString(",")
	buf.WriteString(",")
	buf.WriteString(",")
	buf.WriteString(fmt.Sprintf("%-20.20v/", h.Description))

	return buf.String()
}

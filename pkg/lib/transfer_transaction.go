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
	transferTransactionLength = 68
)

func NewTransferTransaction() *TransferTransaction {
	return &TransferTransaction{
		RecordCode: "16",
		TypeCode:   "108",
	}

}

type TransferTransaction struct {
	RecordCode        string
	TypeCode          string
	Amount            string
	FundsType         string
	ValueDate         string
	BankReference     string
	CustomerReference string
	Description       string
}

func (h *TransferTransaction) Validate() error {
	if h.RecordCode != "16" {
		return fmt.Errorf("TransferTransaction: invalid record code")
	}
	if h.TypeCode != "108" && h.TypeCode != "409" {
		return fmt.Errorf("TransferTransaction: invalid type code")
	}

	return nil
}

func (h *TransferTransaction) Parse(line string) error {
	if n := utf8.RuneCountInString(line); n < transferTransactionLength {
		return fmt.Errorf("TransferTransaction: length %d is too short", n)
	}

	h.RecordCode, _ = util.EntryParser(line[0:3], ",")
	h.TypeCode, _ = util.EntryParser(line[3:7], ",")
	h.Amount, _ = util.EntryParser(line[7:21], ",")
	h.FundsType, _ = util.EntryParser(line[21:23], ",")
	h.ValueDate, _ = util.EntryParser(line[23:30], ",")
	h.BankReference, _ = util.EntryParser(line[31:37], ",")
	h.CustomerReference, _ = util.EntryParser(line[37:47], ",")
	h.Description, _ = util.EntryParser(line[47:68], "/")

	return nil
}

func (h *TransferTransaction) String() string {
	var buf bytes.Buffer

	buf.WriteString(fmt.Sprintf("%-2.2v,", h.RecordCode))
	buf.WriteString(fmt.Sprintf("%-3.3v,", h.TypeCode))
	buf.WriteString(fmt.Sprintf("%-13.13v,", h.Amount))
	buf.WriteString(fmt.Sprintf("%-1.1v,", h.FundsType))
	buf.WriteString(fmt.Sprintf("%-6.6v,", h.ValueDate))
	buf.WriteString(",")
	buf.WriteString(fmt.Sprintf("%-5.5v,", h.BankReference))
	buf.WriteString(fmt.Sprintf("%-9.9v,", h.CustomerReference))
	buf.WriteString(fmt.Sprintf("%-20.20v/", h.Description))

	return buf.String()
}

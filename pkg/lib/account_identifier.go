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
	aiParseErrorFmt    = "AccountIdentifier: unable to parse %s"
	aiValidateErrorFmt = "AccountIdentifierCurrent: invalid %s"
)

/*

Account Identifier and Summary/Status for Current (CDA), Personal (PDA), and Loan accounts

CDA and PDA accounts
This record contains information on opening and closing balances for CDA and PDA accounts. It always
has a record code of 03.

*/

// Creating CDA and PDA accounts
func NewAccountIdentifier() *AccountIdentifier {
	return &AccountIdentifier{
		RecordCode: "03",
	}
}

// CDA and PDA accounts
type AccountIdentifier struct {
	RecordCode    string
	AccountNumber string
	CurrencyCode  string   `json:",omitempty"`
	TypeCode      string   `json:",omitempty"`
	Amount        string   `json:",omitempty"`
	ItemCount     string   `json:",omitempty"`
	FundsType     string   `json:",omitempty"`
	Composite     []string `json:",omitempty"`
}

func (h *AccountIdentifier) Validate() error {
	if h.RecordCode != "03" {
		return fmt.Errorf(fmt.Sprintf(aiValidateErrorFmt, "RecordCode"))
	}
	if h.AccountNumber == "" {
		return fmt.Errorf(fmt.Sprintf(aiValidateErrorFmt, "AccountNumber"))
	}
	if h.Amount != "" && !util.ValidateAmount(h.Amount) {
		return fmt.Errorf(fmt.Sprintf(aiValidateErrorFmt, "Amount"))
	}
	if h.CurrencyCode != "" && !util.ValidateCurrencyCode(h.CurrencyCode) {
		return fmt.Errorf(fmt.Sprintf(aiValidateErrorFmt, "CurrencyCode"))
	}
	if h.TypeCode != "" && !util.ValidateTypeCode(h.TypeCode) {
		return fmt.Errorf(fmt.Sprintf(aiValidateErrorFmt, "TypeCode"))
	}
	if h.FundsType != "" && !util.ValidateFundsType(h.FundsType) {
		return fmt.Errorf(fmt.Sprintf(aiValidateErrorFmt, "FundsType"))
	}

	return nil
}

func (h *AccountIdentifier) Parse(data string) (int, error) {

	var line string
	var err error
	var size, read int

	length := util.GetSize(data)
	if length < 2 {
		return 0, fmt.Errorf(fmt.Sprintf(aiParseErrorFmt, "record"))
	} else {
		line = data[:length]
	}

	// RecordCode
	if h.RecordCode, size, err = util.ReadField(line[read:]); err != nil {
		return 0, fmt.Errorf(fmt.Sprintf(aiParseErrorFmt, "RecordCode"))
	} else {
		read += size
	}

	// AccountNumber
	if h.AccountNumber, size, err = util.ReadField(line[read:]); err != nil {
		return 0, fmt.Errorf(fmt.Sprintf(aiParseErrorFmt, "AccountNumber"))
	} else {
		read += size
	}

	// CurrencyCode
	if h.CurrencyCode, size, err = util.ReadField(line[read:]); err != nil {
		return 0, fmt.Errorf(fmt.Sprintf(aiParseErrorFmt, "CurrencyCode"))
	} else {
		read += size
	}

	// TypeCode
	if h.TypeCode, size, err = util.ReadField(line[read:]); err != nil {
		return 0, fmt.Errorf(fmt.Sprintf(aiParseErrorFmt, "TypeCode"))
	} else {
		read += size
	}

	// Amount
	if h.Amount, size, err = util.ReadField(line[read:]); err != nil {
		return 0, fmt.Errorf(fmt.Sprintf(aiParseErrorFmt, "Amount"))
	} else {
		read += size
	}

	// ItemCount
	if h.ItemCount, size, err = util.ReadField(line[read:]); err != nil {
		return 0, fmt.Errorf(fmt.Sprintf(aiParseErrorFmt, "ItemCount"))
	} else {
		read += size
	}

	// FundsType
	if h.FundsType, size, err = util.ReadField(line[read:]); err != nil {
		return 0, fmt.Errorf(fmt.Sprintf(aiParseErrorFmt, "FundsType"))
	} else {
		read += size
	}

	for int64(read) < length {
		var composite string
		if composite, size, err = util.ReadField(line[read:]); err != nil {
			return 0, fmt.Errorf(fmt.Sprintf(aiParseErrorFmt, "ExtAmount1"))
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

func (h *AccountIdentifier) String() string {
	var buf bytes.Buffer

	buf.WriteString(fmt.Sprintf("%s,", h.RecordCode))
	buf.WriteString(fmt.Sprintf("%s,", h.AccountNumber))
	buf.WriteString(fmt.Sprintf("%s,", h.CurrencyCode))
	buf.WriteString(fmt.Sprintf("%s,", h.TypeCode))
	buf.WriteString(fmt.Sprintf("%s,", h.Amount))
	buf.WriteString(fmt.Sprintf("%s,", h.ItemCount))
	buf.WriteString(fmt.Sprintf("%s", h.FundsType))

	for _, composite := range h.Composite {
		buf.WriteString(fmt.Sprintf(",%s", composite))
	}

	buf.WriteString("/")
	return buf.String()
}

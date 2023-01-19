// Copyright 2022 The Moov Authors
// Use of this source code is governed by an Apache License
// license that can be found in the LICENSE file.

package lib

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/moov-io/bai2/pkg/util"
)

/*

FILE STRUCTURE

To simplify processing, balance reporting transmission files are divided into “envelopes” of data.
These envelopes organize data at the following levels:
• Account
• Group
• File

Account:
	The first level of organization is the account. An account envelope includes balance and transaction data.
	Example: Account #1256793 at Last National Bank, previous day information as of midnight.

*/

// Creating new account object
func NewAccount() *Account {
	return &Account{}
}

// Account Format
type Account struct {
	// Account Identifier
	AccountNumber string   `json:"accountNumber"`
	CurrencyCode  string   `json:"currencyCode,omitempty"`
	TypeCode      string   `json:"typeCode,omitempty"`
	Amount        string   `json:"amount,omitempty"`
	ItemCount     int64    `json:"itemCount,omitempty"`
	FundsType     string   `json:"fundsType,omitempty"`
	Composite     []string `json:"composite,omitempty"`

	// Account Trailer
	AccountControlTotal string `json:"accountControlTotal"`
	NumberRecords       int64  `json:"numberRecords"`

	Details []TransactionDetail

	header  accountIdentifier
	trailer accountTrailer
}

func (r *Account) copyRecords() {

	r.header = accountIdentifier{
		AccountNumber: r.AccountNumber,
		CurrencyCode:  r.CurrencyCode,
		TypeCode:      r.TypeCode,
		Amount:        r.Amount,
		ItemCount:     r.ItemCount,
		FundsType:     r.FundsType,
		Composite:     r.Composite,
	}

	r.trailer = accountTrailer{
		AccountControlTotal: r.AccountControlTotal,
		NumberRecords:       r.NumberRecords,
	}

}

func (r *Account) String() string {

	r.copyRecords()

	var buf bytes.Buffer
	buf.WriteString(r.header.string() + "\n")
	for i := range r.Details {
		buf.WriteString(r.Details[i].string() + "\n")
	}
	buf.WriteString(r.trailer.string())

	return buf.String()
}

func (r *Account) Validate() error {

	r.copyRecords()

	if err := r.header.validate(); err != nil {
		return err
	}

	for i := range r.Details {
		if err := r.Details[i].validate(); err != nil {
			return err
		}
	}

	if err := r.trailer.validate(); err != nil {
		return err
	}

	return nil
}

func (r *Account) Read(scan Bai2Scanner, input string, lineNum int) (int, string, error) {

	var detail *TransactionDetail

	for line := scan.ScanLine(input); line != ""; line = scan.ScanLine(input) {

		input = ""

		// don't expect new line
		line = strings.TrimSpace(strings.ReplaceAll(line, "\n", ""))

		// find record code
		if len(line) < 3 {
			lineNum++
			continue
		}

		switch line[:2] {
		case util.AccountIdentifierCode:

			lineNum++
			newRecord := accountIdentifier{}
			_, err := newRecord.parse(line)
			if err != nil {
				return lineNum, line, fmt.Errorf("ERROR parsing account identifier on line %d - %v", lineNum, err)
			}

			r.AccountNumber = newRecord.AccountNumber
			r.CurrencyCode = newRecord.CurrencyCode
			r.TypeCode = newRecord.TypeCode
			r.Amount = newRecord.Amount
			r.ItemCount = newRecord.ItemCount
			r.FundsType = newRecord.FundsType
			r.Composite = newRecord.Composite

		case util.AccountTrailerCode:

			lineNum++
			newRecord := accountTrailer{}
			_, err := newRecord.parse(line)
			if err != nil {
				return lineNum, line, fmt.Errorf("ERROR parsing account trailer on line %d - %v", lineNum, err)
			}

			r.AccountControlTotal = newRecord.AccountControlTotal
			r.NumberRecords = newRecord.NumberRecords

			if detail != nil {
				r.Details = append(r.Details, *detail)
			}

			return lineNum, "", nil

		case util.TransactionDetailCode:

			lineNum++
			if detail != nil {
				r.Details = append(r.Details, *detail)
				detail = nil
			}

			detail = NewTransactionDetail()
			_, err := detail.parse(line)
			if err != nil {
				return lineNum, line, fmt.Errorf("ERROR parsing transaction detail on line %d - %v", lineNum, err)
			}

		case util.ContinuationCode:

			lineNum++
			newRecord := continuationRecord{}
			_, err := newRecord.parse(line)
			if err != nil {
				return lineNum, line, fmt.Errorf("ERROR parsing continuation on line %d - %v", lineNum, err)
			}

			if detail == nil {
				r.Composite = append(r.Composite, newRecord.Composite...)
			} else {
				detail.Composite = append(detail.Composite, newRecord.Composite...)
			}

		default:

			return lineNum, line, nil

		}
	}

	return lineNum, "", nil
}

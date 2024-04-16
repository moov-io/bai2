// Copyright 2022 The Moov Authors
// Use of this source code is governed by an Apache License
// license that can be found in the LICENSE file.

package lib

import (
	"bytes"
	"errors"
	"fmt"
	"regexp"
	"strconv"

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
	AccountNumber string           `json:"accountNumber"`
	CurrencyCode  string           `json:"currencyCode,omitempty"`
	Summaries     []AccountSummary `json:"summaries,omitempty"`

	// Account Trailer
	AccountControlTotal string `json:"accountControlTotal"`
	NumberRecords       int64  `json:"numberRecords"`

	Details []Detail

	header  accountIdentifier
	trailer accountTrailer
}

func (r *Account) copyRecords() {

	r.header = accountIdentifier{
		AccountNumber: r.AccountNumber,
		CurrencyCode:  r.CurrencyCode,
		Summaries:     r.Summaries,
	}

	r.trailer = accountTrailer{
		AccountControlTotal: r.AccountControlTotal,
		NumberRecords:       r.NumberRecords,
	}

}

var accountIdentifierCountExpression = regexp.MustCompile(`(?m:^(?:(?:03)|(?:16)|(?:49)|(?:88)))`)

// Sums the number of 03,16,88,49 records in the account. Maps to the NumberRecords field
func (a *Account) SumRecords(opts ...int64) int64 {
	acctString := a.String(opts...)
	result := accountIdentifierCountExpression.FindAllStringSubmatch(acctString, 10000)
	return int64(len(result))
}

// Sums the Amount fields from all 03 and 16 records. Maps to the AccountControlTotal field
func (a *Account) SumDetailAmounts() (string, error) {
	if err := a.Validate(); err != nil {
		return "0", err
	}
	var sum int64
	for _, detail := range a.Details {
		amt, err := strconv.ParseInt(detail.Amount, 10, 64)
		if err != nil {
			return "0", err
		}
		switch string(detail.TypeCode[0]) {
		case "1", "2", "3":
			sum += amt

		case "4", "5", "6":
			sum -= amt
		default:
			return "0", fmt.Errorf("TypeCode %v is invalid for transaction detail", detail.TypeCode)
		}
	}
	for _, summary := range a.Summaries {
		amt, err := strconv.ParseInt(summary.Amount, 10, 64)
		if err != nil {
			return "0", err
		}
		sum += amt
	}
	return fmt.Sprint(sum), nil
}

func (r *Account) String(opts ...int64) string {

	r.copyRecords()

	var buf bytes.Buffer
	buf.WriteString(r.header.string(opts...) + "\n")
	for i := range r.Details {
		buf.WriteString(r.Details[i].String(opts...) + "\n")
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
		if err := r.Details[i].Validate(); err != nil {
			return err
		}
	}

	if err := r.trailer.validate(); err != nil {
		return err
	}

	return nil
}

func (r *Account) Read(scan *Bai2Scanner, useCurrentLine bool) error {
	if scan == nil {
		return errors.New("invalid bai2 scanner")
	}

	parseAccountIdentifier := func(raw string) error {
		if raw == "" {
			return nil
		}

		newRecord := accountIdentifier{}
		_, err := newRecord.parse(raw)
		if err != nil {
			return fmt.Errorf("ERROR parsing account identifier on line %d (%v)", scan.GetLineIndex(), err)
		}

		r.AccountNumber = newRecord.AccountNumber
		r.CurrencyCode = newRecord.CurrencyCode
		r.Summaries = newRecord.Summaries

		return nil
	}

	var rawData string
	find := false

	for line := scan.ScanLine(useCurrentLine); line != ""; line = scan.ScanLine(useCurrentLine) {
		// find record code
		if len(line) < 3 {
			continue
		}

		useCurrentLine = false
		switch line[:2] {
		case util.AccountIdentifierCode:
			if find {
				break
			}

			rawData = line
			find = true

		case util.ContinuationCode:
			if len(rawData) > 0 {
				rawData = rawData[:len(rawData)-1] + "," + line[3:]
			}

		case util.AccountTrailerCode:
			if err := parseAccountIdentifier(rawData); err != nil {
				return err
			}

			newRecord := accountTrailer{}
			_, err := newRecord.parse(line)
			if err != nil {
				return fmt.Errorf("ERROR parsing account trailer on line %d (%v)", scan.GetLineIndex(), err)
			}

			r.AccountControlTotal = newRecord.AccountControlTotal
			r.NumberRecords = newRecord.NumberRecords

			return nil

		case util.TransactionDetailCode:
			if err := parseAccountIdentifier(rawData); err != nil {
				return err
			} else {
				rawData = ""
			}

			detail := NewDetail()
			err := detail.Read(scan, true)
			if err != nil {
				return err
			}

			r.Details = append(r.Details, *detail)
			useCurrentLine = true
		default:
			return fmt.Errorf("ERROR parsing account on line %d (unable to read record type %s)", scan.GetLineIndex(), line[0:2])

		}
	}

	return nil
}

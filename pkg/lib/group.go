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

Group:
	The next level of organization is the group. A group includes one or more account envelopes, all of which represent accounts at the same financial institution.
	All information in a group is for the same date and time.
	Example: Several accounts from Last National Bank to XYZ Reporting Service, sameday information as of 9:00 AM.

*/

// Creating new group object
func NewGroup() *Group {
	return &Group{}
}

// Group Format
type Group struct {
	// Group Header
	Receiver         string `json:"receiver,omitempty"`
	Originator       string `json:"originator"`
	GroupStatus      int64  `json:"groupStatus"`
	AsOfDate         string `json:"asOfDate"`
	AsOfTime         string `json:"asOfTime,omitempty"`
	CurrencyCode     string `json:"currencyCode,omitempty"`
	AsOfDateModifier int64  `json:"asOfDateModifier,omitempty"`

	// Group Trailer
	GroupControlTotal string `json:"groupControlTotal"`
	NumberOfAccounts  int64  `json:"numberOfAccounts"`
	NumberOfRecords   int64  `json:"numberOfRecords"`

	Accounts []Account

	header  groupHeader
	trailer groupTrailer
}

func (r *Group) copyRecords() {

	r.header = groupHeader{
		Receiver:         r.Receiver,
		Originator:       r.Originator,
		GroupStatus:      r.GroupStatus,
		AsOfDate:         r.AsOfDate,
		AsOfTime:         r.AsOfTime,
		CurrencyCode:     r.CurrencyCode,
		AsOfDateModifier: r.AsOfDateModifier,
	}

	r.trailer = groupTrailer{
		GroupControlTotal: r.GroupControlTotal,
		NumberOfAccounts:  r.NumberOfAccounts,
		NumberOfRecords:   r.NumberOfRecords,
	}

}

func (r *Group) String() string {

	r.copyRecords()

	var buf bytes.Buffer
	buf.WriteString(r.header.string() + "\n")
	for i := range r.Accounts {
		buf.WriteString(r.Accounts[i].String() + "\n")
	}
	buf.WriteString(r.trailer.string())

	return buf.String()
}

func (r *Group) Validate() error {

	r.copyRecords()

	if err := r.header.validate(); err != nil {
		return err
	}

	for i := range r.Accounts {
		if err := r.Accounts[i].Validate(); err != nil {
			return err
		}
	}

	if err := r.trailer.validate(); err != nil {
		return err
	}

	return nil
}

func (r *Group) Read(scan Bai2Scanner, input string, lineNum int) (int, string, error) {

	var err error

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
		case util.GroupHeaderCode:

			lineNum++
			newRecord := groupHeader{}
			_, err = newRecord.parse(line)
			if err != nil {
				return lineNum, line, fmt.Errorf("ERROR parsing group header on line %d - %v", lineNum, err)
			}

			r.Receiver = newRecord.Receiver
			r.Originator = newRecord.Originator
			r.GroupStatus = newRecord.GroupStatus
			r.AsOfDate = newRecord.AsOfDate
			r.AsOfTime = newRecord.AsOfTime
			r.CurrencyCode = newRecord.CurrencyCode
			r.AsOfDateModifier = newRecord.AsOfDateModifier

		case util.GroupTrailerCode:

			lineNum++
			newRecord := groupTrailer{}
			_, err = newRecord.parse(line)
			if err != nil {
				return lineNum, line, fmt.Errorf("ERROR parsing group trailer on line %d - %v", lineNum, err)
			}

			r.GroupControlTotal = newRecord.GroupControlTotal
			r.NumberOfAccounts = newRecord.NumberOfAccounts
			r.NumberOfRecords = newRecord.NumberOfRecords

			return lineNum, "", nil

		case util.AccountIdentifierCode:

			newAccount := NewAccount()
			lineNum, input, err = newAccount.Read(scan, line, lineNum)
			if err != nil {
				return lineNum, input, err
			}

			r.Accounts = append(r.Accounts, *newAccount)

		default:

			return lineNum, line, err

		}
	}

	return lineNum, "", nil
}

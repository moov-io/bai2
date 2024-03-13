// Copyright 2022 The Moov Authors
// Use of this source code is governed by an Apache License
// license that can be found in the LICENSE file.

package lib

import (
	"bytes"
	"errors"
	"fmt"

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

func (r *Group) String(opts ...int64) string {

	r.copyRecords()

	var buf bytes.Buffer
	buf.WriteString(r.header.string() + "\n")
	for i := range r.Accounts {
		buf.WriteString(r.Accounts[i].String(opts...) + "\n")
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

func (r *Group) Read(scan *Bai2Scanner, useCurrentLine bool) error {
	if scan == nil {
		return errors.New("invalid bai2 scanner")
	}

	var err error
	for line := scan.ScanLine(useCurrentLine); line != ""; line = scan.ScanLine(useCurrentLine) {
		useCurrentLine = false

		// find record code
		if len(line) < 3 {
			continue
		}

		switch line[:2] {
		case util.GroupHeaderCode:
			newRecord := groupHeader{}
			_, err = newRecord.parse(line)
			if err != nil {
				return fmt.Errorf("ERROR parsing group header on line %d (%v)", scan.GetLineIndex(), err)
			}

			r.Receiver = newRecord.Receiver
			r.Originator = newRecord.Originator
			r.GroupStatus = newRecord.GroupStatus
			r.AsOfDate = newRecord.AsOfDate
			r.AsOfTime = newRecord.AsOfTime
			r.CurrencyCode = newRecord.CurrencyCode
			r.AsOfDateModifier = newRecord.AsOfDateModifier

		case util.AccountIdentifierCode:
			newAccount := NewAccount()
			err = newAccount.Read(scan, true)
			if err != nil {
				return err
			}

			r.Accounts = append(r.Accounts, *newAccount)

		case util.GroupTrailerCode:
			newRecord := groupTrailer{}
			_, err = newRecord.parse(line)
			if err != nil {
				return fmt.Errorf("ERROR parsing group trailer on line %d (%v)", scan.GetLineIndex(), err)
			}

			r.GroupControlTotal = newRecord.GroupControlTotal
			r.NumberOfAccounts = newRecord.NumberOfAccounts
			r.NumberOfRecords = newRecord.NumberOfRecords

			return nil

		default:
			return fmt.Errorf("ERROR parsing group on line %d (unable to read record type %s)", scan.GetLineIndex(), line[0:2])
		}
	}

	return nil
}

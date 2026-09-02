// Copyright 2022 The Moov Authors
// Use of this source code is governed by an Apache License
// license that can be found in the LICENSE file.

package bai3

import (
	"bytes"
	"errors"
	"fmt"
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

Group:
	The next level of organization is the group. A group includes one or more account envelopes, all of which represent accounts at the same financial institution.
	All information in a group is for the same date and time.
	Example: Several accounts from Last National Bank to XYZ Reporting Service, sameday information as of 9:00 AM.

File:
	Groups are combined into files for transmission. A file includes data transmitted between one sender and one receiver.
	A file can include groups of data from any number of originating financial institutions destined for any number of ultimate receivers.
	The groups within a file may contain different As-of-Dates and times.

The following list shows multiple groups within a file and multiple accounts within a group:
  Record
   Code 		Record Name 		Purpose
01 			File Header 		Begins File
02 			Group Header 		Begins First Group
03 			Account Identifier 	First Account
16 			Transaction Detail 	First Account Detail
49 			Account Trailer 	Ends First Account
03 			Account Identifier 	Second Account
49 			Account Trailer Ends Second Account
98 			Group Trailer 		Ends First Group
02 			Group Header Begins Second Group
03 			Account Identifier 	Begins Third Account
88 			Continuation Continues Previous 03 Record
49 			Account Trailer Ends Third Account
98 			Group Trailer Ends Second Group
99 			File Trailer Ends File
The preceding example included two groups. The first group included two accounts, the second
included one account. Only the first account of the first group included transaction detail.

*/

// NewFile returns an empty BTR3 / BAI3 file.
func NewFile() *File {
	return &File{VersionNumber: VersionBTR3}
}

// NewFileWith returns a BTR3 / BAI3 file with the specified options.
func NewFileWith(options Options) *File {
	return &File{VersionNumber: VersionBTR3, options: options}
}

// File is a BTR3 (X9.121 version 3) balance reporting file.
type File struct {
	// File Header
	Sender               string `json:"sender"`
	Receiver             string `json:"receiver"`
	FileCreatedDate      string `json:"fileCreatedDate"`
	FileCreatedTime      string `json:"fileCreatedTime"`
	FileIdNumber         string `json:"fileIdNumber"`
	PhysicalRecordLength int64  `json:"physicalRecordLength,omitempty"`
	BlockSize            int64  `json:"blockSize,omitempty"`
	VersionNumber        int64  `json:"versionNumber"`

	// File trailer
	FileControlTotal string `json:"fileControlTotal"`
	NumberOfBanks    int64  `json:"numberOfBanks"`
	NumberOfRecords  int64  `json:"numberOfRecords"`

	// Groups
	Banks []Bank

	header  fileHeader
	trailer fileTrailer

	options Options
}

func (r *File) SetOptions(options Options) {
	r.options = options
}

func (r *File) copyRecords() {
	r.header = fileHeader{
		Sender:               r.Sender,
		Receiver:             r.Receiver,
		FileCreatedDate:      r.FileCreatedDate,
		FileCreatedTime:      r.FileCreatedTime,
		FileIdNumber:         r.FileIdNumber,
		PhysicalRecordLength: r.PhysicalRecordLength,
		BlockSize:            r.BlockSize,
		VersionNumber:        r.VersionNumber,
	}

	r.trailer = fileTrailer{
		FileControlTotal: r.FileControlTotal,
		NumberOfBanks:    r.NumberOfBanks,
		NumberOfRecords:  r.NumberOfRecords,
	}
}

// SumRecords is the number of records in the file: the sum of each group's
// NumberOfRecords plus the file header (01) and file trailer (99).
func (f *File) SumRecords() int64 {
	var sum int64
	for _, group := range f.Banks {
		sum += group.NumberOfRecords
	}
	return sum + 2
}

// Sums the number of groups. Maps to the NumberOfBanks field.
func (g *File) SumNumberOfBanks() int64 {
	return int64(len(g.Banks))
}

func (f *File) Version() int64 {
	if f.VersionNumber != 0 {
		return f.VersionNumber
	}
	return VersionBTR3
}

// Sums the Group Control Totals. Maps to the FileControlTotal field.
func (a *File) SumBankControlTotals() (string, error) {
	if err := a.Validate(); err != nil {
		return "0", err
	}
	var sum int64
	for _, group := range a.Banks {
		amt, err := strconv.ParseInt(group.BankControlTotal, 10, 64)
		if err != nil {
			return "0", err
		}
		sum += amt
	}
	return fmt.Sprint(sum), nil
}

func (r *File) String() string {
	r.copyRecords()

	var buf bytes.Buffer
	buf.WriteString(r.header.string() + "\n")
	for i := range r.Banks {
		buf.WriteString(r.Banks[i].String(r.PhysicalRecordLength) + "\n")
	}
	buf.WriteString(r.trailer.string())

	return buf.String()
}

func (r *File) Validate() error {
	r.copyRecords()

	if err := r.header.validate(r.options); err != nil {
		return err
	}

	for i := range r.Banks {
		if r.options.StrictBankHeader {
			if r.Banks[i].CurrencyCode != "" {
				return fmt.Errorf("BankHeader: currency must be empty in BTR3")
			}
			if r.Banks[i].GroupStatus != 1 {
				return fmt.Errorf("BankHeader: group status must be 1 (Update) in BTR3")
			}
		}
		if err := r.Banks[i].Validate(); err != nil {
			return err
		}
	}

	if err := r.trailer.validate(); err != nil {
		return err
	}

	return nil
}

func (r *File) Read(scan *util.Scanner) error {
	if scan == nil {
		return errors.New("invalid bai3 scanner")
	}

	var err error
	sawHeader := false
	sawTrailer := false

	for line := scan.ScanLine(); line != ""; line = scan.ScanLine() {
		if err := scan.Err(); err != nil {
			return err
		}

		// find record code
		if len(line) < 3 {
			continue
		}

		switch line[0:2] {
		case util.FileHeaderCode:
			if sawHeader {
				return fmt.Errorf("ERROR parsing file on line %d (duplicate file header)", scan.GetLineIndex())
			}

			newRecord := fileHeader{}
			_, err = newRecord.parse(line, r.options)
			if err != nil {
				return fmt.Errorf("ERROR parsing file header on line %d (%v)", scan.GetLineIndex(), err)
			}

			r.Sender = newRecord.Sender
			r.Receiver = newRecord.Receiver
			r.FileCreatedDate = newRecord.FileCreatedDate
			r.FileCreatedTime = newRecord.FileCreatedTime
			r.FileIdNumber = newRecord.FileIdNumber
			r.PhysicalRecordLength = newRecord.PhysicalRecordLength
			r.BlockSize = newRecord.BlockSize
			r.VersionNumber = newRecord.VersionNumber
			sawHeader = true
			scan.SetPhysicalRecordLength(r.PhysicalRecordLength)

		case util.GroupHeaderCode:
			if !sawHeader {
				return fmt.Errorf("ERROR parsing file on line %d (file header is required before groups)", scan.GetLineIndex())
			}

			newBank := NewBank()
			err = newBank.Read(scan, true)
			if err != nil {
				return err
			}

			r.Banks = append(r.Banks, *newBank)

		case util.FileTrailerCode:
			newRecord := fileTrailer{}
			_, err = newRecord.parse(line)
			if err != nil {
				return fmt.Errorf("ERROR parsing file trailer on line %d (%v)", scan.GetLineIndex(), err)
			}

			r.FileControlTotal = newRecord.FileControlTotal
			r.NumberOfBanks = newRecord.NumberOfBanks
			r.NumberOfRecords = newRecord.NumberOfRecords
			sawTrailer = true

			return r.finishRead(sawHeader, sawTrailer)

		default:
			return fmt.Errorf("ERROR parsing file on line %d (unsupported record type %s)", scan.GetLineIndex(), line[0:2])
		}
	}

	if err := scan.Err(); err != nil {
		return err
	}

	return r.finishRead(sawHeader, sawTrailer)
}

func (r *File) finishRead(sawHeader, sawTrailer bool) error {
	if !sawHeader {
		return errors.New("missing file header (01)")
	}
	if !sawTrailer {
		return errors.New("missing file trailer (99)")
	}
	if r.options.StrictControlTotals {
		return r.validateControlTotals()
	}
	return nil
}

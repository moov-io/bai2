// Copyright 2022 The Moov Authors
// Use of this source code is governed by an Apache License
// license that can be found in the LICENSE file.

package lib

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

// Creating new file object
func NewBai2() *Bai2 {
	return &Bai2{}
}

// FILE with BAI Format
type Bai2 struct {
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
	NumberOfGroups   int64  `json:"numberOfGroups"`
	NumberOfRecords  int64  `json:"numberOfRecords"`

	// Groups
	Groups []Group

	header  fileHeader
	trailer fileTrailer
}

func (r *Bai2) copyRecords() {

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
		NumberOfGroups:   r.NumberOfGroups,
		NumberOfRecords:  r.NumberOfRecords,
	}

}

func (f *Bai2) SumRecords() int64 {
	var sum int64
	for _, group := range f.Groups {
		sum += group.NumberOfRecords
	}
	// Add two for the group header and trailer records
	return sum + 2
}

func (g *Bai2) SumNumberOfGroups() int64 {
	return int64(len(g.Groups))
}

func (a *Bai2) SumGroupControlTotals() (string, error) {
	if err := a.Validate(); err != nil {
		return "0", err
	}
	var sum int64
	for _, group := range a.Groups {
		amt, err := strconv.ParseInt(group.GroupControlTotal, 10, 64)
		if err != nil {
			return "0", err
		}
		sum += amt
	}
	return fmt.Sprint(sum), nil
}

func (r *Bai2) String() string {

	r.copyRecords()

	var buf bytes.Buffer
	buf.WriteString(r.header.string() + "\n")
	for i := range r.Groups {
		buf.WriteString(r.Groups[i].String(r.PhysicalRecordLength) + "\n")
	}
	buf.WriteString(r.trailer.string())

	return buf.String()
}

func (r *Bai2) Validate() error {

	r.copyRecords()

	if err := r.header.validate(); err != nil {
		return err
	}

	for i := range r.Groups {
		if err := r.Groups[i].Validate(); err != nil {
			return err
		}
	}

	if err := r.trailer.validate(); err != nil {
		return err
	}

	return nil
}

func (r *Bai2) Read(scan *Bai2Scanner) error {

	if scan == nil {
		return errors.New("invalid bai2 scanner")
	}

	var err error
	for line := scan.ScanLine(); line != ""; line = scan.ScanLine() {

		// find record code
		if len(line) < 3 {
			continue
		}

		switch line[0:2] {
		case util.FileHeaderCode:

			newRecord := fileHeader{}
			_, err = newRecord.parse(line)
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

		case util.GroupHeaderCode:

			newGroup := NewGroup()
			err = newGroup.Read(scan, true)
			if err != nil {
				return err
			}

			r.Groups = append(r.Groups, *newGroup)

		case util.FileTrailerCode:

			newRecord := fileTrailer{}
			_, err = newRecord.parse(line)
			if err != nil {
				return fmt.Errorf("ERROR parsing file trailer on line %d (%v)", scan.GetLineIndex(), err)
			}

			r.FileControlTotal = newRecord.FileControlTotal
			r.NumberOfGroups = newRecord.NumberOfGroups
			r.NumberOfRecords = newRecord.NumberOfRecords

			return nil

		default:
			return fmt.Errorf("ERROR parsing file on line %d (unsupported record type %s)", scan.GetLineIndex(), line[0:2])
		}
	}

	return nil
}

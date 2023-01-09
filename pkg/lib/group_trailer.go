// Copyright 2022 The Moov Authors
// Use of this source code is governed by an Apache License
// license that can be found in the LICENSE file.

package lib

import (
	"bytes"
	"fmt"

	"github.com/moov-io/bai2/pkg/util"
)

/*
Group Trailer

The group trailer is the second last record in a BAI format file. This record contains information on the
group control total, the number of accounts and the number of records. The group control total is the
sum of the account control totals in the group. The number of records is the total of all type 02, 03, 16, 49,
88 and 98 records in the group.

*/

const (
	gtParseErrorFmt    = "GroupTrailer: unable to parse %s"
	gtValidateErrorFmt = "GroupTrailer: invalid %s"
)

// Creating Group Trailer
func NewGroupTrailer() *GroupTrailer {
	return &GroupTrailer{
		RecordCode: "98",
	}

}

// Group Trailer
type GroupTrailer struct {
	RecordCode        string
	GroupControlTotal string
	NumberOfAccounts  int64
	NumberOfRecords   int64
}

func (h *GroupTrailer) Validate() error {
	if h.RecordCode != "98" {
		return fmt.Errorf(fmt.Sprintf(gtValidateErrorFmt, "RecordCode"))
	}
	if h.GroupControlTotal != "" && !util.ValidateAmount(h.GroupControlTotal) {
		return fmt.Errorf(fmt.Sprintf(gtValidateErrorFmt, "GroupControlTotal"))
	}

	return nil
}

func (h *GroupTrailer) Parse(data string) (int, error) {

	var line string
	var err error
	var size, read int

	if length := util.GetSize(data); length < 2 {
		return 0, fmt.Errorf(fmt.Sprintf(gtParseErrorFmt, "record"))
	} else {
		line = data[:length]
	}

	// RecordCode
	if h.RecordCode, size, err = util.ReadField(util.GetField(line, read)); err != nil {
		return 0, fmt.Errorf(fmt.Sprintf(gtParseErrorFmt, "RecordCode"))
	} else {
		read += size
	}

	// GroupControlTotal
	if h.GroupControlTotal, size, err = util.ReadField(util.GetField(line, read)); err != nil {
		return 0, fmt.Errorf(fmt.Sprintf(gtParseErrorFmt, "GroupControlTotal"))
	} else {
		read += size
	}

	// NumberOfAccounts
	if h.NumberOfAccounts, size, err = util.ReadFieldAsInt(util.GetField(line, read)); err != nil {
		return 0, fmt.Errorf(fmt.Sprintf(gtParseErrorFmt, "NumberOfAccounts"))
	} else {
		read += size
	}

	// NumberOfRecords
	if h.NumberOfRecords, size, err = util.ReadFieldAsInt(util.GetField(line, read)); err != nil {
		return 0, fmt.Errorf(fmt.Sprintf(gtParseErrorFmt, "NumberOfRecords"))
	} else {
		read += size
	}

	if err = h.Validate(); err != nil {
		return 0, err
	}

	return read, nil
}

func (h *GroupTrailer) String() string {
	var buf bytes.Buffer

	buf.WriteString(fmt.Sprintf("%s,", h.RecordCode))
	buf.WriteString(fmt.Sprintf("%s,", h.GroupControlTotal))
	buf.WriteString(fmt.Sprintf("%d,", h.NumberOfAccounts))
	buf.WriteString(fmt.Sprintf("%d/", h.NumberOfRecords))

	return buf.String()
}

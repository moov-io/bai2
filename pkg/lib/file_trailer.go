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

File Trailer

The file trailer is the last record in a BAI format file. This record contains information on the file control
total, the number of groups and the number of records. The file control total is the sum of the group
control totals in the file. The number of groups is the sum of type 02 records in the file. The number of
records is the sum of all records in the file, including the file trailer (type 99) record.

*/

const (
	ftParseErrorFmt    = "FileTrailer: unable to parse %s"
	ftValidateErrorFmt = "FileTrailer: invalid %s"
)

// Creating File Trailer
func NewFileTrailer() *FileTrailer {
	return &FileTrailer{
		RecordCode: "99",
	}

}

// File Trailer
type FileTrailer struct {
	RecordCode       string
	FileControlTotal string
	NumberOfGroups   int64
	NumberOfRecords  int64
}

func (h *FileTrailer) Validate() error {
	if h.RecordCode != "99" {
		return fmt.Errorf(fmt.Sprintf(ftValidateErrorFmt, "RecordCode"))
	}
	if h.FileControlTotal != "" && !util.ValidateAmount(h.FileControlTotal) {
		return fmt.Errorf(fmt.Sprintf(ftValidateErrorFmt, "FileControlTotal"))
	}

	return nil
}

func (h *FileTrailer) Parse(data string) (int, error) {

	var line string
	var err error
	var size, read int

	if length := util.GetSize(data); length < 2 {
		return 0, fmt.Errorf(fmt.Sprintf(ftParseErrorFmt, "record"))
	} else {
		line = data[:length]
	}

	// RecordCode
	if h.RecordCode, size, err = util.ReadField(line, read); err != nil {
		return 0, fmt.Errorf(fmt.Sprintf(ftParseErrorFmt, "RecordCode"))
	} else {
		read += size
	}

	// GroupControlTotal
	if h.FileControlTotal, size, err = util.ReadField(line, read); err != nil {
		return 0, fmt.Errorf(fmt.Sprintf(ftParseErrorFmt, "GroupControlTotal"))
	} else {
		read += size
	}

	// NumberOfGroups
	if h.NumberOfGroups, size, err = util.ReadFieldAsInt(line, read); err != nil {
		return 0, fmt.Errorf(fmt.Sprintf(ftParseErrorFmt, "NumberOfGroups"))
	} else {
		read += size
	}

	// NumberOfRecords
	if h.NumberOfRecords, size, err = util.ReadFieldAsInt(line, read); err != nil {
		return 0, fmt.Errorf(fmt.Sprintf(ftParseErrorFmt, "NumberOfRecords"))
	} else {
		read += size
	}

	if err = h.Validate(); err != nil {
		return 0, err
	}

	return read, nil
}

func (h *FileTrailer) String() string {
	var buf bytes.Buffer

	buf.WriteString(fmt.Sprintf("%s,", h.RecordCode))
	buf.WriteString(fmt.Sprintf("%s,", h.FileControlTotal))
	buf.WriteString(fmt.Sprintf("%d,", h.NumberOfGroups))
	buf.WriteString(fmt.Sprintf("%d/", h.NumberOfRecords))

	return buf.String()
}

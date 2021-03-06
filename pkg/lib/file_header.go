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

/*

File Header

The File Header is the first record in a BAI format file. It always has a record code of 01.

*/

const (
	fileHeaderLength = 37
)

// Creating File Header
func NewFileHeader() *FileHeader {
	return &FileHeader{
		RecordCode:           "01",
		Sender:               "0004",
		PhysicalRecordLength: 80,
		BlockSize:            1,
		VersionNumber:        2,
	}

}

// File Header
type FileHeader struct {
	RecordCode           string
	Sender               string
	Receiver             string
	FileCreatedDate      string
	FileCreatedTime      string
	FileIdNumber         int64
	PhysicalRecordLength int64
	BlockSize            int64
	VersionNumber        int64
}

func (h *FileHeader) Validate() error {
	if h.RecordCode != "01" {
		return fmt.Errorf("FileHeader: invalid record code")
	}
	if h.Sender != "0004" {
		return fmt.Errorf("FileHeader: invalid sender")
	}
	if h.PhysicalRecordLength != 80 {
		return fmt.Errorf("FileHeader: invalid physical record length")
	}
	if h.BlockSize != 1 {
		return fmt.Errorf("FileHeader: invalid block size")
	}
	if h.VersionNumber != 2 {
		return fmt.Errorf("FileHeader: invalid version number")
	}

	return nil
}

func (h *FileHeader) Parse(line string) error {
	if n := utf8.RuneCountInString(line); n < fileHeaderLength {
		return fmt.Errorf("FileHeader: length %d is too short", n)
	}

	h.RecordCode, _ = util.EntryParser(line[0:3], ",")
	h.Sender, _ = util.EntryParser(line[3:8], ",")
	h.Receiver, _ = util.EntryParser(line[8:14], ",")
	h.FileCreatedDate, _ = util.EntryParser(line[14:21], ",")
	h.FileCreatedTime, _ = util.EntryParser(line[21:26], ",")
	h.FileIdNumber, _ = util.EntryParserToInt(line[26:30], ",")
	h.PhysicalRecordLength, _ = util.EntryParserToInt(line[30:33], ",")
	h.BlockSize, _ = util.EntryParserToInt(line[33:35], ",")
	h.VersionNumber, _ = util.EntryParserToInt(line[35:37], "/")

	return nil
}

func (h *FileHeader) String() string {
	var buf bytes.Buffer

	buf.WriteString(fmt.Sprintf("%-2.2v,", h.RecordCode))
	buf.WriteString(fmt.Sprintf("%-4.4v,", h.Sender))
	buf.WriteString(fmt.Sprintf("%-5.5v,", h.Receiver))
	buf.WriteString(fmt.Sprintf("%-6.6v,", h.FileCreatedDate))
	buf.WriteString(fmt.Sprintf("%-4.4v,", h.FileCreatedTime))
	buf.WriteString(fmt.Sprintf("%03.3v,", h.FileIdNumber))
	buf.WriteString(fmt.Sprintf("%02.2v,", h.PhysicalRecordLength))
	buf.WriteString(fmt.Sprintf("%01.1v,", h.BlockSize))
	buf.WriteString(fmt.Sprintf("%01.1v/", h.VersionNumber))

	return buf.String()
}

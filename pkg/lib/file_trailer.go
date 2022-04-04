// Copyright 2022 The Moov Authors
// Use of this source code is governed by an Apache License
// license that can be found in the LICENSE file.

package lib

import (
	"bytes"
	"fmt"
	"github.com/moov-io/bai2/pkg/util"
	"unicode/utf8"
)

const (
	fileTrailerLength = 42
)

func NewFileTrailer() *FileTrailer {
	return &FileTrailer{
		RecordCode: "99",
	}

}

type FileTrailer struct {
	RecordCode        string
	GroupControlTotal string
	NumberOfGroups    int64
	NumberOfRecords   int64
}

func (h *FileTrailer) Validate() error {
	if h.RecordCode != "99" {
		return fmt.Errorf("FileTrailer: invalid record code")
	}

	return nil
}

func (h *FileTrailer) Parse(line string) error {
	if n := utf8.RuneCountInString(line); n < fileTrailerLength {
		return fmt.Errorf("FileTrailer: length %d is too short", n)
	}

	h.RecordCode, _ = util.EntryParser(line[0:3], ",")
	h.GroupControlTotal, _ = util.EntryParser(line[3:22], ",")
	h.NumberOfGroups, _ = util.EntryParserToInt(line[22:32], ",")
	h.NumberOfRecords, _ = util.EntryParserToInt(line[32:42], "/")

	return nil
}

func (h *FileTrailer) String() string {
	var buf bytes.Buffer

	buf.WriteString(fmt.Sprintf("%-2.2v,", h.RecordCode))
	buf.WriteString(fmt.Sprintf("%-18.18v,", h.GroupControlTotal))
	buf.WriteString(fmt.Sprintf("%09.9v,", h.NumberOfGroups))
	buf.WriteString(fmt.Sprintf("%09.9v/", h.NumberOfRecords))

	return buf.String()
}

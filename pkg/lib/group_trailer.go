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

const (
	groupTrailerLength = 42
)

func NewGroupTrailer() *GroupTrailer {
	return &GroupTrailer{
		RecordCode: "98",
	}

}

type GroupTrailer struct {
	RecordCode        string
	GroupControlTotal string
	NumberOfAccounts  int64
	NumberOfRecords   int64
}

func (h *GroupTrailer) Validate() error {
	if h.RecordCode != "98" {
		return fmt.Errorf("GroupTrailer: invalid record code")
	}

	return nil
}

func (h *GroupTrailer) Parse(line string) error {
	if n := utf8.RuneCountInString(line); n < groupTrailerLength {
		return fmt.Errorf("GroupTrailer: length %d is too short", n)
	}

	h.RecordCode, _ = util.EntryParser(line[0:3], ",")
	h.GroupControlTotal, _ = util.EntryParser(line[3:22], ",")
	h.NumberOfAccounts, _ = util.EntryParserToInt(line[22:32], ",")
	h.NumberOfRecords, _ = util.EntryParserToInt(line[32:42], "/")

	return nil
}

func (h *GroupTrailer) String() string {
	var buf bytes.Buffer

	buf.WriteString(fmt.Sprintf("%-2.2v,", h.RecordCode))
	buf.WriteString(fmt.Sprintf("%-18.18v,", h.GroupControlTotal))
	buf.WriteString(fmt.Sprintf("%09.9v,", h.NumberOfAccounts))
	buf.WriteString(fmt.Sprintf("%09.9v/", h.NumberOfRecords))

	return buf.String()
}

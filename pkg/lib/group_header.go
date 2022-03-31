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
	fileGroupLength = 29
)

func NewGroupHeader() *GroupHeader {
	return &GroupHeader{
		RecordCode:   "02",
		Sender:       "0004",
		CurrencyCode: "USD",
	}

}

type GroupHeader struct {
	RecordCode   string
	Receiver     string
	Sender       string
	GroupStatus  string
	AsOfDate     string
	CurrencyCode string
}

func (h *GroupHeader) Validate() error {
	if h.RecordCode != "02" {
		return fmt.Errorf("GroupHeader: invalid record code")
	}
	if h.Sender != "0004" {
		return fmt.Errorf("GroupHeader: invalid sender")
	}
	if h.CurrencyCode != "USD" && h.CurrencyCode != "CAD" {
		return fmt.Errorf("GroupHeader: invalid currency code")
	}

	return nil
}

func (h *GroupHeader) Parse(line string) error {
	if n := utf8.RuneCountInString(line); n < fileGroupLength {
		return fmt.Errorf("GroupHeader: length %d is too short", n)
	}

	h.RecordCode, _ = util.EntryParser(line[0:3], ",")
	h.Receiver, _ = util.EntryParser(line[3:9], ",")
	h.Sender, _ = util.EntryParser(line[9:14], ",")
	h.GroupStatus, _ = util.EntryParser(line[14:16], ",")
	h.AsOfDate, _ = util.EntryParser(line[16:23], ",")
	h.CurrencyCode, _ = util.EntryParser(line[24:28], ",")

	return nil
}

func (h *GroupHeader) String() string {
	var buf bytes.Buffer

	buf.WriteString(fmt.Sprintf("%2.2v,", h.RecordCode))
	buf.WriteString(fmt.Sprintf("%5.5v,", h.Receiver))
	buf.WriteString(fmt.Sprintf("%4.4v,", h.Sender))
	buf.WriteString(fmt.Sprintf("%1.1v,", h.GroupStatus))
	buf.WriteString(fmt.Sprintf("%6.6v,", h.AsOfDate))
	buf.WriteString(fmt.Sprintf(","))
	buf.WriteString(fmt.Sprintf("%3.3v,", h.CurrencyCode))
	buf.WriteString(fmt.Sprintf("/"))

	return buf.String()
}

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

Group Header

The Group Header is the second record in a BAI format file. It always has a record code of 02.

*/

const (
	ghParseErrorFmt    = "GroupHeader: unable to parse %s"
	ghValidateErrorFmt = "GroupHeader: invalid %s"
)

// Creating Group Header
func NewGroupHeader() *GroupHeader {
	return &GroupHeader{
		RecordCode: "02",
	}
}

// Group Header
type GroupHeader struct {
	RecordCode       string
	Receiver         string `json:",omitempty"`
	Originator       string
	GroupStatus      int64
	AsOfDate         string
	AsOfTime         string `json:",omitempty"`
	CurrencyCode     string `json:",omitempty"`
	AsOfDateModifier int64  `json:",omitempty"`
}

func (h *GroupHeader) Validate() error {
	if h.RecordCode != "02" {
		return fmt.Errorf(fmt.Sprintf(ghValidateErrorFmt, "RecordCode"))
	}
	if h.Originator == "" {
		return fmt.Errorf(fmt.Sprintf(ghValidateErrorFmt, "Originator"))
	}
	if h.GroupStatus < 0 || h.GroupStatus > 4 {
		return fmt.Errorf(fmt.Sprintf(ghValidateErrorFmt, "GroupStatus"))
	}
	if h.AsOfDate == "" {
		return fmt.Errorf(fmt.Sprintf(ghValidateErrorFmt, "AsOfDate"))
	} else if !util.ValidateData(h.AsOfDate) {
		return fmt.Errorf(fmt.Sprintf(ghValidateErrorFmt, "AsOfDate"))
	}
	if h.AsOfTime != "" && !util.ValidateTime(h.AsOfTime) {
		return fmt.Errorf(fmt.Sprintf(ghValidateErrorFmt, "AsOfTime"))
	}
	if h.CurrencyCode != "" && !util.ValidateCurrencyCode(h.CurrencyCode) {
		return fmt.Errorf(fmt.Sprintf(ghValidateErrorFmt, "CurrencyCode"))
	}
	if h.AsOfDateModifier < 0 || h.AsOfDateModifier > 4 {
		return fmt.Errorf(fmt.Sprintf(ghValidateErrorFmt, "AsOfDateModifier"))
	}

	return nil
}

func (h *GroupHeader) Parse(data string) (int, error) {

	var line string
	var err error
	var size, read int

	if length := util.GetSize(data); length < 2 {
		return 0, fmt.Errorf(fmt.Sprintf(ghParseErrorFmt, "record"))
	} else {
		line = data[:length]
	}

	// RecordCode
	if h.RecordCode, size, err = util.ReadField(util.GetField(line, read)); err != nil {
		return 0, fmt.Errorf(fmt.Sprintf(ghParseErrorFmt, "RecordCode"))
	} else {
		read += size
	}

	// Receiver
	if h.Receiver, size, err = util.ReadField(util.GetField(line, read)); err != nil {
		return 0, fmt.Errorf(fmt.Sprintf(ghParseErrorFmt, "Receiver"))
	} else {
		read += size
	}

	// Originator
	if h.Originator, size, err = util.ReadField(util.GetField(line, read)); err != nil {
		return 0, fmt.Errorf(fmt.Sprintf(ghParseErrorFmt, "Originator"))
	} else {
		read += size
	}

	// GroupStatus
	if h.GroupStatus, size, err = util.ReadFieldAsInt(util.GetField(line, read)); err != nil {
		return 0, fmt.Errorf(fmt.Sprintf(ghParseErrorFmt, "GroupStatus"))
	} else {
		read += size
	}

	// AsOfDate
	if h.AsOfDate, size, err = util.ReadField(util.GetField(line, read)); err != nil {
		return 0, fmt.Errorf(fmt.Sprintf(ghParseErrorFmt, "AsOfDate"))
	} else {
		read += size
	}

	// AsOfTime
	if h.AsOfTime, size, err = util.ReadField(util.GetField(line, read)); err != nil {
		return 0, fmt.Errorf(fmt.Sprintf(ghParseErrorFmt, "AsOfTime"))
	} else {
		read += size
	}

	// CurrencyCode
	if h.CurrencyCode, size, err = util.ReadField(util.GetField(line, read)); err != nil {
		return 0, fmt.Errorf(fmt.Sprintf(ghParseErrorFmt, "CurrencyCode"))
	} else {
		read += size
	}

	// AsOfDateModifier
	if h.AsOfDateModifier, size, err = util.ReadFieldAsInt(util.GetField(line, read)); err != nil {
		return 0, fmt.Errorf(fmt.Sprintf(ghParseErrorFmt, "AsOfDateModifier"))
	} else {
		read += size
	}

	if err = h.Validate(); err != nil {
		return 0, err
	}

	return read, nil
}

func (h *GroupHeader) String() string {
	var buf bytes.Buffer

	buf.WriteString(fmt.Sprintf("%s,", h.RecordCode))
	buf.WriteString(fmt.Sprintf("%s,", h.Receiver))
	buf.WriteString(fmt.Sprintf("%s,", h.Originator))
	buf.WriteString(fmt.Sprintf("%d,", h.GroupStatus))
	buf.WriteString(fmt.Sprintf("%s,", h.AsOfDate))
	buf.WriteString(fmt.Sprintf("%s,", h.AsOfTime))
	buf.WriteString(fmt.Sprintf("%s,", h.CurrencyCode))
	if h.AsOfDateModifier > 0 {
		buf.WriteString(fmt.Sprintf("%d/", h.AsOfDateModifier))
	} else {
		buf.WriteString("/")
	}

	return buf.String()
}

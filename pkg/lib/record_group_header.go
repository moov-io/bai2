// Copyright 2022 The Moov Authors
// Use of this source code is governed by an Apache License
// license that can be found in the LICENSE file.

package lib

import (
	"bytes"
	"fmt"

	"github.com/moov-io/bai2/pkg/util"
)

const (
	ghParseErrorFmt    = "GroupHeader: unable to parse %s"
	ghValidateErrorFmt = "GroupHeader: invalid %s"
)

type groupHeader struct {
	Receiver         string `json:",omitempty"`
	Originator       string
	GroupStatus      int64
	AsOfDate         string
	AsOfTime         string `json:",omitempty"`
	CurrencyCode     string `json:",omitempty"`
	AsOfDateModifier int64  `json:",omitempty"`
}

func (h *groupHeader) validate() error {
	if h.Originator == "" {
		return fmt.Errorf(ghValidateErrorFmt, "Originator")
	}
	if h.GroupStatus < 1 || h.GroupStatus > 4 {
		return fmt.Errorf(ghValidateErrorFmt, "GroupStatus")
	}
	if h.AsOfDate == "" {
		return fmt.Errorf(ghValidateErrorFmt, "AsOfDate")
	} else if !util.ValidateDate(h.AsOfDate) {
		return fmt.Errorf(ghValidateErrorFmt, "AsOfDate")
	}
	if h.AsOfTime != "" && !util.ValidateTime(h.AsOfTime) {
		return fmt.Errorf(ghValidateErrorFmt, "AsOfTime")
	}
	if h.CurrencyCode != "" && !util.ValidateCurrencyCode(h.CurrencyCode) {
		return fmt.Errorf(ghValidateErrorFmt, "CurrencyCode")
	}
	if h.AsOfDateModifier < 0 || h.AsOfDateModifier > 4 {
		return fmt.Errorf(ghValidateErrorFmt, "AsOfDateModifier")
	}

	return nil
}

func (h *groupHeader) parse(data string) (int, error) {

	var line string
	var err error
	var size, read int

	if length := util.GetSize(data); length < 3 {
		return 0, fmt.Errorf(ghParseErrorFmt, "record")
	} else {
		line = data[:length]
	}

	// RecordCode
	if util.GroupHeaderCode != data[:2] {
		return 0, fmt.Errorf(ghParseErrorFmt, "RecordCode")
	}
	read += 3

	// Receiver
	if h.Receiver, size, err = util.ReadField(line, read); err != nil {
		return 0, fmt.Errorf(ghParseErrorFmt, "Receiver")
	} else {
		read += size
	}

	// Originator
	if h.Originator, size, err = util.ReadField(line, read); err != nil {
		return 0, fmt.Errorf(ghParseErrorFmt, "Originator")
	} else {
		read += size
	}

	// GroupStatus
	if h.GroupStatus, size, err = util.ReadFieldAsInt(line, read); err != nil {
		return 0, fmt.Errorf(ghParseErrorFmt, "GroupStatus")
	} else {
		read += size
	}

	// AsOfDate
	if h.AsOfDate, size, err = util.ReadField(line, read); err != nil {
		return 0, fmt.Errorf(ghParseErrorFmt, "AsOfDate")
	} else {
		read += size
	}

	// AsOfTime
	if h.AsOfTime, size, err = util.ReadField(line, read); err != nil {
		return 0, fmt.Errorf(ghParseErrorFmt, "AsOfTime")
	} else {
		read += size
	}

	// CurrencyCode
	if h.CurrencyCode, size, err = util.ReadField(line, read); err != nil {
		return 0, fmt.Errorf(ghParseErrorFmt, "CurrencyCode")
	} else {
		read += size
	}

	// AsOfDateModifier
	if h.AsOfDateModifier, size, err = util.ReadFieldAsInt(line, read); err != nil {
		return 0, fmt.Errorf(ghParseErrorFmt, "AsOfDateModifier")
	} else {
		read += size
	}

	if err = h.validate(); err != nil {
		return 0, err
	}

	return read, nil
}

func (h *groupHeader) string() string {
	var buf bytes.Buffer

	buf.WriteString(fmt.Sprintf("%s,", util.GroupHeaderCode))
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

// Copyright 2022 The Moov Authors
// Use of this source code is governed by an Apache License
// license that can be found in the LICENSE file.

package bai3

import (
	"bytes"
	"fmt"

	"github.com/moov-io/bai2/pkg/util"
)

const (
	bhParseErrorFmt    = "BankHeader: unable to parse %s"
	bhValidateErrorFmt = "BankHeader: invalid %s"
)

type bankHeader struct {
	Receiver           string `json:",omitempty"`
	BankIdentification string
	GroupStatus        int64
	AsOfDate           string
	AsOfTime           string `json:",omitempty"`
	CurrencyCode       string `json:",omitempty"`
	AsOfDateModifier   int64  `json:",omitempty"`
}

func (h *bankHeader) validate() error {
	if h.BankIdentification == "" {
		return fmt.Errorf(bhValidateErrorFmt, "BankIdentification")
	}
	if h.GroupStatus < 1 || h.GroupStatus > 4 {
		return fmt.Errorf(bhValidateErrorFmt, "GroupStatus")
	}
	if h.AsOfDate == "" {
		return fmt.Errorf(bhValidateErrorFmt, "AsOfDate")
	} else if !util.ValidateDate(h.AsOfDate) {
		return fmt.Errorf(bhValidateErrorFmt, "AsOfDate")
	}
	if h.AsOfTime != "" && !util.ValidateTime(h.AsOfTime) {
		return fmt.Errorf(bhValidateErrorFmt, "AsOfTime")
	}
	if h.CurrencyCode != "" && !util.ValidateCurrencyCode(h.CurrencyCode) {
		return fmt.Errorf(bhValidateErrorFmt, "CurrencyCode")
	}
	if h.AsOfDateModifier < 0 || h.AsOfDateModifier > 4 {
		return fmt.Errorf(bhValidateErrorFmt, "AsOfDateModifier")
	}

	return nil
}

func (h *bankHeader) parse(data string) (int, error) {

	var line string
	var err error
	var size, read int

	if length := util.GetSize(data); length < 3 {
		return 0, fmt.Errorf(bhParseErrorFmt, "record")
	} else {
		line = data[:length]
	}

	// RecordCode
	if util.GroupHeaderCode != data[:2] {
		return 0, fmt.Errorf(bhParseErrorFmt, "RecordCode")
	}
	read += 3

	// Receiver
	if h.Receiver, size, err = util.ReadField(line, read); err != nil {
		return 0, fmt.Errorf(bhParseErrorFmt, "Receiver")
	} else {
		read += size
	}

	// Originator
	if h.BankIdentification, size, err = util.ReadField(line, read); err != nil {
		return 0, fmt.Errorf(bhParseErrorFmt, "BankIdentification")
	} else {
		read += size
	}

	// GroupStatus
	if h.GroupStatus, size, err = util.ReadFieldAsInt(line, read); err != nil {
		return 0, fmt.Errorf(bhParseErrorFmt, "GroupStatus")
	} else {
		read += size
	}

	// AsOfDate
	if h.AsOfDate, size, err = util.ReadField(line, read); err != nil {
		return 0, fmt.Errorf(bhParseErrorFmt, "AsOfDate")
	} else {
		read += size
	}

	// AsOfTime
	if h.AsOfTime, size, err = util.ReadField(line, read); err != nil {
		return 0, fmt.Errorf(bhParseErrorFmt, "AsOfTime")
	} else {
		read += size
	}

	// CurrencyCode
	if h.CurrencyCode, size, err = util.ReadField(line, read); err != nil {
		return 0, fmt.Errorf(bhParseErrorFmt, "CurrencyCode")
	} else {
		read += size
	}

	// AsOfDateModifier
	if h.AsOfDateModifier, size, err = util.ReadFieldAsInt(line, read); err != nil {
		return 0, fmt.Errorf(bhParseErrorFmt, "AsOfDateModifier")
	} else {
		read += size
	}

	if h.GroupStatus == 0 {
		h.GroupStatus = 1
	}
	// BTR3: currency is positional-null on record 02; do not default it.

	if err = h.validate(); err != nil {
		return 0, err
	}

	return read, nil
}

func (h *bankHeader) string() string {
	var buf bytes.Buffer

	buf.WriteString(fmt.Sprintf("%s,", util.GroupHeaderCode))
	buf.WriteString(fmt.Sprintf("%s,", h.Receiver))
	buf.WriteString(fmt.Sprintf("%s,", h.BankIdentification))
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

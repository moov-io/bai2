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
	gtParseErrorFmt    = "GroupTrailer: unable to parse %s"
	gtValidateErrorFmt = "GroupTrailer: invalid %s"
)

type groupTrailer struct {
	GroupControlTotal string
	NumberOfAccounts  int64
	NumberOfRecords   int64
}

func (h *groupTrailer) validate() error {
	if h.GroupControlTotal != "" && !util.ValidateAmount(h.GroupControlTotal) {
		return fmt.Errorf(fmt.Sprintf(gtValidateErrorFmt, "GroupControlTotal"))
	}

	return nil
}

func (h *groupTrailer) parse(data string) (int, error) {

	var line string
	var err error
	var size, read int

	if length := util.GetSize(data); length < 3 {
		return 0, fmt.Errorf(fmt.Sprintf(gtParseErrorFmt, "record"))
	} else {
		line = data[:length]
	}

	// RecordCode
	if util.GroupTrailerCode != data[:2] {
		return 0, fmt.Errorf(fmt.Sprintf(fhParseErrorFmt, "RecordCode"))
	}
	read += 3

	// GroupControlTotal
	if h.GroupControlTotal, size, err = util.ReadField(line, read); err != nil {
		return 0, fmt.Errorf(fmt.Sprintf(gtParseErrorFmt, "GroupControlTotal"))
	} else {
		read += size
	}

	// NumberOfAccounts
	if h.NumberOfAccounts, size, err = util.ReadFieldAsInt(line, read); err != nil {
		return 0, fmt.Errorf(fmt.Sprintf(gtParseErrorFmt, "NumberOfAccounts"))
	} else {
		read += size
	}

	// NumberOfRecords
	if h.NumberOfRecords, size, err = util.ReadFieldAsInt(line, read); err != nil {
		return 0, fmt.Errorf(fmt.Sprintf(gtParseErrorFmt, "NumberOfRecords"))
	} else {
		read += size
	}

	if err = h.validate(); err != nil {
		return 0, err
	}

	return read, nil
}

func (h *groupTrailer) string() string {
	var buf bytes.Buffer

	buf.WriteString(fmt.Sprintf("%s,", util.GroupTrailerCode))
	buf.WriteString(fmt.Sprintf("%s,", h.GroupControlTotal))
	buf.WriteString(fmt.Sprintf("%d,", h.NumberOfAccounts))
	buf.WriteString(fmt.Sprintf("%d/", h.NumberOfRecords))

	return buf.String()
}

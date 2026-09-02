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
	ftParseErrorFmt    = "FileTrailer: unable to parse %s"
	ftValidateErrorFmt = "FileTrailer: invalid %s"
)

type fileTrailer struct {
	FileControlTotal string
	NumberOfBanks    int64
	NumberOfRecords  int64
}

func (h *fileTrailer) validate() error {
	if h.FileControlTotal != "" && !util.ValidateAmount(h.FileControlTotal) {
		return fmt.Errorf(ftValidateErrorFmt, "FileControlTotal")
	}

	return nil
}

func (h *fileTrailer) parse(data string) (int, error) {

	var line string
	var err error
	var size, read int

	if length := util.GetSize(data); length < 3 {
		return 0, fmt.Errorf(ftParseErrorFmt, "record")
	} else {
		line = data[:length]
	}

	// RecordCode
	if util.FileTrailerCode != line[:2] {
		return 0, fmt.Errorf(ftParseErrorFmt, "RecordCode")
	}
	read += 3

	// BankControlTotal
	if h.FileControlTotal, size, err = util.ReadField(line, read); err != nil {
		return 0, fmt.Errorf(ftParseErrorFmt, "FileControlTotal")
	} else {
		read += size
	}

	// NumberOfBanks
	if h.NumberOfBanks, size, err = util.ReadFieldAsInt(line, read); err != nil {
		return 0, fmt.Errorf(ftParseErrorFmt, "NumberOfBanks")
	} else {
		read += size
	}

	// NumberOfRecords
	if h.NumberOfRecords, size, err = util.ReadFieldAsInt(line, read); err != nil {
		return 0, fmt.Errorf(ftParseErrorFmt, "NumberOfRecords")
	} else {
		read += size
	}

	if err = h.validate(); err != nil {
		return 0, err
	}

	return read, nil
}

func (h *fileTrailer) string() string {
	var buf bytes.Buffer

	buf.WriteString(fmt.Sprintf("%s,", util.FileTrailerCode))
	buf.WriteString(fmt.Sprintf("%s,", h.FileControlTotal))
	buf.WriteString(fmt.Sprintf("%d,", h.NumberOfBanks))
	buf.WriteString(fmt.Sprintf("%d/", h.NumberOfRecords))

	return buf.String()
}

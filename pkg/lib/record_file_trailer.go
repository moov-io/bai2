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
	ftParseErrorFmt    = "FileTrailer: unable to parse %s"
	ftValidateErrorFmt = "FileTrailer: invalid %s"
)

type fileTrailer struct {
	FileControlTotal string
	NumberOfGroups   int64
	NumberOfRecords  int64
}

func (h *fileTrailer) validate() error {
	if h.FileControlTotal != "" && !util.ValidateAmount(h.FileControlTotal) {
		return fmt.Errorf(fmt.Sprintf(ftValidateErrorFmt, "FileControlTotal"))
	}

	return nil
}

func (h *fileTrailer) parse(data string) (int, error) {

	var line string
	var err error
	var size, read int

	if length := util.GetSize(data); length < 3 {
		return 0, fmt.Errorf(fmt.Sprintf(ftParseErrorFmt, "record"))
	} else {
		line = data[:length]
	}

	// RecordCode
	if util.FileTrailerCode != line[:2] {
		return 0, fmt.Errorf(fmt.Sprintf(fhParseErrorFmt, "RecordCode"))
	}
	read += 3

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

	if err = h.validate(); err != nil {
		return 0, err
	}

	return read, nil
}

func (h *fileTrailer) string() string {
	var buf bytes.Buffer

	buf.WriteString(fmt.Sprintf("%s,", util.FileTrailerCode))
	buf.WriteString(fmt.Sprintf("%s,", h.FileControlTotal))
	buf.WriteString(fmt.Sprintf("%d,", h.NumberOfGroups))
	buf.WriteString(fmt.Sprintf("%d/", h.NumberOfRecords))

	return buf.String()
}

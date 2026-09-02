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
	atParseErrorFmt    = "AccountTrailer: unable to parse %s"
	atValidateErrorFmt = "AccountTrailer: invalid %s"
)

type accountTrailer struct {
	AccountControlTotal string
	NumberRecords       int64
}

func (h *accountTrailer) validate() error {
	if h.AccountControlTotal != "" && !util.ValidateAmount(h.AccountControlTotal) {
		return fmt.Errorf(atValidateErrorFmt, "Amount")
	}

	return nil
}

func (h *accountTrailer) parse(data string) (int, error) {

	var line string
	var err error
	var size, read int

	length := util.GetSize(data)
	if length < 3 {
		return 0, fmt.Errorf(atParseErrorFmt, "record")
	} else {
		line = data[:length]
	}

	// RecordCode
	if util.AccountTrailerCode != data[:2] {
		return 0, fmt.Errorf(atParseErrorFmt, "RecordCode")
	}
	read += 3

	// AccountControlTotal
	if h.AccountControlTotal, size, err = util.ReadField(line, read); err != nil {
		return 0, fmt.Errorf(atParseErrorFmt, "AccountControlTotal")
	} else {
		read += size
	}

	// NumberRecords
	if h.NumberRecords, size, err = util.ReadFieldAsInt(line, read); err != nil {
		return 0, fmt.Errorf(atParseErrorFmt, "NumberRecords")
	} else {
		read += size
	}

	if err = h.validate(); err != nil {
		return 0, err
	}

	return read, nil
}

func (h *accountTrailer) string() string {
	var buf bytes.Buffer

	buf.WriteString(fmt.Sprintf("%s,", util.AccountTrailerCode))
	buf.WriteString(fmt.Sprintf("%s,", h.AccountControlTotal))
	buf.WriteString(fmt.Sprintf("%d/", h.NumberRecords))

	return buf.String()
}

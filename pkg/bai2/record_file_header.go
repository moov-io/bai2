// Copyright 2022 The Moov Authors
// Use of this source code is governed by an Apache License
// license that can be found in the LICENSE file.

package bai2

import (
	"bytes"
	"fmt"

	"github.com/moov-io/bai2/pkg/util"
)

const (
	fhParseErrorFmt    = "FileHeader: unable to parse %s"
	fhValidateErrorFmt = "FileHeader: invalid %s"
)

type fileHeader struct {
	Sender               string
	Receiver             string
	FileCreatedDate      string
	FileCreatedTime      string
	FileIdNumber         string
	PhysicalRecordLength int64 `json:",omitempty"`
	BlockSize            int64 `json:",omitempty"`
	VersionNumber        int64
}

func (h *fileHeader) validate(options Options) error {
	if h.Sender == "" {
		return fmt.Errorf(fhValidateErrorFmt, "Sender")
	}
	if h.Receiver == "" {
		return fmt.Errorf(fhValidateErrorFmt, "Receiver")
	}
	if h.FileCreatedDate == "" {
		return fmt.Errorf(fhValidateErrorFmt, "FileCreatedDate")
	} else if !util.ValidateDate(h.FileCreatedDate) {
		return fmt.Errorf(fhValidateErrorFmt, "FileCreatedDate")
	}
	if h.FileCreatedTime == "" {
		return fmt.Errorf(fhValidateErrorFmt, "FileCreatedTime")
	} else if !util.ValidateTime(h.FileCreatedTime) {
		return fmt.Errorf(fhValidateErrorFmt, "FileCreatedTime")
	}
	if h.FileIdNumber == "" {
		return fmt.Errorf(fhValidateErrorFmt, "FileIdNumber")
	}
	if h.VersionNumber != VersionBAI2 && !options.IgnoreVersion {
		if h.VersionNumber == VersionBTR3 {
			return fmt.Errorf("FileHeader: version 3 (BTR3) is not supported by the BAI2 reader; set Options.IgnoreVersion for BAI2 files stamped with version 3")
		}
		return fmt.Errorf(fhValidateErrorFmt, "VersionNumber")
	}

	return nil
}

func (h *fileHeader) parse(data string, options Options) (int, error) {

	var line string
	var err error
	var size, read int

	if length := util.GetSize(data); length < 3 {
		return 0, fmt.Errorf(fhParseErrorFmt, "record")
	} else {
		line = data[:length]
	}

	// RecordCode
	if util.FileHeaderCode != line[:2] {
		return 0, fmt.Errorf(fhParseErrorFmt, "RecordCode")
	}
	read += 3

	// Sender
	if h.Sender, size, err = util.ReadField(line, read); err != nil {
		return 0, fmt.Errorf(fhParseErrorFmt, "Sender")
	} else {
		read += size
	}

	// Receiver
	if h.Receiver, size, err = util.ReadField(line, read); err != nil {
		return 0, fmt.Errorf(fhParseErrorFmt, "Receiver")
	} else {
		read += size
	}

	// FileCreatedDate
	if h.FileCreatedDate, size, err = util.ReadField(line, read); err != nil {
		return 0, fmt.Errorf(fhParseErrorFmt, "FileCreatedDate")
	} else {
		read += size
	}

	// FileCreatedTime
	if h.FileCreatedTime, size, err = util.ReadField(line, read); err != nil {
		return 0, fmt.Errorf(fhParseErrorFmt, "FileCreatedTime")
	} else {
		read += size
	}

	// FileIdNumber
	if h.FileIdNumber, size, err = util.ReadField(line, read); err != nil {
		return 0, fmt.Errorf(fhParseErrorFmt, "FileIdNumber")
	} else {
		read += size
	}

	// PhysicalRecordLength
	if h.PhysicalRecordLength, size, err = util.ReadFieldAsInt(line, read); err != nil {
		return 0, fmt.Errorf(fhParseErrorFmt, "PhysicalRecordLength")
	} else {
		read += size
	}

	// BlockSize
	if h.BlockSize, size, err = util.ReadFieldAsInt(line, read); err != nil {
		return 0, fmt.Errorf(fhParseErrorFmt, "BlockSize")
	} else {
		read += size
	}

	// VersionNumber
	if h.VersionNumber, size, err = util.ReadFieldAsInt(line, read); err != nil {
		return 0, fmt.Errorf(fhParseErrorFmt, "VersionNumber")
	} else {
		read += size
	}

	if err = h.validate(options); err != nil {
		return 0, err
	}

	return read, nil
}

func (h *fileHeader) string() string {
	var buf bytes.Buffer

	buf.WriteString(fmt.Sprintf("%s,", util.FileHeaderCode))
	buf.WriteString(fmt.Sprintf("%s,", h.Sender))
	buf.WriteString(fmt.Sprintf("%s,", h.Receiver))
	buf.WriteString(fmt.Sprintf("%s,", h.FileCreatedDate))
	buf.WriteString(fmt.Sprintf("%s,", h.FileCreatedTime))
	buf.WriteString(fmt.Sprintf("%s,", h.FileIdNumber))
	if h.PhysicalRecordLength > 0 {
		buf.WriteString(fmt.Sprintf("%d,", h.PhysicalRecordLength))
	} else {
		buf.WriteString(",")
	}
	if h.BlockSize > 0 {
		buf.WriteString(fmt.Sprintf("%d,", h.BlockSize))
	} else {
		buf.WriteString(",")
	}
	buf.WriteString(fmt.Sprintf("%d/", h.VersionNumber))

	return buf.String()
}

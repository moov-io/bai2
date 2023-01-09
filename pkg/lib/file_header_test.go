// Copyright 2022 The Moov Authors
// Use of this source code is governed by an Apache License
// license that can be found in the LICENSE file.

package lib

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func mockFileHeader() *FileHeader {
	return &FileHeader{
		RecordCode:           "01",
		Sender:               "0004",
		Receiver:             "12345",
		FileCreatedDate:      "060321",
		FileCreatedTime:      "0829",
		FileIdNumber:         "001",
		PhysicalRecordLength: 80,
		BlockSize:            1,
		VersionNumber:        2,
	}
}

func TestFileHeader(t *testing.T) {

	record := mockFileHeader()
	require.NoError(t, record.Validate())

	record.VersionNumber = 0
	require.Error(t, record.Validate())
	require.Equal(t, "FileHeader: invalid VersionNumber", record.Validate().Error())

	record.FileIdNumber = ""
	require.Error(t, record.Validate())
	require.Equal(t, "FileHeader: invalid FileIdNumber", record.Validate().Error())

	record.FileCreatedTime = ""
	require.Error(t, record.Validate())
	require.Equal(t, "FileHeader: invalid FileCreatedTime", record.Validate().Error())

	record.FileCreatedDate = ""
	require.Error(t, record.Validate())
	require.Equal(t, "FileHeader: invalid FileCreatedDate", record.Validate().Error())

	record.Receiver = ""
	require.Error(t, record.Validate())
	require.Equal(t, "FileHeader: invalid Receiver", record.Validate().Error())

	record.Sender = ""
	require.Error(t, record.Validate())
	require.Equal(t, "FileHeader: invalid Sender", record.Validate().Error())

	record.RecordCode = ""
	require.Error(t, record.Validate())
	require.Equal(t, "FileHeader: invalid RecordCode", record.Validate().Error())

}

func TestFileHeaderWithOptional(t *testing.T) {

	sample := "01,0004,12345,060321,0829,001,80,1,2/"
	record := NewFileHeader()

	size, err := record.Parse(sample)
	require.NoError(t, err)
	require.Equal(t, 37, size)

	require.Equal(t, "01", record.RecordCode)
	require.Equal(t, "0004", record.Sender)
	require.Equal(t, "12345", record.Receiver)
	require.Equal(t, "060321", record.FileCreatedDate)
	require.Equal(t, "0829", record.FileCreatedTime)
	require.Equal(t, "001", record.FileIdNumber)
	require.Equal(t, int64(80), record.PhysicalRecordLength)
	require.Equal(t, int64(1), record.BlockSize)
	require.Equal(t, int64(2), record.VersionNumber)

	require.Equal(t, sample, record.String())
}

func TestFileHeaderWithoutOptional(t *testing.T) {

	sample := "01,2,12345,060321,0829,1,,,2/"
	record := NewFileHeader()

	size, err := record.Parse(sample)
	require.NoError(t, err)
	require.Equal(t, 29, size)

	require.Equal(t, "01", record.RecordCode)
	require.Equal(t, "2", record.Sender)
	require.Equal(t, "12345", record.Receiver)
	require.Equal(t, "060321", record.FileCreatedDate)
	require.Equal(t, "0829", record.FileCreatedTime)
	require.Equal(t, "1", record.FileIdNumber)
	require.Equal(t, int64(0), record.PhysicalRecordLength)
	require.Equal(t, int64(0), record.BlockSize)
	require.Equal(t, int64(2), record.VersionNumber)

	require.Equal(t, sample, record.String())
}

func TestFileHeaderWithInvalidSample(t *testing.T) {

	record := NewFileHeader()
	_, err := record.Parse("01,2,12345,06032,0829,1,,,2/")
	require.Error(t, err)

	_, err = record.Parse("01,2,12345,060321,082,1,,,2/")
	require.Error(t, err)

	_, err = record.Parse("01,2,12345,060321,082a,1,,,2/")
	require.Error(t, err)
}

func TestFileHeaderWithInvalidSample2(t *testing.T) {

	sample := "01,2,12345,06032,0829,1"
	record := NewAccountIdentifier()

	size, err := record.Parse(sample)
	require.Equal(t, "AccountIdentifier: unable to parse record", err.Error())
	require.Equal(t, 0, size)

	sample = "01,2,12345/"
	size, err = record.Parse(sample)
	require.Equal(t, "AccountIdentifier: unable to parse Amount", err.Error())
	require.Equal(t, 0, size)

}

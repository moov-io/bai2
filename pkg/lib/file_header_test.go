// Copyright 2022 The Moov Authors
// Use of this source code is governed by an Apache License
// license that can be found in the LICENSE file.

package lib

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestFileHeader(t *testing.T) {

	header := NewFileHeader()
	require.NoError(t, header.Validate())

	header.VersionNumber = ""
	require.Error(t, header.Validate())
	require.Equal(t, "FileHeader: invalid version number", header.Validate().Error())

	header.BlockSize = ""
	require.Error(t, header.Validate())
	require.Equal(t, "FileHeader: invalid block size", header.Validate().Error())

	header.PhysicalRecordLength = ""
	require.Error(t, header.Validate())
	require.Equal(t, "FileHeader: invalid physical record length", header.Validate().Error())

	header.Sender = ""
	require.Error(t, header.Validate())
	require.Equal(t, "FileHeader: invalid sender", header.Validate().Error())

	header.RecordCode = ""
	require.Error(t, header.Validate())
	require.Equal(t, "FileHeader: invalid record code", header.Validate().Error())

}

func TestFileHeaderWithSample(t *testing.T) {

	sample := "01,0004,12345,060321,0829,001,80,1,2/"
	header := NewFileHeader()

	err := header.Parse(sample)
	require.NoError(t, err)

	require.Equal(t, "01", header.RecordCode)
	require.Equal(t, "0004", header.Sender)
	require.Equal(t, "12345", header.Receiver)
	require.Equal(t, "060321", header.FileCreatedDate)
	require.Equal(t, "0829", header.FileCreatedTime)
	require.Equal(t, "001", header.FileIdNumber)
	require.Equal(t, "80", header.PhysicalRecordLength)
	require.Equal(t, "1", header.BlockSize)
	require.Equal(t, "2", header.VersionNumber)

	require.Equal(t, sample, header.String())

	header = &FileHeader{}
	require.Error(t, header.Validate())

	err = header.Parse(sample[:30])
	require.Equal(t, "FileHeader: length 30 is too short", err.Error())
}

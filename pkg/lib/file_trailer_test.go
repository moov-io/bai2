// Copyright 2022 The Moov Authors
// Use of this source code is governed by an Apache License
// license that can be found in the LICENSE file.

package lib

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestFileTrailer(t *testing.T) {

	record := NewFileTrailer()
	require.NoError(t, record.Validate())

	record.RecordCode = ""
	require.Error(t, record.Validate())
	require.Equal(t, "FileTrailer: invalid RecordCode", record.Validate().Error())

}

func TestFileTrailerWithSample(t *testing.T) {

	sample := "99,+00000000001280000,1,27/"
	record := NewFileTrailer()

	size, err := record.Parse(sample)
	require.NoError(t, err)
	require.Equal(t, 27, size)

	require.Equal(t, "99", record.RecordCode)
	require.Equal(t, "+00000000001280000", record.FileControlTotal)
	require.Equal(t, int64(1), record.NumberOfGroups)
	require.Equal(t, int64(27), record.NumberOfRecords)

	require.Equal(t, sample, record.String())
}

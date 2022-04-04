// Copyright 2022 The Moov Authors
// Use of this source code is governed by an Apache License
// license that can be found in the LICENSE file.

package lib

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestFileTrailer(t *testing.T) {

	header := NewFileTrailer()
	require.NoError(t, header.Validate())

	header.RecordCode = ""
	require.Error(t, header.Validate())
	require.Equal(t, "FileTrailer: invalid record code", header.Validate().Error())

}

func TestFileTrailerWithSample(t *testing.T) {

	sample := "99,+00000000001280000,000000001,000000027/"
	header := NewFileTrailer()

	err := header.Parse(sample)
	require.NoError(t, err)

	require.Equal(t, "99", header.RecordCode)
	require.Equal(t, "+00000000001280000", header.GroupControlTotal)
	require.Equal(t, int64(1), header.NumberOfGroups)
	require.Equal(t, int64(27), header.NumberOfRecords)

	require.Equal(t, sample, header.String())

	header = &FileTrailer{}
	require.Error(t, header.Validate())

	err = header.Parse(sample[:20])
	require.Equal(t, "FileTrailer: length 20 is too short", err.Error())
}

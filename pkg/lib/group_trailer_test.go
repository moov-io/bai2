// Copyright 2022 The Moov Authors
// Use of this source code is governed by an Apache License
// license that can be found in the LICENSE file.

package lib

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGroupTrailer(t *testing.T) {

	record := NewGroupTrailer()
	require.NoError(t, record.Validate())

	record.RecordCode = ""
	require.Error(t, record.Validate())
	require.Equal(t, "GroupTrailer: invalid RecordCode", record.Validate().Error())

}

func TestGroupTrailerWithSample(t *testing.T) {

	sample := "98,+00000000001280000,2,25/"
	record := NewGroupTrailer()

	size, err := record.Parse(sample)
	require.NoError(t, err)
	require.Equal(t, 27, size)

	require.Equal(t, "98", record.RecordCode)
	require.Equal(t, "+00000000001280000", record.GroupControlTotal)
	require.Equal(t, int64(2), record.NumberOfAccounts)
	require.Equal(t, int64(25), record.NumberOfRecords)

	require.Equal(t, sample, record.String())
}

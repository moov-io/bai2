// Copyright 2022 The Moov Authors
// Use of this source code is governed by an Apache License
// license that can be found in the LICENSE file.

package lib

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGroupTrailer(t *testing.T) {

	header := NewGroupTrailer()
	require.NoError(t, header.Validate())

	header.RecordCode = ""
	require.Error(t, header.Validate())
	require.Equal(t, "GroupTrailer: invalid record code", header.Validate().Error())

}

func TestGroupTrailerWithSample(t *testing.T) {

	sample := "98,+00000000001280000,000000002,000000025/"
	header := NewGroupTrailer()

	err := header.Parse(sample)
	require.NoError(t, err)

	require.Equal(t, "98", header.RecordCode)
	require.Equal(t, "+00000000001280000", header.GroupControlTotal)
	require.Equal(t, int64(2), header.NumberOfAccounts)
	require.Equal(t, int64(25), header.NumberOfRecords)

	require.Equal(t, sample, header.String())

	header = &GroupTrailer{}
	require.Error(t, header.Validate())

	err = header.Parse(sample[:20])
	require.Equal(t, "GroupTrailer: length 20 is too short", err.Error())
}

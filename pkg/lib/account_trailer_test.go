// Copyright 2022 The Moov Authors
// Use of this source code is governed by an Apache License
// license that can be found in the LICENSE file.

package lib

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestAccountTrailer(t *testing.T) {

	header := NewAccountTrailer()
	require.NoError(t, header.Validate())

	header.RecordCode = ""
	require.Error(t, header.Validate())
	require.Equal(t, "AccountTrailer: invalid record code", header.Validate().Error())

}

func TestAccountTrailerWithSample(t *testing.T) {

	sample := "49,+00000000000446000,000000009/"
	header := NewAccountTrailer()

	err := header.Parse(sample)
	require.NoError(t, err)

	require.Equal(t, "49", header.RecordCode)
	require.Equal(t, "+00000000000446000", header.AccountControlTotal)
	require.Equal(t, int64(9), header.NumberRecords)

	require.Equal(t, sample, header.String())

	header = &AccountTrailer{}
	require.Error(t, header.Validate())

	err = header.Parse(sample[:20])
	require.Equal(t, "AccountTrailer: length 20 is too short", err.Error())
}

// Copyright 2022 The Moov Authors
// Use of this source code is governed by an Apache License
// license that can be found in the LICENSE file.

package lib

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestAccountTrailer(t *testing.T) {

	record := NewAccountTrailer()
	require.NoError(t, record.Validate())

	record.RecordCode = ""
	require.Error(t, record.Validate())
	require.Equal(t, "AccountTrailer: invalid RecordCode", record.Validate().Error())

}

func TestAccountTrailerWithSample(t *testing.T) {

	sample := "49,+00000000000446000,9/"
	record := NewAccountTrailer()

	size, err := record.Parse(sample)
	require.NoError(t, err)
	require.Equal(t, 24, size)

	require.Equal(t, "49", record.RecordCode)
	require.Equal(t, "+00000000000446000", record.AccountControlTotal)
	require.Equal(t, int64(9), record.NumberRecords)

	require.Equal(t, sample, record.String())
}

func TestAccountTrailerWithSample2(t *testing.T) {

	sample := "49,+00000000000446000"
	record := NewAccountIdentifier()

	size, err := record.Parse(sample)
	require.Equal(t, "AccountIdentifier: unable to parse record", err.Error())
	require.Equal(t, 0, size)

	sample = "49,+00000000000446000/"
	size, err = record.Parse(sample)
	require.Equal(t, "AccountIdentifier: unable to parse CurrencyCode", err.Error())
	require.Equal(t, 0, size)

}

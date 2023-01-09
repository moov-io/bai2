// Copyright 2022 The Moov Authors
// Use of this source code is governed by an Apache License
// license that can be found in the LICENSE file.

package lib

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func mockAccountIdentifier() *AccountIdentifier {
	return &AccountIdentifier{
		RecordCode:    "03",
		AccountNumber: "0004",
	}
}

func TestAccountIdentifierCurrent(t *testing.T) {

	record := mockAccountIdentifier()
	require.NoError(t, record.Validate())

	record.AccountNumber = ""
	require.Error(t, record.Validate())
	require.Equal(t, "AccountIdentifierCurrent: invalid AccountNumber", record.Validate().Error())

	record.RecordCode = ""
	require.Error(t, record.Validate())
	require.Equal(t, "AccountIdentifierCurrent: invalid RecordCode", record.Validate().Error())

}

func TestAccountIdentifierCurrentWithSample1(t *testing.T) {

	sample := "03,10200123456,CAD,040,+000000000000,,,045,+000000000000,4,0/"
	record := NewAccountIdentifier()

	size, err := record.Parse(sample)
	require.NoError(t, err)
	require.Equal(t, 61, size)

	require.Equal(t, "03", record.RecordCode)
	require.Equal(t, "10200123456", record.AccountNumber)
	require.Equal(t, "CAD", record.CurrencyCode)
	require.Equal(t, "040", record.TypeCode)
	require.Equal(t, "+000000000000", record.Amount)
	require.Equal(t, "045", record.Composite[0])
	require.Equal(t, "+000000000000", record.Composite[1])
	require.Equal(t, "4", record.Composite[2])
	require.Equal(t, "0", record.Composite[3])

	require.Equal(t, sample, record.String())
}

func TestAccountIdentifierCurrentWithSample2(t *testing.T) {

	sample := "03,5765432,,,,,/"
	record := NewAccountIdentifier()

	size, err := record.Parse(sample)
	require.NoError(t, err)
	require.Equal(t, 16, size)

	require.Equal(t, "03", record.RecordCode)
	require.Equal(t, "5765432", record.AccountNumber)

	require.Equal(t, sample, record.String())
}

func TestAccountIdentifierCurrentWithSample3(t *testing.T) {

	sample := "03,5765432,,,,,,/"
	record := NewAccountIdentifier()

	size, err := record.Parse(sample)
	require.NoError(t, err)
	require.Equal(t, 17, size)

	require.Equal(t, "03", record.RecordCode)
	require.Equal(t, "5765432", record.AccountNumber)
	require.Equal(t, 1, len(record.Composite))

	require.Equal(t, sample, record.String())
}

func TestAccountIdentifierCurrentWithSample4(t *testing.T) {

	sample := "03,5765432,"
	record := NewAccountIdentifier()

	size, err := record.Parse(sample)
	require.Equal(t, "AccountIdentifier: unable to parse record", err.Error())
	require.Equal(t, 0, size)

	sample = "03,5765432,/"
	size, err = record.Parse(sample)
	require.Equal(t, "AccountIdentifier: unable to parse TypeCode", err.Error())
	require.Equal(t, 0, size)

}

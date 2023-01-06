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

	header := mockAccountIdentifier()
	require.NoError(t, header.Validate())

	header.AccountNumber = ""
	require.Error(t, header.Validate())
	require.Equal(t, "AccountIdentifierCurrent: invalid AccountNumber", header.Validate().Error())

	header.RecordCode = ""
	require.Error(t, header.Validate())
	require.Equal(t, "AccountIdentifierCurrent: invalid RecordCode", header.Validate().Error())

}

func TestAccountIdentifierCurrentWithSample1(t *testing.T) {

	sample := "03,10200123456,CAD,040,+000000000000,,,045,+000000000000,4,0/"
	header := NewAccountIdentifier()

	size, err := header.Parse(sample)
	require.NoError(t, err)
	require.Equal(t, 61, size)

	require.Equal(t, "03", header.RecordCode)
	require.Equal(t, "10200123456", header.AccountNumber)
	require.Equal(t, "CAD", header.CurrencyCode)
	require.Equal(t, "040", header.TypeCode)
	require.Equal(t, "+000000000000", header.Amount)
	require.Equal(t, "045", header.Composite[0])
	require.Equal(t, "+000000000000", header.Composite[1])
	require.Equal(t, "4", header.Composite[2])
	require.Equal(t, "0", header.Composite[3])

	require.Equal(t, sample, header.String())
}

func TestAccountIdentifierCurrentWithSample2(t *testing.T) {

	sample := "03,5765432,,,,,/"
	header := NewAccountIdentifier()

	size, err := header.Parse(sample)
	require.NoError(t, err)
	require.Equal(t, 16, size)

	require.Equal(t, "03", header.RecordCode)
	require.Equal(t, "5765432", header.AccountNumber)

	require.Equal(t, sample, header.String())
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

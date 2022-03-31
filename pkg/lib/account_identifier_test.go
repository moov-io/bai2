// Copyright 2022 The Moov Authors
// Use of this source code is governed by an Apache License
// license that can be found in the LICENSE file.

package lib

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestAccountIdentifier(t *testing.T) {

	header := NewAccountIdentifier()
	require.NoError(t, header.Validate())

	header.TypeCode2 = ""
	require.Error(t, header.Validate())
	require.Equal(t, "AccountIdentifier: invalid type code", header.Validate().Error())

	header.TypeCode1 = ""
	require.Error(t, header.Validate())
	require.Equal(t, "AccountIdentifier: invalid type code", header.Validate().Error())

	header.CurrencyCode = ""
	require.Error(t, header.Validate())
	require.Equal(t, "AccountIdentifier: invalid currency code", header.Validate().Error())

	header.RecordCode = ""
	require.Error(t, header.Validate())
	require.Equal(t, "AccountIdentifier: invalid record code", header.Validate().Error())

}

func TestAccountIdentifierWithSample(t *testing.T) {

	sample := "03,10200123456,CAD,040,+000000000000,,,045,+000000000000,,/"
	header := NewAccountIdentifier()

	err := header.Parse(sample)
	require.NoError(t, err)

	require.Equal(t, "03", header.RecordCode)
	require.Equal(t, "10200123456", header.AccountNumber)
	require.Equal(t, "CAD", header.CurrencyCode)
	require.Equal(t, "040", header.TypeCode1)
	require.Equal(t, "+000000000000", header.OpeningBalance)
	require.Equal(t, "045", header.TypeCode2)
	require.Equal(t, "+000000000000", header.ClosingBalance)

	require.Equal(t, sample, header.String())

	header = &AccountIdentifier{}
	require.Error(t, header.Validate())

	err = header.Parse(sample[:20])
	require.Equal(t, "AccountIdentifier: length 20 is too short", err.Error())
}

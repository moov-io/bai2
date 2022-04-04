// Copyright 2022 The Moov Authors
// Use of this source code is governed by an Apache License
// license that can be found in the LICENSE file.

package lib

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestAccountIdentifierLoan(t *testing.T) {

	header := NewAccountIdentifierLoan()
	require.NoError(t, header.Validate())

	header.TypeCode2 = ""
	require.Error(t, header.Validate())
	require.Equal(t, "AccountIdentifierLoan: invalid type code", header.Validate().Error())

	header.TypeCode1 = ""
	require.Error(t, header.Validate())
	require.Equal(t, "AccountIdentifierLoan: invalid type code", header.Validate().Error())

	header.CurrencyCode = ""
	require.Error(t, header.Validate())
	require.Equal(t, "AccountIdentifierLoan: invalid currency code", header.Validate().Error())

	header.RecordCode = ""
	require.Error(t, header.Validate())
	require.Equal(t, "AccountIdentifierLoan: invalid record code", header.Validate().Error())

}

func TestAccountIdentifierLoanWithSample(t *testing.T) {

	sample := "03,1020123456701,CAD,056,+000000000000,,V,060201,,056,+000000000000,,V,060228,/"
	header := NewAccountIdentifierLoan()

	err := header.Parse(sample)
	require.NoError(t, err)

	require.Equal(t, "03", header.RecordCode)
	require.Equal(t, "1020123456701", header.AccountNumber)
	require.Equal(t, "CAD", header.CurrencyCode)
	require.Equal(t, "056", header.TypeCode1)
	require.Equal(t, "+000000000000", header.OpeningBalance)
	require.Equal(t, "060201", header.ValueDate1)
	require.Equal(t, "V", header.FundsType1)
	require.Equal(t, "056", header.TypeCode2)
	require.Equal(t, "+000000000000", header.ClosingBalance)
	require.Equal(t, "060228", header.ValueDate2)
	require.Equal(t, "V", header.FundsType2)

	require.Equal(t, sample, header.String())

	header = &AccountIdentifierLoan{}
	require.Error(t, header.Validate())

	err = header.Parse(sample[:20])
	require.Equal(t, "AccountIdentifierLoan: length 20 is too short", err.Error())
}

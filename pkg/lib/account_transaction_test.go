// Copyright 2022 The Moov Authors
// Use of this source code is governed by an Apache License
// license that can be found in the LICENSE file.

package lib

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestAccountTransaction(t *testing.T) {

	header := NewAccountTransaction()
	require.NoError(t, header.Validate())

	header.TypeCode = ""
	require.Error(t, header.Validate())
	require.Equal(t, "AccountTransaction: invalid type code", header.Validate().Error())

	header.RecordCode = ""
	require.Error(t, header.Validate())
	require.Equal(t, "AccountTransaction: invalid record code", header.Validate().Error())

}

func TestAccountTransactionWithSample(t *testing.T) {

	sample := "16,409,000000000002500,V,060316,,,,RETURNED CHEQUE     /"
	header := NewAccountTransaction()

	err := header.Parse(sample)
	require.NoError(t, err)

	require.Equal(t, "16", header.RecordCode)
	require.Equal(t, "409", header.TypeCode)
	require.Equal(t, "000000000002500", header.Amount)
	require.Equal(t, "V", header.FundsType)
	require.Equal(t, "060316", header.ValueDate)
	require.Equal(t, "RETURNED CHEQUE     ", header.Description)

	require.Equal(t, sample, header.String())

	header = &AccountTransaction{}
	require.Error(t, header.Validate())

	err = header.Parse(sample[:20])
	require.Equal(t, "AccountTransaction: length 20 is too short", err.Error())
}

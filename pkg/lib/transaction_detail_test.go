// Copyright 2022 The Moov Authors
// Use of this source code is governed by an Apache License
// license that can be found in the LICENSE file.

package lib

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestTransactionDetail(t *testing.T) {

	record := TransactionDetail{
		RecordCode: "16",
		TypeCode:   "890",
	}
	require.NoError(t, record.Validate())

	record.TypeCode = "AAA"
	require.Error(t, record.Validate())
	require.Equal(t, "AccountTransaction: invalid TypeCode", record.Validate().Error())

	record.RecordCode = ""
	require.Error(t, record.Validate())
	require.Equal(t, "AccountTransaction: invalid RecordCode", record.Validate().Error())

}

func TestTransactionDetailWithSample(t *testing.T) {

	sample := "16,409,000000000002500,V,060316,,,,RETURNED CHEQUE     /"
	record := NewTransactionDetail()

	size, err := record.Parse(sample)
	require.NoError(t, err)
	require.Equal(t, 56, size)

	require.Equal(t, "16", record.RecordCode)
	require.Equal(t, "409", record.TypeCode)
	require.Equal(t, "000000000002500", record.Amount)
	require.Equal(t, "V", record.FundsType)
	require.Equal(t, 5, len(record.Composite))
	require.Equal(t, "060316", record.Composite[0])
	require.Equal(t, "RETURNED CHEQUE     ", record.Composite[4])

	require.Equal(t, sample, record.String())
}

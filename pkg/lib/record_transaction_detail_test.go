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
		TypeCode: "890",
	}
	require.NoError(t, record.validate())

	record.TypeCode = "AAA"
	require.Error(t, record.validate())
	require.Equal(t, "AccountTransaction: invalid TypeCode", record.validate().Error())

}

func TestTransactionDetailWithSample(t *testing.T) {

	sample := "16,409,000000000002500,V,060316,,,,RETURNED CHEQUE     /"
	record := NewTransactionDetail()

	size, err := record.parse(sample)
	require.NoError(t, err)
	require.Equal(t, 56, size)

	require.Equal(t, "409", record.TypeCode)
	require.Equal(t, "000000000002500", record.Amount)
	require.Equal(t, "V", record.FundsType)
	require.Equal(t, 5, len(record.Composite))
	require.Equal(t, "060316", record.Composite[0])
	require.Equal(t, "RETURNED CHEQUE     ", record.Composite[4])

	require.Equal(t, sample, record.string())
}

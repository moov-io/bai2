// Copyright 2022 The Moov Authors
// Use of this source code is governed by an Apache License
// license that can be found in the LICENSE file.

package lib

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestTransferTransaction(t *testing.T) {

	header := NewTransferTransaction()
	require.NoError(t, header.Validate())

	header.TypeCode = ""
	require.Error(t, header.Validate())
	require.Equal(t, "TransferTransaction: invalid type code", header.Validate().Error())

	header.RecordCode = ""
	require.Error(t, header.Validate())
	require.Equal(t, "TransferTransaction: invalid record code", header.Validate().Error())

}

func TestTransferTransactionWithSample(t *testing.T) {

	sample := "16,108,0000000254671,V,060316,,01023,123456789,STORE NAME          /"
	header := NewTransferTransaction()

	err := header.Parse(sample)
	require.NoError(t, err)

	require.Equal(t, "16", header.RecordCode)
	require.Equal(t, "108", header.TypeCode)
	require.Equal(t, "0000000254671", header.Amount)
	require.Equal(t, "V", header.FundsType)
	require.Equal(t, "060316", header.ValueDate)
	require.Equal(t, "01023", header.BankReference)
	require.Equal(t, "123456789", header.CustomerReference)
	require.Equal(t, "STORE NAME          ", header.Description)

	require.Equal(t, sample, header.String())

	header = &TransferTransaction{}
	require.Error(t, header.Validate())

	err = header.Parse(sample[:20])
	require.Equal(t, "TransferTransaction: length 20 is too short", err.Error())
}

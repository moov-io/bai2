// Copyright 2022 The Moov Authors
// Use of this source code is governed by an Apache License
// license that can be found in the LICENSE file.

package lib

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestAccountIdentifierContinuation(t *testing.T) {

	header := NewAccountIdentifierContinuation()
	require.NoError(t, header.Validate())

	header.TypeCode2 = ""
	require.Error(t, header.Validate())
	require.Equal(t, "AccountIdentifierContinuation: invalid type code", header.Validate().Error())

	header.TypeCode1 = ""
	require.Error(t, header.Validate())
	require.Equal(t, "AccountIdentifierContinuation: invalid type code", header.Validate().Error())

	header.RecordCode = ""
	require.Error(t, header.Validate())
	require.Equal(t, "AccountIdentifierContinuation: invalid record code", header.Validate().Error())

}

func TestAccountIdentifierContinuationWithSample(t *testing.T) {

	sample := "88,100,000000000208500,00003,V,060316,,400,000000000208500,00008,V,060316,/"
	header := NewAccountIdentifierContinuation()

	err := header.Parse(sample)
	require.NoError(t, err)

	require.Equal(t, "88", header.RecordCode)
	require.Equal(t, "100", header.TypeCode1)
	require.Equal(t, "000000000208500", header.TotalCreditAmount1)
	require.Equal(t, int64(3), header.TotalOfCredits1)
	require.Equal(t, "V", header.FundsType1)
	require.Equal(t, "060316", header.ValueDate1)
	require.Equal(t, "400", header.TypeCode2)
	require.Equal(t, "000000000208500", header.TotalCreditAmount2)
	require.Equal(t, int64(8), header.TotalOfCredits2)
	require.Equal(t, "V", header.FundsType2)
	require.Equal(t, "060316", header.ValueDate2)

	require.Equal(t, sample, header.String())

	header = &AccountIdentifierContinuation{}
	require.Error(t, header.Validate())

	err = header.Parse(sample[:20])
	require.Equal(t, "AccountIdentifierContinuation: length 20 is too short", err.Error())
}

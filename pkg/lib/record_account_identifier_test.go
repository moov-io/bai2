// Copyright 2022 The Moov Authors
// Use of this source code is governed by an Apache License
// license that can be found in the LICENSE file.

package lib

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func mockAccountIdentifier() *accountIdentifier {
	return &accountIdentifier{
		AccountNumber: "0004",
	}
}

func TestAccountIdentifierCurrent(t *testing.T) {

	record := mockAccountIdentifier()
	require.NoError(t, record.validate())

	record.AccountNumber = ""
	require.Error(t, record.validate())
	require.Equal(t, "AccountIdentifierCurrent: invalid AccountNumber", record.validate().Error())

}

func TestAccountIdentifierCurrentWithSample1(t *testing.T) {

	sample := "03,10200123456,CAD,040,+000000000000,,,045,+000000000000,4,0/"
	record := accountIdentifier{}

	size, err := record.parse(sample)
	require.NoError(t, err)
	require.Equal(t, 61, size)

	require.Equal(t, "10200123456", record.AccountNumber)
	require.Equal(t, "CAD", record.CurrencyCode)
	require.Equal(t, 2, len(record.Summaries))

	summary := record.Summaries[0]
	require.Equal(t, "040", summary.TypeCode)
	require.Equal(t, "+000000000000", summary.Amount)
	require.Equal(t, int64(0), summary.ItemCount)
	require.Equal(t, "", string(summary.FundsType.TypeCode))

	summary = record.Summaries[1]

	require.Equal(t, "045", summary.TypeCode)
	require.Equal(t, "+000000000000", summary.Amount)
	require.Equal(t, int64(4), summary.ItemCount)
	require.Equal(t, FundsType0, string(summary.FundsType.TypeCode))

	require.Equal(t, sample, record.string())
}

func TestAccountIdentifierCurrentWithSample2(t *testing.T) {

	sample := "03,5765432,,,,,/"
	record := accountIdentifier{}

	size, err := record.parse(sample)
	require.NoError(t, err)
	require.Equal(t, 16, size)

	require.Equal(t, "5765432", record.AccountNumber)
	require.Equal(t, 1, len(record.Summaries))

	require.Equal(t, sample, record.string())
}

func TestAccountIdentifierOutputWithContinuationRecords(t *testing.T) {

	record := accountIdentifier{
		AccountNumber: "10200123456",
		CurrencyCode:  "CAD",
	}

	for i := 0; i < 10; i++ {
		record.Summaries = append(record.Summaries,
			AccountSummary{
				TypeCode:  "040",
				Amount:    "+000000000000",
				ItemCount: 10,
			})
	}

	result := record.string()
	expectResult := `03,10200123456,CAD,040,+000000000000,10,,040,+000000000000,10,,040,+000000000000,10,,040,+000000000000,10,,040,+000000000000,10,,040,+000000000000,10,,040,+000000000000,10,,040,+000000000000,10,,040,+000000000000,10,,040,+000000000000,10,/`
	require.Equal(t, expectResult, result)
	require.Equal(t, len(expectResult), len(result))

	result = record.string(80)
	expectResult = `03,10200123456,CAD,040,+000000000000,10,,040,+000000000000,10,,040/
88,+000000000000,10,,040,+000000000000,10,,040,+000000000000,10,,040/
88,+000000000000,10,,040,+000000000000,10,,040,+000000000000,10,,040/
88,+000000000000,10,,040,+000000000000,10,/`
	require.Equal(t, expectResult, result)
	require.Equal(t, len(expectResult), len(result))

	result = record.string(50)
	expectResult = `03,10200123456,CAD,040,+000000000000,10,,040/
88,+000000000000,10,,040,+000000000000,10,,040/
88,+000000000000,10,,040,+000000000000,10,,040/
88,+000000000000,10,,040,+000000000000,10,,040/
88,+000000000000,10,,040,+000000000000,10,,040/
88,+000000000000,10,/`
	require.Equal(t, expectResult, result)
	require.Equal(t, len(expectResult), len(result))
}

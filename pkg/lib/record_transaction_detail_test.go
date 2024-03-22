// Copyright 2022 The Moov Authors
// Use of this source code is governed by an Apache License
// license that can be found in the LICENSE file.

package lib

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestTransactionDetail(t *testing.T) {

	record := transactionDetail{
		TypeCode: "890",
	}
	require.NoError(t, record.validate())

	record.TypeCode = "AAA"
	require.Error(t, record.validate())
	require.Equal(t, "TransactionDetail: invalid TypeCode", record.validate().Error())

}

func TestTransactionDetailWithSample(t *testing.T) {

	sample := "16,409,000000000002500,V,060316,,,,RETURNED CHEQUE     /"
	record := transactionDetail{
		TypeCode: "890",
	}

	size, err := record.parse(sample)
	require.NoError(t, err)
	require.Equal(t, 56, size)

	require.Equal(t, "409", record.TypeCode)
	require.Equal(t, "000000000002500", record.Amount)
	require.Equal(t, "V", string(record.FundsType.TypeCode))
	require.Equal(t, "060316", record.FundsType.Date)
	require.Equal(t, "", record.FundsType.Time)
	require.Equal(t, "", record.BankReferenceNumber)
	require.Equal(t, "", record.CustomerReferenceNumber)
	require.Equal(t, "RETURNED CHEQUE     ", record.Text)

	require.Equal(t, sample, record.string())
}

func TestTransactionDetailOutputWithContinuationRecords(t *testing.T) {

	record := transactionDetail{
		TypeCode:                "409",
		Amount:                  "111111111111111",
		BankReferenceNumber:     "222222222222222",
		CustomerReferenceNumber: "333333333333333",
		Text:                    "RETURNED CHEQUE     444444444444444",
		FundsType: FundsType{
			TypeCode:           FundsTypeD,
			DistributionNumber: 5,
			Distributions: []Distribution{
				{
					Day:    1,
					Amount: 1000000000,
				},
				{
					Day:    2,
					Amount: 2000000000,
				},
				{
					Day:    3,
					Amount: 3000000000,
				},
				{
					Day:    4,
					Amount: 4000000000,
				},
				{
					Day:    5,
					Amount: 5000000000,
				},
				{
					Day:    6,
					Amount: 6000000000,
				},
				{
					Day:    7,
					Amount: 7000000000,
				},
			},
		},
	}

	result := record.string()
	expectResult := `16,409,111111111111111,D,5,1,1000000000,2,2000000000,3,3000000000,4,4000000000,5,5000000000,6,6000000000,7,7000000000,222222222222222,333333333333333,RETURNED CHEQUE     444444444444444/`
	require.Equal(t, expectResult, result)
	require.Equal(t, len(expectResult), len(result))

	result = record.string(80)
	expectResult = `16,409,111111111111111,D,5,1,1000000000,2,2000000000,3,3000000000,4,4000000000/
88,5,5000000000,6,6000000000,7,7000000000,222222222222222,333333333333333/
88,RETURNED CHEQUE     444444444444444/`
	require.Equal(t, expectResult, result)
	require.Equal(t, len(expectResult), len(result))

	result = record.string(50)
	expectResult = `16,409,111111111111111,D,5,1,1000000000,2/
88,2000000000,3,3000000000,4,4000000000,5/
88,5000000000,6,6000000000,7,7000000000/
88,222222222222222,333333333333333/
88,RETURNED CHEQUE     444444444444444/`
	require.Equal(t, expectResult, result)
	require.Equal(t, len(expectResult), len(result))

}

/**
 * Outlines the behavior of a Detail record when the Detail and Continuations for the detail are terminated
 * by a newline character ("\n") rather than a slash ("/").
 *
 * Note: continuation parsing is implemented in `detail.go`, which is why this particular test doesn't parse
 * all of the continuation lines.
 */
 func TestTransactionDetailOutput_ContinuationRecordWithNewlineDelimiter(t *testing.T) {
	data := `16,266,1912,,GI2118700002010,20210706MMQFMPU8000001,Outgoing Wire Return,-
88,CREF: 20210706MMQFMPU8000001
88,EREF: 20210706MMQFMPU8000001
88,DBIC: GSCRUS33
88,CRNM: ABC Company
88,DBNM: SAMPLE INC.`

	record := transactionDetail{}

	size, err := record.parse(data)
	require.NoError(t, err)

	require.Equal(t, "266", record.TypeCode)
	require.Equal(t, "1912", record.Amount)
	require.Equal(t, "", string(record.FundsType.TypeCode))
	require.Equal(t, "", record.FundsType.Date)
	require.Equal(t, "", record.FundsType.Time)
	require.Equal(t, "GI2118700002010", record.BankReferenceNumber)
	require.Equal(t, "20210706MMQFMPU8000001", record.CustomerReferenceNumber)
	require.Equal(t, "Outgoing Wire Return,-", record.Text)
	require.Equal(t, 75, size)

	result := record.string()
	expectResult := `16,266,1912,,GI2118700002010,20210706MMQFMPU8000001,Outgoing Wire Return,-/`
	require.Equal(t, expectResult, result)
}
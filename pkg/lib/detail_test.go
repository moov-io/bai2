// Copyright 2022 The Moov Authors
// Use of this source code is governed by an Apache License
// license that can be found in the LICENSE file.

package lib

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestDetail(t *testing.T) {
	raw := `16,115,10000000,S,5000000,4000000,1000000/
88,AX13612,B096132,AMALGAMATED CORP. LOCKBOX/
88,DEPOSIT-MISC. RECEIVABLES/`

	scan := NewBai2Scanner(bytes.NewReader([]byte(raw)))
	detail := NewDetail()

	err := detail.Read(&scan, false)
	require.NoError(t, err)

	require.Equal(t, "115", detail.TypeCode)
	require.Equal(t, "10000000", detail.Amount)
	require.Equal(t, "S", string(detail.FundsType.TypeCode))
	require.Equal(t, int64(5000000), detail.FundsType.ImmediateAmount)
	require.Equal(t, int64(4000000), detail.FundsType.OneDayAmount)
	require.Equal(t, int64(1000000), detail.FundsType.TwoDayAmount)

	expect := `16,115,10000000,S,5000000,4000000,1000000,AX13612,B096132,AMALGAMATED CORP. LOCKBOX,DEPOSIT-MISC. RECEIVABLES/`
	require.Equal(t, expect, detail.String())

	expect = `16,115,10000000,S,5000000,4000000,1000000/
88,AX13612,B096132,AMALGAMATED CORP. LOCKBOX/
88,DEPOSIT-MISC. RECEIVABLES/`
	require.Equal(t, expect, detail.String(50))
}

/**
 * Outlines the behavior of a Detail record when the Detail and Continuations for the detail are terminated
 * by a newline character ("\n") rather than a slash ("/").
 */
 func TestDetail_ContinuationRecordWithNewlineDelimiter(t *testing.T) {
	data := `16,266,1912,,GI2118700002010,20210706MMQFMPU8000001,Outgoing Wire Return,-
88,CREF: 20210706MMQFMPU8000001
88,EREF: 20210706MMQFMPU8000001
88,DBIC: GSCRUS33
88,CRNM: ABC Company
88,DBNM: SAMPLE INC.`

	scan := NewBai2Scanner(bytes.NewReader([]byte(data)))
	detail := NewDetail()

	err := detail.Read(&scan, false)
	require.NoError(t, err)

	require.Equal(t, "266", detail.TypeCode)
	require.Equal(t, "1912", detail.Amount)
	require.Equal(t, "", string(detail.FundsType.TypeCode))
	require.Equal(t, "", detail.FundsType.Date)
	require.Equal(t, "", detail.FundsType.Time)
	require.Equal(t, "GI2118700002010", detail.BankReferenceNumber)
	require.Equal(t, "20210706MMQFMPU8000001", detail.CustomerReferenceNumber)
	require.Equal(t, "Outgoing Wire Return,,CREF: 20210706MMQFMPU800000,EREF: 20210706MMQFMPU800000,DBIC: GSCRUS3,CRNM: ABC Compan,DBNM: SAMPLE INC.", detail.Text)

	expectResult := `16,266,1912,,GI2118700002010,20210706MMQFMPU8000001,Outgoing Wire Return,,CREF: 20210706MMQFMPU800000,EREF: 20210706MMQFMPU800000,DBIC: GSCRUS3,CRNM: ABC Compan,DBNM: SAMPLE INC./`
	require.Equal(t, expectResult, detail.String())

	expectResult = `16,266,1912,,GI2118700002010/
88,20210706MMQFMPU8000001,Outgoing Wire Return,/
88,CREF: 20210706MMQFMPU800000/
88,EREF: 20210706MMQFMPU800000,DBIC: GSCRUS3/
88,CRNM: ABC Compan,DBNM: SAMPLE INC./`
	require.Equal(t, expectResult, detail.String(50))
}
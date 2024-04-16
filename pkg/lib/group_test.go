// Copyright 2022 The Moov Authors
// Use of this source code is governed by an Apache License
// license that can be found in the LICENSE file.

package lib

import (
	"bytes"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestGroupWithSampleData1(t *testing.T) {

	raw := `
02,12345,0004,1,060317,,CAD,/
03,10200123456,CAD,040,+000000000000,,,045,+000000000000,,/
88,100,000000000208500,00003,V,060316,,400,000000000208500,00008,V,060316,/
16,409,000000000002500,V,060316,,,,RETURNED CHEQUE     /
49,+00000000000834000,14/
03,10200123456,CAD,040,+000000000000,,,045,+000000000000,,/
16,409,000000000000500,V,060317,,,,GALERIES RICHELIEU  /
88,100,000000000208500,00003,V,060316,,400,000000000208500,00008,V,060316,/
49,+00000000000446000,9/
98,+00000000001280000,2,25/
`

	group := Group{}
	scan := NewBai2Scanner(bytes.NewReader([]byte(raw)))

	err := group.Read(&scan, false)
	require.NoError(t, err)
	require.NoError(t, group.Validate())
	require.Equal(t, 10, scan.GetLineIndex())
}

func TestGroupWithSampleData2(t *testing.T) {

	raw := `
02,12345,0004,1,060317,,CAD,/
03,10200123456,CAD,040,+000000000000,,,045,+000000000000,,/
88,100,000000000208500,00003,V,060316,,400,000000000208500,00008,V,060316,/
16,409,000000000002500,V,060316,,,,RETURNED CHEQUE     /
49,+00000000000834000,14/
03,10200123456,CAD,040,+000000000000,,,045,+000000000000,,/
16,409,000000000000500,V,060317,,,,GALERIES RICHELIEU  /
88,100,000000000208500,00003,V,060316,,400,000000000208500,00008,V,060316,/
49,+00000000000446000,9/
98,+00000000001280000,2,25/
99,+00000000001280000,1,27/
`

	group := Group{}
	scan := NewBai2Scanner(bytes.NewReader([]byte(raw)))

	err := group.Read(&scan, false)
	require.NoError(t, err)
	require.NoError(t, group.Validate())
	require.Equal(t, 10, scan.GetLineIndex())
	require.Equal(t, "98,+00000000001280000,2,25/", scan.GetLine())
}

func TestSumGroupRecords(t *testing.T) {
	group := Group{}
	group.Accounts = []Account{
		{NumberRecords: 2},
		{NumberRecords: 3},
		{NumberRecords: 4},
	}
	require.Equal(t, int64(11), group.SumRecords())
}

func TestSumNumberOfAccounts(t *testing.T) {
	group := Group{}
	group.Accounts = []Account{
		{NumberRecords: 2},
		{NumberRecords: 3},
		{NumberRecords: 4},
	}
	require.Equal(t, int64(3), group.SumNumberOfAccounts())
}

func TestSumAccountControlTotals(t *testing.T) {
	group := Group{}
	group.Receiver = "121000358"
	group.Originator = "121000358"
	group.GroupStatus = 1
	group.AsOfDate = time.Now().Format("060102")
	group.AsOfTime = time.Now().Format("1504")
	group.AsOfDateModifier = 2
	group.Accounts = []Account{
		{AccountControlTotal: "100", AccountNumber: "9876543210"},
		{AccountControlTotal: "-100", AccountNumber: "9876543210"},
		{AccountControlTotal: "200", AccountNumber: "9876543210"},
	}
	total, err := group.SumAccountControlTotals()
	require.NoError(t, err)
	require.Equal(t, "200", total)
}

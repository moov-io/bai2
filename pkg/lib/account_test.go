// Copyright 2022 The Moov Authors
// Use of this source code is governed by an Apache License
// license that can be found in the LICENSE file.

package lib

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/require"
)

/*
func TestAccountWithSampleData1(t *testing.T) {

	raw := `
03,9876543210,,010,-500000,,,100,1000000,,,400,2000000,,,190/
88,500000,,,110,1000000,,,072,500000,,,074,500000,,,040/
88,-1500000,,/
16,115,500000,S,,200000,300000,,,LOCK BOX NO.68751/
49,4000000,5/
98,+00000000001280000,2,25/
`

	scan := NewBai2Scanner(bytes.NewReader([]byte(raw)))
	account := Account{}
	err := account.Read(&scan, false)
	require.NoError(t, err)
	require.NoError(t, account.Validate())
	require.Equal(t, 5, scan.GetLineIndex())
	require.Equal(t, "", scan.GetLine())
}

func TestAccountWithSampleData2(t *testing.T) {

	raw := `
03,10200123456,CAD,040,+000000000000,,,045,+000000000000,,/
88,100,000000000111500,00002,V,060317,,400,000000000111500,00004,V,060317,/
16,108,000000000011500,V,060317,,,,TFR 1020 0345678    /
16,108,000000000100000,V,060317,,,,MONTREAL            /
98,+00000000001280000,2,25/
`

	scan := NewBai2Scanner(bytes.NewReader([]byte(raw)))
	account := Account{}
	err := account.Read(&scan, false)
	require.NoError(t, err)
	require.NoError(t, account.Validate())
	require.Equal(t, 5, scan.GetLineIndex())
}
*/

func TestAccountOutputWithContinuationRecord(t *testing.T) {

	raw := `
03,9876543210,,010,-500000,,,100,1000000,,,400,2000000,,,190/
88,500000,,,110,1000000,,,072,500000,,,074,500000,,,040/
88,-1500000,,/
16,115,500000,S,,200000,300000,,,LOCK BOX NO.68751/
49,4000000,5/
`

	scan := NewBai2Scanner(bytes.NewReader([]byte(raw)))
	account := Account{}
	err := account.Read(&scan, false)
	require.NoError(t, err)
	require.NoError(t, account.Validate())
	require.Equal(t, 5, scan.GetLineIndex())
	require.Equal(t, "49,4000000,5/", scan.GetLine())

	result := account.String()
	expectedResult := `03,9876543210,,010,-500000,,,100,1000000,,,400,2000000,,,190,500000,,,110,1000000,,,072,500000,,,074,500000,,,040,-1500000,,/
16,115,500000,S,0,200000,300000,,,LOCK BOX NO.68751/
49,4000000,5/`
	require.Equal(t, expectedResult, result)

	result = account.String(50)
	expectedResult = `03,9876543210,,010,-500000,,,100,1000000,,,400/
88,2000000,,,190,500000,,,110,1000000,,,072/
88,500000,,,074,500000,,,040,-1500000,,/
16,115,500000,S,0,200000,300000,,/
88,LOCK BOX NO.68751/
49,4000000,5/`
	require.Equal(t, expectedResult, result)

}

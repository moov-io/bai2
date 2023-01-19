// Copyright 2022 The Moov Authors
// Use of this source code is governed by an Apache License
// license that can be found in the LICENSE file.

package lib

import (
	"bytes"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestAccountWithSampleData1(t *testing.T) {

	raw := `
03,10200123456,CAD,040,+000000000000,,,045,+000000000000,,/
88,100,000000000111500,00002,V,060317,,400,000000000111500,00004,V,060317,/
16,108,000000000011500,V,060317,,,,TFR 1020 0345678    /
16,108,000000000100000,V,060317,,,,MONTREAL            /
49,+00000000000446000,9/
98,+00000000001280000,2,25/
`

	account := Account{}
	lineNum, line, err := account.Read(NewBai2Scanner(bytes.NewReader([]byte(raw))), "", 0)
	require.NoError(t, err)
	require.NoError(t, account.Validate())
	require.Equal(t, 5, lineNum)
	require.Equal(t, "", line)
}

func TestAccountWithSampleData2(t *testing.T) {

	raw := `
03,10200123456,CAD,040,+000000000000,,,045,+000000000000,,/
88,100,000000000111500,00002,V,060317,,400,000000000111500,00004,V,060317,/
16,108,000000000011500,V,060317,,,,TFR 1020 0345678    /
16,108,000000000100000,V,060317,,,,MONTREAL            /
98,+00000000001280000,2,25/
`

	account := Account{}
	lineNum, line, err := account.Read(NewBai2Scanner(bytes.NewReader([]byte(raw))), "", 0)
	require.NoError(t, err)
	require.NoError(t, account.Validate())
	require.Equal(t, 4, lineNum)
	require.Equal(t, "98,+00000000001280000,2,25/", line)
}

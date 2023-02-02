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

	expect := `16,115,10000000,S,5000000,4000000,1000000,AX13612,B096132,AMALGAMATED CORP. LOCKBOX/`
	require.Equal(t, expect, detail.String())

	expect = `16,115,10000000,S,5000000,4000000,1000000/
88,AX13612,B096132,AMALGAMATED CORP. LOCKBOX/`
	require.Equal(t, expect, detail.String(50))
}

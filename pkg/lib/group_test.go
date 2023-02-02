// Copyright 2022 The Moov Authors
// Use of this source code is governed by an Apache License
// license that can be found in the LICENSE file.

package lib

import (
	"bytes"
	"github.com/stretchr/testify/require"
	"testing"
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

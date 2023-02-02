// Copyright 2022 The Moov Authors
// Use of this source code is governed by an Apache License
// license that can be found in the LICENSE file.

package lib

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestFileWithSampleData(t *testing.T) {
	paths := []string{
		"sample1.txt",
		"sample2.txt",
		"sample3.txt",
	}

	for _, path := range paths {
		samplePath := filepath.Join("..", "..", "test", "testdata", path)
		fd, err := os.Open(samplePath)
		require.NoError(t, err)

		scan := NewBai2Scanner(fd)
		f := NewBai2()
		err = f.Read(&scan)

		require.NoError(t, err)
		require.NoError(t, f.Validate())
	}
}

func TestFileWithContinuationRecord(t *testing.T) {

	raw := `01,0004,12345,060321,0829,001,80,1,2/
02,12345,0004,1,060317,,CAD,/
03,10200123456,CAD,040,+000000000000,,,045,+000000000000,,/
88,046,+000000000000,,,047,+000000000000,,,048,+000000000000,,,049,+000000000000,,/
88,050,+000000000000,,,051,+000000000000,,,052,+000000000000,,,053,+000000000000,,/
16,409,000000000002500,V,060316,1300,,,RETURNED CHEQUE     /
16,409,000000000090000,V,060316,1300,,,RTN-UNKNOWN         /
49,+00000000000834000,14/
98,+00000000001280000,2,25/
99,+00000000001280000,1,27/`

	scan := NewBai2Scanner(strings.NewReader(raw))
	f := NewBai2()
	err := f.Read(&scan)
	require.NoError(t, err)
	require.NoError(t, f.Validate())

	expected := `01,0004,12345,060321,0829,001,80,1,2/
02,12345,0004,1,060317,,CAD,/
03,10200123456,CAD,040,+000000000000,,,045,+000000000000,,,046,+000000000000,,/
88,047,+000000000000,,,048,+000000000000,,,049,+000000000000,,,050/
88,+000000000000,,,051,+000000000000,,,052,+000000000000,,,053,+000000000000,,/
16,409,000000000002500,V,060316,1300,,,RETURNED CHEQUE     /
49,+00000000000834000,14/
98,+00000000001280000,2,25/
99,+00000000001280000,1,27/`
	require.Equal(t, expected, f.String())
}

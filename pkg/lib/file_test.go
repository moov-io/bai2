// Copyright 2022 The Moov Authors
// Use of this source code is governed by an Apache License
// license that can be found in the LICENSE file.

package lib

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestWithSampleData(t *testing.T) {
	paths := []string{
		"sample1.txt",
		"sample2.txt",
		"sample3.txt",
	}

	for _, path := range paths {
		samplePath := filepath.Join("..", "..", "test", "testdata", path)
		fd, err := os.Open(samplePath)
		require.NoError(t, err)

		f := NewBai2()
		err = f.Read(NewBai2Scanner(fd))
		require.NoError(t, err)
		require.NoError(t, f.Validate())
	}
}

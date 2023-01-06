// Copyright 2022 The Moov Authors
// Use of this source code is governed by an Apache License
// license that can be found in the LICENSE file.

package file

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestWithSampleData(t *testing.T) {
	type testData struct {
		FileName string
		EqualRaw bool
	}
	samples := []testData{
		{"sample1.txt", true}, {"sample2.txt", true}, {"sample3.txt", false}}

	for _, sample := range samples {
		samplePath := filepath.Join("..", "..", "data", sample.FileName)
		fd, err := os.Open(samplePath)
		require.NoError(t, err)

		bai, err := Parse(fd)
		require.NoError(t, err)
		require.NoError(t, bai.Validate())

		if sample.EqualRaw {
			raw, err := os.ReadFile(samplePath)
			require.NoError(t, err)

			rawStr := strings.ReplaceAll(string(raw), "\r\n", "\n")
			require.Equal(t, rawStr, bai.String())
		}
	}
}

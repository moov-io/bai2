// Copyright 2020 The Moov Authors
// Use of this source code is governed by an Apache License
// license that can be found in the LICENSE file.

package fuzz

import (
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/moov-io/bai2/pkg/lib"

	"github.com/stretchr/testify/require"
)

func FuzzReaderWriter(f *testing.F) {
	populateCorpus(f)

	f.Fuzz(func(t *testing.T, contents string) {
		scan := lib.NewBai2Scanner(strings.NewReader(contents))
		file := lib.NewBai2()
		file.Read(&scan)
		file.Validate()

		out := file.String()
		require.Greater(t, len(out), 0)
	})
}

func populateCorpus(f *testing.F) {
	f.Helper()

	err := filepath.Walk(filepath.Join("..", "testdata"), func(path string, info fs.FileInfo, _ error) error {
		if info.IsDir() {
			return nil
		}

		bs, err := os.ReadFile(path)
		if err != nil {
			f.Fatal(err)
		}
		f.Add(string(bs))
		return nil
	})
	if err != nil {
		f.Fatal(err)
	}
}

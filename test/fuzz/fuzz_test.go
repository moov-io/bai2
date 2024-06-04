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

func FuzzReaderWriter_ValidFiles(f *testing.F) {
	populateCorpus(f, false)

	f.Fuzz(func(t *testing.T, contents string) {
		scan := lib.NewBai2Scanner(strings.NewReader(contents))
		file := lib.NewBai2()
		require.NotPanics(t, func() { file.Read(&scan, false) })
		require.NotPanics(t, func() { file.Validate(false) })

		out := file.String()
		require.Greater(t, len(out), 0)
	})
}

func FuzzReaderWriter_ErrorFiles(f *testing.F) {
	populateCorpus(f, true)

	f.Fuzz(func(t *testing.T, contents string) {
		scan := lib.NewBai2Scanner(strings.NewReader(contents))
		file := lib.NewBai2()
		require.NotPanics(t, func() { file.Read(&scan, false) })
		require.NotPanics(t, func() { file.Validate(false) })

		out := file.String()
		require.Greater(t, len(out), 0)
	})
}

func populateCorpus(f *testing.F, errorFiles bool) {
	f.Helper()

	err := filepath.Walk(filepath.Join("..", "testdata"), func(path string, info fs.FileInfo, _ error) error {
		path = filepath.ToSlash(path)

		// Skip directories and some files
		if info.IsDir() {
			return nil
		}
		if strings.HasSuffix(path, ".output") {
			return nil // skip
		}
		if errorFiles && !strings.Contains(path, "errors/") {
			f.Logf("skipping %s", path)
			return nil
		}
		if !errorFiles && strings.Contains(path, "errors/") {
			f.Logf("skipping %s", path)
			return nil
		}

		f.Logf("adding %s", path)

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

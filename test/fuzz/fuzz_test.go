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
)

func FuzzReaderWriter(f *testing.F) {
	populateCorpus(f)

	f.Fuzz(func(t *testing.T, contents string) {
		// Bound input size so pathological cases don't OOM CI.
		if len(contents) > 1<<20 {
			t.Skip()
		}

		scan := lib.NewBai2Scanner(strings.NewReader(contents))
		file := lib.NewBai2()

		// Read/Validate/String must never panic on arbitrary input.
		_ = file.Read(&scan)
		_ = file.Validate()
		_ = file.String()
	})
}

func populateCorpus(f *testing.F) {
	f.Helper()

	// Always seed empty / tiny inputs.
	f.Add("")
	f.Add("01,")
	f.Add("01,SENDER,RECEIVER,250101,1200,1,,,2/\n99,0,0,0/")

	err := filepath.Walk(filepath.Join("..", "testdata"), func(path string, info fs.FileInfo, _ error) error {
		path = filepath.ToSlash(path)

		if info.IsDir() {
			return nil
		}
		if strings.HasSuffix(path, ".output") {
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

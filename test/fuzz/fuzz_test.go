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

	moovbai "github.com/moov-io/bai2"
	"github.com/moov-io/bai2/pkg/bai2"
)

func FuzzReaderWriter(f *testing.F) {
	populateCorpus(f)

	f.Fuzz(func(t *testing.T, contents string) {
		// Bound input size so pathological cases don't OOM CI.
		if len(contents) > 1<<20 {
			t.Skip()
		}

		scan := bai2.NewBai2Scanner(strings.NewReader(contents))
		file := bai2.NewBai2()

		// Read/Validate/String must never panic on arbitrary input.
		_ = file.Read(&scan)
		_ = file.Validate()
		_ = file.String()

		got, _ := moovbai.Read(strings.NewReader(contents))
		if got != nil {
			_ = got.Validate()
			_ = got.String()
		}
	})
}

func populateCorpus(f *testing.F) {
	f.Helper()

	// Always seed empty / tiny inputs.
	f.Add("")
	f.Add("01,")
	f.Add("01,SENDER,RECEIVER,250101,1200,1,,,2/\n99,0,0,0/")
	f.Add("01,122099999,123456789,040621,0200,1,,,2/\n02,031001234,122099999,1,040620,2359,,2/\n03,5765432,,,,,/\n49,0,2/\n98,0,1,4/\n99,0,1,6/")
	f.Add("16,890,,,,,detail reports will be delayed until 11:00 AM.")

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

// Copyright 2022 The Moov Authors
// Use of this source code is governed by an Apache License
// license that can be found in the LICENSE file.

// Package bai2 reads Cash Management Balance Reporting files.
// Read inspects the 01 Version Number and returns a BAI2 (pkg/bai2) or
// BAI3 / BTR3 (pkg/bai3) document. Import pkg/bai2 or pkg/bai3 directly
// when the version is already known.
package bai2

import (
	"bytes"
	"fmt"
	"io"
	"strconv"
	"strings"

	bai2pkg "github.com/moov-io/bai2/pkg/bai2"
	"github.com/moov-io/bai2/pkg/bai3"
	"github.com/moov-io/bai2/pkg/util"
)

// File is implemented by BAI2 and BAI3 parsed documents.
type File interface {
	Version() int64
	Validate() error
	String() string
}

// ReadOptions control version detection and parser strictness.
type ReadOptions struct {
	IgnoreVersion       bool
	StrictControlTotals bool
	// ForceVersion skips detection and uses 2 or 3.
	ForceVersion int64
}

// Read detects the file header version and returns a BAI2 or BAI3 document.
func Read(r io.Reader) (File, error) {
	return ReadWithOptions(r, ReadOptions{})
}

// ReadWithOptions is Read with explicit parser options.
func ReadWithOptions(r io.Reader, opt ReadOptions) (File, error) {
	buf, err := io.ReadAll(r)
	if err != nil {
		return nil, err
	}

	version := opt.ForceVersion
	if version == 0 {
		version, err = detectVersion(buf)
		if err != nil {
			return nil, err
		}
	}

	scan := util.NewScanner(bytes.NewReader(buf))
	useBai3 := opt.ForceVersion == bai3.VersionBTR3 || (opt.ForceVersion == 0 && !opt.IgnoreVersion && version == bai3.VersionBTR3)
	if useBai3 {
		f := bai3.NewFileWith(bai3.Options{
			IgnoreVersion:       opt.IgnoreVersion,
			StrictControlTotals: opt.StrictControlTotals,
		})
		if err := f.Read(&scan); err != nil {
			return nil, err
		}
		return f, nil
	}

	if version != bai2pkg.VersionBAI2 && !opt.IgnoreVersion && opt.ForceVersion != bai2pkg.VersionBAI2 {
		return nil, fmt.Errorf("unsupported BAI version %d (want 2 or 3)", version)
	}

	f := bai2pkg.NewBai2With(bai2pkg.Options{
		IgnoreVersion:       opt.IgnoreVersion || opt.ForceVersion == bai2pkg.VersionBAI2,
		StrictControlTotals: opt.StrictControlTotals,
	})
	if err := f.Read(&scan); err != nil {
		return nil, err
	}
	return f, nil
}

func detectVersion(buf []byte) (int64, error) {
	scan := util.NewScanner(bytes.NewReader(buf))
	line := scan.ScanLine()
	if err := scan.Err(); err != nil {
		return 0, err
	}
	if line == "" {
		return 0, fmt.Errorf("missing file header (01)")
	}
	if len(line) < 3 || line[:2] != util.FileHeaderCode {
		return 0, fmt.Errorf("first record is %q, want 01 file header", line)
	}
	fields := strings.Split(strings.TrimSuffix(line, "/"), ",")
	if len(fields) < 9 {
		return 0, fmt.Errorf("file header too short to read version")
	}
	v, err := strconv.ParseInt(strings.TrimSpace(fields[8]), 10, 64)
	if err != nil {
		return 0, fmt.Errorf("invalid version number %q", fields[8])
	}
	return v, nil
}

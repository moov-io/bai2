// Copyright 2022 The Moov Authors
// Use of this source code is governed by an Apache License
// license that can be found in the LICENSE file.

package bai2

import (
	"os"
	"strings"
	"testing"

	bai2pkg "github.com/moov-io/bai2/pkg/bai2"
	"github.com/moov-io/bai2/pkg/bai3"
	"github.com/stretchr/testify/require"
)

func TestReadDetectsBAI2(t *testing.T) {
	raw, err := os.ReadFile("test/testdata/spec-section3.txt")
	require.NoError(t, err)
	f, err := Read(strings.NewReader(string(raw)))
	require.NoError(t, err)
	require.Equal(t, bai2pkg.VersionBAI2, f.Version())
	_, ok := f.(*bai2pkg.Bai2)
	require.True(t, ok)
	require.NoError(t, f.Validate())
}

func TestReadDetectsBAI3(t *testing.T) {
	raw, err := os.ReadFile("test/testdata/bai3/x9-mandatory.txt")
	require.NoError(t, err)
	f, err := Read(strings.NewReader(string(raw)))
	require.NoError(t, err)
	require.Equal(t, bai3.VersionBTR3, f.Version())
	got, ok := f.(*bai3.File)
	require.True(t, ok)
	require.Equal(t, "GBP", got.Banks[0].Accounts[0].CurrencyCode)
	require.NoError(t, f.Validate())
}

func TestReadForceVersion(t *testing.T) {
	raw, err := os.ReadFile("test/testdata/bai3/x9-mandatory.txt")
	require.NoError(t, err)
	_, err = ReadWithOptions(strings.NewReader(string(raw)), ReadOptions{ForceVersion: 2, IgnoreVersion: true})
	require.NoError(t, err)
}

func TestReadIgnoreVersionParsesAsBAI2(t *testing.T) {
	// BAI2 layout stamped as version 3, with omitted account currency (invalid BAI3).
	raw := "01,122099999,123456789,040621,0200,1,,,3/\n" +
		"02,031001234,122099999,1,040620,2359,,2/\n" +
		"03,5765432,,,,,/\n" +
		"49,0,2/\n98,0,1,4/\n99,0,1,6/\n"

	_, err := Read(strings.NewReader(raw))
	require.Error(t, err)

	f, err := ReadWithOptions(strings.NewReader(raw), ReadOptions{IgnoreVersion: true})
	require.NoError(t, err)
	got, ok := f.(*bai2pkg.Bai2)
	require.True(t, ok)
	require.Equal(t, int64(3), got.VersionNumber)
}

func TestReadForceVersion3(t *testing.T) {
	raw, err := os.ReadFile("test/testdata/bai3/x9-mandatory.txt")
	require.NoError(t, err)
	f, err := ReadWithOptions(strings.NewReader(string(raw)), ReadOptions{ForceVersion: 3})
	require.NoError(t, err)
	_, ok := f.(*bai3.File)
	require.True(t, ok)
}

func TestReadUnsupportedVersion(t *testing.T) {
	raw := "01,122099999,123456789,150623,0200,1,,,4/\n" +
		"02,,122099999,1,150622,,,2/\n" +
		"03,0987654321,GBP,,,,/\n" +
		"49,0,2/\n98,0,1,4/\n99,0,1,6/\n"
	_, err := Read(strings.NewReader(raw))
	require.Error(t, err)
	require.Contains(t, err.Error(), "unsupported BAI version 4")
}

// Copyright 2022 The Moov Authors
// Use of this source code is governed by an Apache License
// license that can be found in the LICENSE file.

package bai2

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestStrictControlTotals_Sample2PhysicalRecords(t *testing.T) {
	path := filepath.Join("..", "..", "test", "testdata", "sample2.txt")
	raw, err := os.ReadFile(path)
	require.NoError(t, err)

	scan := NewBai2Scanner(strings.NewReader(string(raw)))
	file := NewBai2With(Options{StrictControlTotals: true})
	require.NoError(t, file.Read(&scan))

	var found bool
	for _, group := range file.Groups {
		for _, account := range group.Accounts {
			if account.AccountNumber == "7890654321" {
				found = true
				require.Equal(t, int64(3), account.NumberRecords)
				require.Equal(t, int64(3), account.physicalRecordCount)
				require.Equal(t, int64(3), account.recordCount(file.PhysicalRecordLength))
				// String() wrapping does not preserve the original 88 record.
				require.Equal(t, int64(2), account.SumRecords(file.PhysicalRecordLength))
			}
		}
	}
	require.True(t, found)
}

func TestStrictControlTotals_WrongRecordCount(t *testing.T) {
	path := filepath.Join("..", "..", "test", "testdata", "sample2.txt")
	raw, err := os.ReadFile(path)
	require.NoError(t, err)

	tampered := strings.Replace(string(raw), "49,11800000,3/", "49,11800000,2/", 1)
	scan := NewBai2Scanner(strings.NewReader(tampered))
	file := NewBai2With(Options{StrictControlTotals: true})
	err = file.Read(&scan)
	require.Error(t, err)
	require.Contains(t, err.Error(), "account 7890654321 record count 2 does not match computed 3")
}

func TestStrictControlTotals_CreateRoundTrip(t *testing.T) {
	file := NewBai2()
	file.Sender = "122099999"
	file.Receiver = "123456789"
	file.FileCreatedDate = "040621"
	file.FileCreatedTime = "0200"
	file.FileIdNumber = "1"
	file.VersionNumber = VersionBAI2

	group := NewGroup()
	group.Receiver = "031001234"
	group.Originator = "122099999"
	group.GroupStatus = 1
	group.AsOfDate = "040620"
	group.AsOfTime = "2359"
	group.AsOfDateModifier = 2

	account := NewAccount()
	account.AccountNumber = "0123456789"
	account.Summaries = []AccountSummary{{
		TypeCode: "010",
		Amount:   "100",
	}}
	group.Accounts = []Account{*account}
	file.Groups = []Group{*group}

	require.NoError(t, file.Create())

	scan := NewBai2Scanner(strings.NewReader(file.String()))
	parsed := NewBai2With(Options{StrictControlTotals: true})
	require.NoError(t, parsed.Read(&scan))
}

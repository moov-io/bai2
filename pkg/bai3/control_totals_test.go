// Copyright 2022 The Moov Authors
// Use of this source code is governed by an Apache License
// license that can be found in the LICENSE file.

package bai3

import (
	"os"
	"strings"
	"testing"

	"github.com/moov-io/bai2/pkg/util"
	"github.com/stretchr/testify/require"
)

func TestStrictControlTotals_ContinuationPhysicalRecords(t *testing.T) {
	raw, err := os.ReadFile(testdata("x9-continuation-16.txt"))
	require.NoError(t, err)

	// The account has five physical records (03, 16, 88, 88, 49), but the 49
	// trailer says 4. Strict validation reports that mismatch. Without a
	// parse-time count, SumRecords would re-serialize the account and report 3.
	scan := util.NewScanner(strings.NewReader(string(raw)))
	file := NewFile()
	require.NoError(t, file.Read(&scan))

	require.Len(t, file.Banks, 1)
	require.Len(t, file.Banks[0].Accounts, 1)
	account := file.Banks[0].Accounts[0]
	require.Equal(t, int64(4), account.NumberRecords)
	require.Equal(t, int64(5), account.physicalRecordCount)
	require.Equal(t, int64(5), account.recordCount(file.PhysicalRecordLength))
	require.Equal(t, int64(3), account.SumRecords(file.PhysicalRecordLength))

	strictScan := util.NewScanner(strings.NewReader(string(raw)))
	strict := NewFileWith(Options{StrictControlTotals: true})
	err = strict.Read(&strictScan)
	require.Error(t, err)
	require.Contains(t, err.Error(), "account 1234567890 record count 4 does not match computed 5")
}

func TestStrictControlTotals_StatusSummaryDetail(t *testing.T) {
	raw, err := os.ReadFile(testdata("x9-status-summary-detail.txt"))
	require.NoError(t, err)

	scan := util.NewScanner(strings.NewReader(string(raw)))
	file := NewFileWith(Options{StrictControlTotals: true})
	require.NoError(t, file.Read(&scan))

	tampered := strings.Replace(string(raw), "49,70525000,3/", "49,70525000,2/", 1)
	if tampered == string(raw) {
		t.Fatal("expected to find account trailer 49,70525000,3/")
	}
	tamperScan := util.NewScanner(strings.NewReader(tampered))
	err = NewFileWith(Options{StrictControlTotals: true}).Read(&tamperScan)
	require.Error(t, err)
	require.Contains(t, err.Error(), "record count 2 does not match computed 3")
}

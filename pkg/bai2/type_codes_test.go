// Copyright 2022 The Moov Authors
// Use of this source code is governed by an Apache License
// license that can be found in the LICENSE file.

package bai2

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestLookupTypeCodeAppendixA(t *testing.T) {
	cases := []struct {
		code        string
		level       TypeLevel
		kind        TransactionKind
		description string
	}{
		{"010", TypeLevelStatus, TransactionNA, "Opening Ledger"},
		{"040", TypeLevelStatus, TransactionNA, "Opening Available"},
		{"045", TypeLevelStatus, TransactionNA, "Closing Available"},
		{"100", TypeLevelSummary, TransactionCR, "Total Credits"},
		{"108", TypeLevelDetail, TransactionCR, "Credit (Any Type)"},
		{"165", TypeLevelDetail, TransactionCR, "Preauthorized ACH Credit"},
		{"190", TypeLevelSummary, TransactionCR, "Total Incoming Money Transfers"},
		{"400", TypeLevelSummary, TransactionDB, "Total Debits"},
		{"409", TypeLevelDetail, TransactionDB, "Debit (Any Type)"},
		{"890", TypeLevelDetail, TransactionNA, "Contains Non-monetary Information"},
	}
	for _, tc := range cases {
		got, ok := LookupTypeCode(tc.code)
		require.True(t, ok, tc.code)
		require.Equal(t, tc.level, got.Level, tc.code)
		require.Equal(t, tc.kind, got.Transaction, tc.code)
		require.Equal(t, tc.description, got.Description, tc.code)
	}
}

func TestClassifyTypeCodeRanges(t *testing.T) {
	level, kind := ClassifyTypeCode("010")
	require.Equal(t, TypeLevelStatus, level)
	require.Equal(t, TransactionNA, kind)

	level, kind = ClassifyTypeCode("901") // custom status 900–919
	require.Equal(t, TypeLevelStatus, level)
	require.Equal(t, TransactionNA, kind)

	level, kind = ClassifyTypeCode("930") // custom credit 920–959
	require.Equal(t, TypeLevelSummary, level)
	require.Equal(t, TransactionCR, kind)

	level, kind = ClassifyTypeCode("970") // custom debit 960–999
	require.Equal(t, TypeLevelSummary, level)
	require.Equal(t, TransactionDB, kind)

	require.True(t, IsStatusTypeCode("072"))
	require.False(t, IsStatusTypeCode("190"))
}

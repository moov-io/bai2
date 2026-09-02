// Copyright 2022 The Moov Authors
// Use of this source code is governed by an Apache License
// license that can be found in the LICENSE file.

package bai2

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestProductionUTF8BOM(t *testing.T) {
	body := "01,122099999,123456789,040621,0200,1,,,2/\n" +
		"02,031001234,122099999,1,040620,2359,,2/\n" +
		"03,0975312468,,010,500000,,,/\n" +
		"49,500000,2/\n" +
		"98,500000,1,4/\n" +
		"99,500000,1,6/\n"
	raw := "\xef\xbb\xbf" + body

	scan := NewBai2Scanner(strings.NewReader(raw))
	file := NewBai2()
	require.NoError(t, file.Read(&scan))
	require.NoError(t, file.Validate())
	require.Equal(t, "122099999", file.Sender)
}

func TestProductionSpacesAroundFields(t *testing.T) {
	raw := "01, 122099999 ,123456789, 040621,0200,1,,,2/\n" +
		"02, 031001234, 122099999, 1, 040620, 2359, , 2/\n" +
		"03, 0975312468 , , 010, 500000, , , /\n" +
		"49,500000,2/\n" +
		"98,500000,1,4/\n" +
		"99,500000,1,6/\n"

	scan := NewBai2Scanner(strings.NewReader(raw))
	file := NewBai2()
	require.NoError(t, file.Read(&scan))
	require.NoError(t, file.Validate())
	require.Equal(t, "122099999", file.Sender)
	require.Equal(t, "031001234", file.Groups[0].Receiver)
	require.Equal(t, "0975312468", file.Groups[0].Accounts[0].AccountNumber)
	require.Equal(t, "010", file.Groups[0].Accounts[0].Summaries[0].TypeCode)
}

func TestProductionOmittedCurrencyAndGroupStatus(t *testing.T) {
	// Spec: omitted group currency is USD; omitted status is treated as Update (1).
	raw := "01,122099999,123456789,040621,0200,1,,,2/\n" +
		"02,,122099999,,040620,,,/\n" +
		"03,0975312468,,,,,/\n" +
		"49,0,2/\n" +
		"98,0,1,4/\n" +
		"99,0,1,6/\n"

	scan := NewBai2Scanner(strings.NewReader(raw))
	file := NewBai2()
	require.NoError(t, file.Read(&scan))
	require.NoError(t, file.Validate())
	require.Equal(t, DefaultGroupStatus, file.Groups[0].GroupStatus)
	require.Equal(t, DefaultCurrency, file.Groups[0].CurrencyCode)
	require.Equal(t, DefaultCurrency, file.Groups[0].Accounts[0].CurrencyCode)
}

func TestProductionTypeCodeMetadata(t *testing.T) {
	raw := "01,122099999,123456789,040621,0200,1,,,2/\n" +
		"02,031001234,122099999,1,040620,2359,,2/\n" +
		"03,0975312468,,010,500000,,,190,70000000,4,0/\n" +
		"16,165,1500000,1,DD1620,, DEALER PAYMENTS\n" +
		"49,72000000,3/\n" +
		"98,72000000,1,5/\n" +
		"99,72000000,1,7/\n"

	scan := NewBai2Scanner(strings.NewReader(raw))
	file := NewBai2()
	require.NoError(t, file.Read(&scan))
	require.NoError(t, file.Validate())

	opening := file.Groups[0].Accounts[0].Summaries[0]
	require.Equal(t, TypeLevelStatus, opening.Level)
	require.Equal(t, TransactionNA, opening.Transaction)
	require.Equal(t, "Opening Ledger", opening.Description)
	info, ok := opening.TypeInfo()
	require.True(t, ok)
	require.Equal(t, "Opening Ledger", info.Description)

	detail := file.Groups[0].Accounts[0].Details[0]
	require.Equal(t, TypeLevelDetail, detail.Level)
	require.Equal(t, TransactionCR, detail.Transaction)
	require.Equal(t, "Preauthorized ACH Credit", detail.Description)
}

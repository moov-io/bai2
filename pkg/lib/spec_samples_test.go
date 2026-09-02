// Copyright 2022 The Moov Authors
// Use of this source code is governed by an Apache License
// license that can be found in the LICENSE file.

package lib

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

// Samples are taken from Cash Management Balance Reporting Specifications
// Version 2, section 3 Record Formats.

func TestSpecFileHeaderSample(t *testing.T) {
	// Section 3, 01 – FILE HEADER
	sample := "01,122099999,123456789,040621,0200,1,55,,2/"
	var header fileHeader
	_, err := header.parse(sample, Options{})
	require.NoError(t, err)
	require.Equal(t, "122099999", header.Sender)
	require.Equal(t, "123456789", header.Receiver)
	require.Equal(t, "040621", header.FileCreatedDate)
	require.Equal(t, "0200", header.FileCreatedTime)
	require.Equal(t, "1", header.FileIdNumber)
	require.Equal(t, int64(55), header.PhysicalRecordLength)
	require.Equal(t, int64(0), header.BlockSize)
	require.Equal(t, int64(2), header.VersionNumber)
	require.Equal(t, sample, header.string())
}

func TestSpecGroupHeaderSample(t *testing.T) {
	// Section 3, 02 – GROUP HEADER
	sample := "02,031001234,122099999,1,040620,2359,,2/"
	var header groupHeader
	_, err := header.parse(sample)
	require.NoError(t, err)
	require.Equal(t, "031001234", header.Receiver)
	require.Equal(t, "122099999", header.Originator)
	require.Equal(t, int64(1), header.GroupStatus)
	require.Equal(t, "040620", header.AsOfDate)
	require.Equal(t, "2359", header.AsOfTime)
	require.Equal(t, "", header.CurrencyCode)
	require.Equal(t, int64(2), header.AsOfDateModifier)
	require.Equal(t, sample, header.string())
}

func TestSpecAccountIdentifierSample(t *testing.T) {
	// Section 3, 03 – ACCOUNT IDENTIFIER AND SUMMARY STATUS
	sample := "03,0975312468,,010,500000,,,190,70000000,4,0/"
	var rec accountIdentifier
	_, err := rec.parse(sample)
	require.NoError(t, err)
	require.Equal(t, "0975312468", rec.AccountNumber)
	require.Equal(t, "", rec.CurrencyCode)
	require.Len(t, rec.Summaries, 2)

	require.Equal(t, "010", rec.Summaries[0].TypeCode)
	require.Equal(t, "500000", rec.Summaries[0].Amount)
	require.Equal(t, int64(0), rec.Summaries[0].ItemCount)
	require.True(t, IsStatusTypeCode("010"))

	require.Equal(t, "190", rec.Summaries[1].TypeCode)
	require.Equal(t, "70000000", rec.Summaries[1].Amount)
	require.Equal(t, int64(4), rec.Summaries[1].ItemCount)
	require.Equal(t, FundsType0, string(rec.Summaries[1].FundsType.TypeCode))
}

func TestSpecAccountIdentifierNoStatusOrSummary(t *testing.T) {
	// Section 3, 03 note: account number only, five defaulted fields.
	sample := "03,5765432,,,,,/"
	var rec accountIdentifier
	_, err := rec.parse(sample)
	require.NoError(t, err)
	require.Equal(t, "5765432", rec.AccountNumber)
	require.Empty(t, rec.Summaries)
	require.Equal(t, sample, rec.string())
}

func TestSpecTransactionDetailSample(t *testing.T) {
	// Section 3, 16 – TRANSACTION DETAIL
	// Text is not slash-delimited.
	sample := "16,165,1500000,1,DD1620,, DEALER PAYMENTS"
	var rec transactionDetail
	_, err := rec.parse(sample)
	require.NoError(t, err)
	require.Equal(t, "165", rec.TypeCode)
	require.Equal(t, "1500000", rec.Amount)
	require.Equal(t, FundsType1, string(rec.FundsType.TypeCode))
	require.Equal(t, "DD1620", rec.BankReferenceNumber)
	require.Equal(t, "", rec.CustomerReferenceNumber)
	require.Equal(t, " DEALER PAYMENTS", rec.Text)
	require.Equal(t, sample, rec.string())
}

func TestSpecNonMonetaryType890(t *testing.T) {
	// Section 2 / Appendix A: type code 890
	sample := "16,890,,,,,detail reports will be delayed until 11:00 AM."
	var rec transactionDetail
	_, err := rec.parse(sample)
	require.NoError(t, err)
	require.Equal(t, "890", rec.TypeCode)
	require.Equal(t, "", rec.Amount)
	require.Equal(t, "detail reports will be delayed until 11:00 AM.", rec.Text)

	tc, ok := LookupTypeCode("890")
	require.True(t, ok)
	require.Equal(t, TypeLevelDetail, tc.Level)
	require.Equal(t, TransactionNA, tc.Transaction)
}

func TestSpecContinuationSample(t *testing.T) {
	// Section 3, 88 – CONTINUATION RECORD
	raw := "16,115,10000000,S,5000000,4000000,1000000/\n88,AX13612,B096132,AMALGAMATED CORP. LOCKBOX\n88,DEPOSIT-MISC. RECEIVABLES\n"
	scan := NewBai2Scanner(strings.NewReader(raw))
	detail := NewDetail()
	require.NoError(t, detail.Read(&scan, false))
	require.Equal(t, "115", detail.TypeCode)
	require.Equal(t, "10000000", detail.Amount)
	require.Equal(t, FundsTypeS, string(detail.FundsType.TypeCode))
	require.Equal(t, int64(5000000), detail.FundsType.ImmediateAmount)
	require.Equal(t, int64(4000000), detail.FundsType.OneDayAmount)
	require.Equal(t, int64(1000000), detail.FundsType.TwoDayAmount)
	require.Equal(t, "AX13612", detail.BankReferenceNumber)
	require.Equal(t, "B096132", detail.CustomerReferenceNumber)
	require.Equal(t, "AMALGAMATED CORP. LOCKBOX,DEPOSIT-MISC. RECEIVABLES", detail.Text)
}

func TestSpecAccountTrailerSample(t *testing.T) {
	// Section 3, 49 – ACCOUNT TRAILER
	sample := "49,18650000,3/"
	var rec accountTrailer
	_, err := rec.parse(sample)
	require.NoError(t, err)
	require.Equal(t, "18650000", rec.AccountControlTotal)
	require.Equal(t, int64(3), rec.NumberRecords)
	require.Equal(t, sample, rec.string())
}

func TestSpecGroupTrailerSample(t *testing.T) {
	// Section 3, 98 – GROUP TRAILER
	sample := "98,11800000,2,6/"
	var rec groupTrailer
	_, err := rec.parse(sample)
	require.NoError(t, err)
	require.Equal(t, "11800000", rec.GroupControlTotal)
	require.Equal(t, int64(2), rec.NumberOfAccounts)
	require.Equal(t, int64(6), rec.NumberOfRecords)
	require.Equal(t, sample, rec.string())
}

func TestSpecTimes0000And2400And9999(t *testing.T) {
	// Section 3: 0000 beginning of day, 2400 end of day, 9999 also end of day.
	for _, tm := range []string{"0000", "2400", "9999"} {
		sample := "02,031001234,122099999,1,040620," + tm + ",,2/"
		var header groupHeader
		_, err := header.parse(sample)
		require.NoError(t, err, tm)
		require.Equal(t, tm, header.AsOfTime)
	}
}

func TestSpecCreateControlTotals(t *testing.T) {
	// Algebraic sum of Amount fields on 03 + 16, per 49 definition.
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

	acct := NewAccount()
	acct.AccountNumber = "0975312468"
	acct.Summaries = []AccountSummary{
		{TypeCode: "010", Amount: "500000"},
		{TypeCode: "190", Amount: "70000000", ItemCount: 4, FundsType: FundsType{TypeCode: FundsType0}},
	}
	detail := NewDetail()
	detail.TypeCode = "165"
	detail.Amount = "1500000"
	detail.FundsType = FundsType{TypeCode: FundsType1}
	detail.BankReferenceNumber = "DD1620"
	detail.Text = " DEALER PAYMENTS"
	acct.Details = []Detail{*detail}
	group.Accounts = []Account{*acct}
	file.Groups = []Group{*group}

	require.NoError(t, file.Create())
	require.NoError(t, file.Validate())

	require.Equal(t, "72000000", file.Groups[0].Accounts[0].AccountControlTotal)
	require.Equal(t, int64(3), file.Groups[0].Accounts[0].NumberRecords) // 03 + 16 + 49
	require.Equal(t, "72000000", file.Groups[0].GroupControlTotal)
	require.Equal(t, int64(1), file.Groups[0].NumberOfAccounts)
	require.Equal(t, int64(5), file.Groups[0].NumberOfRecords) // 02+03+16+49+98
	require.Equal(t, "72000000", file.FileControlTotal)
	require.Equal(t, int64(1), file.NumberOfGroups)
	require.Equal(t, int64(7), file.NumberOfRecords) // +01+99

	file.options.StrictControlTotals = true
	require.NoError(t, file.Validate())
}

func TestSpecPhysicalRecordLengthBlockedFile(t *testing.T) {
	// Section 2: physical record occupies N characters; trailing positions are padding.
	pad := func(s string, n int) string {
		if len(s) > n {
			return s[:n]
		}
		return s + strings.Repeat(" ", n-len(s))
	}
	n := 80
	raw := pad("01,122099999,123456789,040621,0200,1,80,,2/", n) +
		pad("02,031001234,122099999,1,040620,2359,,2/", n) +
		pad("03,0975312468,,010,500000,,,/", n) +
		pad("16,165,1500000,1,DD1620,, DEALER PAYMENTS", n) +
		pad("49,2000000,3/", n) +
		pad("98,2000000,1,5/", n) +
		pad("99,2000000,1,7/", n)

	scan := NewBai2Scanner(strings.NewReader(raw))
	file := NewBai2()
	require.NoError(t, file.Read(&scan))
	require.NoError(t, file.Validate())
	require.Equal(t, int64(80), file.PhysicalRecordLength)
	require.Equal(t, "122099999", file.Sender)
	require.Len(t, file.Groups, 1)
	require.Equal(t, "0975312468", file.Groups[0].Accounts[0].AccountNumber)
	require.Equal(t, " DEALER PAYMENTS", file.Groups[0].Accounts[0].Details[0].Text)
}

func TestMissingFileTrailer(t *testing.T) {
	raw := "01,122099999,123456789,040621,0200,1,,,2/\n02,031001234,122099999,1,040620,2359,,2/\n"
	scan := NewBai2Scanner(strings.NewReader(raw))
	file := NewBai2()
	err := file.Read(&scan)
	require.Error(t, err)
	require.Contains(t, err.Error(), "missing")
}

func TestVersion3RejectedWithoutIgnoreVersion(t *testing.T) {
	sample := "01,122099999,123456789,040621,0200,1,,,3/"
	var header fileHeader
	_, err := header.parse(sample, Options{})
	require.Error(t, err)
	require.Contains(t, err.Error(), "version 3")

	_, err = header.parse(sample, Options{IgnoreVersion: true})
	require.NoError(t, err)
	require.Equal(t, int64(3), header.VersionNumber)
}

func TestGroupStatusZeroInvalid(t *testing.T) {
	var header groupHeader
	_, err := header.parse("02,031001234,122099999,0,040620,2359,,2/")
	require.Error(t, err)
	require.Contains(t, err.Error(), "GroupStatus")
}

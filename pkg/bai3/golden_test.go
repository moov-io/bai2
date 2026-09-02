// Copyright 2022 The Moov Authors
// Use of this source code is governed by an Apache License
// license that can be found in the LICENSE file.

package bai3

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/moov-io/bai2/pkg/util"
	"github.com/stretchr/testify/require"
)

func testdata(name string) string {
	return filepath.Join("..", "..", "test", "testdata", "bai3", name)
}

func TestPublicGoldenFiles(t *testing.T) {
	files := []string{
		"x9-mandatory.txt",
		"x9-status-summary-detail.txt",
		"x9-bank-header-all-fields.txt",
		"x9-continuation-16.txt",
		"x9-03-split-summary.txt",
	}
	for _, name := range files {
		t.Run(name, func(t *testing.T) {
			raw, err := os.ReadFile(testdata(name))
			require.NoError(t, err)
			scan := util.NewScanner(strings.NewReader(string(raw)))
			f := NewFile()
			require.NoError(t, f.Read(&scan), name)
			require.NoError(t, f.Validate(), name)
			require.Equal(t, VersionBTR3, f.Version())
			require.NotEmpty(t, f.Banks)
			require.NotEmpty(t, f.Banks[0].Accounts)
			require.NotEmpty(t, f.Banks[0].Accounts[0].CurrencyCode)
			require.Empty(t, f.Banks[0].CurrencyCode)
		})
	}
}

func TestX9MandatoryFields(t *testing.T) {
	raw, err := os.ReadFile(testdata("x9-mandatory.txt"))
	require.NoError(t, err)
	scan := util.NewScanner(strings.NewReader(string(raw)))
	f := NewFile()
	require.NoError(t, f.Read(&scan))

	require.Equal(t, "122099999", f.Sender)
	require.Equal(t, "123456789", f.Receiver)
	require.Equal(t, "150623", f.FileCreatedDate)
	require.Equal(t, "0200", f.FileCreatedTime)
	require.Equal(t, "1", f.FileIdNumber)
	require.Equal(t, int64(3), f.VersionNumber)

	bank := f.Banks[0]
	require.Equal(t, "122099999", bank.BankIdentification)
	require.Equal(t, int64(1), bank.GroupStatus)
	require.Equal(t, "150622", bank.AsOfDate)
	require.Equal(t, int64(2), bank.AsOfDateModifier)
	require.Empty(t, bank.CurrencyCode)

	acct := bank.Accounts[0]
	require.Equal(t, "0987654321", acct.AccountNumber)
	require.Equal(t, "GBP", acct.CurrencyCode)
	require.Empty(t, acct.Summaries)
}

func TestX9StatusSummaryAndDetail(t *testing.T) {
	raw, err := os.ReadFile(testdata("x9-status-summary-detail.txt"))
	require.NoError(t, err)
	scan := util.NewScanner(strings.NewReader(string(raw)))
	f := NewFile()
	require.NoError(t, f.Read(&scan))

	acct := f.Banks[0].Accounts[0]
	require.Equal(t, "CAD", acct.CurrencyCode)
	require.Len(t, acct.Summaries, 2)
	require.Equal(t, "010", acct.Summaries[0].TypeCode)
	require.Equal(t, "500000", acct.Summaries[0].Amount)
	require.Equal(t, TypeLevelStatus, acct.Summaries[0].Level)
	require.Equal(t, "190", acct.Summaries[1].TypeCode)
	require.Equal(t, "70000000", acct.Summaries[1].Amount)
	require.Equal(t, int64(4), acct.Summaries[1].ItemCount)
	require.Equal(t, FundsType0, string(acct.Summaries[1].FundsType.TypeCode))

	require.Len(t, acct.Details, 1)
	require.Equal(t, "399", acct.Details[0].TypeCode)
	require.Equal(t, "25000", acct.Details[0].Amount)
}

func TestX9BankHeaderAllFields(t *testing.T) {
	raw, err := os.ReadFile(testdata("x9-bank-header-all-fields.txt"))
	require.NoError(t, err)
	scan := util.NewScanner(strings.NewReader(string(raw)))
	f := NewFile()
	require.NoError(t, f.Read(&scan))
	require.Equal(t, "SWXXXXXX", f.Banks[0].Receiver)
	require.Equal(t, "2359", f.Banks[0].AsOfTime)
	require.Empty(t, f.Banks[0].CurrencyCode)
	require.Equal(t, FundsTypeS, string(f.Banks[0].Accounts[0].Details[0].FundsType.TypeCode))
}

func TestX9ContinuationDetail(t *testing.T) {
	raw, err := os.ReadFile(testdata("x9-continuation-16.txt"))
	require.NoError(t, err)
	scan := util.NewScanner(strings.NewReader(string(raw)))
	f := NewFile()
	require.NoError(t, f.Read(&scan))
	d := f.Banks[0].Accounts[0].Details[0]
	require.Equal(t, "165", d.TypeCode)
	require.Equal(t, "TRACE12345678900", d.BankReferenceNumber)
	require.Equal(t, "CUST REF 79", d.CustomerReferenceNumber)
	require.Equal(t, "THIS IS THE TEXT HERE", d.Text)
}

func TestX9AccountCurrencyRequired(t *testing.T) {
	raw := "01,122099999,123456789,150623,0200,1,,,3/\n" +
		"02,,122099999,1,150622,,,2/\n" +
		"03,0987654321,,,,,/\n" +
		"49,0,2/\n98,0,1,4/\n99,0,1,6/\n"
	scan := util.NewScanner(strings.NewReader(raw))
	f := NewFile()
	err := f.Read(&scan)
	require.Error(t, err)
	require.Contains(t, err.Error(), "CurrencyCode")
}

func TestX9SplitSummaryContinuation(t *testing.T) {
	raw, err := os.ReadFile(testdata("x9-03-split-summary.txt"))
	require.NoError(t, err)
	scan := util.NewScanner(strings.NewReader(string(raw)))
	f := NewFile()
	require.NoError(t, f.Read(&scan))

	codes := make([]string, 0, len(f.Banks[0].Accounts[0].Summaries))
	for _, s := range f.Banks[0].Accounts[0].Summaries {
		codes = append(codes, s.TypeCode)
	}
	require.Equal(t, []string{"010", "015", "040", "045", "050", "055", "072", "073", "074", "100", "190", "400", "570"}, codes)
	require.Equal(t, "99666666", f.Banks[0].Accounts[0].Summaries[9].Amount)
	require.Equal(t, int64(3), f.Banks[0].Accounts[0].Summaries[9].ItemCount)
}

func TestCreateFillsTrailers(t *testing.T) {
	f := NewFile()
	f.Sender = "122099999"
	f.Receiver = "123456789"
	f.FileCreatedDate = "150623"
	f.FileCreatedTime = "0200"
	f.FileIdNumber = "1"
	f.Banks = []Bank{{
		BankIdentification: "122099999",
		GroupStatus:        1,
		AsOfDate:           "150622",
		AsOfDateModifier:   2,
		Accounts: []Account{{
			AccountNumber: "0987654321",
			CurrencyCode:  "GBP",
			Summaries:     []AccountSummary{{TypeCode: "010", Amount: "500000"}},
		}},
	}}
	require.NoError(t, f.Create())
	require.Equal(t, "500000", f.Banks[0].Accounts[0].AccountControlTotal)
	require.Equal(t, int64(2), f.Banks[0].Accounts[0].NumberRecords)
	require.Equal(t, "500000", f.Banks[0].BankControlTotal)
	require.Equal(t, int64(1), f.Banks[0].NumberOfAccounts)
	require.Equal(t, int64(4), f.Banks[0].NumberOfRecords)
	require.Equal(t, "500000", f.FileControlTotal)
	require.Equal(t, int64(1), f.NumberOfBanks)
	require.Equal(t, int64(6), f.NumberOfRecords)
	require.NoError(t, f.Validate())
}

func TestFundsTypeDStillParsed(t *testing.T) {
	raw := "01,122099999,123456789,150623,0200,1,,,3/\n" +
		"02,,122099999,1,150622,,,2/\n" +
		"03,0987654321,GBP,,,,/\n" +
		"16,399,25000,D,1,1,100,,,/\n" +
		"49,25000,3/\n98,25000,1,5/\n99,25000,1,7/\n"
	scan := util.NewScanner(strings.NewReader(raw))
	f := NewFile()
	require.NoError(t, f.Read(&scan))
	require.Equal(t, FundsTypeD, string(f.Banks[0].Accounts[0].Details[0].FundsType.TypeCode))
	require.Equal(t, int64(1), f.Banks[0].Accounts[0].Details[0].FundsType.DistributionNumber)
}

func TestStrictBankHeaderRejectsCurrency(t *testing.T) {
	raw := "01,122099999,123456789,150623,0200,1,,,3/\n" +
		"02,,122099999,1,150622,,USD,2/\n" +
		"03,0987654321,GBP,,,,/\n" +
		"49,0,2/\n98,0,1,4/\n99,0,1,6/\n"
	scan := util.NewScanner(strings.NewReader(raw))
	f := NewFileWith(Options{StrictBankHeader: true})
	require.NoError(t, f.Read(&scan))
	err := f.Validate()
	require.Error(t, err)
	require.Contains(t, err.Error(), "currency must be empty")
}

func TestGoldenRoundTrip(t *testing.T) {
	raw, err := os.ReadFile(testdata("x9-mandatory.txt"))
	require.NoError(t, err)
	scan := util.NewScanner(strings.NewReader(string(raw)))
	f := NewFile()
	require.NoError(t, f.Read(&scan))
	scan2 := util.NewScanner(strings.NewReader(f.String()))
	f2 := NewFile()
	require.NoError(t, f2.Read(&scan2))
	require.Equal(t, f.Sender, f2.Sender)
	require.Equal(t, f.Banks[0].BankIdentification, f2.Banks[0].BankIdentification)
	require.Equal(t, f.Banks[0].Accounts[0].CurrencyCode, f2.Banks[0].Accounts[0].CurrencyCode)
}

func TestX9RecordSamples(t *testing.T) {
	var hdr fileHeader
	_, err := hdr.parse("01,122099999,123456789,150623,0200,1,,,3/", Options{})
	require.NoError(t, err)
	require.Equal(t, int64(3), hdr.VersionNumber)

	var bank bankHeader
	_, err = bank.parse("02,,122099999,1,150622,,,2/")
	require.NoError(t, err)
	require.Equal(t, "122099999", bank.BankIdentification)
	require.Empty(t, bank.CurrencyCode)

	_, err = bank.parse("02,SWXXXXXX,122099999,1,150622,2359,,2/")
	require.NoError(t, err)
	require.Equal(t, "SWXXXXXX", bank.Receiver)
	require.Equal(t, "2359", bank.AsOfTime)

	var acct accountIdentifier
	_, err = acct.parse("03,0987654321,GBP,,,,/")
	require.NoError(t, err)
	require.Equal(t, "GBP", acct.CurrencyCode)
	require.Empty(t, acct.Summaries)

	var det transactionDetail
	_, err = det.parse("16,399,25000,0,,,/")
	require.NoError(t, err)
	require.Equal(t, "399", det.TypeCode)
	require.Equal(t, FundsType0, string(det.FundsType.TypeCode))
}

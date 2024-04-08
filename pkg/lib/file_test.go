// Copyright 2022 The Moov Authors
// Use of this source code is governed by an Apache License
// license that can be found in the LICENSE file.

package lib

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestFileWithSampleData(t *testing.T) {
	paths := []string{
		"sample1.txt",
		"sample2.txt",
		"sample3.txt",
		"sample4-continuations-newline-delimited.txt",
	}

	for _, path := range paths {
		samplePath := filepath.Join("..", "..", "test", "testdata", path)
		fd, err := os.Open(samplePath)
		require.NoError(t, err)

		scan := NewBai2Scanner(fd)
		f := NewBai2()
		err = f.Read(&scan)

		require.NoError(t, err)
		require.NoError(t, f.Validate())
	}
}

func TestFileWithContinuationRecord(t *testing.T) {

	raw := `01,0004,12345,060321,0829,001,80,1,2/
02,12345,0004,1,060317,,CAD,/
03,10200123456,CAD,040,+000000000000,,,045,+000000000000,,/
88,046,+000000000000,,,047,+000000000000,,,048,+000000000000,,,049,+000000000000,,/
88,050,+000000000000,,,051,+000000000000,,,052,+000000000000,,,053,+000000000000,,/
16,409,000000000002500,V,060316,1300,,,RETURNED CHEQUE     /
16,409,000000000090000,V,060316,1300,,,RTN-UNKNOWN         /
49,+00000000000834000,14/
98,+00000000001280000,2,25/
99,+00000000001280000,1,27/`

	scan := NewBai2Scanner(strings.NewReader(raw))
	f := NewBai2()
	err := f.Read(&scan)
	require.NoError(t, err)
	require.NoError(t, f.Validate())

	expected := `01,0004,12345,060321,0829,001,80,1,2/
02,12345,0004,1,060317,,CAD,/
03,10200123456,CAD,040,+000000000000,,,045,+000000000000,,,046,+000000000000,,/
88,047,+000000000000,,,048,+000000000000,,,049,+000000000000,,,050/
88,+000000000000,,,051,+000000000000,,,052,+000000000000,,,053,+000000000000,,/
16,409,000000000002500,V,060316,1300,,,RETURNED CHEQUE     /
16,409,000000000090000,V,060316,1300,,,RTN-UNKNOWN         /
49,+00000000000834000,14/
98,+00000000001280000,2,25/
99,+00000000001280000,1,27/`
	require.Equal(t, expected, f.String())
}

func TestSumFileRecords(t *testing.T) {
	file := Bai2{}
	file.Groups = []Group{
		{NumberOfRecords: 27},
	}
	require.Equal(t, int64(29), file.SumRecords())
}

func TestScannedFileTrailerRecordCount(t *testing.T) {

	raw := `01,0004,12345,060321,0829,001,80,1,2/
02,12345,0004,1,060317,,CAD,/
03,10200123456,CAD,040,+000000000000,,,045,+000000000000,,/
88,046,+000000000000,,,047,+000000000000,,,048,+000000000000,,,049,+000000000000,,/
88,050,+000000000000,,,051,+000000000000,,,052,+000000000000,,,053,+000000000000,,/
16,409,000000000002500,V,060316,1300,,,RETURNED CHEQUE     /
16,409,000000000090000,V,060316,1300,,,RTN-UNKNOWN         /
49,+00000000000834000,6/
98,+00000000001280000,2,8/
99,+00000000001280000,1,10/`

	scan := NewBai2Scanner(bytes.NewReader([]byte(raw)))
	file := NewBai2()
	err := file.Read(&scan)
	require.NoError(t, err)
	require.Equal(t, int64(10), file.SumRecords())
}

func TestSumNumberOfGroups(t *testing.T) {
	file := Bai2{}
	file.Groups = []Group{
		{NumberOfRecords: 27},
		{NumberOfRecords: 27},
		{NumberOfRecords: 27},
	}
	require.Equal(t, int64(3), file.SumNumberOfGroups())
}

func TestSumGroupControlTotals(t *testing.T) {
	group := Group{}
	group.Receiver = "121000358"
	group.Originator = "121000358"
	group.GroupStatus = 1
	group.AsOfDate = time.Now().Format("060102")
	group.AsOfTime = time.Now().Format("1504")
	group.AsOfDateModifier = 2
	group.GroupControlTotal = "200"

	file := NewBai2()
	file.Sender = "121000358"
	file.Receiver = "121000358"
	file.FileCreatedDate = time.Now().Format("060102")
	file.FileCreatedTime = time.Now().Format("1504")
	file.FileIdNumber = "01"
	file.PhysicalRecordLength = 80
	file.BlockSize = 1
	file.VersionNumber = 2
	file.Groups = append(file.Groups, group)
	file.NumberOfRecords = file.SumRecords()

	total, err := file.SumGroupControlTotals()
	require.NoError(t, err)
	require.Equal(t, "200", total)
}

func TestBuildFileAggregates(t *testing.T) {

	recordLength := int64(80)

	details := []Detail{}
	for i := 0; i < 10; i++ {
		detail := NewDetail()
		detail.TypeCode = "409"
		detail.Amount = "274006"
		detail.BankReferenceNumber = "1234567"
		detail.Text = "TV Purchase"
		details = append(details, *detail)
	}

	account1 := Account{}
	account1.AccountNumber = "1234567"
	account1.CurrencyCode = "USD"
	account1.Details  = append(account1.Details, details...)
	account1.Summaries = []AccountSummary{
		{TypeCode: "040", Amount: "+000000000000"},
		{TypeCode: "045", Amount: "+000000000000"},
		{TypeCode: "046", Amount: "+000000000000"},
		{TypeCode: "047", Amount: "+000000000000"},
		{TypeCode: "048", Amount: "+000000000000"},
		{TypeCode: "049", Amount: "+000000000000"},
		{TypeCode: "050", Amount: "+000000000000"},
		{TypeCode: "051", Amount: "+000000000000"},
		{TypeCode: "052", Amount: "+000000000000"},
		{TypeCode: "053", Amount: "+000000000000"},
	}
	controlTotal, _ := account1.SumDetailAmounts()
	account1.AccountControlTotal = controlTotal
	account1.NumberRecords = account1.SumRecords(recordLength)

	account2 := Account{}
	account2.AccountNumber = "1234567"
	account2.CurrencyCode = "USD"
	account2.Details  = append(account2.Details, details...)
	controlTotal, _ = account2.SumDetailAmounts()
	account2.AccountControlTotal = controlTotal
	account2.NumberRecords = account2.SumRecords(recordLength)

	group := Group{}
	group.Receiver = "121000358"
	group.Originator = "121000358"
	group.GroupStatus = 1
	group.AsOfDate = time.Now().Format("060102")
	group.AsOfTime = time.Now().Format("1504")
	group.AsOfDateModifier = 2
	group.Accounts = append(group.Accounts, account1, account2)
	controlTotal, _ = group.SumAccountControlTotals()
	group.GroupControlTotal = controlTotal
	group.NumberOfAccounts = group.SumNumberOfAccounts()
	group.NumberOfRecords = group.SumRecords()

	file := NewBai2()
	file.Sender = "121000358"
	file.Receiver = "121000358"
	file.FileCreatedDate = time.Now().Format("060102")
	file.FileCreatedTime = time.Now().Format("1504")
	file.FileIdNumber = "01"
	file.PhysicalRecordLength = recordLength
	file.BlockSize = 1
	file.VersionNumber = 2
	file.Groups = append(file.Groups, group)
	controlTotal, _ = file.SumGroupControlTotals()
	file.FileControlTotal = controlTotal
	file.NumberOfGroups = file.SumNumberOfGroups()
	file.NumberOfRecords = file.SumRecords()

	require.Equal(t, "-5480120", file.FileControlTotal)
	require.Equal(t, int64(1), file.NumberOfGroups)
	require.Equal(t, int64(30), file.NumberOfRecords)

	fmt.Print(file.String())

}

// Copyright 2022 The Moov Authors
// Use of this source code is governed by an Apache License
// license that can be found in the LICENSE file.

package bai3

import "fmt"

// Create fills trailer aggregates (control totals and record counts) from the
// file body, using PhysicalRecordLength to count continuation records the same
// way String() writes them.
func (f *File) Create() error {
	for i := range f.Banks {
		if err := f.Banks[i].create(f.PhysicalRecordLength); err != nil {
			return err
		}
	}

	total, err := f.sumBankControlTotals()
	if err != nil {
		return err
	}
	f.FileControlTotal = total
	f.NumberOfBanks = f.SumNumberOfBanks()
	f.NumberOfRecords = f.SumRecords()
	return nil
}

func (g *Bank) create(recordLength int64) error {
	for i := range g.Accounts {
		total, err := g.Accounts[i].sumAmounts()
		if err != nil {
			return err
		}
		g.Accounts[i].AccountControlTotal = total
		g.Accounts[i].NumberRecords = g.Accounts[i].SumRecords(recordLength)
	}

	total, err := g.sumAccountControlTotals()
	if err != nil {
		return err
	}
	g.BankControlTotal = total
	g.NumberOfAccounts = g.SumNumberOfAccounts()
	g.NumberOfRecords = g.SumRecords()
	return nil
}

func (g *Bank) sumAccountControlTotals() (string, error) {
	var sum int64
	for _, account := range g.Accounts {
		amt, err := parseAmount(account.AccountControlTotal)
		if err != nil {
			return "0", err
		}
		sum += amt
	}
	return formatAmount(sum), nil
}

func (f *File) sumBankControlTotals() (string, error) {
	var sum int64
	for _, group := range f.Banks {
		amt, err := parseAmount(group.BankControlTotal)
		if err != nil {
			return "0", err
		}
		sum += amt
	}
	return formatAmount(sum), nil
}

func (f *File) validateControlTotals() error {
	for gi := range f.Banks {
		g := &f.Banks[gi]
		for ai := range g.Accounts {
			a := &g.Accounts[ai]
			got, err := a.sumAmounts()
			if err != nil {
				return err
			}
			if !amountsEqual(got, a.AccountControlTotal) {
				return fmt.Errorf("account %s control total %s does not match computed %s", a.AccountNumber, a.AccountControlTotal, got)
			}
			n := a.SumRecords(f.PhysicalRecordLength)
			if a.NumberRecords != n {
				return fmt.Errorf("account %s record count %d does not match computed %d", a.AccountNumber, a.NumberRecords, n)
			}
		}

		got, err := g.sumAccountControlTotals()
		if err != nil {
			return err
		}
		if !amountsEqual(got, g.BankControlTotal) {
			return fmt.Errorf("group control total %s does not match computed %s", g.BankControlTotal, got)
		}
		if g.NumberOfAccounts != g.SumNumberOfAccounts() {
			return fmt.Errorf("group number of accounts %d does not match computed %d", g.NumberOfAccounts, g.SumNumberOfAccounts())
		}
		if g.NumberOfRecords != g.SumRecords() {
			return fmt.Errorf("group number of records %d does not match computed %d", g.NumberOfRecords, g.SumRecords())
		}
	}

	got, err := f.sumBankControlTotals()
	if err != nil {
		return err
	}
	if !amountsEqual(got, f.FileControlTotal) {
		return fmt.Errorf("file control total %s does not match computed %s", f.FileControlTotal, got)
	}
	if f.NumberOfBanks != f.SumNumberOfBanks() {
		return fmt.Errorf("file number of groups %d does not match computed %d", f.NumberOfBanks, f.SumNumberOfBanks())
	}
	if f.NumberOfRecords != f.SumRecords() {
		return fmt.Errorf("file number of records %d does not match computed %d", f.NumberOfRecords, f.SumRecords())
	}
	return nil
}

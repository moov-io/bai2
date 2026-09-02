// Copyright 2022 The Moov Authors
// Use of this source code is governed by an Apache License
// license that can be found in the LICENSE file.

package lib

import "fmt"

// Create fills trailer aggregates (control totals and record counts) from the
// file body, using PhysicalRecordLength to count continuation records the same
// way String() writes them.
func (f *Bai2) Create() error {
	for i := range f.Groups {
		if err := f.Groups[i].create(f.PhysicalRecordLength); err != nil {
			return err
		}
	}

	total, err := f.sumGroupControlTotals()
	if err != nil {
		return err
	}
	f.FileControlTotal = total
	f.NumberOfGroups = f.SumNumberOfGroups()
	f.NumberOfRecords = f.SumRecords()
	return nil
}

func (g *Group) create(recordLength int64) error {
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
	g.GroupControlTotal = total
	g.NumberOfAccounts = g.SumNumberOfAccounts()
	g.NumberOfRecords = g.SumRecords()
	return nil
}

func (g *Group) sumAccountControlTotals() (string, error) {
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

func (f *Bai2) sumGroupControlTotals() (string, error) {
	var sum int64
	for _, group := range f.Groups {
		amt, err := parseAmount(group.GroupControlTotal)
		if err != nil {
			return "0", err
		}
		sum += amt
	}
	return formatAmount(sum), nil
}

func (f *Bai2) validateControlTotals() error {
	for gi := range f.Groups {
		g := &f.Groups[gi]
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
		if !amountsEqual(got, g.GroupControlTotal) {
			return fmt.Errorf("group control total %s does not match computed %s", g.GroupControlTotal, got)
		}
		if g.NumberOfAccounts != g.SumNumberOfAccounts() {
			return fmt.Errorf("group number of accounts %d does not match computed %d", g.NumberOfAccounts, g.SumNumberOfAccounts())
		}
		if g.NumberOfRecords != g.SumRecords() {
			return fmt.Errorf("group number of records %d does not match computed %d", g.NumberOfRecords, g.SumRecords())
		}
	}

	got, err := f.sumGroupControlTotals()
	if err != nil {
		return err
	}
	if !amountsEqual(got, f.FileControlTotal) {
		return fmt.Errorf("file control total %s does not match computed %s", f.FileControlTotal, got)
	}
	if f.NumberOfGroups != f.SumNumberOfGroups() {
		return fmt.Errorf("file number of groups %d does not match computed %d", f.NumberOfGroups, f.SumNumberOfGroups())
	}
	if f.NumberOfRecords != f.SumRecords() {
		return fmt.Errorf("file number of records %d does not match computed %d", f.NumberOfRecords, f.SumRecords())
	}
	return nil
}

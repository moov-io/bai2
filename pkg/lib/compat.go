// Copyright 2022 The Moov Authors
// Use of this source code is governed by an Apache License
// license that can be found in the LICENSE file.

// Package lib is a compatibility shim for github.com/moov-io/bai2/pkg/bai2.
//
// Deprecated: import github.com/moov-io/bai2/pkg/bai2 instead. See docs/MIGRATING.md.
package lib

import (
	"io"

	"github.com/moov-io/bai2/pkg/bai2"
)

type (
	Bai2            = bai2.Bai2
	Options         = bai2.Options
	Group           = bai2.Group
	Account         = bai2.Account
	AccountSummary  = bai2.AccountSummary
	Detail          = bai2.Detail
	FundsType       = bai2.FundsType
	FundsTypeCode   = bai2.FundsTypeCode
	Distribution    = bai2.Distribution
	TypeCode        = bai2.TypeCode
	TypeLevel       = bai2.TypeLevel
	TransactionKind = bai2.TransactionKind
	Bai2Scanner     = bai2.Bai2Scanner
)

const (
	VersionBAI2        = bai2.VersionBAI2
	VersionBTR3        = bai2.VersionBTR3
	DefaultCurrency    = bai2.DefaultCurrency
	DefaultGroupStatus = bai2.DefaultGroupStatus
	TypeLevelStatus    = bai2.TypeLevelStatus
	TypeLevelSummary   = bai2.TypeLevelSummary
	TypeLevelDetail    = bai2.TypeLevelDetail
	TransactionNA      = bai2.TransactionNA
	TransactionCR      = bai2.TransactionCR
	TransactionDB      = bai2.TransactionDB
	FundsType0         = bai2.FundsType0
	FundsType1         = bai2.FundsType1
	FundsType2         = bai2.FundsType2
	FundsTypeS         = bai2.FundsTypeS
	FundsTypeV         = bai2.FundsTypeV
	FundsTypeD         = bai2.FundsTypeD
	FundsTypeZ         = bai2.FundsTypeZ
)

var (
	NewBai2           = bai2.NewBai2
	NewBai2With       = bai2.NewBai2With
	NewBai2Scanner    = bai2.NewBai2Scanner
	NewGroup          = bai2.NewGroup
	NewAccount        = bai2.NewAccount
	NewDetail         = bai2.NewDetail
	LookupTypeCode    = bai2.LookupTypeCode
	ClassifyTypeCode  = bai2.ClassifyTypeCode
	IsStatusTypeCode  = bai2.IsStatusTypeCode
	IsSummaryTypeCode = bai2.IsSummaryTypeCode
)

// NewScanner is an alias for NewBai2Scanner kept for older call sites.
func NewScanner(r io.Reader) Bai2Scanner {
	return NewBai2Scanner(r)
}

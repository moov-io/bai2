// Copyright 2022 The Moov Authors
// Use of this source code is governed by an Apache License
// license that can be found in the LICENSE file.

package bai2

import "strconv"

// TypeLevel is the BAI2 type-code level from Appendix A.
type TypeLevel string

const (
	TypeLevelStatus  TypeLevel = "Status"
	TypeLevelSummary TypeLevel = "Summary"
	TypeLevelDetail  TypeLevel = "Detail"
)

// TransactionKind is CR, DB, or NA as defined in Appendix A.
type TransactionKind string

const (
	TransactionNA TransactionKind = "NA"
	TransactionCR TransactionKind = "CR"
	TransactionDB TransactionKind = "DB"
)

// TypeCode describes a uniform BAI2 type code from Appendix A.
type TypeCode struct {
	Code        string
	Transaction TransactionKind
	Level       TypeLevel
	Description string
}

// LookupTypeCode returns the Appendix A catalog entry for code, if present.
func LookupTypeCode(code string) (TypeCode, bool) {
	tc, ok := typeCodeCatalog[code]
	return tc, ok
}

func typeMeta(code string) (TransactionKind, TypeLevel, string) {
	if tc, ok := LookupTypeCode(code); ok {
		return tc.Transaction, tc.Level, tc.Description
	}
	level, kind := ClassifyTypeCode(code)
	return kind, level, ""
}

// ClassifyTypeCode returns level and transaction kind from the spec ranges
// (Appendix A, Type Code Ranges). Unknown 3-digit codes still classify by range
// so custom 900–999 codes work. Empty or malformed codes return zero values.
func ClassifyTypeCode(code string) (TypeLevel, TransactionKind) {
	if tc, ok := typeCodeCatalog[code]; ok {
		return tc.Level, tc.Transaction
	}

	n, err := strconv.Atoi(code)
	if err != nil || n < 0 || n > 999 || len(code) != 3 {
		return "", ""
	}

	switch {
	case n <= 99:
		return TypeLevelStatus, TransactionNA
	case n == 100:
		return TypeLevelSummary, TransactionCR
	case n <= 399:
		return TypeLevelSummary, TransactionCR
	case n == 400:
		return TypeLevelSummary, TransactionDB
	case n <= 699:
		return TypeLevelSummary, TransactionDB
	case n <= 799:
		return TypeLevelSummary, TransactionNA
	case n == 890:
		return TypeLevelDetail, TransactionNA
	case n >= 900 && n <= 919:
		return TypeLevelStatus, TransactionNA
	case n >= 920 && n <= 959:
		return TypeLevelSummary, TransactionCR
	case n >= 960:
		return TypeLevelSummary, TransactionDB
	default:
		return "", ""
	}
}

// IsStatusTypeCode reports whether code is an account-status type (03 only,
// no item count or funds type).
func IsStatusTypeCode(code string) bool {
	level, _ := ClassifyTypeCode(code)
	return level == TypeLevelStatus
}

// IsSummaryTypeCode reports whether code is an activity-summary type (03).
func IsSummaryTypeCode(code string) bool {
	level, _ := ClassifyTypeCode(code)
	if level == TypeLevelSummary {
		return true
	}
	// Codes in 101–399 / 401–699 are both summary (03) and detail (16).
	// When used on a 03 they are summaries.
	n, err := strconv.Atoi(code)
	if err != nil {
		return false
	}
	return (n >= 100 && n <= 399) || (n >= 400 && n <= 699) || (n >= 700 && n <= 799) || (n >= 920 && n <= 999)
}

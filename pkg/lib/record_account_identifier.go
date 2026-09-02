// Copyright 2022 The Moov Authors
// Use of this source code is governed by an Apache License
// license that can be found in the LICENSE file.

package lib

import (
	"bytes"
	"fmt"
	"unicode"

	"github.com/moov-io/bai2/pkg/util"
)

const (
	aiParseErrorFmt    = "AccountIdentifier: unable to parse %s"
	aiValidateErrorFmt = "AccountIdentifier: invalid %s"
)

type AccountSummary struct {
	TypeCode    string          `json:"TypeCode"`
	Amount      string          `json:"Amount"`
	ItemCount   int64           `json:"ItemCount"`
	FundsType   FundsType       `json:"FundsType"`
	Transaction TransactionKind `json:"Transaction,omitempty"`
	Level       TypeLevel       `json:"Level,omitempty"`
	Description string          `json:"Description,omitempty"`
}

// TypeInfo returns the Appendix A catalog entry for this summary's type code.
func (s AccountSummary) TypeInfo() (TypeCode, bool) {
	return LookupTypeCode(s.TypeCode)
}

type accountIdentifier struct {
	AccountNumber string
	CurrencyCode  string

	Summaries []AccountSummary
}

func (r *accountIdentifier) validate() error {

	if r.AccountNumber == "" {
		return fmt.Errorf(aiValidateErrorFmt, "AccountNumber")
	}

	if r.CurrencyCode != "" && !util.ValidateCurrencyCode(r.CurrencyCode) {
		return fmt.Errorf(aiValidateErrorFmt, "CurrencyCode")
	}

	for _, summary := range r.Summaries {
		if summary.Amount != "" && !util.ValidateAmount(summary.Amount) {
			return fmt.Errorf(aiValidateErrorFmt, "Amount")
		}
		if summary.TypeCode != "" && !util.ValidateTypeCode(summary.TypeCode) {
			return fmt.Errorf(aiValidateErrorFmt, "TypeCode")
		}
		if summary.FundsType.Validate() != nil {
			return fmt.Errorf(aiValidateErrorFmt, "FundsType")
		}
	}

	return nil
}

func (r *accountIdentifier) parse(data string) (int, error) {

	var line string
	var err error
	var size, read int

	length := util.GetSize(data)
	if length < 3 {
		return 0, fmt.Errorf(aiParseErrorFmt, "record")
	} else {
		line = data[:length]
	}

	// RecordCode
	if util.AccountIdentifierCode != data[:2] {
		return 0, fmt.Errorf(aiParseErrorFmt, "RecordCode")
	}
	read += 3

	// AccountNumber
	if r.AccountNumber, size, err = util.ReadField(line, read); err != nil {
		return 0, fmt.Errorf(aiParseErrorFmt, "AccountNumber")
	} else {
		read += size
	}

	// CurrencyCode
	if r.CurrencyCode, size, err = util.ReadField(line, read); err != nil {
		return 0, fmt.Errorf(aiParseErrorFmt, "CurrencyCode")
	} else {
		read += size
	}

	for read < len(line) {
		if remainingFieldsEmpty(line[read:]) {
			break
		}

		var summary AccountSummary

		// TypeCode
		if summary.TypeCode, size, err = util.ReadField(line, read); err != nil {
			return 0, fmt.Errorf(aiParseErrorFmt, "TypeCode")
		} else {
			read += size
		}

		// Amount
		if summary.Amount, size, err = util.ReadField(line, read); err != nil {
			return 0, fmt.Errorf(aiParseErrorFmt, "Amount")
		} else {
			read += size
		}

		// ItemCount
		if summary.ItemCount, size, err = util.ReadFieldAsInt(line, read); err != nil {
			return 0, fmt.Errorf(aiParseErrorFmt, "ItemCount")
		} else {
			read += size
		}

		if size, err = summary.FundsType.parse(line[read:]); err != nil {
			return 0, fmt.Errorf(aiParseErrorFmt, "FundsType")
		} else {
			read += size
		}

		// Spec: a defaulted type code means no status/summary is reported for this slot.
		if summary.TypeCode == "" && summary.Amount == "" {
			continue
		}

		summary.Transaction, summary.Level, summary.Description = typeMeta(summary.TypeCode)
		r.Summaries = append(r.Summaries, summary)
	}

	if remainingFieldsEmpty(line[read:]) {
		read = len(line)
	}

	if err = r.validate(); err != nil {
		return 0, err
	}

	return read, nil
}

func (r *accountIdentifier) string(opts ...int64) string {

	var maxLen int64
	if len(opts) > 0 {
		maxLen = opts[0]
	}

	var total, buf bytes.Buffer

	buf.WriteString(fmt.Sprintf("%s,", util.AccountIdentifierCode))
	buf.WriteString(fmt.Sprintf("%s,", r.AccountNumber))
	buf.WriteString(fmt.Sprintf("%s,", r.CurrencyCode))

	if len(r.Summaries) == 0 {
		buf.WriteString(",,,")
	} else {
		for index, summary := range r.Summaries {

			util.WriteBuffer(&total, &buf, summary.TypeCode, maxLen)
			buf.WriteString(",")

			util.WriteBuffer(&total, &buf, summary.Amount, maxLen)
			buf.WriteString(",")

			if summary.ItemCount == 0 {
				buf.WriteString(",")
			} else {
				util.WriteBuffer(&total, &buf, fmt.Sprintf("%d", summary.ItemCount), maxLen)
				buf.WriteString(",")
			}

			util.WriteBuffer(&total, &buf, summary.FundsType.String(), maxLen)

			if index < len(r.Summaries)-1 {
				buf.WriteString(",")
			}
		}
	}

	buf.WriteString("/")
	total.WriteString(buf.String())

	return total.String()
}

func remainingFieldsEmpty(s string) bool {
	for _, r := range s {
		if r != ',' && r != '/' && !unicode.IsSpace(r) {
			return false
		}
	}
	return true
}

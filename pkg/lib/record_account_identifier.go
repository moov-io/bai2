// Copyright 2022 The Moov Authors
// Use of this source code is governed by an Apache License
// license that can be found in the LICENSE file.

package lib

import (
	"bytes"
	"fmt"

	"github.com/moov-io/bai2/pkg/util"
)

const (
	aiParseErrorFmt    = "AccountIdentifier: unable to parse %s"
	aiValidateErrorFmt = "AccountIdentifierCurrent: invalid %s"
)

type AccountSummary struct {
	TypeCode  string
	Amount    string
	ItemCount int64
	FundsType FundsType
}

type accountIdentifier struct {
	AccountNumber string
	CurrencyCode  string

	Summaries []AccountSummary
}

func (r *accountIdentifier) validate() error {

	if r.AccountNumber == "" {
		return fmt.Errorf(fmt.Sprintf(aiValidateErrorFmt, "AccountNumber"))
	}

	if r.CurrencyCode != "" && !util.ValidateCurrencyCode(r.CurrencyCode) {
		return fmt.Errorf(fmt.Sprintf(aiValidateErrorFmt, "CurrencyCode"))
	}

	for _, summary := range r.Summaries {
		if summary.Amount != "" && !util.ValidateAmount(summary.Amount) {
			return fmt.Errorf(fmt.Sprintf(aiValidateErrorFmt, "Amount"))
		}
		if summary.TypeCode != "" && !util.ValidateTypeCode(summary.TypeCode) {
			return fmt.Errorf(fmt.Sprintf(aiValidateErrorFmt, "TypeCode"))
		}
		if summary.FundsType.Validate() != nil {
			return fmt.Errorf(fmt.Sprintf(aiValidateErrorFmt, "FundsType"))
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
		return 0, fmt.Errorf(fmt.Sprintf(aiParseErrorFmt, "record"))
	} else {
		line = data[:length]
	}

	// RecordCode
	if util.AccountIdentifierCode != data[:2] {
		return 0, fmt.Errorf(fmt.Sprintf(aiParseErrorFmt, "RecordCode"))
	}
	read += 3

	// AccountNumber
	if r.AccountNumber, size, err = util.ReadField(line, read); err != nil {
		return 0, fmt.Errorf(fmt.Sprintf(aiParseErrorFmt, "AccountNumber"))
	} else {
		read += size
	}

	// CurrencyCode
	if r.CurrencyCode, size, err = util.ReadField(line, read); err != nil {
		return 0, fmt.Errorf(fmt.Sprintf(aiParseErrorFmt, "CurrencyCode"))
	} else {
		read += size
	}

	for read < len(data) {

		var summary AccountSummary

		// TypeCode
		if summary.TypeCode, size, err = util.ReadField(line, read); err != nil {
			return 0, fmt.Errorf(fmt.Sprintf(aiParseErrorFmt, "TypeCode"))
		} else {
			read += size
		}

		// Amount
		if summary.Amount, size, err = util.ReadField(line, read); err != nil {
			return 0, fmt.Errorf(fmt.Sprintf(aiParseErrorFmt, "Amount"))
		} else {
			read += size
		}

		// ItemCount
		if summary.ItemCount, size, err = util.ReadFieldAsInt(line, read); err != nil {
			return 0, fmt.Errorf(fmt.Sprintf(aiParseErrorFmt, "ItemCount"))
		} else {
			read += size
		}

		if size, err = summary.FundsType.parse(line[read:]); err != nil {
			return 0, fmt.Errorf(fmt.Sprintf(aiParseErrorFmt, "FundsType"))
		} else {
			read += size
		}

		r.Summaries = append(r.Summaries, summary)
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

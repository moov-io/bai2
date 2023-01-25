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

type accountIdentifier struct {
	AccountNumber string
	CurrencyCode  string   `json:",omitempty"`
	TypeCode      string   `json:",omitempty"`
	Amount        string   `json:",omitempty"`
	ItemCount     int64    `json:",omitempty"`
	FundsType     string   `json:",omitempty"`
	Composite     []string `json:",omitempty"`
}

func (h *accountIdentifier) validate() error {
	if h.AccountNumber == "" {
		return fmt.Errorf(fmt.Sprintf(aiValidateErrorFmt, "AccountNumber"))
	}
	if h.Amount != "" && !util.ValidateAmount(h.Amount) {
		return fmt.Errorf(fmt.Sprintf(aiValidateErrorFmt, "Amount"))
	}
	if h.CurrencyCode != "" && !util.ValidateCurrencyCode(h.CurrencyCode) {
		return fmt.Errorf(fmt.Sprintf(aiValidateErrorFmt, "CurrencyCode"))
	}
	if h.TypeCode != "" && !util.ValidateTypeCode(h.TypeCode) {
		return fmt.Errorf(fmt.Sprintf(aiValidateErrorFmt, "TypeCode"))
	}
	if h.FundsType != "" && !util.ValidateFundsType(h.FundsType) {
		return fmt.Errorf(fmt.Sprintf(aiValidateErrorFmt, "FundsType"))
	}

	return nil
}

func (h *accountIdentifier) parse(data string) (int, error) {

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
	if h.AccountNumber, size, err = util.ReadField(line, read); err != nil {
		return 0, fmt.Errorf(fmt.Sprintf(aiParseErrorFmt, "AccountNumber"))
	} else {
		read += size
	}

	// CurrencyCode
	if h.CurrencyCode, size, err = util.ReadField(line, read); err != nil {
		return 0, fmt.Errorf(fmt.Sprintf(aiParseErrorFmt, "CurrencyCode"))
	} else {
		read += size
	}

	// TypeCode
	if h.TypeCode, size, err = util.ReadField(line, read); err != nil {
		return 0, fmt.Errorf(fmt.Sprintf(aiParseErrorFmt, "TypeCode"))
	} else {
		read += size
	}

	// Amount
	if h.Amount, size, err = util.ReadField(line, read); err != nil {
		return 0, fmt.Errorf(fmt.Sprintf(aiParseErrorFmt, "Amount"))
	} else {
		read += size
	}

	// ItemCount
	if h.ItemCount, size, err = util.ReadFieldAsInt(line, read); err != nil {
		return 0, fmt.Errorf(fmt.Sprintf(aiParseErrorFmt, "ItemCount"))
	} else {
		read += size
	}

	// FundsType
	if h.FundsType, size, err = util.ReadField(line, read); err != nil {
		return 0, fmt.Errorf(fmt.Sprintf(aiParseErrorFmt, "FundsType"))
	} else {
		read += size
	}

	for int64(read) < length {
		var composite string
		if composite, size, err = util.ReadField(line, read); err != nil {
			return 0, fmt.Errorf(fmt.Sprintf(aiParseErrorFmt, "Composite"))
		} else {
			read += size
		}
		h.Composite = append(h.Composite, composite)
	}

	if err = h.validate(); err != nil {
		return 0, err
	}

	return read, nil
}

func (h *accountIdentifier) string(opts ...int64) string {

	var totalBuf bytes.Buffer
	var buf bytes.Buffer

	buf.WriteString(fmt.Sprintf("%s,", util.AccountIdentifierCode))
	buf.WriteString(fmt.Sprintf("%s,", h.AccountNumber))
	buf.WriteString(fmt.Sprintf("%s,", h.CurrencyCode))
	buf.WriteString(fmt.Sprintf("%s,", h.TypeCode))
	buf.WriteString(fmt.Sprintf("%s,", h.Amount))
	if h.ItemCount > 0 {
		buf.WriteString(fmt.Sprintf("%d,", h.ItemCount))
	} else {
		buf.WriteString(",")
	}
	buf.WriteString(h.FundsType)

	var maxLen int64
	if len(opts) > 0 {
		maxLen = opts[0]
	}

	for _, composite := range h.Composite {
		if maxLen > 0 {
			if int64(buf.Len()+len(composite)+2) > maxLen {
				// refresh buf
				buf.WriteString("/" + "\n") // added new line
				totalBuf.WriteString(buf.String())

				// new buf
				buf = bytes.Buffer{}
				buf.WriteString(util.ContinuationCode)
			}
		}

		buf.WriteString(fmt.Sprintf(",%s", composite))
	}

	buf.WriteString("/")
	totalBuf.WriteString(buf.String())

	return totalBuf.String()
}

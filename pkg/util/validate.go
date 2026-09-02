// Copyright 2022 The Moov Authors
// Use of this source code is governed by an Apache License
// license that can be found in the LICENSE file.

package util

import (
	"regexp"
	"strconv"
)

var dateYYMMDDTypeRegex = regexp.MustCompile(`^[0-9]{2}(0[1-9]|1[0-2])(0[1-9]|[12][0-9]|3[01])$`)
var singedNumber = regexp.MustCompile(`^(-|\+|)?[0-9]\d*$`)
var currencyCodeRegex = regexp.MustCompile(`^[a-zA-Z]{3}$`)
var typeCodeRegex = regexp.MustCompile(`^[0-9]{3}$`)

func ValidateDate(input string) bool {
	return dateYYMMDDTypeRegex.MatchString(input)
}

// ValidateTime reports whether input is a BAI2 military time.
// Spec: 0000 through 2400 (0000 = beginning of day, 2400 = end of day).
// Some processors send 9999 for end of day; that is accepted.
func ValidateTime(input string) bool {
	if input == "9999" || input == "2400" {
		return true
	}
	if len(input) != 4 {
		return false
	}
	hh, err1 := strconv.Atoi(input[:2])
	mm, err2 := strconv.Atoi(input[2:])
	if err1 != nil || err2 != nil {
		return false
	}
	if hh == 24 {
		return mm == 0
	}
	return hh >= 0 && hh <= 23 && mm >= 0 && mm <= 59
}

func ValidateFundsType(input string) bool {
	if input == "0" || input == "1" || input == "2" || input == "Z" || input == "V" ||
		input == "S" || input == "D" {
		return true
	}
	return false
}

func ValidateAmount(input string) bool {
	return singedNumber.MatchString(input)
}

func ValidateCurrencyCode(input string) bool {
	return currencyCodeRegex.MatchString(input)
}

func ValidateTypeCode(input string) bool {
	return typeCodeRegex.MatchString(input)
}

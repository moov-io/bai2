// Copyright 2022 The Moov Authors
// Use of this source code is governed by an Apache License
// license that can be found in the LICENSE file.

package util

import "regexp"

var dateYYMMDDTypeRegex = regexp.MustCompile(`[0-9][0-9](0[1-9]|1[0-2])(0[1-9]|1[0-9]|2[0-9]|3[01])`)
var timeTypeRegex = regexp.MustCompile(`[0-9][0-9][0-9][0-9]`)
var singedNumber = regexp.MustCompile(`^(-|\+|)?[0-9]\d*$`)
var currencyCodeRegex = regexp.MustCompile(`^[a-zA-Z]{3}$`)
var typeCodeRegex = regexp.MustCompile(`^[0-9]{3}$`)

func ValidateData(input string) bool {
	return dateYYMMDDTypeRegex.MatchString(input)
}

func ValidateTime(input string) bool {
	return timeTypeRegex.MatchString(input)
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

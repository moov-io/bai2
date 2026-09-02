// Copyright 2022 The Moov Authors
// Use of this source code is governed by an Apache License
// license that can be found in the LICENSE file.

package bai2

import (
	"strconv"
	"strings"
)

func parseAmount(s string) (int64, error) {
	s = strings.TrimSpace(s)
	if s == "" || s == "+" || s == "-" {
		return 0, nil
	}
	return strconv.ParseInt(s, 10, 64)
}

func formatAmount(n int64) string {
	return strconv.FormatInt(n, 10)
}

func amountsEqual(a, b string) bool {
	ai, err1 := parseAmount(a)
	bi, err2 := parseAmount(b)
	if err1 != nil || err2 != nil {
		return a == b
	}
	return ai == bi
}

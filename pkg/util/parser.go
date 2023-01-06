// Copyright 2022 The Moov Authors
// Use of this source code is governed by an Apache License
// license that can be found in the LICENSE file.

package util

import (
	"fmt"
	"strconv"
	"strings"
)

func getIndex(input string) int {
	idx1 := strings.Index(input, ",")
	idx2 := strings.Index(input, "/")

	if idx1 == -1 {
		return idx2
	}
	return idx1
}

func ReadField(input string) (string, int, error) {

	idx := getIndex(input)
	if idx == -1 {
		return "", 0, fmt.Errorf("doesn't have valid delimiter")
	}

	return input[:idx], idx + 1, nil
}

func ReadFieldAsInt(input string) (int64, int, error) {

	idx := getIndex(input)
	if idx == -1 {
		return 0, 0, fmt.Errorf("doesn't have valid delimiter")
	}

	value, _ := strconv.ParseInt(input[:idx], 10, 64)
	return value, idx + 1, nil
}

func GetSize(line string) int64 {

	size := strings.Index(line, "/")
	if size > 0 {
		return int64(size + 1)
	}

	return int64(size)
}

// Copyright 2022 The Moov Authors
// Use of this source code is governed by an Apache License
// license that can be found in the LICENSE file.

package util

import (
	"fmt"
	"math"
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

func ReadField(input string, start int) (string, int, error) {

	data := ""

	if start < len(input) {
		data = input[start:]
	}

	if data == "" {
		return "", 0, fmt.Errorf("doesn't enough input string")
	}

	idx := getIndex(data)
	if idx == -1 {
		return "", 0, fmt.Errorf("doesn't have valid delimiter")
	}

	return data[:idx], idx + 1, nil
}

func ReadFieldAsInt(input string, start int) (int64, int, error) {

	data := ""

	if start < len(input) {
		data = input[start:]
	}

	if data == "" {
		return 0, 0, fmt.Errorf("doesn't enough input string")
	}

	idx := getIndex(data)
	if idx == -1 {
		return 0, 0, fmt.Errorf("doesn't have valid delimiter")
	}

	if data[:idx] == "" {
		return 0, 1, nil
	}

	value, err := strconv.ParseInt(data[:idx], 10, 64)
	if err != nil {
		return 0, 0, fmt.Errorf("doesn't have valid value")
	}

	// check for lower and uppper bounds
	if value > 0 && value <= math.MaxInt64 {
		return value, idx + 1, nil
	}

	return 0, idx + 1, nil
}

func GetSize(line string) int64 {

	size := strings.Index(line, "/")
	if size >= 0 {
		return int64(size + 1)
	}

	return int64(size)
}

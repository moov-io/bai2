// Copyright 2022 The Moov Authors
// Use of this source code is governed by an Apache License
// license that can be found in the LICENSE file.

package util

import (
	"fmt"
	"strconv"
	"strings"
)

func EntryParser(entry, delimiter string) (string, error) {

	if len(entry) < 1 {
		return "", fmt.Errorf("invalid length")
	}

	size := len(entry) - 1
	if entry[size:] != delimiter {
		return "", fmt.Errorf("contains invalid delimiter")
	}

	return entry[:size], nil
}

func EntryParserToInt(entry, delimiter string) (int64, error) {

	if len(entry) < 1 {
		return 0, fmt.Errorf("invalid length")
	}

	size := len(entry) - 1
	if entry[size:] != delimiter {
		return 0, fmt.Errorf("contains invalid delimiter")
	}

	value, _ := strconv.ParseInt(entry[:size], 10, 64)

	return value, nil
}

func GetSize(line string) int64 {

	size := strings.Index(line, "/")
	if size > 0 {
		return int64(size + 1)
	}

	return int64(size)
}

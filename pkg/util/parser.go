// Copyright 2022 The Moov Authors
// Use of this source code is governed by an Apache License
// license that can be found in the LICENSE file.

package util

import (
	"fmt"
	"strconv"
	"strings"
)

func getIndex(input string, opts ...bool) int {
	index_comma := strings.Index(input, ",")
	index_slash := strings.Index(input, "/")
	index_newline := strings.Index(input, "\n")
	// NB. `opts[0]`: if true, returns the index of the last character of the line. this will lead the parser to read the
	// remainder of the line. note, if line is terminated with a `/` character, we will read up to that index insted of
	// reading `len(input)`
	read_remainder_of_line := len(opts) > 0 && opts[0]

	// If there is no `,` separator in the input, return either the index of the next explicit terminating character (`/`)
	// or the index of the next newline character, if no terminating character is present.
	if read_remainder_of_line || index_comma == -1 {
		if index_slash != -1 {
			return index_slash
		}
		if index_newline != -1 {
			return index_newline
		}
		return len(input)
	}

	// If a line is terminated with a `/` character and the terminator is BEFORE the next `,` character, return
	// the index of the `/` character.
	if index_slash > -1 && index_slash < index_comma {
		return index_slash
	}

	// If a line is terminated with a `\n` character (and is NOT terminated with a / character) and the terminator is
	// BEFORE the next `,` character, return the index of the `\n` character.
	if index_slash < 0 && index_newline > -1 && index_newline < index_comma {
		return index_newline
	}

	// Otherwise, return the index of the next `,` character. Value will not be `-1` due to earlier function logic.
	return index_comma
}

func ReadField(input string, start int, opts ...bool) (string, int, error) {

	data := ""

	if start < len(input) {
		data = input[start:]
	}

	if data == "" {
		return "", 0, fmt.Errorf("doesn't enough input string")
	}

	idx := getIndex(data, opts...)
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

	return value, idx + 1, nil
}

func GetSize(line string) int64 {

	size := strings.Index(line, "/")
	if size >= 0 {
		return int64(size + 1)
	}

	size = strings.Index(line, "\n")
	if size >= 0 {
		return int64(size + 1)
	}

	return int64(len(line))
}

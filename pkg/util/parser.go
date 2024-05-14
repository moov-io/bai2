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

	// NB. `opts[0]`: if true, slash character is allowed in text and does not signify a line termination. this is the
	// case for Transaction Detail and Continuation records.
	// `opts[1]`: if true, returns the index of the first newline, or the full length of the input if no newline exists.
	allow_slash_as_character := len(opts) > 0 && opts[0]
	read_remainder_of_line := len(opts) > 1 && opts[1]

	if read_remainder_of_line {
		if index_newline != -1 {
			return index_newline
		}
		return len(input)
	}

	// If there is no `,` separator in the input, return either the index of the next explicit terminating character (`/`)
	// or the index of the next newline character, if no terminating character is present.
	//
	// If slash is allowed as a non-terminating character, only newlines are respected here.
	if index_comma == -1 {
		if !allow_slash_as_character && index_slash != -1 {
			return index_slash
		}
		if index_newline != -1 {
			return index_newline
		}
		return len(input)
	}

	// If a line is terminated with a `/` character and the terminator is BEFORE the next `,` character, return
	// the index of the `/` character.
	//
	// If slash is allowed as a non-terminating character, this check is skipped.
	if !allow_slash_as_character && index_slash > -1 && index_slash < index_comma {
		return index_slash
	}

	// If a line is terminated with a `\n` character (and is NOT terminated with a `/` character, or if `/` is an allowed character)
	// and the `\n` is BEFORE the next `,` character, return the index of the `\n` character.
	if (index_slash < 0 || allow_slash_as_character) && index_newline > -1 && index_newline < index_comma {
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

func GetSize(line string, opts ...bool) int64 {
	allow_slash_as_character := len(opts) > 0 && opts[0]
	read_remainder_of_line := len(opts) > 1 && opts[1]
	if read_remainder_of_line {
		return int64(len(line))
	}

	size := strings.Index(line, "/")
	if !allow_slash_as_character && size >= 0 {
		return int64(size + 1)
	}

	size = strings.Index(line, "\n")
	if size >= 0 {
		return int64(size + 1)
	}

	return int64(len(line))
}

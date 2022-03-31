// Copyright 2022 The Moov Authors
// Use of this source code is governed by an Apache License
// license that can be found in the LICENSE file.

package util

import "fmt"

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

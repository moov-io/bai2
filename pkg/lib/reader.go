// Copyright 2022 The Moov Authors
// Use of this source code is governed by an Apache License
// license that can be found in the LICENSE file.

package lib

import (
	"bufio"
	"io"

	"github.com/moov-io/bai2/pkg/util"
)

type Bai2Scanner struct {
	scan *bufio.Scanner
}

func NewBai2Scanner(fd io.Reader) Bai2Scanner {
	scan := bufio.NewScanner(fd)
	scan.Split(scanRecord)
	return Bai2Scanner{scan: scan}
}

func (b *Bai2Scanner) ScanLine(line string) string {

	if len(line) > 0 {
		return line
	}

	if !b.scan.Scan() {
		return ""
	}

	return b.scan.Text()
}

// scanRecord allows Reader to read each segment
func scanRecord(data []byte, atEOF bool) (advance int, token []byte, err error) {

	if atEOF && len(data) == 0 {
		return 0, nil, nil
	}

	index := util.GetSize(string(data))
	if index < 1 || !atEOF {
		// need more data
		return 0, nil, nil
	}

	return int(index), data[:int(index)], nil
}

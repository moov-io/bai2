// Copyright 2022 The Moov Authors
// Use of this source code is governed by an Apache License
// license that can be found in the LICENSE file.

package lib

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"log"
	"strings"
	"unicode"

	"github.com/moov-io/bai2/pkg/util"
)

type Bai2Scanner struct {
	reader      *bufio.Reader
	currentLine *bytes.Buffer
	index       int
}

func NewBai2Scanner(fd io.Reader) Bai2Scanner {
	reader := bufio.NewReader(fd)
	currentLine := new(bytes.Buffer)
	return Bai2Scanner{reader: reader, currentLine: currentLine}
}

func (b *Bai2Scanner) GetLineIndex() int {
	return b.index
}

func (b *Bai2Scanner) GetLine() string {
	return strings.TrimSpace(b.currentLine.String())
}

// ScanLine returns a line from the underlying reader
// arg[0]: useCurrentLine (if false read a new line)
func (b *Bai2Scanner) ScanLine(arg ...bool) string {

	useCurrentLine := false
	if len(arg) > 0 {
		useCurrentLine = arg[0]
	}

	if useCurrentLine {
		return b.GetLine()
	}

	// Reset the read buffer every time we read a new line.
	b.currentLine.Reset()

	for {
		// Read each rune in the file until a newline or a `/` or EOF.
		rune, _, err := b.reader.ReadRune()
		if err != nil {
			if err != io.EOF {
				log.Fatal(err)
			}
			break
		}

		char := string(rune)
		switch char {
		case "/":
			// Add `/` to line if it exists. Parsers use this to help internally represent the delineation
			// between records. On observing a `/` character, check to see if we have a full record available
			// for processing.
			b.currentLine.WriteString(char)
			goto fullLine
		case "\n", "\r":
			// On observing a newline character, check to see if we have a full record available for processing.
			goto fullLine
		default:
			b.currentLine.WriteString(char)
		}

		continue

		// This routine processes a "full line". In the context of a BAI2 file, a line is a single record
		// and may be terminated either by a `/` or a newline character. In specific circumstances, a logical record
		// ("line") may continue onto the next line, and in that event, processing should read the contents of
		// the following line before considering the record "complete".
	fullLine:
		line := strings.TrimSpace(b.currentLine.String())
		// If the current line has only white space, ignore it and continue reading.
		if blankLine(line) {
			b.currentLine.Reset()
			continue
		}

		// If the line ends with a `/` delimiter, treat it as a complete record and process it as is.
		if strings.HasSuffix(line, "/") {
			break
		}

		// If a line ends with a newline character, look ahead to the next three bytes. If the next line
		// is a new record, it will have a defined and valid record code. If a valid record code is not
		// observed, continue parsing lines until a distinct record is observed.
		bytes, err := b.reader.Peek(3)
		if err != nil && err != io.EOF {
			log.Fatal(err)
		}

		// If the next three bytes are any of the defined BAI2 record codes (followed by a comma), we consider the next line
		// as a new record and process the current line up to this point.
		nextThreeBytes := string(bytes)
		headerCodes := []string{util.FileHeaderCode, util.GroupHeaderCode, util.AccountIdentifierCode, util.TransactionDetailCode, util.ContinuationCode, util.AccountTrailerCode, util.GroupTrailerCode, util.FileTrailerCode}
		nextLineHasNewRecord := false
		for _, header := range headerCodes {
			if nextThreeBytes == fmt.Sprintf("%s,", header) {
				b.currentLine.WriteString("/")
				nextLineHasNewRecord = true
				break
			}
		}

		if nextLineHasNewRecord {
			break
		}

		// Here, the current line "continued" onto the next line without a delimiter and without a new record code on
		// the subsequent line. Parse the next line as though it is a continuation of the current line.
		continue
	}

	b.index++
	return b.GetLine()
}

func blankLine(line string) bool {
	for _, r := range line {
		if !unicode.IsSpace(r) {
			return false
		}
	}
	return true
}

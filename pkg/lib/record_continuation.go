// Copyright 2022 The Moov Authors
// Use of this source code is governed by an Apache License
// license that can be found in the LICENSE file.

package lib

import (
	"bytes"
	"fmt"

	"github.com/moov-io/bai2/pkg/util"
)

const (
	crParseErrorFmt = "Continuation: unable to parse %s"
)

type continuationRecord struct {
	Composite []string `json:",omitempty"`
}

func (h *continuationRecord) validate() error {
	return nil
}

func (h *continuationRecord) parse(data string) (int, error) {
	var line string
	var err error
	var size, read int

	length := util.GetSize(data)
	if length < 3 {
		return 0, fmt.Errorf(fmt.Sprintf(crParseErrorFmt, "record"))
	} else {
		line = data[:length]
	}

	// RecordCode
	if util.ContinuationCode != data[:2] {
		return 0, fmt.Errorf(fmt.Sprintf(fhParseErrorFmt, "RecordCode"))
	}
	read += 3

	for int64(read) < length {
		var composite string
		if composite, size, err = util.ReadField(line, read); err != nil {
			return 0, fmt.Errorf(fmt.Sprintf(crParseErrorFmt, "Composite"))
		} else {
			read += size
		}
		h.Composite = append(h.Composite, composite)
	}

	if err = h.validate(); err != nil {
		return 0, err
	}

	return read, nil
}

func (h *continuationRecord) string() string {
	var buf bytes.Buffer

	buf.WriteString(util.ContinuationCode)
	for _, composite := range h.Composite {
		buf.WriteString(fmt.Sprintf(",%s", composite))
	}
	buf.WriteString("/")
	return buf.String()
}

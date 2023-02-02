// Copyright 2022 The Moov Authors
// Use of this source code is governed by an Apache License
// license that can be found in the LICENSE file.

package util

import (
	"bytes"
	"fmt"
	"strings"
)

// WriteBuffer
//
//	Input type (ELM1,EML2,ELM3..,ELMN)
func WriteBuffer(total, buf *bytes.Buffer, input string, maxLen int64) {

	if maxLen > 0 {

		elements := strings.Split(input, ",")
		newInput := ""

		for _, elm := range elements {

			newSize := int64(buf.Len() + len(newInput) + len(elm) + 2)
			if newSize > maxLen {
				if newInput == "" {
					org := buf.String()
					org = org[:len(org)-1] + "/" + "\n"
					total.WriteString(org)
				} else {
					buf.WriteString(newInput + "/" + "\n") // added new line
					total.WriteString(buf.String())
				}

				// refresh buf
				buf.Reset()

				buf.WriteString(fmt.Sprintf("%s,", ContinuationCode))
				newInput = elm
			} else {
				if newInput == "" {
					newInput = elm
				} else {
					newInput = newInput + "," + elm
				}
			}
		}

		if len(newInput) > 0 {
			buf.WriteString(newInput)
		}

	} else {
		buf.WriteString(input)
	}
}

// Copyright 2022 The Moov Authors
// Use of this source code is governed by an Apache License
// license that can be found in the LICENSE file.

package file

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"

	"github.com/moov-io/bai2/pkg/lib"
	"github.com/moov-io/bai2/pkg/record"
	"github.com/moov-io/bai2/pkg/util"
)

/*

The records in a balance reporting transmission file are ordered as follows:

	----------------------------------------------------------------
	Record Code | Record Name  						| Purpose
	----------------------------------------------------------------
			 01 | File Header  						| Begins File
			 02 | Group Header 						| Begins Group
			 03 | Account Identifier 				| Begins Account
			 16 | Transaction Detail (Optional) 	| Within Account
			 49 | Account Trailer 					| Ends Account
			 98 | Group Trailer 					| Ends Group
			 99 | File Trailer 						|Ends File
	----------------------------------------------------------------

*/

// Creating new file object
func NewBai2() Bai2 {
	return Bai2{}
}

// FILE with BAI Format
type Bai2 struct {
	Header  *lib.FileHeader
	Groups  []*Group
	Trailer *lib.FileTrailer
}

func (r *Bai2) String() string {
	var buf bytes.Buffer

	if r.Header != nil {
		buf.WriteString(r.Header.String() + "\n")
	}

	for i := range r.Groups {
		buf.WriteString(r.Groups[i].String())
	}

	if r.Trailer != nil {
		buf.WriteString(r.Trailer.String())
	}

	return buf.String()
}

func (r *Bai2) Validate() error {

	if r.Header != nil {
		if err := r.Header.Validate(); err != nil {
			return err
		}
	}

	for i := range r.Groups {
		if err := r.Groups[i].Validate(); err != nil {
			return err
		}
	}

	if r.Trailer != nil {
		if err := r.Trailer.Validate(); err != nil {
			return err
		}
	}

	return nil
}

func NewGroup() *Group {
	return &Group{}
}

type Group struct {
	Header  *lib.GroupHeader
	Details []record.Record
	Trailer *lib.GroupTrailer
}

func (r *Group) String() string {
	var buf bytes.Buffer

	if r.Header != nil {
		buf.WriteString(r.Header.String() + "\n")
	}

	for i := range r.Details {
		buf.WriteString(r.Details[i].String() + "\n")
	}

	if r.Trailer != nil {
		buf.WriteString(r.Trailer.String() + "\n")
	}

	return buf.String()
}

func (r *Group) Validate() error {

	if r.Header != nil {
		if err := r.Header.Validate(); err != nil {
			return err
		}
	}

	for i := range r.Details {
		if err := r.Details[i].Validate(); err != nil {
			return err
		}
	}

	if r.Trailer != nil {
		if err := r.Trailer.Validate(); err != nil {
			return err
		}
	}

	return nil
}

// Parse will return file object after parse
func Parse(fd io.Reader) (*Bai2, error) {
	file := NewBai2()

	var lineNum int
	var group *Group
	var hasBlock bool

	scan := bufio.NewScanner(fd)
	for scan.Scan() {
		line := scan.Text()
		lineNum++

		if len(line) < 2 {
			continue
		}

		size := util.GetSize(line)
		if size < 0 {
			continue
		}

		switch line[0:2] {
		case "01":

			newRecord := lib.NewFileHeader()
			err := newRecord.Parse(line)
			if err != nil {
				return &file, fmt.Errorf("ERROR parsing file header on line %d - %v", lineNum, err)
			}

			file.Header = newRecord

		case "99":

			newRecord := lib.NewFileTrailer()
			err := newRecord.Parse(line)
			if err != nil {
				return &file, fmt.Errorf("ERROR parsing file trailer on line %d - %v", lineNum, err)
			}

			file.Trailer = newRecord

		case "02":

			// init group
			group = NewGroup()

			newRecord := lib.NewGroupHeader()
			err := newRecord.Parse(line)
			if err != nil {
				return &file, fmt.Errorf("ERROR parsing file header on line %d - %v", lineNum, err)
			}

			group.Header = newRecord

		case "98":

			newRecord := lib.NewGroupTrailer()
			err := newRecord.Parse(line)
			if err != nil {
				return &file, fmt.Errorf("ERROR parsing file trailer on line %d - %v", lineNum, err)
			}

			group.Trailer = newRecord

			// append group
			file.Groups = append(file.Groups, group)

		case "03":

			if size == 59 {
				newRecord := lib.NewAccountIdentifierCurrent()
				err := newRecord.Parse(line)
				if err != nil {
					return &file, fmt.Errorf("ERROR parsing account indentifier on line %d - %v", lineNum, err)
				}

				group.Details = append(group.Details, newRecord)
			} else if size == 79 {
				newRecord := lib.NewAccountIdentifierLoan()
				err := newRecord.Parse(line)
				if err != nil {
					return &file, fmt.Errorf("ERROR parsing account indentifier on line %d - %v", lineNum, err)
				}

				group.Details = append(group.Details, newRecord)
			}

		case "49":

			newRecord := lib.NewAccountTrailer()
			err := newRecord.Parse(line)
			if err != nil {
				return &file, fmt.Errorf("ERROR parsing account trailer on line %d - %v", lineNum, err)
			}

			group.Details = append(group.Details, newRecord)

		case "16":

			if size == 56 {
				newRecord := lib.NewAccountTransaction()
				err := newRecord.Parse(line)
				if err != nil {
					return &file, fmt.Errorf("ERROR parsing account transaction detail on line %d - %v", lineNum, err)
				}

				group.Details = append(group.Details, newRecord)
			} else if size == 68 {
				newRecord := lib.NewTransferTransaction()
				err := newRecord.Parse(line)
				if err != nil {
					return &file, fmt.Errorf("ERROR parsing account transaction detail on line %d - %v", lineNum, err)
				}

				group.Details = append(group.Details, newRecord)
			}

		case "88":

			newRecord := lib.NewAccountIdentifierContinuation()
			err := newRecord.Parse(line)
			if err != nil {
				return &file, fmt.Errorf("ERROR parsing continuation of account summary record on line %d - %v", lineNum, err)
			}

			group.Details = append(group.Details, newRecord)

		default:
			continue

		}

		hasBlock = true

	}

	if !hasBlock {
		return nil, errors.New("invalid file format")
	}

	return &file, nil
}

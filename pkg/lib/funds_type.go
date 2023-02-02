// Copyright 2022 The Moov Authors
// Use of this source code is governed by an Apache License
// license that can be found in the LICENSE file.

package lib

import (
	"bytes"
	"errors"
	"fmt"
	"strings"

	"github.com/moov-io/bai2/pkg/util"
)

const (
	FundsType0 = "0"
	FundsType1 = "1"
	FundsType2 = "2"
	FundsTypeS = "S"
	FundsTypeV = "V"
	FundsTypeD = "D"
	FundsTypeZ = "Z"
)

type FundsType struct {
	TypeCode FundsTypeCode `json:"type_code,omitempty"`

	// Type 0,1,2,S
	ImmediateAmount int64 `json:"immediate_amount,omitempty"` // availability amount
	OneDayAmount    int64 `json:"one_day_amount,omitempty"`   // one-day availability amount
	TwoDayAmount    int64 `json:"two_day_amount,omitempty"`   // more than one-day availability amount

	// Type V
	Date string `json:"date,omitempty"`
	Time string `json:"time,omitempty"`

	// Type D
	DistributionNumber int64          `json:"distribution_number,omitempty"`
	Distributions      []Distribution `json:"distributions,omitempty"`
}

func (f *FundsType) Validate() error {

	if err := f.TypeCode.Validate(); err != nil {
		return err
	}

	if strings.ToUpper(string(f.TypeCode)) == FundsTypeD && int(f.DistributionNumber) != len(f.Distributions) {
		return errors.New("number of distributions is not match")
	}

	if strings.ToUpper(string(f.TypeCode)) == FundsTypeV {
		if f.Date != "" && !util.ValidateData(f.Date) {
			return errors.New("invalid date of fund type V (" + f.Date + ")")
		}
		if f.Time != "" && !util.ValidateTime(f.Time) {
			return errors.New("invalid time of fund type V (" + f.Time + ")")
		}
	}

	return nil
}

func (f *FundsType) String() string {

	fType := strings.ToUpper(string(f.TypeCode))

	var buf bytes.Buffer
	if f.TypeCode == "" || fType == FundsTypeZ {
		buf.WriteString(strings.ToUpper(string(f.TypeCode)))
	} else {

		if fType == FundsType0 || fType == FundsType1 || fType == FundsType2 {
			buf.WriteString(strings.ToUpper(string(f.TypeCode)))
		} else if fType == FundsTypeS {
			buf.WriteString(strings.ToUpper(string(f.TypeCode)) + ",")
			buf.WriteString(fmt.Sprintf("%d,%d,%d", f.ImmediateAmount, f.OneDayAmount, f.TwoDayAmount))
		} else if fType == FundsTypeV {
			buf.WriteString(strings.ToUpper(string(f.TypeCode)) + ",")
			buf.WriteString(fmt.Sprintf("%s,%s", f.Date, f.Time))
		} else if fType == FundsTypeD {
			if len(f.Distributions) > 0 {
				buf.WriteString(fmt.Sprintf("%s,%d,", strings.ToUpper(string(f.TypeCode)), f.DistributionNumber))
				for index, distribution := range f.Distributions {
					if index < len(f.Distributions)-1 {
						buf.WriteString(fmt.Sprintf("%d,%d,", distribution.Day, distribution.Amount))
					} else {
						buf.WriteString(fmt.Sprintf("%d,%d", distribution.Day, distribution.Amount))
					}
				}
			} else {
				buf.WriteString(strings.ToUpper(string(f.TypeCode)) + ",0")
			}
		}
	}

	return buf.String()
}

func (f *FundsType) parse(data string) (int, error) {

	var err error
	var size, read int

	code, size, err := util.ReadField(data, read)
	if err != nil {
		return 0, errors.New("FundsType: unable to parse type code")
	} else {
		read += size
	}

	f.TypeCode = FundsTypeCode(code)

	if f.TypeCode == FundsTypeS {

		f.ImmediateAmount, size, err = util.ReadFieldAsInt(data, read)
		if err != nil {
			return 0, errors.New("FundsType: unable to parse amount")
		} else {
			read += size
		}

		f.OneDayAmount, size, err = util.ReadFieldAsInt(data, read)
		if err != nil {
			return 0, errors.New("FundsType: unable to parse amount")
		} else {
			read += size
		}

		f.TwoDayAmount, size, err = util.ReadFieldAsInt(data, read)
		if err != nil {
			return 0, errors.New("FundsType: unable to parse amount")
		} else {
			read += size
		}

	} else if f.TypeCode == FundsTypeV {
		f.Date, size, err = util.ReadField(data, read)
		if err != nil {
			return 0, errors.New("FundsType: unable to parse date")
		} else {
			read += size
		}

		f.Time, size, err = util.ReadField(data, read)
		if err != nil {
			return 0, errors.New("FundsType: unable to parse time")
		} else {
			read += size
		}
	} else if f.TypeCode == FundsTypeD {
		f.DistributionNumber, size, err = util.ReadFieldAsInt(data, read)
		if err != nil {
			return 0, errors.New("FundsType: unable to parse distribution number")
		} else {
			read += size
		}

		for index := 0; index < int(f.DistributionNumber); index++ {

			var amount, day int64

			day, size, err = util.ReadFieldAsInt(data, read)
			if err != nil {
				return 0, errors.New("FundsType: unable to parse day")
			} else {
				read += size
			}

			amount, size, err = util.ReadFieldAsInt(data, read)
			if err != nil {
				return 0, errors.New("FundsType: unable to parse amount")
			} else {
				read += size
			}

			f.Distributions = append(f.Distributions, Distribution{Day: day, Amount: amount})

		}
	}

	if strings.Contains(data[:read-1], "/") {
		return 0, errors.New("FundsType: unable to parse sub elements")
	}

	if err = f.Validate(); err != nil {
		return 0, err
	}

	return read, nil
}

type Distribution struct {
	Day    int64 `json:"day,omitempty"`    // availability amount
	Amount int64 `json:"amount,omitempty"` // availability amount
}

type FundsTypeCode string

func (c FundsTypeCode) Validate() error {

	str := string(c)
	if len(str) == 0 {
		return nil
	}

	availableTypes := []string{FundsType0, FundsType1, FundsType2, FundsTypeS, FundsTypeV, FundsTypeD, FundsTypeZ}

	for _, t := range availableTypes {
		if strings.ToUpper(str) == t {
			return nil
		}
	}

	return errors.New("invalid fund type")
}

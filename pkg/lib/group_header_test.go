// Copyright 2022 The Moov Authors
// Use of this source code is governed by an Apache License
// license that can be found in the LICENSE file.

package lib

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func mockGroupHeader() *GroupHeader {
	return &GroupHeader{
		RecordCode:       "02",
		Receiver:         "0004",
		Originator:       "12345",
		GroupStatus:      1,
		AsOfDate:         "060321",
		AsOfTime:         "0829",
		CurrencyCode:     "USD",
		AsOfDateModifier: 2,
	}
}

func TestGroupHeader(t *testing.T) {

	record := mockGroupHeader()
	require.NoError(t, record.Validate())

	record.AsOfDateModifier = 5
	require.Error(t, record.Validate())
	require.Equal(t, "GroupHeader: invalid AsOfDateModifier", record.Validate().Error())

	record.CurrencyCode = "A"
	require.Error(t, record.Validate())
	require.Equal(t, "GroupHeader: invalid CurrencyCode", record.Validate().Error())

	record.AsOfTime = "AAA"
	require.Error(t, record.Validate())
	require.Equal(t, "GroupHeader: invalid AsOfTime", record.Validate().Error())

	record.AsOfDate = ""
	require.Error(t, record.Validate())
	require.Equal(t, "GroupHeader: invalid AsOfDate", record.Validate().Error())

	record.GroupStatus = 5
	require.Error(t, record.Validate())
	require.Equal(t, "GroupHeader: invalid GroupStatus", record.Validate().Error())

	record.Originator = ""
	require.Error(t, record.Validate())
	require.Equal(t, "GroupHeader: invalid Originator", record.Validate().Error())

	record.RecordCode = ""
	require.Error(t, record.Validate())
	require.Equal(t, "GroupHeader: invalid RecordCode", record.Validate().Error())

}

func TestGroupHeaderWithOptional(t *testing.T) {

	sample := "02,12345,0004,1,060317,0000,CAD,2/"
	record := NewGroupHeader()

	size, err := record.Parse(sample)
	require.NoError(t, err)
	require.Equal(t, 34, size)

	require.Equal(t, "02", record.RecordCode)
	require.Equal(t, "0004", record.Originator)
	require.Equal(t, "12345", record.Receiver)
	require.Equal(t, "060317", record.AsOfDate)
	require.Equal(t, "0000", record.AsOfTime)
	require.Equal(t, "CAD", record.CurrencyCode)

	require.Equal(t, sample, record.String())
}

func TestGroupHeaderWithoutOptional(t *testing.T) {

	sample := "02,,0004,1,060317,,,/"
	record := NewGroupHeader()

	size, err := record.Parse(sample)
	require.NoError(t, err)
	require.Equal(t, 21, size)

	require.Equal(t, "02", record.RecordCode)
	require.Equal(t, "0004", record.Originator)
	require.Equal(t, "060317", record.AsOfDate)

	require.Equal(t, sample, record.String())
}

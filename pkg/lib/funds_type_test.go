// Copyright 2022 The Moov Authors
// Use of this source code is governed by an Apache License
// license that can be found in the LICENSE file.

package lib

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func mockFuncType() FundsType {
	return FundsType{
		ImmediateAmount:    10000,
		OneDayAmount:       -20000,
		TwoDayAmount:       30000,
		Date:               "040701",
		Time:               "1300",
		DistributionNumber: 5,
		Distributions: []Distribution{
			{
				Day:    1,
				Amount: 100,
			},
			{
				Day:    2,
				Amount: 200,
			},
			{
				Day:    3,
				Amount: 300,
			},
			{
				Day:    4,
				Amount: 400,
			},
			{
				Day:    5,
				Amount: 500,
			},
		},
	}
}

func TestFundsType_String(t *testing.T) {
	f := mockFuncType()

	require.NoError(t, f.Validate())
	require.Equal(t, "", f.String())

	f.TypeCode = FundsType0
	require.NoError(t, f.Validate())
	require.Equal(t, "0", f.String())

	f.TypeCode = FundsType1
	require.NoError(t, f.Validate())
	require.Equal(t, "1", f.String())

	f.TypeCode = FundsType2
	require.NoError(t, f.Validate())
	require.Equal(t, "2", f.String())

	f.TypeCode = FundsTypeS
	require.NoError(t, f.Validate())
	require.Equal(t, "S,10000,-20000,30000", f.String())

	f.TypeCode = FundsTypeV
	require.NoError(t, f.Validate())
	require.Equal(t, "V,040701,1300", f.String())

	f.TypeCode = FundsTypeZ
	require.NoError(t, f.Validate())
	require.Equal(t, "Z", f.String())

	f.TypeCode = FundsTypeD
	require.NoError(t, f.Validate())
	require.Equal(t, "D,5,1,100,2,200,3,300,4,400,5,500", f.String())

	f.DistributionNumber = 0
	require.Error(t, f.Validate())
	require.Equal(t, "number of distributions is not match", f.Validate().Error())

	f.Distributions = nil
	require.NoError(t, f.Validate())
	require.Equal(t, "D,0", f.String())
}

func TestFundsType_Parse(t *testing.T) {

	type testsample struct {
		Input         string
		IsReadErr     bool
		IsValidateErr bool
		CodeType      string
		ReadSize      int
		Output        string
	}

	samples := []testsample{
		{
			Input:         ",",
			IsReadErr:     false,
			IsValidateErr: false,
			CodeType:      "",
			ReadSize:      1,
			Output:        "",
		},
		{
			Input:         "Z,",
			IsReadErr:     false,
			IsValidateErr: false,
			CodeType:      FundsTypeZ,
			ReadSize:      2,
			Output:        "Z",
		},
		{
			Input:         "0,10000,",
			IsReadErr:     false,
			IsValidateErr: false,
			CodeType:      FundsType0,
			ReadSize:      2,
			Output:        "0",
		},
		{
			Input:         "1,10000,",
			IsReadErr:     false,
			IsValidateErr: false,
			CodeType:      FundsType1,
			ReadSize:      2,
			Output:        "1",
		},
		{
			Input:         "2,10000,",
			IsReadErr:     false,
			IsValidateErr: false,
			CodeType:      FundsType2,
			ReadSize:      2,
			Output:        "2",
		},
		{
			Input:         "S,10000,-20000,30000,",
			IsReadErr:     false,
			IsValidateErr: false,
			CodeType:      FundsTypeS,
			ReadSize:      21,
			Output:        "S,10000,-20000,30000",
		},
		{
			Input:         "V,040701,1300,",
			IsReadErr:     false,
			IsValidateErr: false,
			CodeType:      FundsTypeV,
			ReadSize:      14,
			Output:        "V,040701,1300",
		},
		{
			Input:         "D,5,1,10000,2,10000,3,10000,4,10000,5,10000,",
			IsReadErr:     false,
			IsValidateErr: false,
			CodeType:      FundsTypeD,
			ReadSize:      44,
			Output:        "D,5,1,10000,2,10000,3,10000,4,10000,5,10000",
		},
		{
			Input:         "D,0,1,10000,2,10000,3,10000,4,10000,5,10000,",
			IsReadErr:     false,
			IsValidateErr: false,
			CodeType:      FundsTypeD,
			ReadSize:      4,
			Output:        "D,0",
		},
		{
			Input:         "D/0,1,10000,2,10000,3,10000,4,10000,5,10000,",
			IsReadErr:     true,
			IsValidateErr: false,
			CodeType:      FundsTypeD,
			ReadSize:      0,
			Output:        "D,0",
		},
	}

	for _, sample := range samples {
		f := FundsType{}

		size, err := f.parse(sample.Input)

		if sample.IsReadErr {
			require.Error(t, err)
		} else {
			require.NoError(t, err)
		}

		if sample.IsValidateErr {
			require.Error(t, f.Validate())
		} else {
			require.NoError(t, f.Validate())
		}

		require.Equal(t, sample.CodeType, string(f.TypeCode))
		require.Equal(t, sample.ReadSize, size)
		require.Equal(t, sample.Output, f.String())
	}

}

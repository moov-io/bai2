// Copyright 2022 The Moov Authors
// Use of this source code is governed by an Apache License
// license that can be found in the LICENSE file.

package util

import (
	"testing"

	"github.com/stretchr/testify/require"
)

type input struct {
	Data  string
	Start int
}

type want struct {
	Error    bool
	ErrorMsg string
	Read     int
	Value    string
	IntValue int64
}

type testSample struct {
	Input input
	Want  want
}

func TestReadField(t *testing.T) {

	samples := []testSample{
		{
			Input: input{
				Data:  "01,",
				Start: 0,
			},
			Want: want{
				Error: false,
				Read:  3,
				Value: "01",
			},
		},
		{
			Input: input{
				Data:  "01/",
				Start: 0,
			},
			Want: want{
				Error: false,
				Read:  3,
				Value: "01",
			},
		},
		{
			Input: input{
				Data:  "ODFI’,",
				Start: 0,
			},
			Want: want{
				Error: false,
				Read:  8,
				Value: "ODFI’",
			},
		},
		{
			Input: input{
				Data:  "ODFI’,",
				Start: 6,
			},
			Want: want{
				Error: false,
				Read:  2,
				Value: "\x99",
			},
		},
		{
			Input: input{
				Data:  "ODFI’,",
				Start: 7,
			},
			Want: want{
				Error: false,
				Read:  1,
				Value: "",
			},
		},
		{
			Input: input{
				Data:  "ODFI’,",
				Start: 8,
			},
			Want: want{
				Error:    true,
				ErrorMsg: "doesn't enough input string",
				Read:     0,
				Value:    "",
			},
		},
		{
			Input: input{
				Data:  "ODFI’,",
				Start: 10,
			},
			Want: want{
				Error:    true,
				ErrorMsg: "doesn't enough input string",
				Read:     0,
				Value:    "",
			},
		},
	}

	for _, sample := range samples {
		value, size, err := ReadField(sample.Input.Data, sample.Input.Start)
		if !sample.Want.Error {
			require.NoError(t, err)
		} else {

			require.Error(t, err)
			require.Equal(t, sample.Want.ErrorMsg, err.Error())
		}
		require.Equal(t, sample.Want.Read, size)
		require.Equal(t, sample.Want.Value, value)
	}
}

func TestReadFieldAsInt(t *testing.T) {

	samples := []testSample{
		{
			Input: input{
				Data:  "11,",
				Start: 0,
			},
			Want: want{
				Error:    false,
				Read:     3,
				IntValue: 11,
			},
		},
		{
			Input: input{
				Data:  "01/",
				Start: 0,
			},
			Want: want{
				Error:    false,
				Read:     3,
				IntValue: 1,
			},
		},
		{
			Input: input{
				Data:  "ODFI’,",
				Start: 0,
			},
			Want: want{
				Error:    true,
				ErrorMsg: "doesn't have valid value",
				Read:     0,
				IntValue: 0,
			},
		},
		{
			Input: input{
				Data:  "ODFI’,",
				Start: 6,
			},
			Want: want{
				Error:    true,
				ErrorMsg: "doesn't have valid value",
				Read:     0,
				IntValue: 0,
			},
		},
		{
			Input: input{
				Data:  "ODFI’,",
				Start: 7,
			},
			Want: want{
				Error:    true,
				ErrorMsg: "doesn't have valid value",
				Read:     0,
				IntValue: 0,
			},
		},
		{
			Input: input{
				Data:  "ODFI’,",
				Start: 8,
			},
			Want: want{
				Error:    true,
				ErrorMsg: "doesn't enough input string",
				Read:     0,
				IntValue: 0,
			},
		},
		{
			Input: input{
				Data:  "ODFI’,",
				Start: 10,
			},
			Want: want{
				Error:    true,
				ErrorMsg: "doesn't enough input string",
				Read:     0,
				IntValue: 0,
			},
		},
	}

	for _, sample := range samples {
		value, size, err := ReadFieldAsInt(sample.Input.Data, sample.Input.Start)
		if !sample.Want.Error {
			require.NoError(t, err)
		} else {

			require.Error(t, err)
			require.Equal(t, sample.Want.ErrorMsg, err.Error())
		}
		require.Equal(t, sample.Want.Read, size)
		require.Equal(t, sample.Want.IntValue, value)
	}
}

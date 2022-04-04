// Copyright 2022 The Moov Authors
// Use of this source code is governed by an Apache License
// license that can be found in the LICENSE file.

package lib

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGroupHeader(t *testing.T) {

	header := NewGroupHeader()
	require.NoError(t, header.Validate())

	header.CurrencyCode = ""
	require.Error(t, header.Validate())
	require.Equal(t, "GroupHeader: invalid currency code", header.Validate().Error())

	header.Sender = ""
	require.Error(t, header.Validate())
	require.Equal(t, "GroupHeader: invalid sender", header.Validate().Error())

	header.RecordCode = ""
	require.Error(t, header.Validate())
	require.Equal(t, "GroupHeader: invalid record code", header.Validate().Error())

}

func TestGroupHeaderWithSample(t *testing.T) {

	sample := "02,12345,0004,1,060317,,CAD,/"
	header := NewGroupHeader()

	err := header.Parse(sample)
	require.NoError(t, err)

	require.Equal(t, "02", header.RecordCode)
	require.Equal(t, "0004", header.Sender)
	require.Equal(t, "12345", header.Receiver)
	require.Equal(t, "060317", header.AsOfDate)
	require.Equal(t, "CAD", header.CurrencyCode)

	require.Equal(t, sample, header.String())

	header = &GroupHeader{}
	require.Error(t, header.Validate())

	err = header.Parse(sample[:20])
	require.Equal(t, "GroupHeader: length 20 is too short", err.Error())
}

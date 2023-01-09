// Copyright 2022 The Moov Authors
// Use of this source code is governed by an Apache License
// license that can be found in the LICENSE file.

package lib

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestContinuationRecord(t *testing.T) {

	record := NewContinuationRecord()
	require.NoError(t, record.Validate())

	record.RecordCode = ""
	require.Error(t, record.Validate())
	require.Equal(t, "ContinuationRecord: invalid RecordCode", record.Validate().Error())

}

func TestContinuationRecordWithSample(t *testing.T) {

	sample := "88,100,000000000208500,00003,V,060316,,400,000000000208500,00008,V,060316,/"
	record := NewContinuationRecord()

	size, err := record.Parse(sample)
	require.NoError(t, err)
	require.Equal(t, 75, size)

	require.Equal(t, "88", record.RecordCode)
	require.Equal(t, 12, len(record.Composite))
	require.Equal(t, "100", record.Composite[0])
	require.Equal(t, "000000000208500", record.Composite[1])

}

func TestContinuationRecordWithSample2(t *testing.T) {

	sample := "88,100,000000000208500"
	record := NewAccountIdentifier()

	size, err := record.Parse(sample)
	require.Equal(t, "AccountIdentifier: unable to parse record", err.Error())
	require.Equal(t, 0, size)

	sample = "88,100,000000000208500,/"
	size, err = record.Parse(sample)
	require.Equal(t, "AccountIdentifier: unable to parse Amount", err.Error())
	require.Equal(t, 0, size)

}

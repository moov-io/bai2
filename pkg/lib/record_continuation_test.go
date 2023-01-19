// Copyright 2022 The Moov Authors
// Use of this source code is governed by an Apache License
// license that can be found in the LICENSE file.

package lib

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestContinuationRecord(t *testing.T) {

	record := continuationRecord{}
	require.NoError(t, record.validate())

}

func TestContinuationRecordWithSample(t *testing.T) {

	sample := "88,100,000000000208500,00003,V,060316,,400,000000000208500,00008,V,060316,/"
	record := continuationRecord{}

	size, err := record.parse(sample)
	require.NoError(t, err)
	require.Equal(t, 75, size)

	require.Equal(t, 12, len(record.Composite))
	require.Equal(t, "100", record.Composite[0])
	require.Equal(t, "000000000208500", record.Composite[1])

}

func TestContinuationRecordWithSample2(t *testing.T) {

	sample := "88,100,000000000208500"
	record := continuationRecord{}

	size, err := record.parse(sample)
	require.Equal(t, "Continuation: unable to parse record", err.Error())
	require.Equal(t, 0, size)
	require.Equal(t, "88/", record.string())

}

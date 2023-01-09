// Copyright 2022 The Moov Authors
// Use of this source code is governed by an Apache License
// license that can be found in the LICENSE file.

package util

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestReadField(t *testing.T) {

	entry, size, err := ReadField("01,")
	require.NoError(t, err)
	require.Equal(t, 3, size)
	require.Equal(t, "01", entry)

	entry, size, err = ReadField("01/")
	require.NoError(t, err)
	require.Equal(t, 3, size)
	require.Equal(t, "01", entry)

	entry, size, err = ReadField("ODFI’,")
	require.NoError(t, err)
	require.Equal(t, 8, size)
	require.Equal(t, "ODFI’", entry)

}

func TestReadFieldAsInt(t *testing.T) {

	entry, size, err := ReadFieldAsInt("11,")
	require.NoError(t, err)
	require.Equal(t, 3, size)
	require.Equal(t, int64(11), entry)

	entry, size, err = ReadFieldAsInt("11/")
	require.NoError(t, err)
	require.Equal(t, 3, size)
	require.Equal(t, int64(11), entry)

	entry, size, err = ReadFieldAsInt("ODFI’,")
	require.NoError(t, err)
	require.Equal(t, 8, size)
	require.Equal(t, int64(0), entry)

}

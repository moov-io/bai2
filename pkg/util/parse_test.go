// Copyright 2022 The Moov Authors
// Use of this source code is governed by an Apache License
// license that can be found in the LICENSE file.

package util

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestEntryParser(t *testing.T) {

	entry, err := EntryParser("01,", ",")
	require.NoError(t, err)
	require.Equal(t, "01", entry)

	entry, err = EntryParser("2/", "/")
	require.NoError(t, err)
	require.Equal(t, "2", entry)

}

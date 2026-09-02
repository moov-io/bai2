// Copyright 2022 The Moov Authors
// Use of this source code is governed by an Apache License
// license that can be found in the LICENSE file.

package util

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestValidateDate(t *testing.T) {
	require.True(t, ValidateDate("040621"))
	require.True(t, ValidateDate("060317"))
	require.False(t, ValidateDate("04062"))
	require.False(t, ValidateDate("xx040621"))
	require.False(t, ValidateDate("040621x"))
	require.False(t, ValidateDate("041321"))
	require.False(t, ValidateDate(""))
}

func TestValidateTime(t *testing.T) {
	// Spec section 3: 0000 beginning of day, 2400 end of day, 9999 also end of day.
	require.True(t, ValidateTime("0000"))
	require.True(t, ValidateTime("0200"))
	require.True(t, ValidateTime("0829"))
	require.True(t, ValidateTime("2359"))
	require.True(t, ValidateTime("2400"))
	require.True(t, ValidateTime("9999"))

	require.False(t, ValidateTime(""))
	require.False(t, ValidateTime("829"))
	require.False(t, ValidateTime("082"))
	require.False(t, ValidateTime("25:00"))
	require.False(t, ValidateTime("2500"))
	require.False(t, ValidateTime("2360"))
	require.False(t, ValidateTime("2401"))
	require.False(t, ValidateTime("9998"))
}

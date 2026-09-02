// Copyright 2022 The Moov Authors
// Use of this source code is governed by an Apache License
// license that can be found in the LICENSE file.

package lib

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCompatShimStillParsesBAI2(t *testing.T) {
	raw := "01,122099999,123456789,040621,0200,1,,,2/\n" +
		"02,031001234,122099999,1,040620,2359,,2/\n" +
		"03,0975312468,,010,500000,,,/\n" +
		"49,500000,2/\n98,500000,1,4/\n99,500000,1,6/\n"
	scan := NewBai2Scanner(strings.NewReader(raw))
	f := NewBai2()
	require.NoError(t, f.Read(&scan))
	require.NoError(t, f.Validate())
	require.Equal(t, int64(2), f.Version())
}

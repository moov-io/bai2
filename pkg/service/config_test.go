// Copyright 2022 The Moov Authors
// Use of this source code is governed by an Apache License
// license that can be found in the LICENSE file.

package service_test

import (
	"testing"

	bai2 "github.com/moov-io/bai2"
	"github.com/moov-io/bai2/pkg/service"
	"github.com/moov-io/base/config"
	"github.com/moov-io/base/log"
	"github.com/stretchr/testify/require"
)

func Test_ConfigLoading(t *testing.T) {
	logger := log.NewNopLogger()

	ConfigService := config.NewService(logger)

	gc := &service.GlobalConfig{}
	err := ConfigService.LoadFromFS(gc, bai2.ConfigDefaults)
	require.Nil(t, err)
	require.Equal(t, ":8208", gc.Bai2.Servers.Public.Bind.Address)
}

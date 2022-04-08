// Copyright 2022 The Moov Authors
// Use of this source code is governed by an Apache License
// license that can be found in the LICENSE file.

package service_test

import (
	"os"
	"testing"

	"github.com/go-kit/log"
	baseLog "github.com/moov-io/base/log"

	"github.com/moov-io/bai2/pkg/service"
	"github.com/stretchr/testify/assert"
)

func Test_Environment_Startup(t *testing.T) {
	a := assert.New(t)

	env := &service.Environment{
		Logger: baseLog.NewLogger(log.NewLogfmtLogger(log.NewSyncWriter(os.Stderr))),
	}

	env, err := service.NewEnvironment(env)
	a.Nil(err)

	shutdown := env.RunServers(false)

	env1, err := service.NewEnvironment(&service.Environment{})
	a.Nil(err)
	env1.Shutdown()

	t.Cleanup(shutdown)
}

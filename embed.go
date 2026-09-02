// Copyright 2022 The Moov Authors
// Use of this source code is governed by an Apache License
// license that can be found in the LICENSE file.

package bai2

import "embed"

// ConfigDefaults is the packaged default YAML loaded when APP_CONFIG is unset.
//
//go:embed configs/config.default.yml
var ConfigDefaults embed.FS

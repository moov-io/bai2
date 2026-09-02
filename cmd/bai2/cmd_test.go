// Copyright 2022 The Moov Authors
// Use of this source code is governed by an Apache License
// license that can be found in the LICENSE file.

package main

import (
	"bytes"
	"github.com/spf13/cobra"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	testFileName       = filepath.Join("..", "..", "test", "testdata", "sample1.txt")
	parseErrorFileName = filepath.Join("..", "..", "test", "testdata", "errors", "sample-parseError.txt")
)

func TestMain(m *testing.M) {
	initRootCmd()
	os.Exit(m.Run())
}

func executeCommandC(root *cobra.Command, args ...string) (c *cobra.Command, output string, err error) {
	buf := new(bytes.Buffer)
	root.SetOutput(buf)
	root.SetArgs(args)

	c, err = root.ExecuteC()

	return c, buf.String(), err
}

func executeCommand(root *cobra.Command, args ...string) (output string, err error) {
	_, output, err = executeCommandC(root, args...)
	return output, err
}

func TestWebTest(t *testing.T) {
	_, err := executeCommand(rootCmd, "web", "--test=true")
	if err != nil {
		t.Errorf("%s", err.Error())
	}
}

func TestPrint(t *testing.T) {
	_, err := executeCommand(rootCmd, "print", "--input", testFileName)
	if err != nil {
		t.Errorf("%s", err.Error())
	}
}

func TestParse(t *testing.T) {
	_, err := executeCommand(rootCmd, "parse", "--input", testFileName)
	if err != nil {
		t.Errorf("%s", err.Error())
	}
}

func TestFormat(t *testing.T) {
	_, err := executeCommand(rootCmd, "format", "--input", testFileName)
	if err != nil {
		t.Errorf("%s", err.Error())
	}
}

func TestPrint_ParseError(t *testing.T) {
	_, err := executeCommand(rootCmd, "print", "--input", parseErrorFileName)
	assert.Equal(t, err.Error(), `first record is "00,0004,1/", want 01 file header`)
}

func TestParse_ParseError(t *testing.T) {
	_, err := executeCommand(rootCmd, "parse", "--input", parseErrorFileName)
	assert.Equal(t, err.Error(), `first record is "00,0004,1/", want 01 file header`)
}

func TestParse_BAI3(t *testing.T) {
	path := filepath.Join("..", "..", "test", "testdata", "bai3", "x9-mandatory.txt")
	_, err := executeCommand(rootCmd, "parse", "--input", path)
	assert.NoError(t, err)
}

func TestParse_IgnoreVersion(t *testing.T) {
	path := filepath.Join("..", "..", "test", "testdata", "spec-section3.txt")
	data, err := os.ReadFile(path)
	assert.NoError(t, err)
	v3 := bytes.Replace(data, []byte(",2/"), []byte(",3/"), 1)
	tmp, err := os.CreateTemp("", "bai2-v3-*.txt")
	assert.NoError(t, err)
	defer os.Remove(tmp.Name())
	_, err = tmp.Write(v3)
	assert.NoError(t, err)
	assert.NoError(t, tmp.Close())

	defer func() {
		ignoreVersion = false
		strictControlTotals = false
	}()

	_, err = executeCommand(rootCmd, "parse", "--input", tmp.Name())
	assert.Error(t, err)

	_, err = executeCommand(rootCmd, "parse", "--input", tmp.Name(), "--ignoreVersion")
	assert.NoError(t, err)
}

func TestFormat_ParseError(t *testing.T) {
	_, err := executeCommand(rootCmd, "format", "--input", parseErrorFileName)
	assert.Equal(t, err.Error(), `first record is "00,0004,1/", want 01 file header`)
}

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
	testFileName = filepath.Join("..", "..", "test", "testdata", "sample1.txt")
	parseErrorFileName = filepath.Join("..", "..", "test", "testdata", "sample-parseError.txt")
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
		t.Errorf(err.Error())
	}
}

func TestPrint(t *testing.T) {
	_, err := executeCommand(rootCmd, "print", "--input", testFileName)
	if err != nil {
		t.Errorf(err.Error())
	}
}

func TestParse(t *testing.T) {
	_, err := executeCommand(rootCmd, "parse", "--input", testFileName)
	if err != nil {
		t.Errorf(err.Error())
	}
}

func TestFormat(t *testing.T) {
	_, err := executeCommand(rootCmd, "format", "--input", testFileName)
	if err != nil {
		t.Errorf(err.Error())
	}
}

func TestPrint_ParseError(t *testing.T) {
	_, err := executeCommand(rootCmd, "print", "--input", parseErrorFileName)
	assert.Equal(t, err.Error(), "ERROR parsing file on line 1 (unsupported record type 00)")
}

func TestParse_ParseError(t *testing.T) {
	_, err := executeCommand(rootCmd, "parse", "--input", parseErrorFileName)
	assert.Equal(t, err.Error(), "ERROR parsing file on line 1 (unsupported record type 00)")
}

func TestFormat_ParseError(t *testing.T) {
	_, err := executeCommand(rootCmd, "format", "--input", parseErrorFileName)
	assert.Equal(t, err.Error(), "ERROR parsing file on line 1 (unsupported record type 00)")
}

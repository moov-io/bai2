// Copyright 2022 The Moov Authors
// Use of this source code is governed by an Apache License
// license that can be found in the LICENSE file.

package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"

	"github.com/moov-io/bai2/pkg/lib"
	"github.com/moov-io/bai2/pkg/service"
	baseLog "github.com/moov-io/base/log"
)

var (
	documentFileName string
	documentBuffer   []byte
)

var WebCmd = &cobra.Command{
	Use:   "web",
	Short: "Launches web server",
	Long:  "Launches web server",
	RunE: func(cmd *cobra.Command, args []string) error {
		env := &service.Environment{
			Logger: baseLog.NewDefaultLogger(),
		}

		env, err := service.NewEnvironment(env)
		if err != nil {
			env.Logger.Fatal().LogErrorf("Error loading up environment.", err).Err()
			os.Exit(1)
		}
		defer env.Shutdown()

		env.Logger.Info().Log("Starting web service")
		test, _ := cmd.Flags().GetBool("test")
		if !test {
			shutdown := env.RunServers(true)
			defer shutdown()
		}
		return nil
	},
}

var Parse = &cobra.Command{
	Use:   "parse",
	Short: "parse bai2 report",
	Long:  "Parse an incoming bai2 report",
	RunE: func(cmd *cobra.Command, args []string) error {

		var err error

		scan := lib.NewBai2Scanner(bytes.NewReader(documentBuffer))
		f := lib.NewBai2()
		err = f.Read(&scan)
		if err != nil {
			return err
		}

		err = f.Validate()
		if err != nil {
			return errors.New("Parsing report was successful, but not valid")
		}

		log.Println("Parsing report was successful and the report is valid")

		return nil
	},
}

var Print = &cobra.Command{
	Use:   "print",
	Short: "Print bai2 report",
	Long:  "Print an incoming bai2 report after parse",
	RunE: func(cmd *cobra.Command, args []string) error {

		var err error

		scan := lib.NewBai2Scanner(bytes.NewReader(documentBuffer))
		f := lib.NewBai2()
		err = f.Read(&scan)
		if err != nil {
			return err
		}

		err = f.Validate()
		if err != nil {
			return err
		}

		fmt.Println(f.String())
		return nil
	},
}

var Format = &cobra.Command{
	Use:   "format",
	Short: "Format bai2 report",
	Long:  "Format an incoming bai2 report after parse",
	RunE: func(cmd *cobra.Command, args []string) error {

		var err error

		scan := lib.NewBai2Scanner(bytes.NewReader(documentBuffer))
		f := lib.NewBai2()
		err = f.Read(&scan)
		if err != nil {
			return err
		}

		err = f.Validate()
		if err != nil {
			return err
		}

		body, ferr := json.Marshal(f)
		if ferr != nil {
			return ferr
		}

		fmt.Println(string(body))
		return nil
	},
}

var rootCmd = &cobra.Command{
	Use:   "",
	Short: "",
	Long:  "",
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		isWeb := false
		cmdNames := make([]string, 0)
		getName := func(c *cobra.Command) {}
		getName = func(c *cobra.Command) {
			if c == nil {
				return
			}
			cmdNames = append([]string{c.Name()}, cmdNames...)
			if c.Name() == "web" {
				isWeb = true
			}
			getName(c.Parent())
		}
		getName(cmd)

		if !isWeb {
			if documentFileName == "" {
				path, err := os.Getwd()
				if err != nil {
					log.Fatal(err)
				}
				documentFileName = filepath.Join(path, "bai2.bin")
			}

			_, err := os.Stat(documentFileName)
			if os.IsNotExist(err) {
				return errors.New("invalid input file")
			}

			documentBuffer, err = os.ReadFile(documentFileName)
			if err != nil {
				return err
			}
		}

		return nil
	},
}

func initRootCmd() {
	WebCmd.Flags().BoolP("test", "t", false, "test server")

	rootCmd.SilenceUsage = true
	rootCmd.PersistentFlags().StringVar(&documentFileName, "input", "", "bai2 report file")
	rootCmd.AddCommand(WebCmd)
	rootCmd.AddCommand(Print)
	rootCmd.AddCommand(Parse)
	rootCmd.AddCommand(Format)
}

func main() {
	initRootCmd()

	rootCmd.Execute()
}

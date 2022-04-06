// Copyright 2022 The Moov Authors
// Use of this source code is governed by an Apache License
// license that can be found in the LICENSE file.

package service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/moov-io/bai2/pkg/file"
)

func outputError(w http.ResponseWriter, code int, err error) {
	w.WriteHeader(code)
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"error": err.Error(),
	})
}

func outputSuccess(w http.ResponseWriter, output string) {
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status": output,
	})
}

func parseInputFromRequest(r *http.Request) (*file.Bai2, error) {
	inputFile, _, err := r.FormFile("input")
	if err != nil {
		return nil, err
	}
	defer inputFile.Close()

	var input bytes.Buffer
	if _, err = io.Copy(&input, inputFile); err != nil {
		return nil, err
	}

	// convert byte slice to io.Reader
	reader := bytes.NewReader(input.Bytes())

	return file.Parse(reader)
}

func outputBufferToWriter(w http.ResponseWriter, f *file.Bai2) {
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	fmt.Fprintf(w, f.String())
}

// parse - parse bai2 report
func parse(w http.ResponseWriter, r *http.Request) {
	f, err := parseInputFromRequest(r)
	if err != nil {
		outputError(w, http.StatusBadRequest, err)
		return
	}

	err = f.Validate()
	if err != nil {
		outputError(w, http.StatusNotImplemented, err)
		return
	}

	outputSuccess(w, "valid file")
}

// print - print bai2 report after parse
func print(w http.ResponseWriter, r *http.Request) {
	f, err := parseInputFromRequest(r)
	if err != nil {
		outputError(w, http.StatusBadRequest, err)
		return
	}

	outputBufferToWriter(w, f)
}

// health - health check
func health(w http.ResponseWriter, r *http.Request) {
	outputSuccess(w, "alive")
}

// configure handlers
func ConfigureHandlers(r *mux.Router) error {

	r.HandleFunc("/health", health).Methods("GET")
	r.HandleFunc("/print", print).Methods("POST")
	r.HandleFunc("/parse", parse).Methods("POST")

	return nil
}

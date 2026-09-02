// Copyright 2022 The Moov Authors
// Use of this source code is governed by an Apache License
// license that can be found in the LICENSE file.

package service_test

import (
	"bytes"
	"encoding/json"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/gorilla/mux"
	"github.com/moov-io/bai2/pkg/lib"
	"github.com/moov-io/bai2/pkg/service"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

var (
	parseErrorFileName                = "errors/sample-parseError.txt"
	testFileName                      = "sample1.txt"
	testDetailsWithNewlineTermination = "sample4-continuations-newline-delimited.txt"
	testDetailsWithSlashInText        = "sample5-issue113.txt"
)

type HandlersTest struct {
	suite.Suite
	testServer *mux.Router
}

func (suite *HandlersTest) makeRequest(method, url, body string) (*httptest.ResponseRecorder, *http.Request) {
	request, err := http.NewRequest(method, url, strings.NewReader(body))
	assert.Equal(suite.T(), nil, err)
	recorder := httptest.NewRecorder()
	return recorder, request
}

func (suite *HandlersTest) getWriter(name string) (*multipart.Writer, *bytes.Buffer) {

	path := filepath.Join("..", "..", "test", "testdata", name)
	file, err := os.Open(path)
	assert.Equal(suite.T(), nil, err)

	defer file.Close()

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile("input", filepath.Base(path))
	assert.Equal(suite.T(), nil, err)

	_, err = io.Copy(part, file)
	assert.Equal(suite.T(), nil, err)
	return writer, body
}

func (suite *HandlersTest) SetupTest() {

	suite.testServer = mux.NewRouter()

	err := service.ConfigureHandlers(suite.testServer)
	assert.Equal(suite.T(), nil, err)
}

func TestHandlersTestSuite(t *testing.T) {
	suite.Run(t, new(HandlersTest))
}

func (suite *HandlersTest) TestUnknownRequest() {
	recorder, request := suite.makeRequest(http.MethodGet, "/unknown", "")
	suite.testServer.ServeHTTP(recorder, request)
	assert.Equal(suite.T(), http.StatusNotFound, recorder.Code)
}

func (suite *HandlersTest) TestHealth() {
	recorder, request := suite.makeRequest(http.MethodGet, "/health", "")
	suite.testServer.ServeHTTP(recorder, request)
	assert.Equal(suite.T(), http.StatusOK, recorder.Code)
}

func (suite *HandlersTest) TestPrint() {

	writer, body := suite.getWriter(testFileName)

	err := writer.Close()
	assert.Equal(suite.T(), nil, err)

	recorder, request := suite.makeRequest(http.MethodPost, "/print", body.String())
	request.Header.Set("Content-Type", writer.FormDataContentType())

	suite.testServer.ServeHTTP(recorder, request)
	assert.Equal(suite.T(), http.StatusOK, recorder.Code)
	suite.requirePrintRoundTrip(recorder.Body.String())
}

func (suite *HandlersTest) requirePrintRoundTrip(body string) {
	scan := lib.NewBai2Scanner(strings.NewReader(body))
	file := lib.NewBai2()
	require.NoError(suite.T(), file.Read(&scan))
	require.NoError(suite.T(), file.Validate())
	require.NotEmpty(suite.T(), file.Groups)
}

func (suite *HandlersTest) TestParse() {

	writer, body := suite.getWriter(testFileName)
	err := writer.Close()
	assert.Equal(suite.T(), nil, err)

	recorder, request := suite.makeRequest(http.MethodPost, "/parse", body.String())
	request.Header.Set("Content-Type", writer.FormDataContentType())

	suite.testServer.ServeHTTP(recorder, request)
	assert.Equal(suite.T(), http.StatusOK, recorder.Code)
}

func (suite *HandlersTest) TestFormat() {

	writer, body := suite.getWriter(testFileName)
	err := writer.Close()
	assert.Equal(suite.T(), nil, err)

	recorder, request := suite.makeRequest(http.MethodPost, "/format", body.String())
	request.Header.Set("Content-Type", writer.FormDataContentType())

	suite.testServer.ServeHTTP(recorder, request)
	assert.Equal(suite.T(), http.StatusOK, recorder.Code)

	var file lib.Bai2
	require.NoError(suite.T(), json.Unmarshal(recorder.Body.Bytes(), &file))
	require.Equal(suite.T(), "0004", file.Sender)
	require.Len(suite.T(), file.Groups, 1)
	require.Len(suite.T(), file.Groups[0].Accounts, 2)
	require.Equal(suite.T(), 11, len(file.Groups[0].Accounts[0].Details))
	for _, acct := range file.Groups[0].Accounts {
		for _, summary := range acct.Summaries {
			require.NotEmpty(suite.T(), summary.TypeCode)
		}
		for _, detail := range acct.Details {
			require.False(suite.T(), strings.HasSuffix(detail.Text, "/"), detail.Text)
		}
	}
}

func (suite *HandlersTest) TestPrint_Bai2FileWithNewlineDelimitedContinuations() {
	writer, body := suite.getWriter(testDetailsWithNewlineTermination)
	err := writer.Close()
	assert.Equal(suite.T(), nil, err)

	recorder, request := suite.makeRequest(http.MethodPost, "/print", body.String())
	request.Header.Set("Content-Type", writer.FormDataContentType())

	suite.testServer.ServeHTTP(recorder, request)
	assert.Equal(suite.T(), http.StatusOK, recorder.Code)
	suite.requirePrintRoundTrip(recorder.Body.String())
}

func (suite *HandlersTest) TestFormat_Bai2FileWithNewlineDelimitedContinuations() {
	writer, body := suite.getWriter(testDetailsWithNewlineTermination)
	err := writer.Close()
	assert.Equal(suite.T(), nil, err)

	recorder, request := suite.makeRequest(http.MethodPost, "/format", body.String())
	request.Header.Set("Content-Type", writer.FormDataContentType())

	suite.testServer.ServeHTTP(recorder, request)
	assert.Equal(suite.T(), http.StatusOK, recorder.Code)

	var file lib.Bai2
	require.NoError(suite.T(), json.Unmarshal(recorder.Body.Bytes(), &file))
	require.Len(suite.T(), file.Groups, 1)
	accounts := file.Groups[0].Accounts
	require.Len(suite.T(), accounts, 5)
	require.Len(suite.T(), accounts[1].Details, 3)
	require.Len(suite.T(), accounts[2].Details, 14)
	for _, acct := range accounts {
		for _, summary := range acct.Summaries {
			require.NotEmpty(suite.T(), summary.TypeCode)
		}
	}
}

func (suite *HandlersTest) TestPrint_Bai2FileWithSlashInText_Issue113() {
	writer, body := suite.getWriter(testDetailsWithSlashInText)
	err := writer.Close()
	assert.Equal(suite.T(), nil, err)

	recorder, request := suite.makeRequest(http.MethodPost, "/print", body.String())
	request.Header.Set("Content-Type", writer.FormDataContentType())

	suite.testServer.ServeHTTP(recorder, request)
	assert.Equal(suite.T(), http.StatusOK, recorder.Code)
	suite.requirePrintRoundTrip(recorder.Body.String())
	require.Contains(suite.T(), recorder.Body.String(), "AB/GS/RPFILERP0001/RPBA0001")
	require.Contains(suite.T(), recorder.Body.String(), "08/18/23 Invoice - Sample")
}

func (suite *HandlersTest) TestFormat_Bai2FileWithSlashInText_Issue113() {
	writer, body := suite.getWriter(testDetailsWithSlashInText)
	err := writer.Close()
	assert.Equal(suite.T(), nil, err)

	recorder, request := suite.makeRequest(http.MethodPost, "/format", body.String())
	request.Header.Set("Content-Type", writer.FormDataContentType())

	suite.testServer.ServeHTTP(recorder, request)
	assert.Equal(suite.T(), http.StatusOK, recorder.Code)

	var file lib.Bai2
	require.NoError(suite.T(), json.Unmarshal(recorder.Body.Bytes(), &file))
	require.Len(suite.T(), file.Groups, 1)
	accounts := file.Groups[0].Accounts
	require.Len(suite.T(), accounts, 5)
	require.Len(suite.T(), accounts[1].Details, 3)
	require.Len(suite.T(), accounts[2].Details, 17)
	require.Equal(suite.T(), "AB/GS/RPFILERP0001/RPBA0001", accounts[1].Details[2].CustomerReferenceNumber)
}

func (suite *HandlersTest) TestParse_BAI3() {
	writer, body := suite.getWriter("bai3/x9-mandatory.txt")
	err := writer.Close()
	assert.Equal(suite.T(), nil, err)

	recorder, request := suite.makeRequest(http.MethodPost, "/parse", body.String())
	request.Header.Set("Content-Type", writer.FormDataContentType())
	suite.testServer.ServeHTTP(recorder, request)
	assert.Equal(suite.T(), http.StatusOK, recorder.Code)
}

func (suite *HandlersTest) TestFormat_BAI3() {
	writer, body := suite.getWriter("bai3/x9-status-summary-detail.txt")
	err := writer.Close()
	assert.Equal(suite.T(), nil, err)

	recorder, request := suite.makeRequest(http.MethodPost, "/format", body.String())
	request.Header.Set("Content-Type", writer.FormDataContentType())
	suite.testServer.ServeHTTP(recorder, request)
	assert.Equal(suite.T(), http.StatusOK, recorder.Code)
	require.Contains(suite.T(), recorder.Body.String(), `"Banks"`)
	require.Contains(suite.T(), recorder.Body.String(), `"CAD"`)
}

func (suite *HandlersTest) TestParse_IgnoreVersion() {
	raw := "01,122099999,123456789,040621,0200,1,,,3/\n02,031001234,122099999,1,040620,2359,,2/\n03,5765432,,,,,/\n49,0,2/\n98,0,1,4/\n99,0,1,6/\n"
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile("input", "v3.bai")
	require.NoError(suite.T(), err)
	_, err = part.Write([]byte(raw))
	require.NoError(suite.T(), err)
	require.NoError(suite.T(), writer.Close())

	recorder, request := suite.makeRequest(http.MethodPost, "/parse", body.String())
	request.Header.Set("Content-Type", writer.FormDataContentType())
	suite.testServer.ServeHTTP(recorder, request)
	assert.Equal(suite.T(), http.StatusBadRequest, recorder.Code)

	recorder, request = suite.makeRequest(http.MethodPost, "/parse?ignoreVersion=true", body.String())
	request.Header.Set("Content-Type", writer.FormDataContentType())
	suite.testServer.ServeHTTP(recorder, request)
	assert.Equal(suite.T(), http.StatusOK, recorder.Code)
}

func (suite *HandlersTest) TestPrint_ParseError() {
	writer, body := suite.getWriter(parseErrorFileName)
	err := writer.Close()
	assert.Equal(suite.T(), nil, err)

	recorder, request := suite.makeRequest(http.MethodPost, "/print", body.String())
	request.Header.Set("Content-Type", writer.FormDataContentType())

	suite.testServer.ServeHTTP(recorder, request)
	assert.Equal(suite.T(), http.StatusBadRequest, recorder.Code)
	assert.Equal(suite.T(), recorder.Body.String(), `{"error":"first record is \"00,0004,1/\", want 01 file header"}
`)
}

func (suite *HandlersTest) TestParse_ParseError() {
	writer, body := suite.getWriter(parseErrorFileName)
	err := writer.Close()
	assert.Equal(suite.T(), nil, err)

	recorder, request := suite.makeRequest(http.MethodPost, "/parse", body.String())
	request.Header.Set("Content-Type", writer.FormDataContentType())

	suite.testServer.ServeHTTP(recorder, request)
	assert.Equal(suite.T(), http.StatusBadRequest, recorder.Code)
	assert.Equal(suite.T(), recorder.Body.String(), `{"error":"first record is \"00,0004,1/\", want 01 file header"}
`)
}

func (suite *HandlersTest) TestFormat_ParseError() {
	writer, body := suite.getWriter(parseErrorFileName)
	err := writer.Close()
	assert.Equal(suite.T(), nil, err)

	recorder, request := suite.makeRequest(http.MethodPost, "/format", body.String())
	request.Header.Set("Content-Type", writer.FormDataContentType())

	suite.testServer.ServeHTTP(recorder, request)
	assert.Equal(suite.T(), http.StatusBadRequest, recorder.Code)
	assert.Equal(suite.T(), recorder.Body.String(), `{"error":"first record is \"00,0004,1/\", want 01 file header"}
`)
}

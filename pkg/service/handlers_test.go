// Copyright 2022 The Moov Authors
// Use of this source code is governed by an Apache License
// license that can be found in the LICENSE file.

package service_test

import (
	"bytes"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/gorilla/mux"
	"github.com/moov-io/bai2/pkg/service"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

var (
	testFileName = "sample1.txt"
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
	assert.Equal(suite.T(), recorder.Body.String(), `{"sender":"0004","receiver":"12345","fileCreatedDate":"060321","fileCreatedTime":"0829","fileIdNumber":"001","physicalRecordLength":80,"blockSize":1,"versionNumber":2,"fileControlTotal":"+00000000001280000","numberOfGroups":1,"numberOfRecords":27,"Groups":[{"receiver":"12345","originator":"0004","groupStatus":1,"asOfDate":"060317","currencyCode":"CAD","groupControlTotal":"+00000000001280000","numberOfAccounts":2,"numberOfRecords":25,"Accounts":[{"accountNumber":"10200123456","currencyCode":"CAD","summaries":[{"TypeCode":"040","Amount":"+000000000000","ItemCount":0,"FundsType":{}},{"TypeCode":"045","Amount":"+000000000000","ItemCount":0,"FundsType":{}},{"TypeCode":"100","Amount":"000000000208500","ItemCount":3,"FundsType":{"type_code":"V","date":"060316"}},{"TypeCode":"400","Amount":"000000000208500","ItemCount":8,"FundsType":{"type_code":"V","date":"060316"}}],"accountControlTotal":"+00000000000834000","numberRecords":14,"Details":[{"TypeCode":"409","Amount":"000000000002500","FundsType":{"type_code":"V","date":"060316"},"BankReferenceNumber":"","CustomerReferenceNumber":"","Text":"RETURNED CHEQUE     "}]},{"accountNumber":"10200123456","currencyCode":"CAD","summaries":[{"TypeCode":"040","Amount":"+000000000000","ItemCount":0,"FundsType":{}},{"TypeCode":"045","Amount":"+000000000000","ItemCount":0,"FundsType":{}},{"TypeCode":"100","Amount":"000000000111500","ItemCount":2,"FundsType":{"type_code":"V","date":"060317"}},{"TypeCode":"400","Amount":"000000000111500","ItemCount":4,"FundsType":{"type_code":"V","date":"060317"}}],"accountControlTotal":"+00000000000446000","numberRecords":9,"Details":[{"TypeCode":"108","Amount":"000000000011500","FundsType":{"type_code":"V","date":"060317"},"BankReferenceNumber":"","CustomerReferenceNumber":"","Text":"TFR 1020 0345678    "}]}]}]}
`)
}

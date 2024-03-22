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
	parseErrorFileName                = "sample-parseError.txt"
	testFileName                      = "sample1.txt"
	testDetailsWithNewlineTermination = "sample4-continuations-newline-delimited.txt"
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

	// Verify that the printed file matches the input file.
	path := filepath.Join("..", "..", "test", "testdata", testFileName)
	fixture, err := os.ReadFile(path)
	assert.Equal(suite.T(), nil, err)

	// NB. Account continuations are currently not written to file exactly as they were read.
	// Because of this behavior, the returned body does NOT strictly match the file data.
	// This test currently asserts on the shape of the file as created by the current return.
	// The difference between this output and the sample file is that a subset of data provided on
	// each account continuation (88) is instead output on the Account record (03).
	assert.NotEqual(suite.T(), recorder.Body.String(), string(fixture))

	expectedFileBody := `01,0004,12345,060321,0829,001,80,1,2/
02,12345,0004,1,060317,,CAD,/
03,10200123456,CAD,040,+000000000000,,,045,+000000000000,,,100,000000000208500/
88,3,V,060316,,400,000000000208500,8,V,060316,/
16,409,000000000002500,V,060316,,,,RETURNED CHEQUE     /
16,409,000000000090000,V,060316,,,,RTN-UNKNOWN         /
16,409,000000000000500,V,060316,,,,RTD CHQ SERVICE CHRG/
16,108,000000000203500,V,060316,,,,TFR 1020 0345678    /
16,108,000000000002500,V,060316,,,,MACLEOD MALL        /
16,108,000000000002500,V,060316,,,,MASCOUCHE QUE       /
16,409,000000000020000,V,060316,,,,1000 ISLANDS MALL   /
16,409,000000000090000,V,060316,,,,PENHORA MALL        /
16,409,000000000002000,V,060316,,,,CAPILANO MALL       /
16,409,000000000002500,V,060316,,,,GALERIES LA CAPITALE/
16,409,000000000001000,V,060316,,,,PLAZA ROCK FOREST   /
49,+00000000000834000,14/
03,10200123456,CAD,040,+000000000000,,,045,+000000000000,,,100,000000000111500/
88,2,V,060317,,400,000000000111500,4,V,060317,/
16,108,000000000011500,V,060317,,,,TFR 1020 0345678    /
16,108,000000000100000,V,060317,,,,MONTREAL            /
16,409,000000000100000,V,060317,,,,GRANDFALL NB        /
16,409,000000000009000,V,060317,,,,HAMILTON ON         /
16,409,000000000002000,V,060317,,,,WOODSTOCK NB        /
16,409,000000000000500,V,060317,,,,GALERIES RICHELIEU  /
49,+00000000000446000,9/
98,+00000000001280000,2,25/
99,+00000000001280000,1,27/`
	assert.Equal(suite.T(), recorder.Body.String(), expectedFileBody)
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
	assert.Equal(suite.T(), recorder.Body.String(), `{"sender":"0004","receiver":"12345","fileCreatedDate":"060321","fileCreatedTime":"0829","fileIdNumber":"001","physicalRecordLength":80,"blockSize":1,"versionNumber":2,"fileControlTotal":"+00000000001280000","numberOfGroups":1,"numberOfRecords":27,"Groups":[{"receiver":"12345","originator":"0004","groupStatus":1,"asOfDate":"060317","currencyCode":"CAD","groupControlTotal":"+00000000001280000","numberOfAccounts":2,"numberOfRecords":25,"Accounts":[{"accountNumber":"10200123456","currencyCode":"CAD","summaries":[{"TypeCode":"040","Amount":"+000000000000","ItemCount":0,"FundsType":{}},{"TypeCode":"045","Amount":"+000000000000","ItemCount":0,"FundsType":{}},{"TypeCode":"100","Amount":"000000000208500","ItemCount":3,"FundsType":{"type_code":"V","date":"060316"}},{"TypeCode":"400","Amount":"000000000208500","ItemCount":8,"FundsType":{"type_code":"V","date":"060316"}}],"accountControlTotal":"+00000000000834000","numberRecords":14,"Details":[{"TypeCode":"409","Amount":"000000000002500","FundsType":{"type_code":"V","date":"060316"},"BankReferenceNumber":"","CustomerReferenceNumber":"","Text":"RETURNED CHEQUE     "},{"TypeCode":"409","Amount":"000000000090000","FundsType":{"type_code":"V","date":"060316"},"BankReferenceNumber":"","CustomerReferenceNumber":"","Text":"RTN-UNKNOWN         "},{"TypeCode":"409","Amount":"000000000000500","FundsType":{"type_code":"V","date":"060316"},"BankReferenceNumber":"","CustomerReferenceNumber":"","Text":"RTD CHQ SERVICE CHRG"},{"TypeCode":"108","Amount":"000000000203500","FundsType":{"type_code":"V","date":"060316"},"BankReferenceNumber":"","CustomerReferenceNumber":"","Text":"TFR 1020 0345678    "},{"TypeCode":"108","Amount":"000000000002500","FundsType":{"type_code":"V","date":"060316"},"BankReferenceNumber":"","CustomerReferenceNumber":"","Text":"MACLEOD MALL        "},{"TypeCode":"108","Amount":"000000000002500","FundsType":{"type_code":"V","date":"060316"},"BankReferenceNumber":"","CustomerReferenceNumber":"","Text":"MASCOUCHE QUE       "},{"TypeCode":"409","Amount":"000000000020000","FundsType":{"type_code":"V","date":"060316"},"BankReferenceNumber":"","CustomerReferenceNumber":"","Text":"1000 ISLANDS MALL   "},{"TypeCode":"409","Amount":"000000000090000","FundsType":{"type_code":"V","date":"060316"},"BankReferenceNumber":"","CustomerReferenceNumber":"","Text":"PENHORA MALL        "},{"TypeCode":"409","Amount":"000000000002000","FundsType":{"type_code":"V","date":"060316"},"BankReferenceNumber":"","CustomerReferenceNumber":"","Text":"CAPILANO MALL       "},{"TypeCode":"409","Amount":"000000000002500","FundsType":{"type_code":"V","date":"060316"},"BankReferenceNumber":"","CustomerReferenceNumber":"","Text":"GALERIES LA CAPITALE"},{"TypeCode":"409","Amount":"000000000001000","FundsType":{"type_code":"V","date":"060316"},"BankReferenceNumber":"","CustomerReferenceNumber":"","Text":"PLAZA ROCK FOREST   "}]},{"accountNumber":"10200123456","currencyCode":"CAD","summaries":[{"TypeCode":"040","Amount":"+000000000000","ItemCount":0,"FundsType":{}},{"TypeCode":"045","Amount":"+000000000000","ItemCount":0,"FundsType":{}},{"TypeCode":"100","Amount":"000000000111500","ItemCount":2,"FundsType":{"type_code":"V","date":"060317"}},{"TypeCode":"400","Amount":"000000000111500","ItemCount":4,"FundsType":{"type_code":"V","date":"060317"}}],"accountControlTotal":"+00000000000446000","numberRecords":9,"Details":[{"TypeCode":"108","Amount":"000000000011500","FundsType":{"type_code":"V","date":"060317"},"BankReferenceNumber":"","CustomerReferenceNumber":"","Text":"TFR 1020 0345678    "},{"TypeCode":"108","Amount":"000000000100000","FundsType":{"type_code":"V","date":"060317"},"BankReferenceNumber":"","CustomerReferenceNumber":"","Text":"MONTREAL            "},{"TypeCode":"409","Amount":"000000000100000","FundsType":{"type_code":"V","date":"060317"},"BankReferenceNumber":"","CustomerReferenceNumber":"","Text":"GRANDFALL NB        "},{"TypeCode":"409","Amount":"000000000009000","FundsType":{"type_code":"V","date":"060317"},"BankReferenceNumber":"","CustomerReferenceNumber":"","Text":"HAMILTON ON         "},{"TypeCode":"409","Amount":"000000000002000","FundsType":{"type_code":"V","date":"060317"},"BankReferenceNumber":"","CustomerReferenceNumber":"","Text":"WOODSTOCK NB        "},{"TypeCode":"409","Amount":"000000000000500","FundsType":{"type_code":"V","date":"060317"},"BankReferenceNumber":"","CustomerReferenceNumber":"","Text":"GALERIES RICHELIEU  "}]}]}]}
`)
}

func (suite *HandlersTest) TestPrint_Bai2FileWithNewlineDelimitedContinuations() {

	writer, body := suite.getWriter(testDetailsWithNewlineTermination)
	err := writer.Close()
	assert.Equal(suite.T(), nil, err)

	recorder, request := suite.makeRequest(http.MethodPost, "/print", body.String())
	request.Header.Set("Content-Type", writer.FormDataContentType())

	suite.testServer.ServeHTTP(recorder, request)
	assert.Equal(suite.T(), http.StatusOK, recorder.Code)

	expectedFileBody := `01,GSBI,cont001,210706,1249,1,,,2/
02,cont001,026015079,1,230906,2000,,/
03,107049924,USD,,,,,060,13053325440,,,100,000,,,400,000,,/
49,13053325440,2/
03,107049932,USD,,,,,060,6865898,,,100,1912,1,,400,000,,/
16,447,60000,,SPB2322984714570,1111,ACH Credit Payment,Entry Description: EXP; -, SEC: CCD, Client Ref ID: 1111, GS ID: SPB2322984714570
88,EREF: 1111
88,DBNM: TEST INC
88,CACT: ACHCONTROLOUTUSD01
16,261,143500,,SB2322600000404,GSQ4FBGFDGWGKY,ACH Credit Reject,From: TEST INC, Remittance Info: "ACH- Test - Addenda Record", Entry Description: TRADE; -, SEC: CTX, Client Ref ID: GSQ4FBGFDGWGKY, GS ID: SB2322600000404
88,CREF: 
88,REMI: ACH- Test - Addenda Record
88,EREF: GSQ4FBGFDGWGKY
88,CRNM: Test
88,DBNM: SAMPLE INC
88,DACT: 101152046
88,DABA: 026015079
16,447,928650,,SPB2322684598521,AB-GS-RPFILERP0001-RPBA0001,ACH Credit Payment,Entry Description: TRADE; -, SEC: CTX, Client Ref ID: AB-GS-TEST0001-RPBA0001, GS ID: SPB2322684598521
88,EREF: AB-GS-RPFILERP0001-RPBA0001
88,DBNM: SAMPLE INC
88,CACT: ACHCONTROLOUTUSD01
49,-1260161341762,26/
16,557,200000,,SB2322600000214,021000080000030,ACH Credit Receipt Return,Return To: Test, Remittance Info: "SB2322300000052", Entry Description: EXP; -, SEC: CCD, Reason: "R02", Return of Client Ref ID: 021000080000030, GS ID: SB2322600000214
88,CREF: 026015076104300
88,IDNM: 1114
88,EREF: 021000080000030
88,CRNM: Test
88,DBNM: SAMPLE INC.
88,CABA: 021000089
16,451,55555,,SB2322600000455,021000020000021,ACH Debit Payment,To: TEST, Entry Description: INVOICES; 210630, SEC: CCD, Client Ref ID: 021000020000021, GS ID: SB2322600000455
88,CREF: 021000020000021
88,IDNM: 2009282
88,EREF: 021000020000021
88,CRNM: TEST
88,DBNM: SAMPLE INC
88,CABA: 021000021
16,266,1912,,GI2118700002010,20210706MMQFMPU8000001,Outgoing Wire Return,-
88,CREF: 20210706MMQFMPU8000001
88,EREF: 20210706MMQFMPU8000001
88,DBIC: GSCRUS33
88,CRNM: ABC Company
88,DBNM: SAMPLE INC.
16,495,50500,,GI2321400000090,GSV0DL6RKT,Outgoing Wire,To: TEST COMPANY, Remittance Info: "QWERTIOP", Client Ref ID: GSV0DL6RKT, GS ID: GI2321400000090, Settled Amt: EUR 322.00, FX Rate: 156.833677
88,REMI: QWERTIOP
88,EREF: GSV0DL6RKT
88,CBIC: COBADEFF
88,CRNM: TEST COMPANY
88,DBNM: SAMPLE TEST
16,195,1125,,GI2229300000187,GS0D9VGMP1IWPLW,Incoming Wire,-
88,EREF: GS0D9VGMP1IWPLW
88,DBIC: CITIUS30XXX
88,CRNM: ABC CORPORATION
88,DACT: 8348572423
88,CHKN: GSIL2X6103UNCRSF
16,257,60000,,SB2225800001203,028000020000335,ACH Debit Payment Return,Return From: Company1, Entry Description: TRADE; -, SEC: CCD, Reason: "R02", Return of Client Ref ID: 028000020000335, GS ID: SB2225800001203
88,IDNM: 1
88,EREF: 028000020000335
88,CRNM: TEST INC
88,DBNM: Company1
88,DABA: 028000024
16,255,931,,SC2134800001999,,Check Return,Return From: Test2 Customer, Check Serial Number: 0009000000, Return Reason: "Payee does not exist", Client Ref ID: 74564762445, GS ID: SC213480000120999
88:EREF: 07370568132
88,CRNM: Test Inc.
88,DBNM: Test2 Customer
88,CABA: 12345
88,CHKN: 0009000000
16,195,50050,,GI2228400005800,RTR60880840833,RTP Incoming,From: SAMPLE INC, Remittance Info: "Test Remittance", Client Ref ID: RTR60880840833, GS ID: GI2228400005800, Clearing Ref: 001
88,REMI: Test Remittance
88,EREF: RTR60880840833
88,CRNM: RTR-CdtrName
88,DBNM: SAMPLE INC
88,DACT: 02122056789012205
88,DABA: 000000010
16,175,527,,SX22293073766088,GS4N04L1COP45VY,Check Deposit,-
88,EREF: GS4N04L1COP45VY
88,DACT: 100168723
16,475,10100,,SC2229300000152,01030340329,Check Paid,-
88,REMI: UAT testing for Checks
88,EREF: 01030340329
88,CRNM: TEST INC
88,DBNM: ABC CORP
88,CABA: 12345
88,CHKN: 006034594478
16,275,337686,,GI2318000014342,e457328416d411eeaf020a58a9feac02,Cash Concentration,From: SAMPLE INC, Account: 290000020437, GS Cash Concentration, "Structure ID: CC0000000", GS ID: GI2318000212121
88,REMI: Structure ID: CC0000082
88,EREF: e123456786d411eeaf020a58a9feac02
88,DBIC: GSCRUS33VIA
88,CRNM: SAMPLE INC
88,DBNM: SAMPLE INC
88,DACT: 290000020437
16,165,5000,,SPB2321284264201,AB-GS-DDFILEAB0001-DDBAB0001,ACH Debit Collection,Entry Description: BILL PMT; -, SEC: CCD, Client Ref ID: AB-GS-DDFILEAB0001-DDBAB0001, GS ID: SPB2321284264201
88,EREF: AB-GS-DDFILEAB0001-DDBAB0001
88,CRNM: SAMPLE LLP
88,DACT: ACHCONTROLINUSD01
16,475,44250,,SC2323300002416,8ce1829175a74ec88d67010dd7fb6132,Check Paid,To: TEST AND COMPANY LLC, Check Serial Number: 24108, GS ID: SC2323300002416
88,EREF: 8ce1829175a74ec88d67010dd7fb6132
88,CRNM: TEST AND COMPANY LLC
88,DBNM: Sample Inc.
88,CABA: 0
88,CHKN: 24108
16,495,30000000,,GI2323300009168,3785726,Outgoing Wire,To: TEST AND COMPANY, Remittance Info: "081823 Invoice - Sample", Client Ref ID: 3785726, GS ID: GI2323300009168, Clearing Ref: 20230821MMQFMPU7004100
88,CREF: 20230821MMQFMPU7004100
88,REMI: 081823 Invoice - Sample
88,EREF: 3785726
88,CRNM: TEST AND COMPANY
88,DBNM: Sample Inc.
88,CACT: 609873838
88,CABA: 021000021
49,6869722,8/
49,000,2/
03,280000010657,USD,,,,,060,000,,,100,000,,,400,000,,/
49,000,2/
98,13060195162,4,16/
99,13060195162,1,18/`
	assert.Equal(suite.T(), recorder.Body.String(), expectedFileBody)
}

func (suite *HandlersTest) TestFormat_Bai2FileWithNewlineDelimitedContinuations() {

	writer, body := suite.getWriter(testDetailsWithNewlineTermination)
	err := writer.Close()
	assert.Equal(suite.T(), nil, err)

	recorder, request := suite.makeRequest(http.MethodPost, "/format", body.String())
	request.Header.Set("Content-Type", writer.FormDataContentType())

	suite.testServer.ServeHTTP(recorder, request)
	assert.Equal(suite.T(), http.StatusOK, recorder.Code)
	assert.Equal(suite.T(), recorder.Body.String(), `{"sender":"GSBI","receiver":"cont001","fileCreatedDate":"210706","fileCreatedTime":"1249","fileIdNumber":"1","versionNumber":2,"fileControlTotal":"13060195162","numberOfGroups":1,"numberOfRecords":18,"Groups":[{"receiver":"cont001","originator":"026015079","groupStatus":1,"asOfDate":"230906","asOfTime":"2000","groupControlTotal":"13060195162","numberOfAccounts":4,"numberOfRecords":16,"Accounts":[{"accountNumber":"107049924","currencyCode":"USD","summaries":[{"TypeCode":"","Amount":"","ItemCount":0,"FundsType":{}},{"TypeCode":"060","Amount":"13053325440","ItemCount":0,"FundsType":{}},{"TypeCode":"100","Amount":"000","ItemCount":0,"FundsType":{}},{"TypeCode":"400","Amount":"000","ItemCount":0,"FundsType":{}}],"accountControlTotal":"13053325440","numberRecords":2,"Details":null},{"accountNumber":"107049932","currencyCode":"USD","summaries":[{"TypeCode":"","Amount":"","ItemCount":0,"FundsType":{}},{"TypeCode":"060","Amount":"6865898","ItemCount":0,"FundsType":{}},{"TypeCode":"100","Amount":"1912","ItemCount":1,"FundsType":{}},{"TypeCode":"400","Amount":"000","ItemCount":0,"FundsType":{}}],"accountControlTotal":"000","numberRecords":2,"Details":[{"TypeCode":"447","Amount":"60000","FundsType":{},"BankReferenceNumber":"SPB2322984714570","CustomerReferenceNumber":"1111","Text":"ACH Credit Payment,Entry Description: EXP; -, SEC: CCD, Client Ref ID: 1111, GS ID: SPB2322984714570\n88,EREF: 1111\n88,DBNM: TEST INC\n88,CACT: ACHCONTROLOUTUSD01\n16,261,143500,,SB2322600000404,GSQ4FBGFDGWGKY,ACH Credit Reject,From: TEST INC, Remittance Info: \"ACH- Test - Addenda Record\", Entry Description: TRADE; -, SEC: CTX, Client Ref ID: GSQ4FBGFDGWGKY, GS ID: SB2322600000404\n88,CREF: \n88,REMI: ACH- Test - Addenda Record\n88,EREF: GSQ4FBGFDGWGKY\n88,CRNM: Test\n88,DBNM: SAMPLE INC\n88,DACT: 101152046\n88,DABA: 026015079\n16,447,928650,,SPB2322684598521,AB-GS-RPFILERP0001-RPBA0001,ACH Credit Payment,Entry Description: TRADE; -, SEC: CTX, Client Ref ID: AB-GS-TEST0001-RPBA0001, GS ID: SPB2322684598521\n88,EREF: AB-GS-RPFILERP0001-RPBA0001\n88,DBNM: SAMPLE INC\n88,CACT: ACHCONTROLOUTUSD01\n49,-1260161341762,26"},{"TypeCode":"557","Amount":"200000","FundsType":{},"BankReferenceNumber":"SB2322600000214","CustomerReferenceNumber":"021000080000030","Text":"ACH Credit Receipt Return,Return To: Test, Remittance Info: \"SB2322300000052\", Entry Description: EXP; -, SEC: CCD, Reason: \"R02\", Return of Client Ref ID: 021000080000030, GS ID: SB2322600000214\n88,CREF: 026015076104300\n88,IDNM: 1114\n88,EREF: 021000080000030\n88,CRNM: Test\n88,DBNM: SAMPLE INC.\n88,CABA: 021000089\n16,451,55555,,SB2322600000455,021000020000021,ACH Debit Payment,To: TEST, Entry Description: INVOICES; 210630, SEC: CCD, Client Ref ID: 021000020000021, GS ID: SB2322600000455\n88,CREF: 021000020000021\n88,IDNM: 2009282\n88,EREF: 021000020000021\n88,CRNM: TEST\n88,DBNM: SAMPLE INC\n88,CABA: 021000021\n16,266,1912,,GI2118700002010,20210706MMQFMPU8000001,Outgoing Wire Return,-\n88,CREF: 20210706MMQFMPU8000001\n88,EREF: 20210706MMQFMPU8000001\n88,DBIC: GSCRUS33\n88,CRNM: ABC Company\n88,DBNM: SAMPLE INC.\n16,495,50500,,GI2321400000090,GSV0DL6RKT,Outgoing Wire,To: TEST COMPANY, Remittance Info: \"QWERTIOP\", Client Ref ID: GSV0DL6RKT, GS ID: GI2321400000090, Settled Amt: EUR 322.00, FX Rate: 156.833677\n88,REMI: QWERTIOP\n88,EREF: GSV0DL6RKT\n88,CBIC: COBADEFF\n88,CRNM: TEST COMPANY\n88,DBNM: SAMPLE TEST\n16,195,1125,,GI2229300000187,GS0D9VGMP1IWPLW,Incoming Wire,-\n88,EREF: GS0D9VGMP1IWPLW\n88,DBIC: CITIUS30XXX\n88,CRNM: ABC CORPORATION\n88,DACT: 8348572423\n88,CHKN: GSIL2X6103UNCRSF\n16,257,60000,,SB2225800001203,028000020000335,ACH Debit Payment Return,Return From: Company1, Entry Description: TRADE; -, SEC: CCD, Reason: \"R02\", Return of Client Ref ID: 028000020000335, GS ID: SB2225800001203\n88,IDNM: 1\n88,EREF: 028000020000335\n88,CRNM: TEST INC\n88,DBNM: Company1\n88,DABA: 028000024\n16,255,931,,SC2134800001999,,Check Return,Return From: Test2 Customer, Check Serial Number: 0009000000, Return Reason: \"Payee does not exist\", Client Ref ID: 74564762445, GS ID: SC213480000120999\n88:EREF: 07370568132\n88,CRNM: Test Inc.\n88,DBNM: Test2 Customer\n88,CABA: 12345\n88,CHKN: 0009000000\n16,195,50050,,GI2228400005800,RTR60880840833,RTP Incoming,From: SAMPLE INC, Remittance Info: \"Test Remittance\", Client Ref ID: RTR60880840833, GS ID: GI2228400005800, Clearing Ref: 001\n88,REMI: Test Remittance\n88,EREF: RTR60880840833\n88,CRNM: RTR-CdtrName\n88,DBNM: SAMPLE INC\n88,DACT: 02122056789012205\n88,DABA: 000000010\n16,175,527,,SX22293073766088,GS4N04L1COP45VY,Check Deposit,-\n88,EREF: GS4N04L1COP45VY\n88,DACT: 100168723\n16,475,10100,,SC2229300000152,01030340329,Check Paid,-\n88,REMI: UAT testing for Checks\n88,EREF: 01030340329\n88,CRNM: TEST INC\n88,DBNM: ABC CORP\n88,CABA: 12345\n88,CHKN: 006034594478\n16,275,337686,,GI2318000014342,e457328416d411eeaf020a58a9feac02,Cash Concentration,From: SAMPLE INC, Account: 290000020437, GS Cash Concentration, \"Structure ID: CC0000000\", GS ID: GI2318000212121\n88,REMI: Structure ID: CC0000082\n88,EREF: e123456786d411eeaf020a58a9feac02\n88,DBIC: GSCRUS33VIA\n88,CRNM: SAMPLE INC\n88,DBNM: SAMPLE INC\n88,DACT: 290000020437\n16,165,5000,,SPB2321284264201,AB-GS-DDFILEAB0001-DDBAB0001,ACH Debit Collection,Entry Description: BILL PMT; -, SEC: CCD, Client Ref ID: AB-GS-DDFILEAB0001-DDBAB0001, GS ID: SPB2321284264201\n88,EREF: AB-GS-DDFILEAB0001-DDBAB0001\n88,CRNM: SAMPLE LLP\n88,DACT: ACHCONTROLINUSD01\n16,475,44250,,SC2323300002416,8ce1829175a74ec88d67010dd7fb6132,Check Paid,To: TEST AND COMPANY LLC, Check Serial Number: 24108, GS ID: SC2323300002416\n88,EREF: 8ce1829175a74ec88d67010dd7fb6132\n88,CRNM: TEST AND COMPANY LLC\n88,DBNM: Sample Inc.\n88,CABA: 0\n88,CHKN: 24108\n16,495,30000000,,GI2323300009168,3785726,Outgoing Wire,To: TEST AND COMPANY, Remittance Info: \"081823 Invoice - Sample\", Client Ref ID: 3785726, GS ID: GI2323300009168, Clearing Ref: 20230821MMQFMPU7004100\n88,CREF: 20230821MMQFMPU7004100\n88,REMI: 081823 Invoice - Sample\n88,EREF: 3785726\n88,CRNM: TEST AND COMPANY\n88,DBNM: Sample Inc.\n88,CACT: 609873838\n88,CABA: 021000021\n49,6869722,8"}]},{"accountNumber":"280000010657","currencyCode":"USD","summaries":[{"TypeCode":"","Amount":"","ItemCount":0,"FundsType":{}},{"TypeCode":"060","Amount":"000","ItemCount":0,"FundsType":{}},{"TypeCode":"100","Amount":"000","ItemCount":0,"FundsType":{}},{"TypeCode":"400","Amount":"000","ItemCount":0,"FundsType":{}}],"accountControlTotal":"000","numberRecords":2,"Details":null}]}]}
`)
}

func (suite *HandlersTest) TestPrint_ParseError() {
	writer, body := suite.getWriter(parseErrorFileName)
	err := writer.Close()
	assert.Equal(suite.T(), nil, err)

	recorder, request := suite.makeRequest(http.MethodPost, "/print", body.String())
	request.Header.Set("Content-Type", writer.FormDataContentType())

	suite.testServer.ServeHTTP(recorder, request)
	assert.Equal(suite.T(), http.StatusBadRequest, recorder.Code)
	assert.Equal(suite.T(), recorder.Body.String(), `{"error":"ERROR parsing file on line 1 (unsupported record type 00)"}
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
	assert.Equal(suite.T(), recorder.Body.String(), `{"error":"ERROR parsing file on line 1 (unsupported record type 00)"}
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
	assert.Equal(suite.T(), recorder.Body.String(), `{"error":"ERROR parsing file on line 1 (unsupported record type 00)"}
`)
}

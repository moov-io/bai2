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
	"github.com/moov-io/bai2/pkg/client"
	"github.com/moov-io/bai2/pkg/service"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

var (
	parseErrorFileName                = "errors/sample-parseError.txt"
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
16,447,60000,,SPB2322984714570,1111,ACH Credit Payment,Entry Description: EXP; -, SEC: CCD, Client Ref ID: 1111, GS ID: SPB2322984714570,EREF: 1111,DBNM: TEST INC,CACT: ACHCONTROLOUTUSD01/
16,261,143500,,SB2322600000404,GSQ4FBGFDGWGKY,ACH Credit Reject,From: TEST INC, Remittance Info: "ACH- Test - Addenda Record", Entry Description: TRADE; -, SEC: CTX, Client Ref ID: GSQ4FBGFDGWGKY, GS ID: SB2322600000404,CREF: ,REMI: ACH- Test - Addenda Record,EREF: GSQ4FBGFDGWGKY,CRNM: Test,DBNM: SAMPLE INC,DACT: 101152046,DABA: 026015079/
16,447,928650,,SPB2322684598521,AB-GS-RPFILERP0001-RPBA0001,ACH Credit Payment,Entry Description: TRADE; -, SEC: CTX, Client Ref ID: AB-GS-TEST0001-RPBA0001, GS ID: SPB2322684598521,EREF: AB-GS-RPFILERP0001-RPBA0001,DBNM: SAMPLE INC,CACT: ACHCONTROLOUTUSD01/
49,-1260161341762,26/
03,104108339,USD,010,159581194,,,015,159381194,,,040,158568897,,,045,158368897,,,100,000,,,400,200000,1,/
16,557,200000,,SB2322600000214,021000080000030,ACH Credit Receipt Return,Return To: Test, Remittance Info: "SB2322300000052", Entry Description: EXP; -, SEC: CCD, Reason: "R02", Return of Client Ref ID: 021000080000030, GS ID: SB2322600000214,CREF: 026015076104300,IDNM: 1114,EREF: 021000080000030,CRNM: Test,DBNM: SAMPLE INC.,CABA: 021000089/
16,451,55555,,SB2322600000455,021000020000021,ACH Debit Payment,To: TEST, Entry Description: INVOICES; 210630, SEC: CCD, Client Ref ID: 021000020000021, GS ID: SB2322600000455,CREF: 021000020000021,IDNM: 2009282,EREF: 021000020000021,CRNM: TEST,DBNM: SAMPLE INC,CABA: 021000021/
16,266,1912,,GI2118700002010,20210706MMQFMPU8000001,Outgoing Wire Return,-,CREF: 20210706MMQFMPU8000001,EREF: 20210706MMQFMPU8000001,DBIC: GSCRUS33,CRNM: ABC Company,DBNM: SAMPLE INC./
16,495,50500,,GI2321400000090,GSV0DL6RKT,Outgoing Wire,To: TEST COMPANY, Remittance Info: "QWERTIOP", Client Ref ID: GSV0DL6RKT, GS ID: GI2321400000090, Settled Amt: EUR 322.00, FX Rate: 156.833677,REMI: QWERTIOP,EREF: GSV0DL6RKT,CBIC: COBADEFF,CRNM: TEST COMPANY,DBNM: SAMPLE TEST/
16,195,1125,,GI2229300000187,GS0D9VGMP1IWPLW,Incoming Wire,-,EREF: GS0D9VGMP1IWPLW,DBIC: CITIUS30XXX,CRNM: ABC CORPORATION,DACT: 8348572423,CHKN: GSIL2X6103UNCRSF/
16,257,60000,,SB2225800001203,028000020000335,ACH Debit Payment Return,Return From: Company1, Entry Description: TRADE; -, SEC: CCD, Reason: "R02", Return of Client Ref ID: 028000020000335, GS ID: SB2225800001203,IDNM: 1,EREF: 028000020000335,CRNM: TEST INC,DBNM: Company1,DABA: 028000024/
16,255,931,,SC2134800001999,,Check Return,Return From: Test2 Customer, Check Serial Number: 0009000000, Return Reason: "Payee does not exist", Client Ref ID: 74564762445, GS ID: SC213480000120999,EREF: 07370568132,CRNM: Test Inc.,DBNM: Test2 Customer,CABA: 12345,CHKN: 0009000000/
16,195,50050,,GI2228400005800,RTR60880840833,RTP Incoming,From: SAMPLE INC, Remittance Info: "Test Remittance", Client Ref ID: RTR60880840833, GS ID: GI2228400005800, Clearing Ref: 001,REMI: Test Remittance,EREF: RTR60880840833,CRNM: RTR-CdtrName,DBNM: SAMPLE INC,DACT: 02122056789012205,DABA: 000000010/
16,175,527,,SX22293073766088,GS4N04L1COP45VY,Check Deposit,-,EREF: GS4N04L1COP45VY,DACT: 100168723/
16,475,10100,,SC2229300000152,01030340329,Check Paid,-,REMI: UAT testing for Checks,EREF: 01030340329,CRNM: TEST INC,DBNM: ABC CORP,CABA: 12345,CHKN: 006034594478/
16,275,337686,,GI2318000014342,e457328416d411eeaf020a58a9feac02,Cash Concentration,From: SAMPLE INC, Account: 290000020437, GS Cash Concentration, "Structure ID: CC0000000", GS ID: GI2318000212121,REMI: Structure ID: CC0000082,EREF: e123456786d411eeaf020a58a9feac02,DBIC: GSCRUS33VIA,CRNM: SAMPLE INC,DBNM: SAMPLE INC,DACT: 290000020437/
16,165,5000,,SPB2321284264201,AB-GS-DDFILEAB0001-DDBAB0001,ACH Debit Collection,Entry Description: BILL PMT; -, SEC: CCD, Client Ref ID: AB-GS-DDFILEAB0001-DDBAB0001, GS ID: SPB2321284264201,EREF: AB-GS-DDFILEAB0001-DDBAB0001,CRNM: SAMPLE LLP,DACT: ACHCONTROLINUSD01/
16,475,44250,,SC2323300002416,8ce1829175a74ec88d67010dd7fb6132,Check Paid,To: TEST AND COMPANY LLC, Check Serial Number: 24108, GS ID: SC2323300002416,EREF: 8ce1829175a74ec88d67010dd7fb6132,CRNM: TEST AND COMPANY LLC,DBNM: Sample Inc.,CABA: 0,CHKN: 24108/
16,495,30000000,,GI2323300009168,3785726,Outgoing Wire,To: TEST AND COMPANY, Remittance Info: "081823 Invoice - Sample", Client Ref ID: 3785726, GS ID: GI2323300009168, Clearing Ref: 20230821MMQFMPU7004100,CREF: 20230821MMQFMPU7004100,REMI: 081823 Invoice - Sample,EREF: 3785726,CRNM: TEST AND COMPANY,DBNM: Sample Inc.,CACT: 609873838,CABA: 021000021/
49,6869722,8/
03,260000033037,USD,,,,,060,000,,,100,000,,,400,000,,/
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
	assert.Equal(suite.T(), recorder.Body.String(), `{"sender":"GSBI","receiver":"cont001","fileCreatedDate":"210706","fileCreatedTime":"1249","fileIdNumber":"1","versionNumber":2,"fileControlTotal":"13060195162","numberOfGroups":1,"numberOfRecords":18,"Groups":[{"receiver":"cont001","originator":"026015079","groupStatus":1,"asOfDate":"230906","asOfTime":"2000","groupControlTotal":"13060195162","numberOfAccounts":4,"numberOfRecords":16,"Accounts":[{"accountNumber":"107049924","currencyCode":"USD","summaries":[{"TypeCode":"","Amount":"","ItemCount":0,"FundsType":{}},{"TypeCode":"060","Amount":"13053325440","ItemCount":0,"FundsType":{}},{"TypeCode":"100","Amount":"000","ItemCount":0,"FundsType":{}},{"TypeCode":"400","Amount":"000","ItemCount":0,"FundsType":{}}],"accountControlTotal":"13053325440","numberRecords":2,"Details":null},{"accountNumber":"107049932","currencyCode":"USD","summaries":[{"TypeCode":"","Amount":"","ItemCount":0,"FundsType":{}},{"TypeCode":"060","Amount":"6865898","ItemCount":0,"FundsType":{}},{"TypeCode":"100","Amount":"1912","ItemCount":1,"FundsType":{}},{"TypeCode":"400","Amount":"000","ItemCount":0,"FundsType":{}}],"accountControlTotal":"-1260161341762","numberRecords":26,"Details":[{"TypeCode":"447","Amount":"60000","FundsType":{},"BankReferenceNumber":"SPB2322984714570","CustomerReferenceNumber":"1111","Text":"ACH Credit Payment,Entry Description: EXP; -, SEC: CCD, Client Ref ID: 1111, GS ID: SPB2322984714570,EREF: 1111,DBNM: TEST INC,CACT: ACHCONTROLOUTUSD01"},{"TypeCode":"261","Amount":"143500","FundsType":{},"BankReferenceNumber":"SB2322600000404","CustomerReferenceNumber":"GSQ4FBGFDGWGKY","Text":"ACH Credit Reject,From: TEST INC, Remittance Info: \"ACH- Test - Addenda Record\", Entry Description: TRADE; -, SEC: CTX, Client Ref ID: GSQ4FBGFDGWGKY, GS ID: SB2322600000404,CREF: ,REMI: ACH- Test - Addenda Record,EREF: GSQ4FBGFDGWGKY,CRNM: Test,DBNM: SAMPLE INC,DACT: 101152046,DABA: 026015079"},{"TypeCode":"447","Amount":"928650","FundsType":{},"BankReferenceNumber":"SPB2322684598521","CustomerReferenceNumber":"AB-GS-RPFILERP0001-RPBA0001","Text":"ACH Credit Payment,Entry Description: TRADE; -, SEC: CTX, Client Ref ID: AB-GS-TEST0001-RPBA0001, GS ID: SPB2322684598521,EREF: AB-GS-RPFILERP0001-RPBA0001,DBNM: SAMPLE INC,CACT: ACHCONTROLOUTUSD01"}]},{"accountNumber":"104108339","currencyCode":"USD","summaries":[{"TypeCode":"010","Amount":"159581194","ItemCount":0,"FundsType":{}},{"TypeCode":"015","Amount":"159381194","ItemCount":0,"FundsType":{}},{"TypeCode":"040","Amount":"158568897","ItemCount":0,"FundsType":{}},{"TypeCode":"045","Amount":"158368897","ItemCount":0,"FundsType":{}},{"TypeCode":"100","Amount":"000","ItemCount":0,"FundsType":{}},{"TypeCode":"400","Amount":"200000","ItemCount":1,"FundsType":{}}],"accountControlTotal":"6869722","numberRecords":8,"Details":[{"TypeCode":"557","Amount":"200000","FundsType":{},"BankReferenceNumber":"SB2322600000214","CustomerReferenceNumber":"021000080000030","Text":"ACH Credit Receipt Return,Return To: Test, Remittance Info: \"SB2322300000052\", Entry Description: EXP; -, SEC: CCD, Reason: \"R02\", Return of Client Ref ID: 021000080000030, GS ID: SB2322600000214,CREF: 026015076104300,IDNM: 1114,EREF: 021000080000030,CRNM: Test,DBNM: SAMPLE INC.,CABA: 021000089"},{"TypeCode":"451","Amount":"55555","FundsType":{},"BankReferenceNumber":"SB2322600000455","CustomerReferenceNumber":"021000020000021","Text":"ACH Debit Payment,To: TEST, Entry Description: INVOICES; 210630, SEC: CCD, Client Ref ID: 021000020000021, GS ID: SB2322600000455,CREF: 021000020000021,IDNM: 2009282,EREF: 021000020000021,CRNM: TEST,DBNM: SAMPLE INC,CABA: 021000021"},{"TypeCode":"266","Amount":"1912","FundsType":{},"BankReferenceNumber":"GI2118700002010","CustomerReferenceNumber":"20210706MMQFMPU8000001","Text":"Outgoing Wire Return,-,CREF: 20210706MMQFMPU8000001,EREF: 20210706MMQFMPU8000001,DBIC: GSCRUS33,CRNM: ABC Company,DBNM: SAMPLE INC."},{"TypeCode":"495","Amount":"50500","FundsType":{},"BankReferenceNumber":"GI2321400000090","CustomerReferenceNumber":"GSV0DL6RKT","Text":"Outgoing Wire,To: TEST COMPANY, Remittance Info: \"QWERTIOP\", Client Ref ID: GSV0DL6RKT, GS ID: GI2321400000090, Settled Amt: EUR 322.00, FX Rate: 156.833677,REMI: QWERTIOP,EREF: GSV0DL6RKT,CBIC: COBADEFF,CRNM: TEST COMPANY,DBNM: SAMPLE TEST"},{"TypeCode":"195","Amount":"1125","FundsType":{},"BankReferenceNumber":"GI2229300000187","CustomerReferenceNumber":"GS0D9VGMP1IWPLW","Text":"Incoming Wire,-,EREF: GS0D9VGMP1IWPLW,DBIC: CITIUS30XXX,CRNM: ABC CORPORATION,DACT: 8348572423,CHKN: GSIL2X6103UNCRSF"},{"TypeCode":"257","Amount":"60000","FundsType":{},"BankReferenceNumber":"SB2225800001203","CustomerReferenceNumber":"028000020000335","Text":"ACH Debit Payment Return,Return From: Company1, Entry Description: TRADE; -, SEC: CCD, Reason: \"R02\", Return of Client Ref ID: 028000020000335, GS ID: SB2225800001203,IDNM: 1,EREF: 028000020000335,CRNM: TEST INC,DBNM: Company1,DABA: 028000024"},{"TypeCode":"255","Amount":"931","FundsType":{},"BankReferenceNumber":"SC2134800001999","CustomerReferenceNumber":"","Text":"Check Return,Return From: Test2 Customer, Check Serial Number: 0009000000, Return Reason: \"Payee does not exist\", Client Ref ID: 74564762445, GS ID: SC213480000120999,EREF: 07370568132,CRNM: Test Inc.,DBNM: Test2 Customer,CABA: 12345,CHKN: 0009000000"},{"TypeCode":"195","Amount":"50050","FundsType":{},"BankReferenceNumber":"GI2228400005800","CustomerReferenceNumber":"RTR60880840833","Text":"RTP Incoming,From: SAMPLE INC, Remittance Info: \"Test Remittance\", Client Ref ID: RTR60880840833, GS ID: GI2228400005800, Clearing Ref: 001,REMI: Test Remittance,EREF: RTR60880840833,CRNM: RTR-CdtrName,DBNM: SAMPLE INC,DACT: 02122056789012205,DABA: 000000010"},{"TypeCode":"175","Amount":"527","FundsType":{},"BankReferenceNumber":"SX22293073766088","CustomerReferenceNumber":"GS4N04L1COP45VY","Text":"Check Deposit,-,EREF: GS4N04L1COP45VY,DACT: 100168723"},{"TypeCode":"475","Amount":"10100","FundsType":{},"BankReferenceNumber":"SC2229300000152","CustomerReferenceNumber":"01030340329","Text":"Check Paid,-,REMI: UAT testing for Checks,EREF: 01030340329,CRNM: TEST INC,DBNM: ABC CORP,CABA: 12345,CHKN: 006034594478"},{"TypeCode":"275","Amount":"337686","FundsType":{},"BankReferenceNumber":"GI2318000014342","CustomerReferenceNumber":"e457328416d411eeaf020a58a9feac02","Text":"Cash Concentration,From: SAMPLE INC, Account: 290000020437, GS Cash Concentration, \"Structure ID: CC0000000\", GS ID: GI2318000212121,REMI: Structure ID: CC0000082,EREF: e123456786d411eeaf020a58a9feac02,DBIC: GSCRUS33VIA,CRNM: SAMPLE INC,DBNM: SAMPLE INC,DACT: 290000020437"},{"TypeCode":"165","Amount":"5000","FundsType":{},"BankReferenceNumber":"SPB2321284264201","CustomerReferenceNumber":"AB-GS-DDFILEAB0001-DDBAB0001","Text":"ACH Debit Collection,Entry Description: BILL PMT; -, SEC: CCD, Client Ref ID: AB-GS-DDFILEAB0001-DDBAB0001, GS ID: SPB2321284264201,EREF: AB-GS-DDFILEAB0001-DDBAB0001,CRNM: SAMPLE LLP,DACT: ACHCONTROLINUSD01"},{"TypeCode":"475","Amount":"44250","FundsType":{},"BankReferenceNumber":"SC2323300002416","CustomerReferenceNumber":"8ce1829175a74ec88d67010dd7fb6132","Text":"Check Paid,To: TEST AND COMPANY LLC, Check Serial Number: 24108, GS ID: SC2323300002416,EREF: 8ce1829175a74ec88d67010dd7fb6132,CRNM: TEST AND COMPANY LLC,DBNM: Sample Inc.,CABA: 0,CHKN: 24108"},{"TypeCode":"495","Amount":"30000000","FundsType":{},"BankReferenceNumber":"GI2323300009168","CustomerReferenceNumber":"3785726","Text":"Outgoing Wire,To: TEST AND COMPANY, Remittance Info: \"081823 Invoice - Sample\", Client Ref ID: 3785726, GS ID: GI2323300009168, Clearing Ref: 20230821MMQFMPU7004100,CREF: 20230821MMQFMPU7004100,REMI: 081823 Invoice - Sample,EREF: 3785726,CRNM: TEST AND COMPANY,DBNM: Sample Inc.,CACT: 609873838,CABA: 021000021"}]},{"accountNumber":"260000033037","currencyCode":"USD","summaries":[{"TypeCode":"","Amount":"","ItemCount":0,"FundsType":{}},{"TypeCode":"060","Amount":"000","ItemCount":0,"FundsType":{}},{"TypeCode":"100","Amount":"000","ItemCount":0,"FundsType":{}},{"TypeCode":"400","Amount":"000","ItemCount":0,"FundsType":{}}],"accountControlTotal":"000","numberRecords":2,"Details":null},{"accountNumber":"280000010657","currencyCode":"USD","summaries":[{"TypeCode":"","Amount":"","ItemCount":0,"FundsType":{}},{"TypeCode":"060","Amount":"000","ItemCount":0,"FundsType":{}},{"TypeCode":"100","Amount":"000","ItemCount":0,"FundsType":{}},{"TypeCode":"400","Amount":"000","ItemCount":0,"FundsType":{}}],"accountControlTotal":"000","numberRecords":2,"Details":null}]}]}
`)
	file := client.NewNullableFile(nil)
	err = file.UnmarshalJSON(recorder.Body.Bytes())
	assert.Equal(suite.T(), nil, err)
	groups := file.Get().GetGroups()
	assert.Equal(suite.T(), len(groups), 1)
	group := groups[0]
	accounts := group.GetAccounts()
	assert.Equal(suite.T(), len(accounts), 5)

	for i := 0; i < len(accounts); i++ {
		if i == 1 {
			assert.Equal(suite.T(), len(accounts[i].GetDetails()), 3)
		} else if i == 2 {
			assert.Equal(suite.T(), len(accounts[i].GetDetails()), 14)
		} else {
			assert.Equal(suite.T(), len(accounts[i].GetDetails()), 0)
		}
	}
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

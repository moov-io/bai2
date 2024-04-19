package v2

import (
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestReader(t *testing.T) {
	report, err := ReadFilepath(filepath.Join("testdata", "example.txt"))
	require.NoError(t, err)
	require.NotNil(t, report)

	fh := FileHeader{
		RecordCode:       "01",
		SenderName:       "122099999",
		SenderID:         "123456789",
		FileCreationDate: "040621",
		FileCreationTime: "0200",
		FileIDModifier:   "1",
	}
	require.Equal(t, fh, report.FileHeader)

	require.Len(t, report.Groups, 1)
	g := report.Groups[0]

	require.Len(t, g.AccountIdentifiers, 1)
	require.Len(t, g.TransactionDetails, 2)
	require.Len(t, g.ContinuationRecords, 2)
	require.Len(t, g.AccountTrailers, 1)

	gh := GroupHeader{
		RecordCode:         "02",
		UltimateReceiverID: "031001234",
		OriginatorID:       "122099999",
		GroupStatus:        "1",
		AsOfDate:           "040620",
		AsOfTime:           "2359",
		CurrencyCode:       "USD",
		AsOfDateModifier:   "2",
	}
	require.Equal(t, gh, g.Header)

	// AccountIdentifier
	ai := AccountIdentifier{
		RecordCode:            "03",
		CustomerAccountNumber: "0975312468",
		CurrencyCode:          "USD",
		TypeCode:              "010",
		Amount:                500000,
		ItemCount:             0,
		FundsType:             "",
	}
	require.Equal(t, ai, g.AccountIdentifiers[0])

	// TODO(adam):
	// Data in this record are for the sending bank’s account number (0975312468). The leading
	// zero in the account number is significant and must be included in the data. The optional
	// currency code is defaulted to the group currency code. The amount for type code (010) is
	// $5,000.00 (500000). The Item Count and Funds Type fields are defaulted to “unknown” as
	// indicated by adjacent delimiters (,,,). The amount for type code (190) is $700,000.00
	// (70000000). The item count for this amount is four (4) and the availability is immediate (0).

	// TransactionDetails
	td1 := TransactionDetail{RecordCode: "16", TypeCode: "165", Amount: 1500000, FundsType: "1", DetailText: "DD1620"}
	require.Equal(t, td1, g.TransactionDetails[0])

	// TODO(adam):
	// This is a Detail Record (16). The amount for type code 165 is $15,000.00 (1500000) and has
	// one-day (1) deferred availability (1). The bank reference number is (DD1620). There is no
	// customer reference number (,,). The text is (DEALER PAYMENTS). The remainder of the
	// field is blank filled if fixed length records are used, and the text field is delimited by the fact
	// that the next record is not “88”.

	td2 := TransactionDetail{RecordCode: "16", TypeCode: "115", Amount: 1e+07, FundsType: "S", DetailText: "5000000"}
	require.Equal(t, td2, g.TransactionDetails[1])

	// ContinuationRecords
	// TODO(adam): need to parse into previous TransactionDetail record
	cr1 := ContinuationRecord{RecordCode: "88", ContinuationText: "AX13612,B096132,AMALGAMATED CORP. LOCKBOX"}
	require.Equal(t, cr1, g.ContinuationRecords[0])
	// TODO(adam):
	cr2 := ContinuationRecord{RecordCode: "88", ContinuationText: "DEPOSIT-MISC. RECEIVABLES"}
	require.Equal(t, cr2, g.ContinuationRecords[1])

	// AccountTrailers
	at := AccountTrailer{RecordCode: "49", ControlTotal: 18650000, NumberOfRecords: 3}
	require.Equal(t, at, g.AccountTrailers[0])

	// GroupTrailer
	gt := GroupTrailer{
		RecordCode:        "98",
		GroupControlTotal: 11800000,
		NumberOfAccounts:  2,
		NumberOfRecords:   6,
	}
	require.Equal(t, gt, g.Trailer)

	ft := FileTrailer{
		RecordCode:      "99",
		ControlTotal:    1215450000,
		NumberOfGroups:  4,
		NumberOfRecords: 36,
	}
	require.Equal(t, ft, report.FileTrailer)
}

func TestParseAccountIdentifier(t *testing.T) {
	// Define test cases in a slice of structs.
	testCases := []struct {
		name    string
		line    string
		want    AccountIdentifier
		wantErr bool
	}{
		{
			name:    "Valid complete data",
			line:    "03,100123456789,USD,11,1200.50,25,A",
			want:    AccountIdentifier{"03", "100123456789", "USD", "11", 1200.50, 25, "A"},
			wantErr: false,
		},
		{
			name:    "Valid data with optional fields missing",
			line:    "03,100123456790,GBP,22,500.00,/,/",
			want:    AccountIdentifier{"03", "100123456790", "GBP", "22", 500.00, 0, ""},
			wantErr: false,
		},
		{
			name:    "Invalid data with missing fields",
			line:    "03,100123456791",
			want:    AccountIdentifier{},
			wantErr: true,
		},
	}

	// Iterate through each test case.
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := parseAccountIdentifier(tc.line)
			// Check for error expectation.
			if (err != nil) != tc.wantErr {
				t.Errorf("parseAccountIdentifier() error = %v, wantErr %v", err, tc.wantErr)
				return
			}
			// Compare the result struct only if no error is expected.
			if !tc.wantErr && got != tc.want {
				t.Errorf("parseAccountIdentifier() got = %v, want %v", got, tc.want)
			}
		})
	}
}

func TestParseTransactionDetail(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected TransactionDetail
		wantErr  bool
	}{
		{
			name:     "valid transaction detail without type 890",
			input:    "16,500,1000.00,C,Invoice Payment",
			expected: TransactionDetail{"16", "500", 1000.00, "C", "Invoice Payment"},
			wantErr:  false,
		},
		{
			name:     "type 890 with no amount needed",
			input:    "16,890,,C,Non-financial adjustment notice",
			expected: TransactionDetail{"16", "890", 0, "C", "Non-financial adjustment notice"},
			wantErr:  false,
		},
		{
			name:     "Regular transaction detail",
			input:    "16,500,1500.00,C,Payment for services/",
			expected: TransactionDetail{"16", "500", 1500.00, "C", "Payment for services/"},
			wantErr:  false,
		},
		{
			name:     "Type 890 transaction detail with leading slash",
			input:    "16,890,,,,/delayed report update",
			expected: TransactionDetail{"16", "890", 0, "", ""},
			wantErr:  false,
		},
		{
			name:    "incorrect format with invalid amount",
			input:   "16,500,one thousand,C,Payment",
			wantErr: true,
		},
		{
			name:    "missing fields",
			input:   "16,500",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := parseTransactionDetail(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("parseTransactionDetail() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				require.Equal(t, tt.expected, result)

			}
		})
	}
}

func TestTransactionAndContinuationParsing(t *testing.T) {
	tests := []struct {
		name        string
		lineTD      string
		lineCR      string
		expectedTD  TransactionDetail
		expectedCR  ContinuationRecord
		expectErrTD bool
		expectErrCR bool
	}{
		{
			name:        "Valid Transaction and Continuation",
			lineTD:      "16,500,1500.00,C,Payment for services/",
			lineCR:      "88, This is continued text from transaction detail/",
			expectedTD:  TransactionDetail{"16", "500", 1500.00, "C", "Payment for services/"},
			expectedCR:  ContinuationRecord{"88", " This is continued text from transaction detail/"},
			expectErrTD: false,
			expectErrCR: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			td, errTD := parseTransactionDetail(tt.lineTD)
			if (errTD != nil) != tt.expectErrTD {
				t.Errorf("TestTransactionAndContinuationParsing(): errTD %v, expectErrTD %v", errTD, tt.expectErrTD)
				return
			}
			if !tt.expectErrTD && td != tt.expectedTD {
				t.Errorf("TestTransactionAndContinuationParsing(): got%+v, want %+v", td, tt.expectedTD)
			}

			cr, errCR := parseContinuationRecord(tt.lineCR)
			if (errCR != nil) != tt.expectErrCR {
				t.Errorf("TestTransactionAndContinuationParsing(): errCR %v, expectErrCR %v", errCR, tt.expectErrCR)
				return
			}
			if !tt.expectErrCR && cr != tt.expectedCR {
				t.Errorf("TestTransactionAndContinuationParsing(): got%+v, want %+v", cr, tt.expectedCR)
			}
		})
	}
}

func TestTransactionDetailValidation(t *testing.T) {
	tests := []struct {
		td       TransactionDetail
		expected bool
	}{
		{TransactionDetail{RecordCode: "16", TypeCode: "200", Amount: 150.00, FundsType: "C", DetailText: "Adjustment for account"}, true},
		{TransactionDetail{RecordCode: "16", TypeCode: "500", Amount: 100.00, FundsType: "D", DetailText: "Unknown code test"}, false},
	}

	for _, test := range tests {
		err := test.td.Validate()
		if (err == nil) != test.expected {
			t.Errorf("For TransactionDetail with type code %s, expected validity %v, but got %v", test.td.TypeCode, test.expected, err == nil)
		}
	}
}

func TestParseFileTrailer(t *testing.T) {
	tests := []struct {
		name        string
		line        string
		expected    FileTrailer
		expectError bool
	}{
		{
			name:        "Positive Control Total",
			line:        "99,+12345.67,100,10",
			expected:    FileTrailer{"99", 12345.67, 100, 10},
			expectError: false,
		},
		{
			name:        "Negative Control Total",
			line:        "99,-12345.67,100,10",
			expected:    FileTrailer{"99", -12345.67, 100, 10},
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual, err := parseFileTrailer(tt.line)
			if (err != nil) != tt.expectError {
				t.Errorf("%s: parseFileTrailer() error = %v, expectError %v", tt.name, err, tt.expectError)
			}
			if !tt.expectError && (actual != tt.expected) {
				t.Errorf("%s: parseFileTrailer() got = %+v, want %+v", tt.name, actual, tt.expected)
			}
		})
	}
}

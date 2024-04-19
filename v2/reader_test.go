package v2

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestParseAccountIdentifier(t *testing.T) {
	tests := []struct {
		name        string
		line        string
		want        AccountIdentifier
		expectError bool
	}{
		{
			name: "Valid data with explicit number of items",
			line: "03,12345,USD,10",
			want: AccountIdentifier{
				RecordCode: "03", CustomerAccountNumber: "12345", CurrencyCode: "USD", NumberOfItems: 10,
			},
			expectError: false,
		},
		{
			name: "Data with default number of items",
			line: "03,12345,USD,,",
			want: AccountIdentifier{
				RecordCode: "03", CustomerAccountNumber: "12345", CurrencyCode: "USD", NumberOfItems: 0,
			},
			expectError: false,
		},
		{
			name:        "Insufficient data fields",
			line:        "03,12345,USD",
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parseAccountIdentifier(tt.line)
			if (err != nil) != tt.expectError {
				t.Errorf("parseAccountIdentifier() error = %v, expectError %v", err, tt.expectError)
				return
			}
			if !tt.expectError && got != tt.want {
				t.Errorf("parseAccountIdentifier() got = %v, want %v", got, tt.want)
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
			line:        "99,+12345.67,100",
			expected:    FileTrailer{"99", 12345.67, 100},
			expectError: false,
		},
		{
			name:        "Negative Control Total",
			line:        "99,-12345.67,100",
			expected:    FileTrailer{"99", -12345.67, 100},
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

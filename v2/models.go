package v2

import (
	"errors"
	"fmt"
)

// Define constants for all record types
const (
	RecordTypeFileHeader        = "01"
	RecordTypeGroupHeader       = "02"
	RecordTypeAccountIdentifier = "03"
	RecordTypeTransactionDetail = "16"
	RecordTypeContinuation      = "88"
	RecordTypeAccountTrailer    = "49"
	RecordTypeGroupTrailer      = "98"
	RecordTypeFileTrailer       = "99"
)

var (
	ErrInvalidRecordCode = errors.New("invalid record code for the given record type")
	ErrInvalidData       = errors.New("provided data is invalid for the given record type")
)

// A BalanceReport encapsulates an entire set of balance report data,
// including groups that contain their relevant transactions.
type BalanceReport struct {
	FileHeader  FileHeader
	Groups      []GroupHeader
	FileTrailer FileTrailer
}

// FileHeader represents the file header details in any transmitted file
type FileHeader struct {
	RecordCode       string
	SenderName       string
	SenderID         string
	FileCreationDate string
	FileCreationTime string
	FileIDModifier   string
}

func (fh FileHeader) Validate() error {
	if fh.RecordCode != RecordTypeFileHeader {
		return ErrInvalidRecordCode
	}
	return nil
}

func (fh FileHeader) Serialize() string {
	return fmt.Sprintf("%2s%-35s%-15s%-8s%-4s%-1s",
		fh.RecordCode,
		fh.SenderName,
		fh.SenderID,
		fh.FileCreationDate,
		fh.FileCreationTime,
		fh.FileIDModifier)
}

// GroupHeader represents the grouping of multiple account transactions
type GroupHeader struct {
	RecordCode         string
	UltimateReceiverID string
	OriginatorID       string
	GroupStatus        string
	AsOfDate           string
	AsOfTime           string // Optional
	CurrencyCode       string // Optional
	AsOfDateModifier   string // Optional

	AccountIdentifiers  []AccountIdentifier
	TransactionDetails  []TransactionDetail
	ContinuationRecords []ContinuationRecord
	AccountTrailers     []AccountTrailer
}

func (gh GroupHeader) Validate() error {
	if gh.RecordCode != RecordTypeGroupHeader {
		return ErrInvalidRecordCode
	}
	return nil
}

func (gh GroupHeader) Serialize() string {
	return fmt.Sprintf("%2s%-15s%-15s%-1s%-6s%-4s%-3s%-1s",
		gh.RecordCode,
		gh.UltimateReceiverID,
		gh.OriginatorID,
		gh.GroupStatus,
		gh.AsOfDate,
		gh.AsOfTime,
		gh.CurrencyCode,
		gh.AsOfDateModifier)
}

// AccountIdentifier represents account details
type AccountIdentifier struct {
	RecordCode            string
	CustomerAccountNumber string
	CurrencyCode          string
	NumberOfItems         int // Field to reflect the number of items in an activity summary.
}

func (ai AccountIdentifier) Validate() error {
	if ai.RecordCode != RecordTypeAccountIdentifier {
		return ErrInvalidRecordCode
	}

	// Check if CurrencyCode (here considered as a Type Code) follows specific rules
	validTypeCodes := map[string]bool{
		"USD": true, // Supposed valid type codes for demonstration
		"EUR": true,
		"JPY": true,
	} // TODO(adam): reference moov-io/iso4712

	if _, ok := validTypeCodes[ai.CurrencyCode]; !ok {
		return fmt.Errorf("invalid type code '%s' for account identifier", ai.CurrencyCode)
	}

	// Additional basic validations could be added here
	if len(ai.CustomerAccountNumber) == 0 {
		return errors.New("customer account number cannot be empty")
	}

	return nil
}

func (ai AccountIdentifier) Serialize() string {
	return fmt.Sprintf("%2s%-34s%-3s",
		ai.RecordCode,
		ai.CustomerAccountNumber,
		ai.CurrencyCode)
}

// TransactionDetail represents details of each transaction
type TransactionDetail struct {
	RecordCode string
	TypeCode   string
	Amount     float64
	FundsType  string
	DetailText string // Attention to initial character restrictions and continuation
}

func (td TransactionDetail) Validate() error {
	if td.RecordCode != RecordTypeTransactionDetail {
		return ErrInvalidRecordCode
	}
	if !td.isValidTypeCode() {
		return fmt.Errorf("invalid type code '%s' for transaction detail", td.TypeCode)
	}
	if td.Amount < 0 {
		return errors.New("transaction amount cannot be negative")
	}
	return nil
}

// isValidTypeCodeForTransactionDetail checks if the given type code is valid
func (td TransactionDetail) isValidTypeCode() bool {
	validTypeCodes := map[string]bool{
		"100": true, // Normal transaction
		"200": true, // Adjustment transaction
		"300": true, // Fee-related transaction
		"400": true, // Interest posting
		"890": true, // Non-monetary transaction
	}

	_, valid := validTypeCodes[td.TypeCode]
	return valid
}

func (td TransactionDetail) Serialize() string {
	return fmt.Sprintf("%2s%-6s%-12.2f%-1s%-34s",
		td.RecordCode,
		td.TypeCode,
		td.Amount,
		td.FundsType,
		td.DetailText)
}

// ContinuationRecord captures continuation details possibly to be appended to transactions
type ContinuationRecord struct {
	RecordCode       string
	ContinuationText string
}

func (cr ContinuationRecord) Validate() error {
	if cr.RecordCode != RecordTypeContinuation {
		return ErrInvalidRecordCode
	}
	return nil
}

func (cr ContinuationRecord) Serialize() string {
	return fmt.Sprintf("%2s%-80s",
		cr.RecordCode,
		cr.ContinuationText)
}

// AccountTrailer represents trailers for individual accounts
type AccountTrailer struct {
	RecordCode      string
	ControlTotal    float64
	NumberOfRecords int
}

func (at AccountTrailer) Validate() error {
	if at.RecordCode != RecordTypeAccountTrailer {
		return ErrInvalidRecordCode
	}
	return nil
}

func (at AccountTrailer) Serialize() string {
	return fmt.Sprintf("%2s%-18.2f%-6d",
		at.RecordCode,
		at.ControlTotal,
		at.NumberOfRecords)
}

// GroupTrailer aggregates the trailers for a group of accounts
type GroupTrailer struct {
	RecordCode        string
	GroupControlTotal float64
	NumberOfAccounts  int
	NumberOfRecords   int
}

func (gt GroupTrailer) Validate() error {
	if gt.RecordCode != RecordTypeGroupTrailer {
		return ErrInvalidRecordCode
	}
	return nil
}

func (gt GroupTrailer) Serialize() string {
	return fmt.Sprintf("%2s%-18.2f%-6d",
		gt.RecordCode,
		gt.GroupControlTotal,
		gt.NumberOfRecords)
}

// Interpretation of control totals with optional '+' or '-' signs in Trailer records
type FileTrailer struct {
	RecordCode      string
	ControlTotal    float64
	NumberOfRecords int
}

func (ft FileTrailer) Validate() error {
	if ft.RecordCode != RecordTypeFileTrailer {
		return ErrInvalidRecordCode
	}
	return nil
}

func (ft FileTrailer) Serialize() string {
	return fmt.Sprintf("%2s%-18.2f%-6d",
		ft.RecordCode,
		ft.ControlTotal,
		ft.NumberOfRecords)
}

package v2

import (
	"bufio"
	"cmp"
	"errors"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
)

// Error definitions for missing required fields or corruption checks
var (
	ErrInvalidRecordType    = errors.New("invalid or unknown record type")
	ErrInvalidDataFormat    = errors.New("data format error")
	ErrIntegrityCheckFailed = errors.New("file integrity check failed")
)

func ReadFilepath(path string) (*BalanceReport, error) {
	fd, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("opening %s failed: %v", path, err)
	}
	defer fd.Close()

	return FromReader(fd)
}

func FromReader(r io.ReadCloser) (*BalanceReport, error) {
	if r == nil {
		return nil, errors.New("nil Reader")
	}
	defer r.Close()

	scanner := bufio.NewScanner(r)
	var currentGroup *Group
	br := &BalanceReport{}
	var err error

	for scanner.Scan() {
		line := scanner.Text()
		recordType := line[:2]

		switch recordType {
		case RecordTypeFileHeader:
			br.FileHeader, err = parseFileHeader(line)
			if err != nil {
				return nil, fmt.Errorf("parsing file header: %w", err)
			}

		case RecordTypeGroupHeader:
			groupHeader, err := parseGroupHeader(line)
			if err != nil {
				return nil, fmt.Errorf("parsing group header: %w", err)
			}
			currentGroup = &Group{
				Header: groupHeader,
			}

		case RecordTypeAccountIdentifier:
			if currentGroup == nil {
				return nil, errors.New("account identifier found outside of a group")
			}
			ai, err := parseAccountIdentifier(line)
			if err != nil {
				return nil, fmt.Errorf("parsing account identifier: %w", err)
			}
			currentGroup.AccountIdentifiers = append(currentGroup.AccountIdentifiers, ai)

		case RecordTypeTransactionDetail:
			if currentGroup == nil {
				return nil, errors.New("transaction detail found outside of a group")
			}
			td, err := parseTransactionDetail(line)
			if err != nil {
				return nil, fmt.Errorf("parsing transaction detail: %w", err)
			}
			currentGroup.TransactionDetails = append(currentGroup.TransactionDetails, td)

		case RecordTypeContinuation:
			if currentGroup == nil {
				return nil, errors.New("continuation record found outside of a group")
			}
			cr, err := parseContinuationRecord(line)
			if err != nil {
				return nil, fmt.Errorf("parsing continuation record: %w", err)
			}
			currentGroup.ContinuationRecords = append(currentGroup.ContinuationRecords, cr)

		case RecordTypeAccountTrailer:
			if currentGroup == nil {
				return nil, errors.New("account trailer found outside of a group")
			}
			at, err := parseAccountTrailer(line)
			if err != nil {
				return nil, fmt.Errorf("parsing account trailer: %w", err)
			}
			currentGroup.AccountTrailers = append(currentGroup.AccountTrailers, at)

		case RecordTypeGroupTrailer:
			if currentGroup == nil {
				return nil, errors.New("group trailer found outside of a group")
			}
			currentGroup.Trailer, err = parseGroupTrailer(line)
			if err != nil {
				return nil, fmt.Errorf("parsing group trailer: %w", err)
			}
			br.Groups = append(br.Groups, *currentGroup)
			currentGroup = nil

		case RecordTypeFileTrailer:
			br.FileTrailer, err = parseFileTrailer(line)
			if err != nil {
				return nil, fmt.Errorf("parsing file trailer: %w", err)
			}

		default:
			return nil, fmt.Errorf("unexpected record type: %s", recordType)
		}
	}

	if err := scanner.Err(); err != nil {
		return br, fmt.Errorf("error during scanning: %v", err)
	}

	return br, nil

}

func parseFileHeader(line string) (FileHeader, error) {
	parts := strings.Split(line, ",")
	if len(parts) < 6 {
		return FileHeader{}, ErrInvalidDataFormat // Ensure there are at least 6 parts.
	}

	return FileHeader{
		RecordCode:       strings.TrimSpace(parts[0]),
		SenderID:         strings.TrimSpace(parts[2]),
		SenderName:       strings.TrimSpace(parts[1]),
		FileCreationDate: strings.TrimSpace(parts[3]),
		FileCreationTime: strings.TrimSpace(parts[4]),
		FileIDModifier:   strings.TrimSpace(parts[5]),
	}, nil
}

func parseGroupHeader(line string) (GroupHeader, error) {
	line = strings.TrimSuffix(line, "/")

	parts := strings.Split(line, ",")

	// Ensure there are at least the minimum required fields.
	if len(parts) < 5 {
		return GroupHeader{}, ErrInvalidDataFormat
	}

	// Trim each part to remove unnecessary whitespaces.
	for i, part := range parts {
		parts[i] = strings.TrimSpace(part)
	}

	gh := GroupHeader{
		RecordCode:         parts[0],
		UltimateReceiverID: parts[1],
		OriginatorID:       parts[2],
		GroupStatus:        parts[3],
		AsOfDate:           parts[4],
	}

	// Handle optional fields
	if len(parts) > 5 && parts[5] != "/" {
		gh.AsOfTime = parts[5]
	}

	if len(parts) > 6 && parts[6] != "/" {
		gh.CurrencyCode = parts[6]
	}
	gh.CurrencyCode = cmp.Or(gh.CurrencyCode, "USD")

	if len(parts) > 7 && parts[7] != "/" {
		gh.AsOfDateModifier = parts[7]
	}

	return gh, nil
}

func parseAccountIdentifier(line string) (AccountIdentifier, error) {
	// Split and trim the line.
	parts := strings.Split(line, ",")
	if len(parts) < 5 { // Minimum parts to form a valid record
		return AccountIdentifier{}, ErrInvalidDataFormat
	}

	// Convert and assign each parts to struct fields, handling optional values.
	asum := AccountIdentifier{
		RecordCode:            strings.TrimSpace(parts[0]),
		CustomerAccountNumber: strings.TrimSpace(parts[1]),
		CurrencyCode:          strings.TrimSpace(parts[2]),
		TypeCode:              strings.TrimSpace(parts[3]),
	}
	asum.CurrencyCode = cmp.Or(asum.CurrencyCode, "USD")

	var err error

	// Parse amount which is a float64.
	parts[4] = strings.TrimSpace(parts[4])
	if parts[4] != "" {
		asum.Amount, err = strconv.ParseFloat(parts[4], 64)
		if err != nil {
			return AccountIdentifier{}, err // Handle amount conversion error.
		}
	}

	// Item Count, where applicable (check if part exists).
	if len(parts) > 5 && parts[5] != "/" {
		if parts[5] != "" {
			asum.ItemCount, err = strconv.Atoi(strings.TrimSpace(parts[5]))
			if err != nil {
				return AccountIdentifier{}, err // Handle item count conversion error.
			}
		}
	}

	// Funds Type, optional last part.
	if len(parts) > 6 && parts[6] != "/" {
		if parts[6] != "" {
			asum.FundsType = strings.TrimSpace(parts[6])
		}
	}

	return asum, nil
}

func parseTransactionDetail(line string) (TransactionDetail, error) {
	// Split the line by comma to read different fields
	fields := strings.Split(line, ",")
	if len(fields) < 4 {
		return TransactionDetail{}, errors.New("insufficient fields")
	}

	// Create initial TransactionDetail struct from split fields
	td := TransactionDetail{
		RecordCode: strings.TrimSpace(fields[0]),
		TypeCode:   strings.TrimSpace(fields[1]),
		FundsType:  strings.TrimSpace(fields[3]),
	}

	// Type code 890 often means not using the Amount field, handled as special case
	if td.TypeCode == "890" {
		if len(fields) > 4 {
			td.DetailText = strings.TrimSpace(fields[4])
		}
		return td, nil
	}

	// Parse Amount field converting string to float
	var err error
	td.Amount, err = strconv.ParseFloat(fields[2], 64)
	if err != nil {
		return TransactionDetail{}, fmt.Errorf("invalid amount format: %v", err)
	}

	// Details text is optional but may not be present; check field availability
	if len(fields) > 4 {
		td.DetailText = strings.TrimSpace(fields[4])
	}

	return td, nil
}

// ParseContinuationRecord parses continuation lines
func parseContinuationRecord(line string) (ContinuationRecord, error) {
	fields := strings.Split(line, ",")
	if len(fields) < 2 {
		return ContinuationRecord{}, fmt.Errorf("invalid continuation record: %s", line)
	}
	cr := ContinuationRecord{
		RecordCode:       fields[0],
		ContinuationText: strings.Join(fields[1:], ","),
	}
	return cr, nil
}

func parseAccountTrailer(line string) (AccountTrailer, error) {
	fields := strings.Split(line, ",")
	if len(fields) != 3 {
		return AccountTrailer{}, fmt.Errorf("invalid account trailer record: %s", line)
	}

	controlTotalStr := fields[1]
	signMultiplier := 1.0
	if strings.HasPrefix(controlTotalStr, "-") {
		signMultiplier = -1.0
		controlTotalStr = controlTotalStr[1:]
	} else if strings.HasPrefix(controlTotalStr, "+") {
		controlTotalStr = controlTotalStr[1:]
	}

	controlTotal, err := strconv.ParseFloat(controlTotalStr, 64)
	if err != nil {
		return AccountTrailer{}, fmt.Errorf("invalid control total: %s", fields[1])
	}

	numRecordsStr := strings.TrimSuffix(fields[2], "/")
	numRecords, err := strconv.Atoi(numRecordsStr)
	if err != nil {
		return AccountTrailer{}, fmt.Errorf("invalid number of records: %s", numRecordsStr)
	}

	at := AccountTrailer{
		RecordCode:      fields[0],
		ControlTotal:    controlTotal * signMultiplier,
		NumberOfRecords: numRecords,
	}
	return at, nil
}

func parseGroupTrailer(line string) (GroupTrailer, error) {
	// Trim the last delimiter if it exists.
	line = strings.TrimSuffix(line, "/")

	parts := strings.Split(line, ",")
	if len(parts) < 4 {
		return GroupTrailer{}, ErrInvalidDataFormat // Ensure there are at least 4 parts.
	}

	// Parse group control total.
	groupControlTotal, errControlTotal := strconv.ParseFloat(strings.TrimSpace(parts[1]), 64)
	if errControlTotal != nil {
		return GroupTrailer{}, errControlTotal
	}

	// Parse number of accounts.
	numAccounts, errNumAccounts := strconv.Atoi(strings.TrimSpace(parts[2]))
	if errNumAccounts != nil {
		return GroupTrailer{}, errNumAccounts
	}

	// Parse number of records, remove any trailing non-numeric characters.
	numRecords, errNumRecords := strconv.Atoi(strings.TrimSpace(parts[3]))
	if errNumRecords != nil {
		return GroupTrailer{}, errNumRecords
	}

	return GroupTrailer{
		RecordCode:        strings.TrimSpace(parts[0]),
		GroupControlTotal: groupControlTotal,
		NumberOfAccounts:  numAccounts,
		NumberOfRecords:   numRecords,
	}, nil
}

// parseFileTrailer translates the file trailer line while validating the number of records and control total
func parseFileTrailer(line string) (FileTrailer, error) {
	fields := strings.Split(line, ",")
	if len(fields) < 4 {
		return FileTrailer{}, ErrInvalidDataFormat
	}

	ft := FileTrailer{
		RecordCode: fields[0],
		NumberOfGroups: func(s string) int {
			if val, err := strconv.Atoi(s); err == nil {
				return val
			}
			return 0
		}(fields[2]), // TODO(adam): why so complex?
	}

	var err error

	// Handle possible sign in the number
	controlTotal := fields[1]
	if controlTotal == "/" {
		ft.ControlTotal = 0
	} else {
		sign := 1.0
		if strings.HasPrefix(controlTotal, "-") {
			sign = -1.0
			controlTotal = controlTotal[1:]
		} else if strings.HasPrefix(controlTotal, "+") {
			controlTotal = controlTotal[1:]
		}
		value, err := strconv.ParseFloat(controlTotal, 64)
		if err != nil {
			return ft, err
		}
		ft.ControlTotal = value * sign
	}

	var recordCount int64
	fields[3] = strings.TrimSuffix(fields[3], "/")
	if fields[3] != "" {
		recordCount, err = strconv.ParseInt(fields[3], 10, 32)
		if err != nil {
			return FileTrailer{}, fmt.Errorf("number of records: %w", err)
		}
		ft.NumberOfRecords = int(recordCount)
	}

	return ft, nil
}

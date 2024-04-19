package v2

import (
	"bufio"
	"errors"
	"fmt"
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

// LoadFromFile processes each line of the file to correctly structure the balance report.
func LoadFromFile(filename string) (*BalanceReport, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %v", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	var currentGroup *GroupHeader
	br := &BalanceReport{}

	for scanner.Scan() {
		line := scanner.Text()
		recordType := line[:2]

		switch recordType {
		case RecordTypeFileHeader:
			br.FileHeader, err = parseFileHeader(line)
			if err != nil {
				return nil, err
			}

		case RecordTypeGroupHeader:
			group, err := parseGroupHeader(line)
			if err != nil {
				return nil, err
			}
			br.Groups = append(br.Groups, group)
			currentGroup = &br.Groups[len(br.Groups)-1]

		case RecordTypeAccountIdentifier:
			if currentGroup == nil {
				return nil, errors.New("account identifier found outside of a group")
			}
			ai, err := parseAccountIdentifier(line)
			if err != nil {
				return nil, err
			}
			currentGroup.AccountIdentifiers = append(currentGroup.AccountIdentifiers, ai)

		case RecordTypeTransactionDetail:
			if currentGroup == nil {
				return nil, errors.New("transaction detail found outside of a group")
			}
			td, err := parseTransactionDetail(line)
			if err != nil {
				return nil, err
			}
			currentGroup.TransactionDetails = append(currentGroup.TransactionDetails, td)

		case RecordTypeContinuation:
			if currentGroup == nil {
				return nil, errors.New("continuation record found outside of a group")
			}
			cr, err := parseContinuationRecord(line)
			if err != nil {
				return nil, err
			}
			currentGroup.ContinuationRecords = append(currentGroup.ContinuationRecords, cr)

		case RecordTypeAccountTrailer:
			if currentGroup == nil {
				return nil, errors.New("account trailer found outside of a group")
			}
			at, err := parseAccountTrailer(line)
			if err != nil {
				return nil, err
			}
			currentGroup.AccountTrailers = append(currentGroup.AccountTrailers, at)

		case RecordTypeFileTrailer:
			br.FileTrailer, err = parseFileTrailer(line)
			if err != nil {
				return nil, err
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

// parseFileHeader parses file header line into the FileHeader struct while validating the data
func parseFileHeader(line string) (FileHeader, error) {
	if len(line) < 65 { // Considering the expected minimal length
		return FileHeader{}, ErrInvalidDataFormat
	}

	return FileHeader{
		RecordCode:       line[:2],
		SenderName:       strings.TrimSpace(line[2:37]),
		SenderID:         strings.TrimSpace(line[37:52]),
		FileCreationDate: line[52:60],
		FileCreationTime: line[60:64],
		FileIDModifier:   line[64:65],
	}, nil
}

// Parse for GroupHeader with enhanced handling for optional fields
func parseGroupHeader(line string) (GroupHeader, error) {
	fields := strings.Split(line, ",")

	gh := GroupHeader{
		RecordCode:   fields[0],
		OriginatorID: fields[1],
		GroupStatus:  fields[2],
		AsOfDate:     fields[3],
	}

	// Handle optional fields
	if len(fields) > 4 && fields[4] != "/" {
		gh.AsOfTime = fields[4]
	}
	if len(fields) > 5 && fields[5] != "/" {
		gh.CurrencyCode = fields[5]
	}
	if len(fields) > 6 && fields[6] != "/" {
		gh.AsOfDateModifier = fields[6]
	}
	return gh, nil
}

func parseAccountIdentifier(line string) (AccountIdentifier, error) {
	fields := strings.Split(line, ",")

	if len(fields) < 4 {
		return AccountIdentifier{}, fmt.Errorf("not enough fields in account identifier record: received %d, require at least 4", len(fields))
	}

	// Extract fields and handle potential empty optional fields with default values or similar logic.
	ai := AccountIdentifier{
		RecordCode:            strings.TrimSpace(fields[0]),
		CustomerAccountNumber: strings.TrimSpace(fields[1]),
	}

	// Example Handle optional third field for currency code in a typical specification
	if len(fields) > 2 && fields[2] != "" {
		ai.CurrencyCode = strings.TrimSpace(fields[2])
	} else {
		// Set default or handle the absence of optional fields if necessary
		ai.CurrencyCode = "USD" // Default currency, or could remain empty if that's acceptable
	}

	// Handle the NumberOfItems with default and unknown cases
	itemsField := fields[3]
	if itemsField == ",," || itemsField == "" {
		ai.NumberOfItems = 0 // Assuming 0 denotes an unknown or not applicable count
	} else {
		var err error
		ai.NumberOfItems, err = strconv.Atoi(itemsField)
		if err != nil {
			return AccountIdentifier{}, fmt.Errorf("error parsing number of items: %v", err)
		}
	}

	return ai, nil
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

// ParseGroupTrailer translates group trailer section
func parseGroupTrailer(line string) GroupTrailer {
	groupControlTotal, _ := strconv.ParseFloat(line[2:20], 64)
	numAccounts, _ := strconv.Atoi(line[20:29])
	numRecords, _ := strconv.Atoi(line[29:38])

	return GroupTrailer{
		RecordCode:        line[:2],
		GroupControlTotal: groupControlTotal,
		NumberOfAccounts:  numAccounts,
		NumberOfRecords:   numRecords,
	}
}

// parseFileTrailer translates the file trailer line while validating the number of records and control total
func parseFileTrailer(line string) (FileTrailer, error) {
	fields := strings.Split(line, ",")
	ft := FileTrailer{
		RecordCode: fields[0],
		NumberOfRecords: func(s string) int {
			if val, err := strconv.Atoi(s); err == nil {
				return val
			}
			return 0
		}(fields[2]), // TODO(adam): why so complex?
	}

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
	return ft, nil
}

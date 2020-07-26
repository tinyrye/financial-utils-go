package account_transactions

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/tinyrye/financial-utils-go/errors"
	"github.com/tinyrye/financial-utils-go/models"

	"github.com/araddon/dateparse"
)

/**
 *
 *
 *
 *
 */
type AmountParser struct {
	Format *regexp.Regexp
	NegativeSubMatches []int
}

/**
 *
 */
type AmountParserError struct {
	AmountText string
	Problem string
}

func (e *AmountParserError) Error() string {
	return fmt.Sprintf("Amount could not be parsed: text=[\"%s\"]; problem=[\"%s\"]", e.AmountText, e.Problem)
}

func (p *AmountParser) IsNegativeSubmatchIndex(subMatchIndex int) bool {
	for _, NegativeSubMatch := range p.NegativeSubMatches {
		if subMatchIndex == NegativeSubMatch {
			return true
		}
	}
	return false
}

func (p *AmountParser) Parse(amountText string) (*float64, error) {
	amountTextTrimmed := strings.Trim(amountText, " ")
	if len(amountTextTrimmed) > 0 {
		amountTextParse := p.Format.FindStringSubmatch(amountTextTrimmed)
		if len(amountTextParse) == 0 {
			return nil, &AmountParserError{amountTextTrimmed, "No valid pattern matched."}
		}
		for subMatchIndex, amountTextParseSubMatch := range amountTextParse[1:] {
			if len(amountTextParseSubMatch) > 0 {
				amount, amountErr := strconv.ParseFloat(amountTextParseSubMatch, 64)
				if p.IsNegativeSubmatchIndex(subMatchIndex) {
					amount *= -1.0
				}
				return &amount, amountErr
			}
		}
	}
	return nil, nil
}

func CreateAmountParser(formatText string, negativeSubMatches []int) *AmountParser {
	amountTextRegex, amountTextRegexErr := regexp.Compile(formatText)
	errors.PanicIf(amountTextRegexErr, "Creating Amount Regex")
	return &AmountParser{amountTextRegex, negativeSubMatches}
}

func ParseOptDate(dateText string) (*time.Time, error) {
	dateTextTrimmed := strings.Trim(dateText, " ")
	if len(dateTextTrimmed) > 0 {
		date, dateErr := dateparse.ParseAny(dateTextTrimmed)
		return &date, dateErr
	} else {
		return nil, nil
	}
}

/**
 *
 *
 *
 *
 */
type InstitutionTransactionReader interface {
	ConvertFromTransactionExport(row map[string]string) (*models.AccountTransaction, error)
}

/**
 *
 *
 *
 *
 */
type UnknownInstitutionByCodeError struct {
	Code string
}

func (e *UnknownInstitutionByCodeError) Error() string {
	return fmt.Sprintf("Unknown financial institution based on code: %s", e.Code)
}

/**
 *
 *
 *
 *
 */
type UwcuTransactionReader struct {
	amountParser *AmountParser
}

func (u *UwcuTransactionReader) ConvertFromTransactionExport(row map[string]string) (*models.AccountTransaction, error) {
	postedAt, postedAtParseErr := ParseOptDate(row["Posted Date"])
	if postedAtParseErr != nil {
		return nil, postedAtParseErr
	}
	amount, amountParseErr := u.amountParser.Parse(row["Amount"])

	if amountParseErr != nil {
		return nil, amountParseErr
	}
	return &models.AccountTransaction {
		InstitutionCode: "uwcu",
		Account: &models.Account {
			InstitutionCode: "uwcu",
			Number: row["AccountNumber"],
			Type: row["AccountType"],
		},
		PostedAt: postedAt,
		Amount: amount,
		Description: row["Description"],
		Category: row["Category"],
		Note: row["Note"],
	}, nil
}

func NewUwcuTransactionReader() *UwcuTransactionReader {
	return &UwcuTransactionReader { CreateAmountParser("\\(\\$(\\d+)\\)|\\$(\\d+)|\\(\\$(\\d+\\.\\d+)\\)|\\$(\\d+\\.\\d+)", []int{0, 2}) }
}

/**
 *
 *
 *
 *
 */
type ChaseTransactionReader struct {
	amountParser *AmountParser
}

func (u *ChaseTransactionReader) ConvertFromTransactionExport(row map[string]string) (*models.AccountTransaction, error) {
	transactedAt, transactedAtParseErr := ParseOptDate(row["Transaction Date"])
	if transactedAtParseErr != nil {
		return nil, transactedAtParseErr
	}
	postedAt, postedAtParseErr := ParseOptDate(row["Post Date"])
	if postedAtParseErr != nil {
		return nil, postedAtParseErr
	}
	amount, amountParseErr := u.amountParser.Parse(row["Amount"])

	if amountParseErr != nil {
		return nil, amountParseErr
	}
	return &models.AccountTransaction {
		InstitutionCode: "chase",
		TransactionType: row["Type"],
		TransactedAt: transactedAt,
		PostedAt: postedAt,
		Amount: amount,
		Description: row["Description"],
		Category: row["Category"],
	}, nil
}

func NewChaseTransactionReader() *ChaseTransactionReader {
	return &ChaseTransactionReader { CreateAmountParser("^-(\\d+)$|^(\\d+)$|^-(\\d+\\.\\d+)$|^(\\d+\\.\\d+)$", []int{0, 2}) }
}

func NewInstitutionTransactionReader(institutionCode string, fileType string) (InstitutionTransactionReader, error) {
	if institutionCode == "uwcu" {
		return NewUwcuTransactionReader(), nil
	} else if institutionCode == "chase" {
		return NewChaseTransactionReader(), nil
	} else {
		return nil, &UnknownInstitutionByCodeError{institutionCode}
	}
}

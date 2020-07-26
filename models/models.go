package models

import (
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/tinyrye/financial-utils-go/errors"

	"github.com/araddon/dateparse"
)

type AccountTransaction struct {
	TransactionId string;
	AccountNumber string;
	AccountType string;
	PostedAt *time.Time;
	Amount *float64;
	Description string;
	Merchant string;
	Category string;
	Note string;
}

type AmountParser struct {
	Format *regexp.Regexp
	NegativeSubMatches []int
}

func CreateAmountParser(formatText string, negativeSubMatches []int) *AmountParser {
	amountTextRegex, amountTextRegexErr := regexp.Compile("\\(\\$(\\d+)\\)|\\$(\\d+)|\\(\\$(\\d+\\.\\d+)\\)|\\$(\\d+\\.\\d+)")
	errors.PanicIf(amountTextRegexErr, "Creating Amount Regex")
	return &AmountParser{amountTextRegex, negativeSubMatches}
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

func ParseOptDate(dateText string) (*time.Time, error) {
	dateTextTrimmed := strings.Trim(dateText, " ")
	if len(dateTextTrimmed) > 0 {
		date, dateErr := dateparse.ParseAny(dateTextTrimmed)
		return &date, dateErr
	} else {
		return nil, nil
	}
}

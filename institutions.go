package account_transactions

import (
	"github.com/tinyrye/financial-utils-go/models"
)

type InstitutionTransactionReader interface {
	ConvertFromTransactionExport(row map[string]string) (*models.AccountTransaction, error)
}

type UwcuTransactionReader struct {
	amountParser *models.AmountParser
}

func (u *UwcuTransactionReader) ConvertFromTransactionExport(row map[string]string) (*models.AccountTransaction, error) {
	postedAt, postedAtParseErr := models.ParseOptDate(row["Posted Date"])
	if postedAtParseErr != nil {
		return nil, postedAtParseErr
	}
	amount, amountParseErr := u.amountParser.Parse(row["Amount"])
	
	if amountParseErr != nil {
		return nil, amountParseErr
	}
	return &models.AccountTransaction {
		AccountNumber: row["AccountNumber"],
		AccountType: row["AccountType"],
		PostedAt: postedAt,
		Amount: amount,
		Description: row["Description"],
		Category: row["Category"],
		Note: row["Note"],
	}, nil
}

func NewUwcuTransactionReader() *UwcuTransactionReader {
	return &UwcuTransactionReader { models.CreateAmountParser("\\(\\$(\\d+)\\)|\\$(\\d+)|\\(\\$(\\d+\\.\\d+)\\)|\\$(\\d+\\.\\d+)", []int{0, 2}) }
}

package account_transactions

import (
	"encoding/json"
	"fmt"
	"io"

	"github.com/tinyrye/financial-utils-go/models"
	"github.com/tinyrye/financial-utils-go/csv"
	"github.com/tinyrye/financial-utils-go/errors"
)

type AccountImportContext struct {
	AssumedAccount *models.Account
}

func (c *AccountImportContext) fillInDefaults(transaction *models.AccountTransaction) {
	if transaction.Account == nil {
		transaction.Account = c.AssumedAccount
	}
}

func ParseAndValidateCsv(institutionType string, filePath string, importContext *AccountImportContext) []models.AccountTransaction {
	transactionReader, readerErr := NewInstitutionTransactionReader(institutionType, "csv")
	errors.PanicIf(readerErr, "NewInstitutionTransactionReader")

	csvDataSet, csvFileErr := csv.ReadDataSet(filePath)
	errors.PanicIf(csvFileErr, "CSV File Reading")
	transactions := make([]models.AccountTransaction, 0)

	csvDataSet.ForEachRow(func(row map[string]string) {
		transaction, transactionConversionErr := transactionReader.ConvertFromTransactionExport(row)
		errors.PanicIf(transactionConversionErr, fmt.Sprintf("Converting transaction: %s", row))
		importContext.fillInDefaults(transaction)
		transactions = append(transactions, *transaction)
	})

	return transactions
}

func WriteTransactionsToJson(transactions []models.AccountTransaction, jsonOutput io.Writer) {
	recordMapsJsonBytes, recordMapsJsonBytesErr := json.Marshal(&transactions)
	errors.PanicIf(recordMapsJsonBytesErr, "JSON Encoding")
	_, recordMapsJsonBytesWriteErr := jsonOutput.Write(recordMapsJsonBytes)
	errors.PanicIf(recordMapsJsonBytesWriteErr, "JSON Std Out Write")
}

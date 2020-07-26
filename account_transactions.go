package account_transactions

import (
	"encoding/json"
	"fmt"
	"io"
	"log"

	"github.com/tinyrye/financial-utils-go/models"
	"github.com/tinyrye/financial-utils-go/csv"
	"github.com/tinyrye/financial-utils-go/errors"
)

func ParseAndValidateCsv(institutionType string, filePath string) []models.AccountTransaction {
	var transactionReader InstitutionTransactionReader

	if institutionType == "uwcu" {
		transactionReader = NewUwcuTransactionReader()
	} else {
		log.Panicf("Unknown transaction file type %s", institutionType)
	}

	csvDataSet, csvFileErr := csv.ReadDataSet(filePath)
	errors.PanicIf(csvFileErr, "CSV File Reading")
	transactions := make([]models.AccountTransaction, 0)

	csvDataSet.ForEachRow(func(row map[string]string) {
		transaction, transactionConversionErr := transactionReader.ConvertFromTransactionExport(row)
		errors.PanicIf(transactionConversionErr, fmt.Sprintf("Converting transaction: %s", row))
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

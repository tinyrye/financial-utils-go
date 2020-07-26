package account_transactions

import (
	"testing"
	"time"

	"github.com/tinyrye/financial-utils-go/errors"

	"github.com/araddon/dateparse"
	"github.com/stretchr/testify/assert"
)

func GetTestDate(text string) *time.Time {
	date, dateErr := dateparse.ParseAny(text)
	errors.PanicIf(dateErr, "Test Date Parse")
	return &date
}
func TestUwcuTransaction(testRun *testing.T) {
	actualTransactions := ParseAndValidateCsv("uwcu", "test-data/test-uwcu-transaction-file.csv")
	assert.Equal(testRun, len(actualTransactions), 8)
	assert.Equal(testRun, actualTransactions[0].AccountNumber, "8675309")
	assert.Equal(testRun, actualTransactions[0].AccountType, "CK")
	assert.Equal(testRun, actualTransactions[0].PostedAt, GetTestDate("7/25/2020"))
	assert.Equal(testRun, *actualTransactions[0].Amount, -200.00)
	assert.Equal(testRun, actualTransactions[0].Description, "Credit Card Payment")
	assert.Equal(testRun, actualTransactions[0].Category, "Credit Card Payment")

	assert.Equal(testRun, actualTransactions[1].AccountNumber, "ABCDEF-123456")
	assert.Equal(testRun, actualTransactions[1].AccountType, "CC")
	assert.Equal(testRun, actualTransactions[1].PostedAt, GetTestDate("7/24/2020"))
	assert.Equal(testRun, *actualTransactions[1].Amount, -6.53)
	assert.Equal(testRun, actualTransactions[1].Description, "My Fav Coffee Shop")
	assert.Equal(testRun, actualTransactions[1].Category, "Dining Out")
}
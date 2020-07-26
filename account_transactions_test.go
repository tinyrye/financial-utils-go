package account_transactions

import (
	"testing"
	"time"

	"github.com/tinyrye/financial-utils-go/errors"
	"github.com/tinyrye/financial-utils-go/models"

	"github.com/araddon/dateparse"
	"github.com/stretchr/testify/assert"
)

func GetTestDate(text string) *time.Time {
	date, dateErr := dateparse.ParseAny(text)
	errors.PanicIf(dateErr, "Test Date Parse")
	return &date
}

func TestUwcuTransaction(testRun *testing.T) {
	actualTransactions := ParseAndValidateCsv("uwcu", "test-data/test-uwcu-transaction-file.csv", nil)
	assert.Equal(testRun, len(actualTransactions), 8)

	assert.Equal(testRun, actualTransactions[0].InstitutionCode, "uwcu")
	assert.Equal(testRun, actualTransactions[0].Account.InstitutionCode, "uwcu")
	assert.Equal(testRun, actualTransactions[0].Account.Number, "8675309")
	assert.Equal(testRun, actualTransactions[0].Account.Type, "CK")
	assert.Equal(testRun, actualTransactions[0].PostedAt, GetTestDate("7/25/2020"))
	assert.Equal(testRun, *actualTransactions[0].Amount, -200.00)
	assert.Equal(testRun, actualTransactions[0].Description, "Credit Card Payment")
	assert.Equal(testRun, actualTransactions[0].Category, "Credit Card Payment")

	assert.Equal(testRun, actualTransactions[1].InstitutionCode, "uwcu")
	assert.Equal(testRun, actualTransactions[1].Account.InstitutionCode, "uwcu")
	assert.Equal(testRun, actualTransactions[1].Account.Number, "ABCDEF-123456")
	assert.Equal(testRun, actualTransactions[1].Account.Type, "CC")
	assert.Equal(testRun, actualTransactions[1].PostedAt, GetTestDate("7/24/2020"))
	assert.Equal(testRun, *actualTransactions[1].Amount, -6.53)
	assert.Equal(testRun, actualTransactions[1].Description, "My Fav Coffee Shop")
	assert.Equal(testRun, actualTransactions[1].Category, "Dining Out")
}

func TestChaseTransaction(testRun *testing.T) {
	assumedContext := &AccountImportContext {
		&models.Account {
			"chase",
			"864210",
			"CC",
		},
	}
	actualTransactions := ParseAndValidateCsv("chase", "test-data/test-chase-transaction-file.csv", assumedContext)
	
	assert.Equal(testRun, len(actualTransactions), 5)

	// 07/10/2020,07/10/2020,PURCHASE INTEREST CHARGE,Fees & Adjustments,Fee,-8
	assert.Equal(testRun, actualTransactions[0].InstitutionCode, "chase")
	assert.Equal(testRun, actualTransactions[0].Account.InstitutionCode, assumedContext.AssumedAccount.InstitutionCode)
	assert.Equal(testRun, actualTransactions[0].Account.Number, assumedContext.AssumedAccount.Number)
	assert.Equal(testRun, actualTransactions[0].Account.Type, assumedContext.AssumedAccount.Type)
	assert.Equal(testRun, actualTransactions[0].TransactionType, "Fee")
	assert.Equal(testRun, actualTransactions[0].TransactedAt, GetTestDate("7/10/2020"))
	assert.Equal(testRun, actualTransactions[0].PostedAt, GetTestDate("7/10/2020"))
	assert.Equal(testRun, *actualTransactions[0].Amount, -8.00)
	assert.Equal(testRun, actualTransactions[0].Description, "PURCHASE INTEREST CHARGE")
	assert.Equal(testRun, actualTransactions[0].Category, "Fees & Adjustments")

	// 06/01/2020,06/03/2020,CKO*Patreon* Membership,Entertainment,Sale,-44.00
	assert.Equal(testRun, actualTransactions[1].InstitutionCode, "chase")
	assert.Equal(testRun, actualTransactions[1].Account.InstitutionCode, assumedContext.AssumedAccount.InstitutionCode)
	assert.Equal(testRun, actualTransactions[1].Account.Number, assumedContext.AssumedAccount.Number)
	assert.Equal(testRun, actualTransactions[1].Account.Type, assumedContext.AssumedAccount.Type)
	assert.Equal(testRun, actualTransactions[1].TransactionType, "Sale")
	assert.Equal(testRun, actualTransactions[1].TransactedAt, GetTestDate("6/01/2020"))
	assert.Equal(testRun, actualTransactions[1].PostedAt, GetTestDate("6/03/2020"))
	assert.Equal(testRun, *actualTransactions[1].Amount, -44.00)
	assert.Equal(testRun, actualTransactions[1].Description, "CKO*Patreon* Membership")
	assert.Equal(testRun, actualTransactions[1].Category, "Entertainment")

	// 01/10/2020,01/12/2020,Shopping Refund,Shopping,Sale,31.17
	assert.Equal(testRun, actualTransactions[4].InstitutionCode, "chase")
	assert.Equal(testRun, actualTransactions[4].Account.InstitutionCode, assumedContext.AssumedAccount.InstitutionCode)
	assert.Equal(testRun, actualTransactions[4].Account.Number, assumedContext.AssumedAccount.Number)
	assert.Equal(testRun, actualTransactions[4].Account.Type, assumedContext.AssumedAccount.Type)
	assert.Equal(testRun, actualTransactions[4].TransactionType, "Sale")
	assert.Equal(testRun, actualTransactions[4].TransactedAt, GetTestDate("01/10/2020"))
	assert.Equal(testRun, actualTransactions[4].PostedAt, GetTestDate("01/12/2020"))
	assert.Equal(testRun, *actualTransactions[4].Amount, 31.17)
	assert.Equal(testRun, actualTransactions[4].Description, "Shopping Rebate")
	assert.Equal(testRun, actualTransactions[4].Category, "Shopping")
}

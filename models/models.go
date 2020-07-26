package models

import (
	"time"
)

type Account struct {
	Id string;
	InstitutionCode string;
	Number string;
	Type string;
}

type AccountTransaction struct {
	Id string;
	InstitutionCode string;
	Account *Account;
	TransactionId string;
	TransactionType string;
	TransactedAt *time.Time;
	PostedAt *time.Time;
	Amount *float64;
	Description string;
	Merchant string;
	Category string;
	Note string;
}

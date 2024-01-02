package domain

import "database/sql"

type CreditCard struct {
	UserID                           uint
	Type, Number, Name, Expired, CVV sql.NullString
}

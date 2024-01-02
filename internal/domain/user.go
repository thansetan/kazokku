package domain

import "database/sql"

type User struct {
	ID                             uint
	Name, Address, Email, Password sql.NullString
	Photos                         []Photo
	CreditCard                     CreditCard
}

func (u User) IsEmpty() bool {
	return u.ID == 0 && u.Name.String == "" && u.Address.String == "" && u.Email.String == "" && u.Password.String == "" && u.Photos == nil && u.CreditCard == (CreditCard{})
}

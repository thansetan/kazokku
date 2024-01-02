package dto

import (
	"strings"
	"time"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
)

type UserRegisterRequest struct {
	Name              string `json:"name" form:"name"`
	Address           string `json:"address" form:"address"`
	Email             string `json:"email" form:"email"`
	Password          string `json:"password" form:"password"`
	CreditCardType    string `json:"creditcard_type" form:"creditcard_type"`
	CreditCardNumber  string `json:"creditcard_number" form:"creditcard_number"`
	CreditCardName    string `json:"creditcard_name" form:"creditcard_name"`
	CreditCardExpired string `json:"creditcard_expired" form:"creditcard_expired"`
	CreditCardCVV     string `json:"creditcard_cvv" form:"creditcard_cvv"`
}

type UserQuery struct {
	Query   string `query:"q"`
	OrderBy string `query:"ob"`
	SortBy  string `query:"sb"`
	Offset  int    `query:"of"`
	Limit   int    `query:"lt"`
}

type UserResponse struct {
	ID         uint               `json:"user_id"`
	Name       string             `json:"name"`
	Email      string             `json:"email"`
	Address    string             `json:"address"`
	Photos     []string           `json:"photos"`
	CreditCard CreditCardResponse `json:"creditcard"`
}

type UserUpdateRequest struct {
	UserID            uint   `json:"user_id" form:"user_id"`
	Name              string `json:"name" form:"name"`
	Address           string `json:"address" form:"address"`
	Email             string `json:"email" form:"email"`
	Password          string `json:"password" form:"password"`
	CreditCardType    string `json:"creditcard_type" form:"creditcard_type"`
	CreditCardNumber  string `json:"creditcard_number" form:"creditcard_number"`
	CreditCardName    string `json:"creditcard_name" form:"creditcard_name"`
	CreditCardExpired string `json:"creditcard_expired" form:"creditcard_expired"`
	CreditCardCVV     string `json:"creditcard_cvv" form:"creditcard_cvv"`
}

var (
	validCCExpDate = validation.NewStringRule(func(s string) bool {
		if len(s) != 5 {
			return false
		}

		exp, err := time.Parse("01/06", s)
		if err != nil {
			return false
		}

		if time.Now().After(exp) {
			return false
		}

		return true
	}, "invalid credit card expired date")

	validOrderBy = validation.NewStringRule(func(s string) bool {
		lowerS := strings.ToLower(s)
		return lowerS == "name" || lowerS == "email"
	}, "invalid order by")

	validSortBy = validation.NewStringRule(func(s string) bool {
		lowerS := strings.ToLower(s)
		return lowerS == "asc" || lowerS == "desc"
	}, "invalid sort by")
)

func (r UserRegisterRequest) Validate() error {
	return validation.ValidateStruct(&r,
		validation.Field(&r.Name, validation.Required),
		validation.Field(&r.Email, validation.Required, is.Email),
		validation.Field(&r.Password, validation.Required),
		validation.Field(&r.Address, validation.Required),
		validation.Field(&r.CreditCardType, validation.Required),
		validation.Field(&r.CreditCardNumber, validation.Required, is.CreditCard),
		validation.Field(&r.CreditCardName, validation.Required),
		validation.Field(&r.CreditCardExpired, validation.Required, validCCExpDate),
		validation.Field(&r.CreditCardCVV, validation.Required),
	)
}

func (r UserUpdateRequest) Validate() error {
	return validation.ValidateStruct(&r,
		validation.Field(&r.UserID, validation.Required),
		validation.Field(&r.Email, is.Email),
		validation.Field(&r.CreditCardNumber, is.CreditCard),
		validation.Field(&r.CreditCardExpired, validCCExpDate),
	)
}

func (q UserQuery) Validate() error {
	return validation.ValidateStruct(&q,
		validation.Field(&q.Offset, validation.Min(0)),
		validation.Field(&q.Limit, validation.Min(30)),
		validation.Field(&q.OrderBy, validOrderBy),
		validation.Field(&q.SortBy, validSortBy),
	)
}

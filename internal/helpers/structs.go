package helpers

import (
	"database/sql"
	"kazokku/internal/app/delivery/dto"
	"kazokku/internal/domain"
)

func UserRegisterDTOtoUserDomain(data dto.UserRegisterRequest) domain.User {
	var user domain.User

	user.Name = sql.NullString{
		String: data.Name,
		Valid:  true,
	}
	user.Address = sql.NullString{
		String: data.Address,
		Valid:  true,
	}
	user.Email = sql.NullString{
		String: data.Email,
		Valid:  true,
	}
	user.Password = sql.NullString{
		String: data.Password,
		Valid:  true,
	}

	return user
}

func UserUpdateDTOtoUserDomain(data dto.UserUpdateRequest) domain.User {
	var user domain.User

	user.ID = data.UserID
	user.Name = sql.NullString{
		String: data.Name,
		Valid:  data.Name != "",
	}
	user.Address = sql.NullString{
		String: data.Address,
		Valid:  data.Address != "",
	}
	user.Email = sql.NullString{
		String: data.Email,
		Valid:  data.Email != "",
	}
	user.Password = sql.NullString{
		String: data.Password,
		Valid:  data.Password != "",
	}

	return user
}

func UserRegisterDTOtoCCDomain(data dto.UserRegisterRequest, userID uint) domain.CreditCard {
	var cc domain.CreditCard

	cc.UserID = userID
	cc.Type = sql.NullString{
		String: data.CreditCardType,
		Valid:  true,
	}
	cc.Number = sql.NullString{
		String: data.CreditCardNumber,
		Valid:  true,
	}
	cc.Name = sql.NullString{
		String: data.CreditCardName,
		Valid:  true,
	}
	cc.Expired = sql.NullString{
		String: data.CreditCardExpired,
		Valid:  true,
	}
	cc.CVV = sql.NullString{
		String: data.CreditCardCVV,
		Valid:  true,
	}

	return cc
}

func UserUpdateDTOtoCCDomain(data dto.UserUpdateRequest, userID uint) domain.CreditCard {
	var cc domain.CreditCard

	cc.UserID = userID
	cc.Type = sql.NullString{
		String: data.CreditCardType,
		Valid:  data.CreditCardType != "",
	}
	cc.Number = sql.NullString{
		String: data.CreditCardNumber,
		Valid:  data.CreditCardNumber != "",
	}
	cc.Name = sql.NullString{
		String: data.CreditCardName,
		Valid:  data.CreditCardName != "",
	}
	cc.Expired = sql.NullString{
		String: data.CreditCardExpired,
		Valid:  data.CreditCardExpired != "",
	}
	cc.CVV = sql.NullString{
		String: data.CreditCardCVV,
		Valid:  data.CreditCardCVV != "",
	}

	return cc
}

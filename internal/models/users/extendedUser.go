package users

import "github.com/e1esm/Effective_Test/internal/models/nationalities"

type ExtendedUser struct {
	User
	Age         int                         `json:"age"`
	Sex         string                      `json:"sex"`
	Nationality []nationalities.Nationality `json:"nationality"`
}

func ExtendedFromRequest(user User) *ExtendedUser {
	return &ExtendedUser{
		User: user,
	}
}

func ExtendedFromEntity(user EntityUser) *ExtendedUser {
	return &ExtendedUser{
		Age:         user.Age,
		Sex:         user.Sex,
		Nationality: user.Nationality,
		User: User{
			user.Name,
			user.Surname,
			user.Patronymic,
		},
	}
}

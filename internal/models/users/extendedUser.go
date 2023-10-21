package users

import "github.com/e1esm/Effective_Test/internal/models/nationalities"

type ExtendedUser struct {
	User
	Age         int                         `json:"age"`
	Sex         string                      `json:"sex"`
	Nationality []nationalities.Nationality `json:"nationality"`
}

func NewExtendedUser(user User) *ExtendedUser {
	return &ExtendedUser{
		User: user,
	}
}

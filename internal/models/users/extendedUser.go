package users

type ExtendedUser struct {
	User
	Age         int      `json:"age"`
	Sex         string   `json:"sex"`
	Nationality []string `json:"nationality"`
}

func NewExtendedUser(user User) *ExtendedUser {
	return &ExtendedUser{
		User: user,
	}
}

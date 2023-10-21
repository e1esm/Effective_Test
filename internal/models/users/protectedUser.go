package users

import (
	"sync"
)

type ProtectedUser struct {
	user ExtendedUser
	mu   *sync.Mutex
}

func NewProtectedUser(user ExtendedUser) *ProtectedUser {
	return &ProtectedUser{
		user,
		&sync.Mutex{},
	}
}

func (pu *ProtectedUser) SetSex(sex string) {
	pu.mu.Lock()
	pu.user.Sex = sex
	pu.mu.Unlock()
}

func (pu *ProtectedUser) SetNationality(nationality []string) {
	pu.mu.Lock()
	pu.user.Nationality = nationality
	pu.mu.Unlock()
}

func (pu *ProtectedUser) SetAge(age int) {
	pu.mu.Lock()
	pu.user.Age = age
	pu.mu.Unlock()
}

func (pu *ProtectedUser) GetUser() ExtendedUser {
	return pu.user
}

func (pu *ProtectedUser) Validate() bool {
	if pu.user.Sex == "" || pu.user.Age == 0 || pu.user.Nationality == nil {
		return false
	}
	return true
}

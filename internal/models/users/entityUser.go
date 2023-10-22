package users

import "github.com/google/uuid"

type EntityUser struct {
	ID uuid.UUID
	ExtendedUser
}

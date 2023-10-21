package identifier

import (
	"context"
	"github.com/e1esm/Effective_Test/internal/models/users"
	"sync"
	"time"
)

type Identifier interface {
	Identify(users.User) *users.ExtendedUser
}

type IdentityService struct {
	timeout time.Duration
}

func NewIdentifyService(timeout time.Duration) *IdentityService {
	return &IdentityService{
		timeout: timeout,
	}
}

func (is *IdentityService) Identify(user users.User) *users.ExtendedUser {
	wg := sync.WaitGroup{}
	protectedUser := users.NewProtectedUser(*users.NewExtendedUser(user))

	ctx, cancel := context.WithTimeout(context.Background(), is.timeout)
	defer cancel()

	wg.Add(3)

	go func() {
		defer wg.Done()
		is.requestAge(ctx, protectedUser)
	}()

	go func() {
		defer wg.Done()
		is.requestSex(ctx, protectedUser)
	}()

	go func() {
		defer wg.Done()
		is.requestNationality(ctx, protectedUser)
	}()

	wg.Wait()

	if isOk := protectedUser.Validate(); isOk {
		builtUser := protectedUser.GetUser()
		return &builtUser
	}
	return nil
}

package identifier

import (
	"context"
	"github.com/e1esm/Effective_Test/internal/models/users"
	"github.com/e1esm/Effective_Test/pkg/utils/logger"
	"go.uber.org/zap"
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
		logger.GetLogger().Info("Successfully identified user",
			zap.String("user", user.Name))
		builtUser := protectedUser.GetUser()
		return &builtUser
	}
	logger.GetLogger().Error("Failed to identify user",
		zap.String("user", user.Name))
	return nil
}

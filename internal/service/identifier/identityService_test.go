package identifier

import (
	"github.com/e1esm/Effective_Test/internal/models/users"
	"testing"
	"time"
)

func TestIdentityService_Identify(t *testing.T) {
	identityService := NewIdentifyService(1 * time.Second)
	table := []struct {
		name string
		user users.User
		isOk bool
	}{
		{
			name: "OK",
			user: users.User{
				Name:       "",
				Surname:    "",
				Patronymic: "",
			},
			isOk: false,
		},
		{
			name: "FAIL",
			user: users.User{
				Name:       "Egor",
				Surname:    "Ivanov",
				Patronymic: "Ivanovich",
			},
			isOk: true,
		},
		{
			name: "FAIL",
			user: users.User{
				Name:       "",
				Surname:    "Ivanov",
				Patronymic: "Ivanovich",
			},
			isOk: false,
		},
	}

	for _, test := range table {
		t.Run(test.name, func(t *testing.T) {
			retrievedUser := identityService.Identify(test.user)
			if retrievedUser != nil && !test.isOk {
				t.Errorf("Invalid result. Want: %v, got: %v",
					test.user, retrievedUser)
			}
			if retrievedUser == nil && test.isOk {
				t.Errorf("Invalid result. Want: %v, got: %v", test.user, retrievedUser)
			}
		})
	}
}

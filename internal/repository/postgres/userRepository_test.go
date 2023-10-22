package postgres

import (
	"context"
	"fmt"
	"github.com/e1esm/Effective_Test/internal/models/nationalities"
	"github.com/e1esm/Effective_Test/internal/models/options"
	"github.com/e1esm/Effective_Test/internal/models/users"
	"github.com/e1esm/Effective_Test/internal/repository/postgres/migrations"
	"github.com/e1esm/Effective_Test/pkg/utils/logger"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/ory/dockertest/v3"
	"github.com/ory/dockertest/v3/docker"
	"go.uber.org/zap"
	"os"
	"testing"
	"time"
)

var testRepo PeopleRepository

func TestMain(m *testing.M) {
	pool, err := dockertest.NewPool("")
	if err != nil {
		logger.GetLogger().Fatal("Couldn't have constructed connection pool", zap.String("err", err.Error()))
	}
	if err := pool.Client.Ping(); err != nil {
		logger.GetLogger().Fatal("Connection's failed", zap.String("err", err.Error()))
	}

	testedResource, err := pool.RunWithOptions(&dockertest.RunOptions{
		Repository: "postgres",
		Tag:        "latest",
		Env: []string{
			"POSTGRES_USER=postgres",
			"POSTGRES_PASSWORD=postgres",
			"POSTGRES_DB=test",
			"POSTGRES_PORT=5432",
			"PGUSER=postgres",
			"listen_addresses = '*'",
		},
	}, func(config *docker.HostConfig) {
		config.AutoRemove = true
		config.RestartPolicy = docker.RestartPolicy{Name: "no"}
	})
	if err := testedResource.Expire(100); err != nil {
		logger.GetLogger().Fatal("Error while setting up resource", zap.String("err", err.Error()))
	}
	logger.GetLogger().Info(testedResource.GetHostPort("5432/tcp"))
	dsn := fmt.Sprintf("postgres://postgres:postgres@%s/test?sslmode=disable", testedResource.GetHostPort("5432/tcp"))

	testRepo = PeopleRepository{}
	if err = pool.Retry(func() error {
		testRepo.db, err = pgxpool.New(context.Background(), dsn)
		if err != nil {
			return err
		}
		return testRepo.db.Ping(context.Background())
	}); err != nil {
		logger.GetLogger().Fatal("Couldn't have connected to the DB")
	}
	pool.MaxWait = 30 * time.Second
	_, err = migrations.ConnectAndRunMigrations(context.Background(), dsn, "file://./migrations", migrations.UP)

	if err != nil {
		logger.GetLogger().Fatal("Couldn't have created new connection pool with DB", zap.String("err", err.Error()))
	}

	code := m.Run()

	if err := pool.Purge(testedResource); err != nil {
		logger.GetLogger().Fatal("Couldn't have purged resource", zap.String("err", err.Error()))
	}
	os.Exit(code)
}

func TestPeopleRepository_Save(t *testing.T) {
	table := []struct {
		name      string
		inputUser users.ExtendedUser
	}{
		{
			name: "Success",
			inputUser: users.ExtendedUser{
				User: users.User{
					Name:       "Egor",
					Surname:    "Ivanov",
					Patronymic: "Ivanovich",
				},
				Age: 20,
				Sex: "male",
				Nationality: []nationalities.Nationality{
					{ID: "RU", Probability: 100},
				},
			},
		},
		{
			name:      "Fail",
			inputUser: users.ExtendedUser{},
		},
	}

	for _, test := range table {
		_, err := testRepo.Save(context.Background(), test.inputUser)
		if err != nil && test.name == "Success" || err == nil && test.name == "Fail" {
			t.Errorf("Invalid result. Got: %v", err)
		}
	}
}

func TestPeopleRepository_Update(t *testing.T) {
	table := []struct {
		name      string
		ID        uuid.UUID
		inputUser users.ExtendedUser
	}{
		{
			name: "Success",
			ID:   uuid.UUID{},
			inputUser: users.ExtendedUser{
				User: users.User{
					Name:       "Egor",
					Surname:    "Ivanov",
					Patronymic: "Ivanovich",
				},
				Age: 20,
				Sex: "male",
				Nationality: []nationalities.Nationality{
					{ID: "RU", Probability: 100},
				},
			},
		},
		{
			name:      "Fail",
			ID:        uuid.New(),
			inputUser: users.ExtendedUser{},
		},
		{
			name: "Success",
			ID:   uuid.UUID{},
			inputUser: users.ExtendedUser{
				User: users.User{
					Name:       "Egor",
					Surname:    "Ivanov",
					Patronymic: "Ivanovich",
				},
				Age: 20,
				Sex: "male",
				Nationality: []nationalities.Nationality{
					{ID: "RU", Probability: 100},
				},
			},
		},
	}

	for _, test := range table {
		_, err := testRepo.Update(context.WithValue(context.Background(), "id", test.ID), test.inputUser)
		if err == nil && test.name == "Fail" || err != nil && test.name == "Success" {
			t.Errorf("Invalid result. Got: %v", err)
		}
	}
}

func TestPeopleRepository_Delete(t *testing.T) {
	table := []struct {
		name      string
		ID        uuid.UUID
		inputUser users.ExtendedUser
	}{
		{
			name: "Success",
			ID:   uuid.UUID{},
			inputUser: users.ExtendedUser{
				User: users.User{
					Name:       "Egor",
					Surname:    "Ivanov",
					Patronymic: "Ivanovich",
				},
				Age: 20,
				Sex: "male",
				Nationality: []nationalities.Nationality{
					{ID: "RU", Probability: 100},
				},
			},
		},
		{
			name:      "Fail",
			ID:        uuid.New(),
			inputUser: users.ExtendedUser{},
		},
		{
			name: "Fail",
			ID:   uuid.UUID{},
			inputUser: users.ExtendedUser{
				User: users.User{
					Name:       "Egor",
					Surname:    "Ivanov",
					Patronymic: "Ivanovich",
				},
				Age: 20,
				Sex: "male",
				Nationality: []nationalities.Nationality{
					{ID: "RU", Probability: 100},
				},
			},
		},
	}
	_, err := testRepo.Save(context.Background(), table[0].inputUser)
	if err != nil {
		t.Errorf("Invalid result of Save method")
	}
	for _, test := range table {
		_, err = testRepo.Delete(context.Background(), test.ID)
		if err == nil && test.name == "Fail" || err != nil && test.name == "Success" {
			t.Errorf("Invalid result. Got: %v", err)
		}

	}
}

func TestPeopleRepository_Get(t *testing.T) {
	input := []struct {
		ID        uuid.UUID
		inputUser users.ExtendedUser
	}{
		{
			ID: uuid.New(),
			inputUser: users.ExtendedUser{
				User: users.User{
					Name:       "Egor",
					Surname:    "Ivanov",
					Patronymic: "Ivanovich",
				},
				Age: 20,
				Sex: "male",
				Nationality: []nationalities.Nationality{
					{ID: "RU", Probability: 100},
				},
			},
		},
		{
			ID: uuid.New(),
			inputUser: users.ExtendedUser{
				User: users.User{
					Name:       "Nikita",
					Surname:    "Ivanov",
					Patronymic: "Ivanovich",
				},
				Age: 25,
				Sex: "male",
				Nationality: []nationalities.Nationality{
					{ID: "ES", Probability: 100},
				},
			},
		},
		{
			ID: uuid.New(),
			inputUser: users.ExtendedUser{
				User: users.User{
					Name:       "Albert",
					Surname:    "Ivanov",
					Patronymic: "Ivanovich",
				},
				Age: 30,
				Sex: "male",
				Nationality: []nationalities.Nationality{
					{ID: "BR", Probability: 100},
				},
			},
		},
	}
	for _, test := range input {
		_, err := testRepo.Save(context.Background(), test.inputUser)
		if err != nil {
			t.Errorf("invalid result. Got: %v", err)
		}
	}

	search := []struct {
		name           string
		opts           options.QueryOptions
		expectedLength int
	}{
		{
			name: "SUCCESS",
			opts: options.QueryOptions{
				options.NewUserOptions("female", "", 0, 2, 0),
				options.NewNationalityOptions([]string{}),
			},
			expectedLength: 0,
		},
		{
			name: "SUCCESS",
			opts: options.QueryOptions{
				options.NewUserOptions("male", "", 0, 1, 0),
				options.NewNationalityOptions([]string{"RU"}),
			},
			expectedLength: 1,
		},
	}
	for _, test := range search {
		res, err := testRepo.Get(context.Background(), &test.opts)
		if err != nil && test.name == "SUCCESS" {
			t.Errorf("Invalid result. Got: %v", err)
		}
		if len(res) != test.expectedLength {
			t.Errorf("Invalid length. Got: %v, want: %v", len(res), test.expectedLength)
		}
	}
}

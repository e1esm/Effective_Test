package postgres

import (
	"context"
	"fmt"
	"github.com/e1esm/Effective_Test/internal/models/users"
	"github.com/e1esm/Effective_Test/internal/repository/postgres/migrations"
	"github.com/e1esm/Effective_Test/pkg/utils/envParser"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"log"
)

const (
	dbURL      = "db_url"
	dbPassword = "db_password"
	dbUsername = "db_user"
	dbPort     = "db_port"
	db         = "db"
)

type Repository interface {
	Save(context.Context, users.ExtendedUser) (uuid.UUID, error)
}

type PeopleRepository struct {
	db *pgxpool.Pool
}

func NewPeopleRepository() *PeopleRepository {
	vars, err := envParser.ParseEnvVariable(dbURL, dbUsername, dbPassword, dbPort, db)
	if err != nil {
		log.Println(err.Error())
		return nil
	}

	connectionURL := fmt.Sprintf("postgres://%v:%v@%v:%v/%v",
		vars[dbUsername],
		vars[dbPassword],
		vars[dbURL],
		vars[dbPort],
		vars[db])

	pool, err := migrations.ConnectAndRunMigrations(context.Background(), connectionURL, migrations.UP)
	log.Println(pool)
	if err != nil {
		log.Println(err.Error())
		return nil
	}

	return &PeopleRepository{
		db: pool,
	}
}

func (pr *PeopleRepository) Save(ctx context.Context, person users.ExtendedUser) (uuid.UUID, error) {

	tx, err := pr.db.BeginTx(ctx, pgx.TxOptions{})
	defer tx.Rollback(ctx)
	if err != nil {
		return uuid.UUID{}, err
	}

	id, err := pr.savePerson(ctx, tx, person)
	if err != nil {
		return uuid.UUID{}, err
	}
	err = pr.saveNationality(context.WithValue(ctx, "id", id), tx, person)
	if err != nil {
		return uuid.UUID{}, err
	}
	if err := tx.Commit(ctx); err != nil {
		return uuid.UUID{}, err
	}

	return uuid.UUID{}, nil
}

func (pr *PeopleRepository) savePerson(ctx context.Context, tx pgx.Tx, person users.ExtendedUser) (uuid.UUID, error) {
	id := uuid.New()
	_, err := tx.Exec(ctx, "INSERT INTO people_info VALUES($1, $2, $3, $4, $5, $6)",
		id,
		person.Name,
		person.Surname,
		person.Patronymic,
		person.Age,
		person.Sex)
	if err != nil {
		return uuid.UUID{}, err
	}
	return id, nil
}

func (pr *PeopleRepository) saveNationality(ctx context.Context, tx pgx.Tx, person users.ExtendedUser) error {
	if userID, ok := ctx.Value("id").(uuid.UUID); ok {
		for i := 0; i < len(person.Nationality); i++ {
			if _, err := tx.Exec(ctx, "INSERT INTO person_nationality VALUES ($1, $2, $3, $4)",
				uuid.New(), person.Nationality[i].ID, person.Nationality[i].Probability, userID); err != nil {
				return err
			}
		}
	}

	return fmt.Errorf("couldn't have casted value to uuid")

}
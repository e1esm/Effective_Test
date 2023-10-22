package postgres

import (
	"context"
	"errors"
	"fmt"
	"github.com/e1esm/Effective_Test/internal/models/nationalities"
	"github.com/e1esm/Effective_Test/internal/models/options"
	"github.com/e1esm/Effective_Test/internal/models/users"
	"github.com/e1esm/Effective_Test/internal/repository/postgres/migrations"
	"github.com/e1esm/Effective_Test/pkg/utils/envParser"
	"github.com/e1esm/Effective_Test/pkg/utils/logger"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
	"log"
	"time"
)

const (
	dbURL      = "db_url"
	dbPassword = "db_password"
	dbUsername = "db_user"
	dbPort     = "db_port"
	db         = "db"
)

var (
	NoRecordsFound = errors.New("no records were found")
)

type Repository interface {
	Save(context.Context, users.ExtendedUser) (uuid.UUID, error)
	Delete(context.Context, uuid.UUID) (uuid.UUID, error)
	Update(context.Context, users.ExtendedUser) (uuid.UUID, error)
	Get(context.Context, options.UserFilter) ([]users.EntityUser, error)
}

type PeopleRepository struct {
	db *pgxpool.Pool
}

func NewPeopleRepository() *PeopleRepository {
	vars, err := envParser.ParseEnvVariable(dbURL, dbUsername, dbPassword, dbPort, db)
	if err != nil {
		logger.GetLogger().Error("Something's missing",
			zap.String("error", err.Error()))
		return nil
	}

	connectionURL := fmt.Sprintf("postgres://%v:%v@%v:%v/%v?sslmode=disable",
		vars[dbUsername],
		vars[dbPassword],
		vars[dbURL],
		vars[dbPort],
		vars[db])

	pool, err := migrations.ConnectAndRunMigrations(context.Background(), connectionURL, "", migrations.UP)
	if err != nil {
		logger.GetLogger().Error("Couldn't have either connected to the DB or run migrations",
			zap.String("err", err.Error()))
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
		log.Println(err.Error())
		return uuid.UUID{}, err
	}
	err = pr.saveNationality(context.WithValue(ctx, "id", id), tx, person)
	if err != nil {
		log.Println(err.Error())
		return uuid.UUID{}, err
	}
	if err := tx.Commit(ctx); err != nil {
		return uuid.UUID{}, err
	}

	logger.GetLogger().Info("Successfully committed transaction for user",
		zap.String("id", id.String()))

	return id, nil
}

func (pr *PeopleRepository) savePerson(ctx context.Context, tx pgx.Tx, person users.ExtendedUser) (uuid.UUID, error) {
	var id uuid.UUID
	id = uuid.New()
	if retrievedID := ctx.Value("id"); retrievedID != nil {
		id = retrievedID.(uuid.UUID)
	}
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
	userID := ctx.Value("id").(uuid.UUID)
	for i := 0; i < len(person.Nationality); i++ {
		if _, err := tx.Exec(ctx, "INSERT INTO person_nationality VALUES ($1, $2, $3, $4)",
			uuid.New(), person.Nationality[i].ID, person.Nationality[i].Probability, userID); err != nil {
			return err
		}
	}
	logger.GetLogger().Info("Successfully inserted nationalities",
		zap.String("user", person.Name))
	return nil

}

func (pr *PeopleRepository) Delete(ctx context.Context, id uuid.UUID) (uuid.UUID, error) {

	tag, err := pr.db.Exec(ctx, "DELETE FROM people_info WHERE id = $1", id)
	if err != nil {
		return uuid.UUID{}, err
	}
	if tag.RowsAffected() == 0 {
		return uuid.UUID{}, NoRecordsFound
	}
	logger.GetLogger().Info("Successfully deleted user", zap.String("ID", id.String()))
	return id, nil
}

func (pr *PeopleRepository) Update(ctx context.Context, user users.ExtendedUser) (uuid.UUID, error) {
	if ctx.Value("id") == nil {
		return uuid.UUID{}, fmt.Errorf("ID is null")
	}
	if _, err := pr.Delete(ctx, ctx.Value("id").(uuid.UUID)); err != nil {
		switch err {
		case NoRecordsFound:
			break
		default:
			logger.GetLogger().Error("Error", zap.String("err", err.Error()))
			return uuid.UUID{}, fmt.Errorf("update error: %v", err)
		}
	}
	id, err := pr.Save(ctx, user)
	if err != nil {
		logger.GetLogger().Error("Updating error", zap.String("err", err.Error()))
		return uuid.UUID{}, fmt.Errorf("update error: %v", err)
	}
	logger.GetLogger().Info("Successfully updated user", zap.String("ID", id.String()))
	return id, nil
}

func (pr *PeopleRepository) Get(ctx context.Context, filter options.UserFilter) ([]users.EntityUser, error) {

	queryString := filter.Build()

	ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()
	rows, err := pr.db.Query(ctx, queryString)

	if err != nil {
		logger.GetLogger().Error("User fetching error", zap.String("err", err.Error()))
		return nil, fmt.Errorf("error while fetching users from the database")
	}

	fetchedUsers := make([]users.EntityUser, 0)

	for rows.Next() {
		user := &users.EntityUser{}

		if err := rows.Scan(&user.ID, &user.Name, &user.Surname, &user.Patronymic, &user.Age, &user.Sex); err != nil {
			logger.GetLogger().Error("Error while scanning row", zap.String("err", err.Error()))
			return nil, fmt.Errorf("error while scanning user: %v", err.Error())
		}
		fetchedUsers = append(fetchedUsers, *user)
	}

	if filter.GetNationalityOptions().Nationalities != nil {
		logger.GetLogger().Info("Nation options're received", zap.String("opts", fmt.Sprintf("%v", filter.GetNationalityOptions().Nationalities)), zap.String("len", fmt.Sprintf("%v", len(
			filter.GetNationalityOptions().Nationalities))))
		fetchedUsers, err = pr.getFilteredByNationalities(ctx, fetchedUsers, filter.GetNationalityOptions())
		if err != nil {
			logger.GetLogger().Error("Error while operating over nationalities of users", zap.String("err", err.Error()))
			return nil, fmt.Errorf("nationality operating error: %v", err.Error())
		}
	}

	return fetchedUsers, nil
}

func (pr *PeopleRepository) getFilteredByNationalities(ctx context.Context, fetched []users.EntityUser, nationalityOptions options.NationalityOptions) ([]users.EntityUser, error) {
	for i := 0; i < len(fetched); i++ {
		var nations []nationalities.Nationality
		rows, err := pr.db.Query(ctx, "SELECT nationality, probability from person_nationality where person_id = $1", fetched[i].ID)
		if err != nil {
			return nil, fmt.Errorf("error while fetching users's nationalities: %v", err.Error())
		}
		for rows.Next() {
			nationality := &nationalities.Nationality{}
			if err := rows.Scan(&nationality.ID, &nationality.Probability); err != nil {
				return nil, fmt.Errorf("error while scanning nationality row: %v", err.Error())
			}
			nations = append(nations, *nationality)
		}
		fetched[i].Nationality = nations
	}

	return nationalityOptions.FilterByNationality(fetched), nil
}

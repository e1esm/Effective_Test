package options

import (
	"fmt"
	"github.com/e1esm/Effective_Test/internal/models/users"
	"reflect"
	"strings"
)

type UserFilter interface {
	Build() string
	GetNationalityOptions() NationalityOptions
}

type QueryOptions struct {
	UserOptions
	NationalityOptions
}

type NationalityOptions struct {
	Nationalities []string `json:"nationality,omitempty"`
}

func (no *NationalityOptions) FilterByNationality(users []users.ExtendedUser) []users.ExtendedUser {
	return nil
}

type UserOptions struct {
	Sex    string `json:"sex,omitempty"`
	Age    int    `json:"age,omitempty"`
	Name   string `json:"name,omitempty"`
	Limit  int    `json:"limit"`
	Offset int    `json:"offset"`
}

func (uo *QueryOptions) Build() string {
	return uo.configureQueryString()
}

func (uo *QueryOptions) GetNationalityOptions() NationalityOptions {
	return uo.NationalityOptions
}

func (uo *QueryOptions) validate() map[string]string {
	availableOptions := make(map[string]string)

	switch {
	case uo.Sex != "":
		availableOptions[strings.ToLower(reflect.ValueOf(uo).Elem().Type().Field(0).Name)] = uo.Sex
	case uo.Age != 0:
		availableOptions[strings.ToLower(reflect.ValueOf(uo).Elem().Type().Field(1).Name)] = fmt.Sprintf("%v", uo.Age)
	case uo.Name != "":
		availableOptions[strings.ToLower(reflect.ValueOf(uo).Elem().Type().Field(2).Name)] = uo.Name
	}

	return availableOptions
}

func (uo *QueryOptions) configureQueryString() string {

	type pair struct {
		key   string
		value string
	}

	availableOptions := uo.validate()

	pairs := make([]pair, 0)

	for k, v := range availableOptions {
		pairs = append(pairs, pair{k, v})
	}

	switch len(pairs) {
	case 0:
		return fmt.Sprintf("SELECT * FROM people_info LIMIT %d OFFSET %d", uo.Limit, uo.Offset)
	case 1:
		return fmt.Sprintf("SELECT * FROM people_info WHERE %s = %v LIMIT %d OFFSET %d", pairs[0].key, pairs[0].value, uo.Limit, uo.Offset)
	case 2:
		return fmt.Sprintf("SELECT * FROM people_info WHERE %s = %v AND %s = %v LIMIT %d OFFSET %d",
			pairs[0].key,
			pairs[0].value,
			pairs[1].key,
			pairs[1].value,
			uo.Limit,
			uo.Offset)
	case 3:
		return fmt.Sprintf("SELECT * FROM people_info WHERE %s = %v AND %s = %v AND %s = %s LIMIT %d OFFSET %d",
			pairs[0].key,
			pairs[0].value,
			pairs[1].key,
			pairs[1].value,
			pairs[2].key,
			pairs[2].value,
			uo.Limit,
			uo.Offset)
	default:
		return ""
	}
}

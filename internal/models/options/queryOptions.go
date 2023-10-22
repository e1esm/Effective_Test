package options

import (
	"fmt"
	"github.com/e1esm/Effective_Test/internal/models/nationalities"
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

func NewQueryOptions(options UserOptions, nationalityOptions NationalityOptions) QueryOptions {
	return QueryOptions{options, nationalityOptions}
}

type NationalityOptions struct {
	Nationalities []string `json:"nationality,omitempty"`
}

func NewNationalityOptions(nations []string) NationalityOptions {
	if nations != nil && (len(nations) > 0 && nations[0] == "") {
		nations = nil
	} else {
		for i := 0; i < len(nations); i++ {
			nations[i] = strings.ToUpper(nations[i])
		}
	}
	return NationalityOptions{Nationalities: nations}
}

func (no *NationalityOptions) FilterByNationality(toBeFiltered []users.EntityUser) []users.EntityUser {
	filtered := make([]users.EntityUser, 0)
	for i := 0; i < len(toBeFiltered); i++ {
		if no.contains(toBeFiltered[i].Nationality) {
			filtered = append(filtered, toBeFiltered[i])
		}
	}

	return filtered
}

func (no *NationalityOptions) contains(nats []nationalities.Nationality) bool {
	for i := 0; i < len(nats); i++ {
		for j := 0; j < len(no.Nationalities); j++ {
			if nats[i].ID == no.Nationalities[j] {
				return true
			}
		}
	}

	return false
}

type UserOptions struct {
	Sex    string `json:"sex,omitempty"`
	Age    int    `json:"age,omitempty"`
	Name   string `json:"name,omitempty"`
	Limit  int    `json:"limit"`
	Offset int    `json:"offset"`
}

func NewUserOptions(sex, name string, age, limit, offset int) UserOptions {
	return UserOptions{sex, age, name, limit, offset}
}

func (uo *UserOptions) Build() string {
	return uo.configureQueryString()
}

func (qo *QueryOptions) GetNationalityOptions() NationalityOptions {
	return qo.NationalityOptions
}

func (uo *UserOptions) validate() map[string]string {
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

func (uo *UserOptions) configureQueryString() string {

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
		return fmt.Sprintf("SELECT uuid(id), name, surname, patronymic, age, sex FROM people_info LIMIT %d OFFSET %d", uo.Limit, uo.Offset)
	case 1:
		return fmt.Sprintf("SELECT uuid(id), name, surname, patronymic, age, sex FROM people_info WHERE %s = '%v' LIMIT %d OFFSET %d", pairs[0].key, pairs[0].value, uo.Limit, uo.Offset)
	case 2:
		return fmt.Sprintf("SELECT uuid(id), name, surname, patronymic, age, sex FROM people_info WHERE %s = '%v' AND %s = '%v' LIMIT %d OFFSET %d",
			pairs[0].key,
			pairs[0].value,
			pairs[1].key,
			pairs[1].value,
			uo.Limit,
			uo.Offset)
	case 3:
		return fmt.Sprintf("SELECT uuid(id), name, surname, patronymic, age, sex FROM people_info WHERE %s = '%v' AND %s = %v AND %s = '%s' LIMIT %d OFFSET %d",
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

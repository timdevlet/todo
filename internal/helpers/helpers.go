package helpers

import (
	"errors"
	"fmt"
	"math/rand"
	"reflect"
	"strings"
	"time"

	"github.com/Masterminds/squirrel"
	"github.com/brianvoe/gofakeit/v6"
	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	en_translations "github.com/go-playground/validator/v10/translations/en"
	"github.com/google/uuid"
	"github.com/timdevlet/todo/pkg/postgres"
	"golang.org/x/exp/constraints"
)

func Map[T any, R any](collection []T, iteratee func(item T, index int) R) []R {
	result := make([]R, len(collection))

	for i, item := range collection {
		result[i] = iteratee(item, i)
	}

	return result
}

func Ptr[T any](a T) *T {
	return &a
}

func Uuid() string {
	id := uuid.New()

	return id.String()
}

func Min[T constraints.Ordered](a, b T) T {
	if a < b {
		return a
	}

	return b
}

func InArray[T comparable](s T, arr []T) bool {
	for _, v := range arr {
		if v == s {
			return true
		}
	}

	return false
}

func Join(arr []string, del string) string {
	return strings.Join(arr, del)
}

//

func ValidationStruct[T any](str T) ([]string, bool) {
	en := en.New()
	uni := ut.New(en, en)

	// this is usually know or extracted from http 'Accept-Language' header
	// also see uni.FindTranslator(...)
	trans, _ := uni.GetTranslator("en")

	validate = validator.New()
	en_translations.RegisterDefaultTranslations(validate, trans)

	validate.RegisterTranslation("gte", trans, func(ut ut.Translator) error {
		return ut.Add("gte", "{0} must be gte then {1}", true) // see universal-translator for details
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T("gte", fe.Field(), fmt.Sprint(fe.Value()))

		return t
	})

	//

	err := validate.Struct(str)
	if err != nil {
		var errs validator.ValidationErrors

		if errors.As(err, &errs) {
			ers := Map(errs, func(e validator.FieldError, _ int) string {

				return fmt.Sprintf("[value:%v][field:%v][tag:%v] %v", e.Value(), e.Field(), e.Tag(), e.Translate(trans))
			})

			return ers, false
		}
	}

	return []string{}, true
}

//

func DateNow() string {
	return time.Now().Format(time.RFC3339)
}

func Default[T any](v *T, def T) T {
	if v == nil {
		return def
	}

	return *v
}

func IsValidUUID(u string) bool {
	_, err := uuid.Parse(u)
	return err == nil
}

//

func FakeSentence(length int) string {
	return gofakeit.Sentence(length)[0 : length-1]
}

func FakeEmail() string {
	return gofakeit.Email()
}

func FakeName() string {
	return gofakeit.Name()
}

//

func RandomFromSlice[T any](arr []T) (bool, T) {
	var result T

	if len(arr) == 0 {
		return false, result
	}

	idx := rand.Intn(len(arr) - 1)

	return true, arr[idx]
}

func UuidByHash(s string) string {
	id := uuid.NewSHA1(uuid.UUID{}, []byte(s))

	return id.String()
}

func UuidByTwoStrings(s1 string, s2 string) string {
	sorted := ""
	if s1 < s2 {
		sorted = s1 + "." + s2
	} else {
		sorted = s2 + "." + s1
	}

	id := uuid.NewSHA1(uuid.UUID{}, []byte(sorted))

	return id.String()
}

// postgres

func FetchManyFromPostgres[T any](postgres *postgres.PDB, table string) ([]T, error) {

	fields := getSqlFields(new(T))

	statement := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar).
		Select(strings.Join(fields, ", ")).
		From(table).
		Limit(uint64(100000))

	sql, values, _ := statement.ToSql()
	rows, err := postgres.DB.Query(sql, values...)
	if err != nil {
		return []T{}, err
	}
	defer rows.Close()

	result := make([]T, 0)
	for rows.Next() {
		var i T

		c := new(T)
		l := getSqlColumnToFieldMap(c)

		err := rows.Scan(l...)
		if err != nil {
			return []T{}, err
		}
		result = append(result, i)
	}

	if err = rows.Err(); err != nil {
		return []T{}, err
	}

	return result, nil
}

func FetchOneFromPostgres[T any](postgres *postgres.PDB, table string, uuid string) (*T, error) {

	fields := getSqlFields(new(T))

	statement := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar).
		Select(strings.Join(fields, ", ")).
		From(table).
		Limit(1)

	sql, values, _ := statement.ToSql()
	rows, err := postgres.DB.Query(sql, values...)
	if err != nil {
		return new(T), err
	}
	defer rows.Close()

	result := []T{}
	for rows.Next() {
		var i T

		c := new(T)
		l := getSqlColumnToFieldMap(c)

		err := rows.Scan(l...)
		if err != nil {
			return new(T), err
		}
		result = append(result, i)
	}

	if len(result) == 0 {
		return new(T), fmt.Errorf("[FetchOneFromPostgres][uuid:%s] not found by uuid", uuid)
	}

	if err = rows.Err(); err != nil {
		return new(T), err
	}

	return &result[0], nil
}

func getSqlFields(model interface{}) []string {
	t := reflect.TypeOf(model).Elem()
	v := reflect.ValueOf(model).Elem()

	r := []string{}

	for i := 0; i < v.NumField(); i++ {
		colName := t.Field(i).Tag.Get("sql")
		r = append(r, colName)
	}
	return r
}

func getSqlColumnToFieldMap(model interface{}) []interface{} {

	v := reflect.ValueOf(model).Elem()
	r := []any{}

	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)
		r = append(r, field.Addr().Interface())
	}
	return r
}

// Validations

var validate *validator.Validate = validator.New()

func ValidateEmail(s string) error {
	err := validate.Var(s, "required,email")
	if err != nil {
		return fmt.Errorf("[email:%s] email is not valid", s)
	}

	return nil
}

//

func GetType(s interface{}) string {
	if t := reflect.TypeOf(s); t.Kind() == reflect.Ptr {
		return t.Elem().Name()
	} else {
		return t.Name()
	}
}

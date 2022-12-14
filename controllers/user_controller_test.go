package controllers

import (
	"bytes"
	"database/sql/driver"
	"encoding/json"
	"learn_testing/config"
	"learn_testing/models"
	"net/http"
	"net/http/httptest"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func TestGetUsersController(t *testing.T) {
	dbFakeGorm, mocked, err := sqlmock.New()

	assert.NoError(t, err)

	dbGorm, _ := gorm.Open(mysql.New(mysql.Config{
		SkipInitializeWithVersion: true,
		Conn:                      dbFakeGorm,
	}))

	config.DB = dbGorm

	row := sqlmock.NewRows([]string{"name", "email"}).
		AddRow("ahmad naufal", "ahmad@gmail.com")

	mocked.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `users` WHERE `users`.`deleted_at` IS NULL")).
		WillReturnRows(row)

	testCase := []struct {
		Name             string
		ExpectStatusCode int
		Method           string
		Body             models.Users
		HasReturnBody    bool
		ExpectBody       models.Users
	}{
		{
			"success",
			http.StatusOK,
			"GET",
			models.Users{},
			false,
			models.Users{},
		},
	}

	for _, val := range testCase {
		t.Run(val.Name, func(t *testing.T) {
			r := httptest.NewRequest(val.Method, "/", nil)
			w := httptest.NewRecorder()

			e := echo.New()
			ctx := e.NewContext(r, w)

			err := GetUsersController(ctx)
			assert.NoError(t, err)

			assert.Equal(t, val.ExpectStatusCode, w.Result().StatusCode)

			if val.HasReturnBody {
				var response map[string]interface{}
				err := json.NewDecoder(w.Result().Body).Decode(&response)

				assert.NoError(t, err)
				assert.Equal(t, val.ExpectBody.Name, response["user"].(map[string]interface{})["name"])
			}
		})
	}
}

func TestGetUserController(t *testing.T) {
	dbFakeGorm, mocked, err := sqlmock.New()

	assert.NoError(t, err)

	dbGorm, _ := gorm.Open(mysql.New(mysql.Config{
		SkipInitializeWithVersion: true,
		Conn:                      dbFakeGorm,
	}))

	config.DB = dbGorm

	row := sqlmock.NewRows([]string{"name", "email", "password"}).
		AddRow("ahmad naufal", "ahmad@gmail.com", "alta@1234")

	mocked.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `users` WHERE id = ? AND `users`.`deleted_at` IS NULL")).
		WithArgs(1).
		WillReturnRows(row)

	testCase := []struct {
		Name             string
		ExpectStatusCode int
		Method           string
		Body             models.Users
		HasReturnBody    bool
		ExpectBody       models.Users
	}{
		{
			"success",
			http.StatusOK,
			"GET",
			models.Users{},
			true,
			models.Users{
				Name: "ahmad naufal",
			},
		},
	}

	for _, val := range testCase {
		t.Run(val.Name, func(t *testing.T) {
			res, _ := json.Marshal(val.Body)
			r := httptest.NewRequest(val.Method, "/", bytes.NewBuffer(res))
			w := httptest.NewRecorder()

			e := echo.New()
			ctx := e.NewContext(r, w)
			ctx.SetPath("/:id")
			ctx.SetParamNames("id")
			ctx.SetParamValues("1")

			err := GetUserController(ctx)
			assert.NoError(t, err)

			assert.Equal(t, val.ExpectStatusCode, w.Result().StatusCode)

			if val.HasReturnBody {
				var response map[string]interface{}
				err := json.NewDecoder(w.Result().Body).Decode(&response)

				assert.NoError(t, err)
				assert.Equal(t, val.ExpectBody.Name, response["user"].(map[string]interface{})["name"])
			}
		})
	}
}

type AnyTime struct{}

func (a AnyTime) Match(v driver.Value) bool {
	_, ok := v.(time.Time)
	return ok
}

func TestCreateUserController(t *testing.T) {
	dbFakeGorm, mocked, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))

	assert.NoError(t, err)

	dbGorm, _ := gorm.Open(mysql.New(mysql.Config{
		SkipInitializeWithVersion: true,
		Conn:                      dbFakeGorm,
	}))

	config.DB = dbGorm

	mocked.ExpectBegin()
	mocked.ExpectExec(regexp.QuoteMeta("INSERT INTO `users` (`created_at`,`updated_at`,`deleted_at`,`name`,`email`,`password`) VALUES (?,?,?,?,?,?)")).
		WithArgs(AnyTime{}, AnyTime{}, "NULL", "ahmad naufal", "ahmad@gmail.com", "alta@1234").
		WillReturnResult(sqlmock.NewErrorResult(nil))

	mocked.ExpectCommit()
	testCase := []struct {
		Name             string
		ExpectStatusCode int
		Method           string
		Body             models.Users
		HasReturnBody    bool
		ExpectBody       string
	}{
		{
			"success",
			http.StatusOK,
			"POST",
			models.Users{
				Name:     "ahmad naufal",
				Email:    "ahmad@gmail.com",
				Password: "alta@1234",
			},
			true,
			"success create new users",
		},
	}

	for _, val := range testCase {
		t.Run(val.Name, func(t *testing.T) {
			res, _ := json.Marshal(val.Body)
			r := httptest.NewRequest(val.Method, "/", bytes.NewBuffer(res))
			w := httptest.NewRecorder()

			e := echo.New()
			ctx := e.NewContext(r, w)

			err := CreateUserController(ctx)
			assert.NoError(t, err)

			assert.Equal(t, val.ExpectStatusCode, w.Result().StatusCode)

			if val.HasReturnBody {
				var response map[string]interface{}
				err := json.NewDecoder(w.Result().Body).Decode(&response)

				assert.NoError(t, err)
				assert.Equal(t, val.ExpectBody, response["message"])
			}
		})
	}
}

func TestDeleteUserController(t *testing.T) {
	// mocking
	dbFakeGorm, mocked, err := sqlmock.New()

	assert.NoError(t, err)

	dbGorm, _ := gorm.Open(mysql.New(mysql.Config{
		SkipInitializeWithVersion: true,
		Conn:                      dbFakeGorm,
	}))

	config.DB = dbGorm

	mocked.ExpectBegin()

	mocked.ExpectExec(regexp.QuoteMeta("DELETE FROM `users` WHERE id = ?")).
		WithArgs(1).
		WillReturnResult(sqlmock.NewResult(1, 0)).
		WillReturnError(nil)

	mocked.ExpectCommit()

	testCase := []struct {
		Name             string
		ExpectStatusCode int
		Method           string
		HasReturnBody    bool
		ExpectBody       string
	}{
		{
			"success",
			http.StatusOK,
			"DELETE",
			true,
			"success deleted user by id",
		},
	}

	for _, val := range testCase {
		t.Run(val.Name, func(t *testing.T) {
			r := httptest.NewRequest(val.Method, "/", nil)
			w := httptest.NewRecorder()

			// handler echo
			e := echo.New()
			ctx := e.NewContext(r, w)
			ctx.SetPath("/:id")
			ctx.SetParamNames("id")
			ctx.SetParamValues("1")

			err := DeleteUserController(ctx)
			assert.NoError(t, err)

			assert.Equal(t, val.ExpectStatusCode, w.Result().StatusCode)

			if val.HasReturnBody {
				var response map[string]interface{}
				err := json.NewDecoder(w.Result().Body).Decode(&response)

				assert.NoError(t, err)
				assert.Equal(t, val.ExpectStatusCode, w.Result().StatusCode)
				assert.Equal(t, val.ExpectBody, response["message"])
			}
		})
	}
}

func TestUpdateUserController(t *testing.T) {
	dbFakeGorm, mocked, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))

	assert.NoError(t, err)

	dbGorm, _ := gorm.Open(mysql.New(mysql.Config{
		SkipInitializeWithVersion: true,
		Conn:                      dbFakeGorm,
	}))

	config.DB = dbGorm

	mocked.ExpectBegin()
	mocked.ExpectExec(regexp.QuoteMeta("UPDATE `users` SET `updated_at`=? WHERE id = ? AND `users`.`deleted_at` IS NULL")).
		WithArgs(AnyTime{}, 1).
		WillReturnResult(sqlmock.NewResult(1, 1))

	mocked.ExpectCommit()
	testCase := []struct {
		Name             string
		ExpectStatusCode int
		Method           string
		Body             models.Users
		HasReturnBody    bool
		ExpectBody       string
	}{
		{
			"success",
			http.StatusOK,
			"POST",
			models.Users{
				Name:     "ahmad naufal",
				Email:    "ahmad@gmail.com",
				Password: "alta@1234",
			},
			true,
			"success updated user by id",
		},
	}

	for _, val := range testCase {
		t.Run(val.Name, func(t *testing.T) {
			res, _ := json.Marshal(val.Body)
			r := httptest.NewRequest(val.Method, "/", bytes.NewBuffer(res))
			w := httptest.NewRecorder()

			e := echo.New()
			ctx := e.NewContext(r, w)
			ctx.SetPath("/:id")
			ctx.SetParamNames("id")
			ctx.SetParamValues("1")

			err := UpdateUserController(ctx)
			assert.NoError(t, err)

			assert.Equal(t, val.ExpectStatusCode, w.Result().StatusCode)

			if val.HasReturnBody {
				var response map[string]interface{}
				err := json.NewDecoder(w.Result().Body).Decode(&response)

				assert.NoError(t, err)
				assert.Equal(t, val.ExpectBody, response["message"])
			}
		})
	}
}

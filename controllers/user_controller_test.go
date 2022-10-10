package controllers

import (
	"bytes"
	"encoding/json"
	"learn_testing/config"
	"learn_testing/models"
	"net/http"
	"net/http/httptest"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func TestGetUsersController(t *testing.T) {
	dbFakeGorm, mock, err := sqlmock.New()

	assert.NoError(t, err)

	dbGorm, _ := gorm.Open(mysql.New(mysql.Config{
		SkipInitializeWithVersion: true,
		Conn:                      dbFakeGorm,
	}))

	config.DB = dbGorm

	row := sqlmock.NewRows([]string{"name", "email"}).
		AddRow("ahmad naufal", "ahmad@gmail.com")

	mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `users` WHERE `users`.`deleted_at` IS NULL")).
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
				var response map[string]models.Users
				err := json.NewDecoder(w.Result().Body).Decode(&response)

				assert.NoError(t, err)
				assert.Equal(t, val.ExpectBody.Name, response["user"].Name)
			}
		})
	}
}

func TestGetUserController(t *testing.T) {
	dbFakeGorm, mock, err := sqlmock.New()

	assert.NoError(t, err)

	dbGorm, _ := gorm.Open(mysql.New(mysql.Config{
		SkipInitializeWithVersion: true,
		Conn:                      dbFakeGorm,
	}))

	config.DB = dbGorm

	row := sqlmock.NewRows([]string{"name", "email"}).
		AddRow("ahmad naufal", "ahmad@gmail.com")

	mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `users` WHERE id = ? AND `users`.`deleted_at` IS NULL")).
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
			models.Users{
				Name: "ahmad",
			},
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
				var response map[string]models.Users
				err := json.NewDecoder(w.Result().Body).Decode(&response)

				assert.NoError(t, err)
				assert.Equal(t, val.ExpectBody.Name, response["user"].Name)
			}
		})
	}
}

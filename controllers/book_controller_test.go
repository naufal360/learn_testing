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
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func TestGetBooksController(t *testing.T) {
	dbFakeGorm, mocked, err := sqlmock.New()

	assert.NoError(t, err)

	dbGorm, _ := gorm.Open(mysql.New(mysql.Config{
		SkipInitializeWithVersion: true,
		Conn:                      dbFakeGorm,
	}))

	config.DB = dbGorm

	row := sqlmock.NewRows([]string{"title", "publisher"}).
		AddRow("jalan jalan", "gramed")

	mocked.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `books` WHERE `books`.`deleted_at` IS NULL")).
		WillReturnRows(row)

	testCase := []struct {
		Name             string
		ExpectStatusCode int
		Method           string
		Body             models.Books
		HasReturnBody    bool
		ExpectBody       models.Books
	}{
		{
			"success",
			http.StatusOK,
			"GET",
			models.Books{},
			false,
			models.Books{},
		},
	}

	for _, val := range testCase {
		t.Run(val.Name, func(t *testing.T) {
			r := httptest.NewRequest(val.Method, "/", nil)
			w := httptest.NewRecorder()

			e := echo.New()
			ctx := e.NewContext(r, w)

			err := GetBooksController(ctx)
			assert.NoError(t, err)

			assert.Equal(t, val.ExpectStatusCode, w.Result().StatusCode)

			if val.HasReturnBody {
				var response map[string]interface{}
				err := json.NewDecoder(w.Result().Body).Decode(&response)

				assert.NoError(t, err)
				assert.Equal(t, val.ExpectBody.Title, response["books"].(map[string]interface{})["title"])
			}
		})
	}
}

func TestGetBookController(t *testing.T) {
	dbFakeGorm, mocked, err := sqlmock.New()

	assert.NoError(t, err)

	dbGorm, _ := gorm.Open(mysql.New(mysql.Config{
		SkipInitializeWithVersion: true,
		Conn:                      dbFakeGorm,
	}))

	config.DB = dbGorm

	row := sqlmock.NewRows([]string{"title", "publisher", "author"}).
		AddRow("jalan jalan", "gramed", "ahmad")

	mocked.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `books` WHERE id = ? AND `books`.`deleted_at` IS NULL")).
		WithArgs(1).
		WillReturnRows(row)

	testCase := []struct {
		Name             string
		ExpectStatusCode int
		Method           string
		Body             models.Books
		HasReturnBody    bool
		ExpectBody       models.Books
	}{
		{
			"success",
			http.StatusOK,
			"GET",
			models.Books{},
			true,
			models.Books{
				Title:     "jalan jalan",
				Publisher: "gramed",
				Author:    "ahmad",
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

			err := GetBookController(ctx)
			assert.NoError(t, err)

			assert.Equal(t, val.ExpectStatusCode, w.Result().StatusCode)

			if val.HasReturnBody {
				var response map[string]interface{}
				err := json.NewDecoder(w.Result().Body).Decode(&response)

				assert.NoError(t, err)
				assert.Equal(t, val.ExpectBody.Title, response["book"].(map[string]interface{})["title"])
			}
		})
	}
}

func TestCreateBookController(t *testing.T) {
	dbFakeGorm, mocked, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))

	assert.NoError(t, err)

	dbGorm, _ := gorm.Open(mysql.New(mysql.Config{
		SkipInitializeWithVersion: true,
		Conn:                      dbFakeGorm,
	}))

	config.DB = dbGorm

	mocked.ExpectBegin()
	mocked.ExpectExec(regexp.QuoteMeta("INSERT INTO `books` (`created_at`,`updated_at`,`deleted_at`,`title`,`publisher`,`author`) VALUES (?,?,?,?,?,?)")).
		WithArgs(time.Now(), time.Now(), "NULL", "jalan jalan", "gramed", "ahmad").
		WillReturnResult(sqlmock.NewErrorResult(nil))

	mocked.ExpectCommit()
	testCase := []struct {
		Name             string
		ExpectStatusCode int
		Method           string
		Body             models.Books
		HasReturnBody    bool
		ExpectBody       string
	}{
		{
			"success",
			http.StatusOK,
			"POST",
			models.Books{
				Title:     "jalan jalan",
				Publisher: "gramed",
				Author:    "ahmad",
			},
			true,
			"success create new books",
		},
	}

	for _, val := range testCase {
		t.Run(val.Name, func(t *testing.T) {
			res, _ := json.Marshal(val.Body)
			r := httptest.NewRequest(val.Method, "/", bytes.NewBuffer(res))
			w := httptest.NewRecorder()

			e := echo.New()
			ctx := e.NewContext(r, w)

			err := CreateBookController(ctx)
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

func TestDeleteBookController(t *testing.T) {
	// mocking
	dbFakeGorm, mocked, err := sqlmock.New()

	assert.NoError(t, err)

	dbGorm, _ := gorm.Open(mysql.New(mysql.Config{
		SkipInitializeWithVersion: true,
		Conn:                      dbFakeGorm,
	}))

	config.DB = dbGorm

	mocked.ExpectBegin()

	mocked.ExpectExec(regexp.QuoteMeta("DELETE FROM `books` WHERE id = ?")).
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
			"success deleted book by id",
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

			err := DeleteBookController(ctx)
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

func TestUpdateBookController(t *testing.T) {
	dbFakeGorm, mocked, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))

	assert.NoError(t, err)

	dbGorm, _ := gorm.Open(mysql.New(mysql.Config{
		SkipInitializeWithVersion: true,
		Conn:                      dbFakeGorm,
	}))

	config.DB = dbGorm

	mocked.ExpectBegin()
	mocked.ExpectExec(regexp.QuoteMeta("UPDATE `book` SET `updated_at`=? WHERE id = ? AND `books`.`deleted_at` IS NULL")).
		WithArgs(time.Now(), 1).
		WillReturnResult(sqlmock.NewResult(1, 1))

	mocked.ExpectCommit()
	testCase := []struct {
		Name             string
		ExpectStatusCode int
		Method           string
		Body             models.Books
		HasReturnBody    bool
		ExpectBody       string
	}{
		{
			"success",
			http.StatusOK,
			"POST",
			models.Books{
				Title:     "jalan jalan",
				Publisher: "gramed",
				Author:    "ahmad",
			},
			true,
			"success updated book by id",
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

			err := UpdateBookController(ctx)
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

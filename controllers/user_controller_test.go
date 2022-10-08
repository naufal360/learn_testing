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
	"github.com/stretchr/testify/suite"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type suiteUsers struct {
	suite.Suite
	mocking sqlmock.Sqlmock
	// testCase []struct{}
}

func (s *suiteUsers) SetupSuite() {
	dbFakeGorm, mocking, err := sqlmock.New()

	s.NoError(err)

	dbGorm, _ := gorm.Open(mysql.New(mysql.Config{
		SkipInitializeWithVersion: true,
		Conn:                      dbFakeGorm,
	}))

	config.DB = dbGorm

	s.mocking = mocking
}

func (s *suiteUsers) TestGetUsersController() {
	row := sqlmock.NewRows([]string{"name", "email"}).
		AddRow("ahmad naufal", "ahmad@gmail.com")

	s.mocking.ExpectQuery(regexp.QuoteMeta("SELECT `user`.`name`,`user`.`email` FROM `user` WHERE id = ? AND `user`.`deleted_at` IS NULL")).
		WithArgs(1).
		WillReturnRows(row)

	testCase := []struct {
		Name               string
		ExpectedStatusCode int
		Method             string
		Body               models.User
		HasReturnBody      bool
		ExpectedBody       models.User
	}{
		{
			"success",
			http.StatusOK,
			http.MethodGet,
			models.User{
				Name: "ahmad",
			},
			true,
			models.User{
				Name: "ahmad",
			},
		},
	}

	for _, v := range testCase {
		s.T().Run(v.Name, func(t *testing.T) {
			res, _ := json.Marshal(v.Body)
			r := httptest.NewRequest(v.Method, "/", bytes.NewBuffer(res))
			w := httptest.NewRecorder()

			e := echo.New()
			ctx := e.NewContext(r, w)
			ctx.SetPath("/")
			ctx.SetPath("/:id")
			ctx.SetParamNames("id")
			ctx.SetParamValues("1")

			err := GetUsersController(ctx)

			s.NoError(err)

			s.Equal(v.ExpectedBody, w.Result().StatusCode)

			if v.HasReturnBody {
				var response map[string]models.User

				err := json.NewDecoder(w.Result().Body).Decode(&response)

				s.NoError(err)

				s.Equal(v.ExpectedBody.Name, response["user"].Name)
			}
		})
	}
}

func (s *suiteUsers) TearDownSuite() {
	config.DB = nil
	s.mocking = nil
}

func TestSuiteUsers(t *testing.T) {
	suite.Run(t, new(suiteUsers))
}

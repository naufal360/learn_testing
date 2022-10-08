package routes

import (
	"learn_testing/config"
	c "learn_testing/controllers"
	m "learn_testing/middleware"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func New() *echo.Echo {

	e := echo.New()

	// ROUTING
	// version
	v1 := e.Group("/v1")
	// Logging
	m.LogMiddleware(e)

	// // routing /users to handler function
	v1.POST("/users", c.CreateUserController)
	v1.POST("/login", c.LoginUserController)

	// // routing /book to handler function
	v1.GET("/books", c.GetBooksController)
	v1.GET("/books/:id", c.GetBookController)

	// JWT AUTH
	jwtAuthV1 := v1.Group("")
	jwtAuthV1.Use(middleware.JWT([]byte(config.ViperEnvVariable("SECRET_KEY"))))

	// // routing /auth/users to handler function
	jwtAuthV1.GET("/users", c.GetUsersController)
	jwtAuthV1.GET("/users/:id", c.GetUserController)
	jwtAuthV1.DELETE("/users/:id", c.DeleteUserController)
	jwtAuthV1.PUT("/users/:id", c.UpdateUserController)

	// routing /auth//books to handler function
	jwtAuthV1.POST("/books", c.CreateBookController)
	jwtAuthV1.DELETE("/books/:id", c.DeleteBookController)
	jwtAuthV1.PUT("/books/:id", c.UpdateBookController)

	return e
}

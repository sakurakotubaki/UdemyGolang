package main

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)


func main()  {
	e := echo.New()

	e.GET("/", func(c echo.Context) error {
		return c.String(200, "Hello, World!")
	})

	e.GET("/users/", func(c echo.Context) error {
		return c.String(200, "Hello, Users!")
	})
	
	e.Use(middleware.Logger())

	e.Start(":8080")
}
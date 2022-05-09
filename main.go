package main

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/simon7426/blog-go-echo-mongo/configs"
	"github.com/simon7426/blog-go-echo-mongo/routes"
)

func main() {
	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.CORS())

	// e.GET("/", func(c echo.Context) error {
	// 	return c.JSON(200, &echo.Map{
	// 		"data": "Hello from echo.",
	// 	})
	// })

	configs.ConnectDB()

	routes.BlogRoutes(e)

	e.Logger.Fatal(e.Start(":7000"))

}

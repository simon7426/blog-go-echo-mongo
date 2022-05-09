package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"time"

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

	go func() {
		if err := e.Start(":7000"); err != nil && err != http.ErrServerClosed {
			e.Logger.Fatal("Shutting down the server...")
		}
	}()

	quit := make(chan os.Signal, 1)

	signal.Notify(quit, os.Interrupt)
	<-quit

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := e.Shutdown(ctx); err != nil {
		e.Logger.Fatal(err)
	}

}

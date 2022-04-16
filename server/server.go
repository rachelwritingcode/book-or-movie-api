package server

import (
	handlers "book-or-movie-api/handlers"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

func Server() {
	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.GET("/getmovie", handlers.GetMovie)
	e.GET("/getbook", handlers.GetBook)
	e.GET("/getrecommendation", handlers.GetReccomendation)

	e.Logger.Fatal(e.Start(":8080"))
}

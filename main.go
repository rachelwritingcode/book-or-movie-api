package main

import (
	"io"
	"log"
	"net/http"
	"os"

	"github.com/jeremywohl/flatten"
	"github.com/joho/godotenv"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

func main() {

	// Web Server
	e := echo.New()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// Routes
	// TODO add functionality for user to pass in a parameter to the route
	e.GET("/book", getBook)

	// Start server
	e.Logger.Fatal(e.Start(":8080"))
}

func getBook(c echo.Context) error {

	// Load Environment Variables
	err := godotenv.Load("local.env")
	if err != nil {
		log.Fatalln(err)
	}

	// Environment Variables
	nyt_books_api := os.Getenv("NYT_BOOKS_API")
	nyt_api_key := os.Getenv("NYT_API_KEY")

	// TODO Pass in the request parameter for &api-key&title= title that user selects
	req, err := http.NewRequest("GET", nyt_books_api, nil)
	if err != nil {
		log.Fatalln(err)
	}
	req.Header.Set("Accept", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalln(err)
	}

	defer resp.Body.Close()

	b, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}
	responseString := string(b)
	flat, err := flatten.FlattenString(responseString, "", flatten.DotStyle)

	return c.String(http.StatusOK, flat)
}

// func getMovieByTitle(c echo.Context) error {
// 	return c.String(http.StatusOK, "Hello, World!")
// }

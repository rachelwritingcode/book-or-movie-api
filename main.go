package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/buger/jsonparser"
	"github.com/jeremywohl/flatten"
	"github.com/joho/godotenv"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

type Movie struct {
	Title string `json:title`
}

type BookResponse struct {
	Title  string `json:title`
	Rating int    `json:rating`
	Plot   string `json:plot`
}

type Book struct {
	Title string `json:title`
}

type MovieResponse struct {
	Title    string   `json:title`
	Rating   int      `json:rating`
	Actors   []string `json:actors`
	Director string   `json:director`
	Plot     string   `json:plot`
}

func main() {

	// Web Server
	e := echo.New()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// Routes
	e.GET("/getmovie", getMovie)
	// e.GET("/getbook", getBook)

	e.GET("/status", func(c echo.Context) error {
		return c.HTML(
			http.StatusOK,
			"<h1>API is running</h1>",
		)
	})

	// TODO Provide a Recommendation based on the Ratings
	// TODO Read the book/movie

	// Start server
	e.Logger.Fatal(e.Start(":8080"))
}

// func getBook(c echo.Context) error {
// }

func getMovie(c echo.Context) error {

	// Title Request Parameter
	title := c.QueryParam("title")
	fmt.Println("Title Query: " + title)

	// Load Environment Variables
	err := godotenv.Load("local.env")
	if err != nil {
		log.Fatalln(err)
	}

	// Environment Variables
	omdb_api := os.Getenv("OMDB_API")
	omdb_key := os.Getenv("OMDB_KEY")

	request_url := omdb_api + "=" + omdb_key + "&t=" + title

	fmt.Println("Request URL: " + request_url)

	req, err := http.NewRequest("GET", request_url, nil)
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

	// Parse out Rating
	director, _ := jsonparser.GetString(b, "Director")
	fmt.Println(director)

	plot, _ := jsonparser.GetString(b, "Plot")
	fmt.Println(plot)

	actors, _ := jsonparser.GetString(b, "Actors")
	fmt.Println(actors)

	rotten_tomatoes, _ := jsonparser.GetString(b, "Ratings.1.Value")
	fmt.Println(rotten_tomatoes)

	metacritic, _ := jsonparser.GetString(b, "Ratings.2.Value")
	fmt.Println(metacritic)

	// TODO Add the data to the custom struct for MovieReview

	// Title
	// Actors
	// Director
	// Plot
	// Return a custom movie response
	// TODO Retrieve Book Review Information

	responseString := string(b)
	flat, err := flatten.FlattenString(responseString, "", flatten.DotStyle)

	return c.String(http.StatusOK, flat)
}

func getBookReview(title string) string {
	review := ""
	// TODO Call the Book Review API
	// Pass in the title request query parameter
	return review
}

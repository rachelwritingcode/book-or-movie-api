package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/buger/jsonparser"
	"github.com/gocolly/colly"
	"github.com/joho/godotenv"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

type BookResponse struct {
	Title  string `json:title`
	Author string `json:author`
	Rating string `json:rating`
	Plot   string `json:plot`
	Year   string `json:year`
}

type MovieResponse struct {
	Title       string `json:title`
	Rating      string `json:rating`
	Actors      string `json:actors`
	Director    string `json:director`
	Plot        string `json:plot`
	ReleaseYear string `json:release year`
	Awards      string `json:awards`
}

func main() {

	// Web Server
	e := echo.New()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// Routes
	e.GET("/getmovie", getMovie)
	e.GET("/getbook", getBook)
	// e.GET("/recommend", getRecommendation)

	e.GET("/status", func(c echo.Context) error {
		return c.HTML(
			http.StatusOK,
			"<h1>API is running</h1>",
		)
	})

	// Start server
	e.Logger.Fatal(e.Start(":8080"))
}

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

	// TODO Remove after debugging
	fmt.Println("\n OMDB API RESULTS " + string(b) + "\n")
	// Parse out Rating
	director, _ := jsonparser.GetString(b, "Director")
	plot, _ := jsonparser.GetString(b, "Plot")
	actors, _ := jsonparser.GetString(b, "Actors")
	rating, _ := jsonparser.GetString(b, "imdbRating")
	year, _ := jsonparser.GetString(b, "Year")
	// TODO Escape special character & for awards field
	awards, _ := jsonparser.GetString(b, "Awards")

	movieResponseData := &MovieResponse{
		title,
		rating,
		actors,
		director,
		plot,
		year,
		awards}

	// Convert struct data to string
	movieResponse, _ := json.Marshal(movieResponseData)
	return c.String(http.StatusOK, string(movieResponse))

}

func getBook(c echo.Context) error {
	// TODO Call the Book Review API
	// Pass in the title request query parameter
	// Title Request Parameter
	title := c.QueryParam("title")
	fmt.Println("Title Query: " + title)
	//TODO if title has two words, then add the ASCII %20 symbol
	//https://openlibrary.org/search.json?title=the%20red%20scrolls%20of%20magic
	request_url := "http://openlibrary.org/search.json?title=" + title
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

	fmt.Println(string(b))

	open_library_id, _ := jsonparser.GetString(b, "docs", "[0]", "key")
	fmt.Println("open_library_id : " + open_library_id)

	// Scrape Open Library for rating and plot information
	rating := ""
	var plot string = ""
	var ratings = []string{}
	web_scraper := colly.NewCollector(
		colly.AllowedDomains("https://openlibrary.org", "openlibrary.org"),
	)

	web_scraper.OnHTML(".readers-stats ", func(e *colly.HTMLElement) {
		e.ForEach("li", func(_ int, elem *colly.HTMLElement) {
			ratings = append(ratings, e.DOM.Find("span").Text())
		})
	})

	web_scraper.OnHTML(".book-description-content", func(e *colly.HTMLElement) {
		e.ForEach("p", func(_ int, elem *colly.HTMLElement) {
			plot = e.Text
		})
	})

	web_scraper.Limit(&colly.LimitRule{
		DomainGlob:  "*",
		RandomDelay: 1 * time.Second,
	})

	web_scraper.OnRequest(func(request *colly.Request) {
		fmt.Println("Visiting...https://openlibrary.org" + open_library_id)
	})

	web_scraper.Visit("https://openlibrary.org" + open_library_id)

	rating = strings.ReplaceAll(ratings[len(ratings)-1], "★", "")
	rating_float, _ := strconv.ParseFloat(rating, 64)
	rating_float = math.Ceil(rating_float*100) / 100
	rating = strconv.FormatFloat(rating_float, 'E', -1, 64)
	rating = strings.ReplaceAll(rating, "E+00", "")

	year, _ := jsonparser.GetString(b, "docs", "[0]", "first_publish_year")
	author, _ := jsonparser.GetString(b, "docs", "[0]", "author_name")

	bookResponse := &BookResponse{
		title,
		author,
		rating,
		plot,
		year,
	}
	bookResponseData, _ := json.Marshal(bookResponse)
	return c.String(http.StatusOK, string(bookResponseData))
	// 	return c.String(http.StatusOK, "Test")
}

// getRecommendation(c echo.Context)error{
// }

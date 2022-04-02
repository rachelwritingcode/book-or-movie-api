package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/buger/jsonparser"
	"github.com/gocolly/colly"
	"github.com/joho/godotenv"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

// type BookResponse struct {
// 	Title  string `json:title`
//	Author string `json:author`
// 	Rating int    `json:rating`
// 	Plot   string `json:plot`
//	Year 	string `json:year`
// }

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

	// TODO Remove after debugging
	open_library_id, _ := jsonparser.GetString(b, "docs", "[0]", "key")

	year, _ := jsonparser.GetString(b, "docs", "[0]", "first_publish_year")
	fmt.Println("Book Publishing Year: " + year)

	author, _ := jsonparser.GetString(b, "docs", "[0]", "author_name")
	fmt.Println("Author: " + author)

	// TODO Web Scraping
	// TODO webscrape overview information for book plot
	// TODO webscrape the book star ratings

	// Instantiate default collector
	web_scraper := colly.NewCollector(
		// Visit only domain
		colly.AllowedDomains("openlibrary.org"),
	)

	// Limit the number of threads started by colly to two
	// when visiting links which domains' matches "*httpbin.*" glob
	web_scraper.Limit(&colly.LimitRule{
		DomainGlob:  "*httpbin.*",
		Parallelism: 1,
		Delay:       20 * time.Second,
	})

	// On every a element which has href attribute call callback
	web_scraper.OnHTML("a[href]", func(e *colly.HTMLElement) {
		link := e.Attr("href")
		// Print link
		fmt.Printf("Link found: %q -> %s\n", e.Text, link)
		// Visit link found on page
		// Only those links are visited which are in AllowedDomains
		web_scraper.Visit(e.Request.AbsoluteURL(link))
	})

	// Before making a request print "Visiting ..."
	web_scraper.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL.String())
	})
	//TODO Set web scraping parameters!

	// Start scraping the website
	web_scraper.Visit("https://openlibrary.org/" + open_library_id)

	// bookResponse, _ := json.Marshal()
	return c.String(http.StatusOK, "Open Library URL ID: "+open_library_id)
}

// getRecommendation(c echo.Context)error{

// }

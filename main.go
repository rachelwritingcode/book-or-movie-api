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

const openLibraryAPI = "http://openlibrary.org/search.json?title="
const openLibraryURL = "http://openlibrary.org"

// Recommendation messages
const noBookFound string = "Watch the movie, there is no book with this title."
const noMovieFound string = "Read the book, there is no movie with this title."
const noBookOrMovieFound string = "There is no book or movie with this title. Try searching another title."
const bookFirst string = "Read the book first, then the movie."
const movieFirst string = "Watch the movie first, then the read book."
const bothFirst string = "Watch the movie and read the book, both are evenly rated."

// Error Messages
const noBookRating string = "No rating information exists for the book."
const noMovieRating string = "No rating information exists for the movie."

// Error response
type Error struct {
	Status  string `json:status`
	Code    string `json:code`
	Message string `json:message`
}

type Recommendation struct {
	Title          string `json:title`
	Recommendation string `json:recommendation`
	BookRating     string `json:bookrating`
	MovieRating    string `json:movierating`
	MoviePlot      string `json:movieplot`
	BookPlot       string `json:bookplot`
	Author         string `json: author`
	Director       string `json:director`
	Actors         string `json:actors`
	BookPublished  string `json: bookpublished`
	MovieReleased  string `json: moviereleased`
	MovieAwards    string `json: movieawards`
}

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

	e.GET("/getmovie", getMovie)
	e.GET("/getbook", getBook)
	e.GET("/recommend", getReccomend)

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

	title := c.QueryParam("title")
	movieData := getMovieData(title)
	return c.String(http.StatusOK, movieData)
}

func getMovieData(title string) string {

	// Load Environment Variables
	err := godotenv.Load("local.env")
	if err != nil {
		log.Fatalln(err)
	}

	// Environment Variables
	omdb_api := os.Getenv("OMDB_API")
	omdb_key := os.Getenv("OMDB_KEY")

	request_url := omdb_api + "=" + omdb_key + "&t=" + title
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

	if strings.Contains(string(b), "Movie not found") {
		errorResponse := &Error{
			"Error",
			"400",
			"Unable to find this title",
		}
		errorResp, _ := json.Marshal(errorResponse)
		errorData := string(errorResp)
		return errorData
	}

	director, _ := jsonparser.GetString(b, "Director")
	plot, _ := jsonparser.GetString(b, "Plot")
	actors, _ := jsonparser.GetString(b, "Actors")
	rating, _ := jsonparser.GetString(b, "imdbRating")
	year, _ := jsonparser.GetString(b, "Year")
	awards, _ := jsonparser.GetString(b, "Awards")
	awards = strings.ReplaceAll(awards, "\u0026", "")

	movieResponseData := &MovieResponse{
		title,
		rating,
		actors,
		director,
		plot,
		year,
		awards}

	movieData, _ := json.Marshal(movieResponseData)
	return string(movieData)
}

func getBookData(title string) string {

	request_url := openLibraryAPI + title
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

	numFound, _ := jsonparser.GetInt(b, "numFound")
	if strings.Contains(string(b), "400 Bad request") || numFound == 0 {
		errorResponse := &Error{
			"Error",
			"400",
			"Unable to find this title",
		}
		errorResp, _ := json.Marshal(errorResponse)
		errorData := string(errorResp)
		return errorData
	}

	open_library_id, _ := jsonparser.GetString(b, "docs", "[0]", "key")

	// Web Scrape Open Library for rating and plot information
	rating := ""
	var plot string = ""
	var ratings = []string{}
	web_scraper := colly.NewCollector(
		colly.AllowedDomains(openLibraryURL, "openlibrary.org"),
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
	})

	web_scraper.Visit(openLibraryURL + open_library_id)

	rating = strings.ReplaceAll(ratings[len(ratings)-1], "★", "")
	rating_float, _ := strconv.ParseFloat(rating, 64)
	rating_float = (math.Ceil(rating_float*100) / 100) / 5 * 10
	rating = strconv.FormatFloat(rating_float, 'E', -1, 64)
	rating = strings.ReplaceAll(rating, "E+00", "")

	data_title, _ := jsonparser.GetString(b, "docs", "[0]", "title")
	yearInt, _ := jsonparser.GetInt(b, "docs", "[0]", "first_publish_year")
	yearString := strconv.FormatInt(int64(yearInt), 10)

	author, _ := jsonparser.GetString(b, "docs", "[0]", "author_name")
	plot = strings.ReplaceAll(plot, "\n", "")
	plot = strings.Trim(plot, " ")
	plot = strings.ReplaceAll(plot, "\\", "")

	bookResponse := &BookResponse{
		data_title,
		author,
		rating,
		plot,
		yearString,
	}
	bookResponseData, _ := json.Marshal(bookResponse)
	bookData := string(bookResponseData)
	return bookData
}

// Add %20 to the request url
func cleanTitleQuery(title string) string {

	query_parameter := ""
	title_slice := strings.Split(title, " ")

	for _, s := range title_slice {
		query_parameter += s + "%20"
	}

	return query_parameter
}

func getBook(c echo.Context) error {

	originalTitle := c.QueryParam("title")
	title := cleanTitleQuery(originalTitle)
	bookData := getBookData(title)
	return c.String(http.StatusOK, bookData)
}

func getReccomend(c echo.Context) error {

	originalTitle := c.QueryParam("title")
	title := cleanTitleQuery(originalTitle)

	bookData := getBookData(title)
	movieData := getMovieData(title)

	recommendation := compareBookAndMovieData(movieData, bookData, title)
	status, _ := jsonparser.GetString([]byte(recommendation), "Message")

	if strings.Contains(status, noBookOrMovieFound) {
		return c.String(http.StatusNotFound, recommendation)
	}
	return c.String(http.StatusOK, recommendation)
}

func compareBookAndMovieData(movieData string, bookData string, title string) string {

	fmt.Println("compareBookAndMovieData - start")
	recommend := &Recommendation{}

	fmt.Println("compareBookAndMovieData - book data")
	fmt.Println(bookData + "\n")
	fmt.Println("compareBookAndMovieData - movie data")
	fmt.Println(movieData + "\n")

	bookStatus, _ := jsonparser.GetString([]byte(bookData), "Message")
	movieStatus, _ := jsonparser.GetString([]byte(movieData), "Message")

	fmt.Println("Book Status " + bookStatus + "\n")
	fmt.Println("Movie Status " + movieStatus + "\n")

	// TODO Handle use case when ratings don't exist for book or movie
	// Both movie and book data cannot be retrieved
	if strings.Contains(bookStatus, "Unable to find this title") && strings.Contains(movieStatus, "Unable to find this title") {
		errorResponse := &Error{
			"Error",
			"400",
			noBookOrMovieFound,
		}
		errorResp, _ := json.Marshal(errorResponse)
		errorData := string(errorResp)
		return errorData

	} else if strings.Contains(movieStatus, "Unable to find this title") { // Only book data can be retrieved

		fmt.Println("Found the book but not the movie data")
		recommend.Recommendation = noMovieFound
		titleString, _ := jsonparser.GetString([]byte(bookData), "Title")
		recommend.Title = titleString
		authorString, _ := jsonparser.GetString([]byte(bookData), "Author")
		recommend.Author = authorString
		bookRatingString, _ := jsonparser.GetString([]byte(bookData), "Rating")
		recommend.BookRating = bookRatingString
		bookPlotString, _ := jsonparser.GetString([]byte(bookData), "Plot")
		recommend.BookPlot = bookPlotString
		publishingYear, _ := jsonparser.GetString([]byte(bookData), "Year")
		recommend.BookPublished = publishingYear

	} else if strings.Contains(bookStatus, "Unable to find this title") { // Only movie data can be retrieved

		fmt.Println("Found the movie but not the book data")
		recommend.Recommendation = noBookFound
		titleString, _ := jsonparser.GetString([]byte(movieData), "Title")
		recommend.Title = titleString
		movieRatingString, _ := jsonparser.GetString([]byte(movieData), "Rating")
		recommend.MovieRating = movieRatingString
		actors, _ := jsonparser.GetString([]byte(movieData), "Actors")
		recommend.Actors = actors
		director, _ := jsonparser.GetString([]byte(movieData), "Director")
		recommend.Director = director
		moviePlotString, _ := jsonparser.GetString([]byte(movieData), "Plot")
		recommend.MoviePlot = moviePlotString
		movieYear, _ := jsonparser.GetString([]byte(movieData), "ReleaseYear")
		recommend.MovieReleased = movieYear
		movieAwards, _ := jsonparser.GetString([]byte(movieData), "Awards")
		recommend.MovieAwards = movieAwards

	} else { // Both book and movie data is retrieved

		fmt.Println("Book and movie data is retrieved")
		// TODO Add a success message in the response
		// TODO Add the MovieTitle and Add the BookTitle Separately
		// Book Data
		titleString, _ := jsonparser.GetString([]byte(bookData), "Title")
		recommend.Title = titleString
		authorString, _ := jsonparser.GetString([]byte(bookData), "Author")
		recommend.Author = authorString
		bookRatingString, _ := jsonparser.GetString([]byte(bookData), "Rating")
		recommend.BookRating = bookRatingString
		bookPlotString, _ := jsonparser.GetString([]byte(bookData), "Plot")
		recommend.BookPlot = bookPlotString
		publishingYear, _ := jsonparser.GetString([]byte(bookData), "Year")
		recommend.BookPublished = publishingYear

		// Movie Data
		movieRatingString, _ := jsonparser.GetString([]byte(movieData), "Rating")
		recommend.MovieRating = movieRatingString
		actors, _ := jsonparser.GetString([]byte(movieData), "Actors")
		recommend.Actors = actors
		director, _ := jsonparser.GetString([]byte(movieData), "Director")
		recommend.Director = director
		moviePlotString, _ := jsonparser.GetString([]byte(movieData), "Plot")
		recommend.MoviePlot = moviePlotString
		movieYear, _ := jsonparser.GetString([]byte(movieData), "ReleaseYear")
		recommend.MovieReleased = movieYear
		movieAwards, _ := jsonparser.GetString([]byte(movieData), "Awards")
		recommend.MovieAwards = movieAwards

		// Compare Movie and Book Ratings
		bookRatingFloat, _ := strconv.ParseFloat(bookRatingString, 64)
		movieRatingFloat, _ := strconv.ParseFloat(movieRatingString, 64)

		fmt.Println("bookRatingFloat \n")
		fmt.Println(bookRatingFloat)
		fmt.Println("movieRatingFloat \n")
		fmt.Println(movieRatingFloat)

		if math.Max(bookRatingFloat, movieRatingFloat) == bookRatingFloat {
			recommend.Recommendation = bookFirst
			fmt.Println("bookRatingFloat is greater")
		} else if math.Max(bookRatingFloat, movieRatingFloat) == movieRatingFloat {
			recommend.Recommendation = movieFirst
			fmt.Println("movieRatingFloat is greater")
		} else {
			recommend.Recommendation = bothFirst
			fmt.Println("movies are evenly matched")
		}

	}

	recommendationData, _ := json.Marshal(recommend)
	reccommendResponse := string(recommendationData)
	return reccommendResponse
}

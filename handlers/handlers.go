package handlers

import (
	"net/http"
	"strings"

	models "book-or-movie-api/models"

	"github.com/buger/jsonparser"
	"github.com/labstack/echo"
)

func GetMovie(c echo.Context) error {
	return c.String(http.StatusOK, GetMovieData(c.QueryParam("title")))
}

func GetBook(c echo.Context) error {
	var title string = cleanTitleParameter(c.QueryParam("title"))
	return c.String(http.StatusOK, GetBookData(title))
}

func GetReccomend(c echo.Context) error {

	title := cleanTitleParameter(c.QueryParam("title"))
	recommendation := CompareBookAndMovieData(GetMovieData(title), GetBookData(title), title)
	status, _ := jsonparser.GetString([]byte(recommendation), "Message")

	if strings.Contains(status, models.NoBookOrMovieFound) {
		return c.String(http.StatusNotFound, recommendation)
	}
	return c.String(http.StatusOK, recommendation)
}

// Add %20 to the request query parameter
func cleanTitleParameter(title string) string {

	var queryParameter string = ""
	title_slice := strings.Split(title, " ")
	for _, s := range title_slice {
		queryParameter += s + "%20"
	}
	return queryParameter
}

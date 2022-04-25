// main_test.go
package main

import (
	handlers "book-or-movie-api/handlers"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo"
	"github.com/stretchr/testify/assert"
)

func TestGetBook(t *testing.T) {

	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/getbook", nil)
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/getbook")
	c.SetParamNames("title")
	c.SetParamValues("Klara and the sun")
	res := rec.Result()
	defer res.Body.Close()

	if assert.NoError(t, handlers.GetBook(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
	}
}

func TestGetMovie(t *testing.T) {

	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/getmovie", nil)
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/getmovie")
	c.SetParamNames("title")
	c.SetParamValues("Iron Man")
	res := rec.Result()
	defer res.Body.Close()

	if assert.NoError(t, handlers.GetMovie(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
	}
}

func TestGetReccomendation(t *testing.T) {

	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/getrecommendation", nil)
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/getrecommendation")
	c.SetParamNames("title")
	c.SetParamValues("the hobbit")
	res := rec.Result()
	defer res.Body.Close()

	if assert.NoError(t, handlers.GetReccomendation(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
	}
}

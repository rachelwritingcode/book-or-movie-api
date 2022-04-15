package models

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

var ErrorResponse = Error{
	Status:  "Error",
	Code:    "400",
	Message: "Unable to find this title",
}

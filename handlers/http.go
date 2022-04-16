package handlers

import (
	"io"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

const openLibraryAPI = "http://openlibrary.org/search.json?title="

func Request(title string, movie bool) []byte {

	var request_url string
	if movie {
		err := godotenv.Load("local.env")
		if err != nil {
			log.Fatalln(err)
		}
		omdb_api := os.Getenv("OMDB_API")
		omdb_key := os.Getenv("OMDB_KEY")

		request_url = omdb_api + "=" + omdb_key + "&t=" + title
	} else {
		request_url = openLibraryAPI + title
	}
	req, err := http.NewRequest("GET", request_url, nil)
	if err != nil {
		log.Fatalln(err)
	}
	req.Header.Set("Accept", "application/json")

	client := &http.Client{}
	response, err := client.Do(req)
	if err != nil {
		log.Fatalln(err)
	}
	defer response.Body.Close()
	responseBytes, err := io.ReadAll(response.Body)
	if err != nil {
		log.Fatalln(err)
	}
	return responseBytes
}

package models

const (
	// Recommendation messages
	NoBookFound        string = "Watch the movie, there is no book for this title."
	NoMovieFound       string = "Read the book, there is no movie for this title."
	NoBookOrMovieFound string = "There is no book or movie for this title. Try searching another title."
	BookFirst          string = "Read the book first, then watch the movie."
	MovieFirst         string = "Watch the movie first, then read the book."
	BothFirst          string = "Watch the movie and read the book, both are evenly rated."

	// Error Messages
	NoBookRating        string = "No rating information exists for the book."
	NoMovieRating       string = "No rating information exists for the movie."
	NoBookNoMovieRating string = "No ratings exist for book or movie, watch both."
	NoTitle             string = "Unable to find this title"
)

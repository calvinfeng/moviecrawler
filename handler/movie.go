package handler

import (
	"encoding/json"
	"moviecrawler/model"
	"net/http"

	"github.com/gorilla/mux"
)

// NewMovieListByYearHandler returns a list of movie of a given year and sorted by IMDB rating in
// descending order.
func NewMovieListByYearHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		year := mux.Vars(r)["year"]
		result, err := model.ListMovieByYear(year)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(err.Error()))
			return
		}

		bytes, err := json.Marshal(result)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Write(bytes)
	}
}

package main

import (
	"moviecrawler/crawler"
	"moviecrawler/handler"
	"moviecrawler/model"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

func init() {
	logrus.SetFormatter(&logrus.TextFormatter{
		DisableColors: false,
		FullTimestamp: true,
	})
}

func main() {
	_, err := model.InitPSQLConnection()
	if err != nil {
		logrus.Error(err)
		return
	}

	go crawler.Crawl()

	r := mux.NewRouter()
	r.Handle("/api/movies/year/{year}/", handler.NewMovieListByYearHandler())

	logrus.Info("http server is listening and serving on 8000")
	s := &http.Server{Addr: ":8000", Handler: r}
	if err := s.ListenAndServe(); err != nil {
		panic(err)
	}
}

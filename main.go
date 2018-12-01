package main

import (
	"moviecrawler/model"

	"github.com/sirupsen/logrus"
)

func init() {
	logrus.SetFormatter(&logrus.TextFormatter{
		DisableColors: false,
		FullTimestamp: true,
	})
}

func main() {
	ms := &MovieState{
		VisitedPage: make(chan string),
		IsVisited:   make(chan StateRequest),
		visited:     make(map[string]struct{}),
	}

	go ms.serve()

	_, err := model.InitPSQLConnection()
	if err != nil {
		logrus.Error(err)
		return
	}

	crawl(src, ms)
}

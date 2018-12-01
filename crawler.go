package main

import (
	"moviecrawler/model"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/gocolly/colly"
	"github.com/sirupsen/logrus"
)

const src = "https://ww.gomovies.sc/movies/page/1/"
const maxPage = 200

func crawl(url string, state *MovieState) {
	c := colly.NewCollector()

	// Movie item page
	c.OnHTML("div#mv-info", func(div *colly.HTMLElement) {
		m := &model.Movie{}
		m.Link = div.ChildAttr("a", "href")
		m.Title = div.DOM.Find("div.mvic-desc > h3").Text()
		m.Description = div.DOM.Find("div.mvic-desc > div.desc").Text()
		div.DOM.Find("div.mvic-info > div.mvici-right").Children().Each(func(idx int, sel *goquery.Selection) {
			texts := strings.Split(sel.Text(), ":")
			switch texts[0] {
			case "Release":
				val := strings.Trim(texts[1], " ")
				if year, err := strconv.ParseInt(val, 10, 64); err == nil {
					m.ReleaseYear = uint(year)
				}
			case "IMDb":
				val := strings.Trim(texts[1], " ")
				if rating, err := strconv.ParseFloat(val, 64); err == nil {
					m.IMDbRating = rating
				}
			}
		})

		if err := m.Create(); err != nil {
			logrus.Error(err)
			return
		}

		logrus.Infof("%s (%d) is saved", m.Title, m.ReleaseYear)
	})

	// Movie index page
	c.OnHTML("a.ml-mask", func(a *colly.HTMLElement) {
		// We should only collect HD movies
		quality := a.ChildText("span.mli-quality")
		if quality == "HD" {
			c.Visit(a.Attr("href"))
		}
	})

	// Movie index page pagination
	c.OnHTML("ul.pagination", func(ul *colly.HTMLElement) {
		ul.ForEach("li", func(idx int, li *colly.HTMLElement) {
			pageLink := li.ChildAttr("a", "href")
			if pageLink == "" {
				return
			}

			resp := make(chan bool)
			state.IsVisited <- StateRequest{Page: pageLink, Response: resp}
			if ok := <-resp; ok {
				return
			}

			state.VisitedPage <- pageLink
			c.Visit(pageLink)
		})
	})

	c.Visit(url)
}

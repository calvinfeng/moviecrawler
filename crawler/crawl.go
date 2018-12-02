package crawler

import (
	"context"
	"moviecrawler/model"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/gocolly/colly"
	"github.com/sirupsen/logrus"
)

const url = "https://ww.gomovies.sc/movies/page/1/"
const maxPage = 1

// Crawl will go through gomovies page by page and grab all the movie links.
func Crawl() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	s := NewState(ctx)
	go s.Activate()

	c := colly.NewCollector()

	// Movie index page
	c.OnHTML("a.ml-mask", func(a *colly.HTMLElement) {
		// We should only collect HD movies
		if a.ChildText("span.mli-quality") == "HD" {
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
			s.IsVisited <- request{page: pageLink, resp: resp}
			if ok := <-resp; ok {
				return
			}

			s.VisitedPage <- pageLink
			c.Visit(pageLink)
		})
	})

	// Movie item page
	c.OnHTML("div#mv-info", handleMovieInformationDiv)

	c.Visit(url)
}

func handleMovieInformationDiv(div *colly.HTMLElement) {
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

	if err := m.Create(); err == nil {
		logrus.Infof("%s (%d) is saved", m.Title, m.ReleaseYear)
	}
}

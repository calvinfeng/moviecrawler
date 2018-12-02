package crawler

import (
	"context"

	"github.com/sirupsen/logrus"
)

// request is to ask state for data.
type request struct {
	page string
	resp chan bool
}

// NewState returns a new crawler state object.
func NewState(ctx context.Context) *State {
	return &State{
		VisitedPage: make(chan string),
		IsVisited:   make(chan request),
		visited:     make(map[string]struct{}),
		ctx:         ctx,
	}
}

// State keeps track of what is collected and what is not.
type State struct {
	VisitedPage chan string
	IsVisited   chan request

	// Private members
	visited map[string]struct{}
	ctx     context.Context
}

// Activate turns the state on.
func (s *State) Activate() {
	for {
		select {
		case <-s.ctx.Done():
			logrus.Info("crawler state is terminated, completed crawling")
			return
		case page := <-s.VisitedPage:
			s.visited[page] = struct{}{}
		case req := <-s.IsVisited:
			if len(s.visited) >= maxPage {
				req.resp <- true
			} else {
				_, ok := s.visited[req.page]
				req.resp <- ok
			}
		}
	}
}

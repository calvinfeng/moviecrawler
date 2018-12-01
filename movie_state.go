package main

// StateRequest is to ask state for data.
type StateRequest struct {
	Page     string
	Response chan bool
}

// MovieState keeps track of what is collected and what is not.
type MovieState struct {
	VisitedPage chan string
	IsVisited   chan StateRequest

	// Private members
	visited map[string]struct{}
}

func (ms *MovieState) serve() {
	for {
		select {
		case page := <-ms.VisitedPage:
			ms.visited[page] = struct{}{}
		case request := <-ms.IsVisited:
			if len(ms.visited) >= maxPage {
				request.Response <- true
			} else {
				_, ok := ms.visited[request.Page]
				request.Response <- ok
			}
		}
	}
}

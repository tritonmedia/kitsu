package kitsu

import (
	"strings"
	"time"
)

// StreamMediaRelations returns a stream of all anime relations (async).
// Be very careful to only use this function once as each
// call will start a new goroutine requesting the whole data.
func StreamMediaRelations() chan *MediaRelation {
	channel := make(chan *MediaRelation)
	url := "media-relationships?include=source,destination"
	ticker := time.NewTicker(500 * time.Millisecond)
	rateLimit := ticker.C

	go func() {
		defer close(channel)
		defer ticker.Stop()

		for {
			page, err := GetMediaRelations(url)

			if err != nil {
				panic(err)
			}

			// Feed media relation data from current page to the stream
			for _, relation := range page.Data {
				channel <- relation
			}

			nextURL := page.Links.Next

			// Did we reach the end?
			if nextURL == "" {
				break
			}

			// Cut off API base URL
			nextURL = strings.TrimPrefix(nextURL, APIBaseURL)

			// Continue with the next page
			url = nextURL

			// Wait for rate limiter to allow the next request
			<-rateLimit
		}
	}()

	return channel
}

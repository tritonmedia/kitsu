package kitsu

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/pkg/errors"
)

const (
	pageOffsetKey = "page[offset]"
)

// GetAnime returns the anime for the given ID.
func GetAnime(id string) (*Anime, error) {
	response, requestError := Get("anime/" + id)

	if requestError != nil {
		return nil, requestError
	}

	anime := new(AnimeResponse)
	decodeError := response.Unmarshal(anime)

	return anime.Data, decodeError
}

func getAnimeListPage(id string, next string, e []Episode) ([]Episode, error) {
	if next == "" {
		next = fmt.Sprintf("episodes?filter[mediaId]=%s&page[limit]=20&page[offset]=0", id)
	} else { // make it relatvie
		next = strings.Replace(next, APIBaseURL, "", 1)
	}

	nextURL, err := url.Parse(next)
	if err != nil {
		return e, errors.Wrap(err, "failed to parse next url")
	}
	offset := nextURL.Query().Get(pageOffsetKey)

	res, err := Get(next)
	if err != nil {
		return e, err
	}

	eps := EpisodeListResponse{}
	if err := res.Unmarshal(&eps); err != nil {
		return e, err
	}

	e = append(e, eps.Data...)

	newNextURL, err := url.Parse(eps.Links.Next)
	if err != nil {
		return e, errors.Wrap(err, "failed to parse next url")
	}
	newOffset := newNextURL.Query().Get(pageOffsetKey)

	if newOffset == offset || eps.Links.Next == "" {
		return e, nil
	}

	return getAnimeListPage(id, eps.Links.Next, e)
}

// GetAnimeEpisodes returns a list of episodes for a given ID
func GetAnimeEpisodes(id string) ([]Episode, error) {
	return getAnimeListPage(id, "", make([]Episode, 0))
}

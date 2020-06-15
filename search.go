package twitterscraper

import (
	"context"
	"fmt"
	"net/url"
	"strconv"
)

const ajaxSearchURL = "https://twitter.com/i/search/timeline?q=%s"

// SearchTweets returns channel with tweets for a given search query
func SearchTweets(ctx context.Context, query string, maxTweetsNbr int) <-chan *Result {
	channel := make(chan *Result)
	go func(query string) {
		defer close(channel)
		var maxId string
		tweetsNbr := 0
		for tweetsNbr < maxTweetsNbr {
			select {
			case <-ctx.Done():
				channel <- &Result{Error: ctx.Err()}
				return
			default:
			}

			tweets, err := FetchSearchTweets(query, maxId)
			if err != nil {
				channel <- &Result{Error: err}
				return
			}

			if len(tweets) == 0 {
				break
			}

			for _, tweet := range tweets {
				select {
				case <-ctx.Done():
					channel <- &Result{Error: ctx.Err()}
					return
				default:
				}

				if tweetsNbr < maxTweetsNbr {
					lastId, _ := strconv.ParseInt(tweet.ID, 10, 64)
					maxId = strconv.FormatInt(lastId-1, 10)
					channel <- &Result{Tweet: *tweet}
				}
				tweetsNbr++
			}
		}
	}(query)
	return channel
}

// FetchSearchTweets gets tweets for a given search query, via the Twitter frontend API
func FetchSearchTweets(query, maxId string) ([]*Tweet, error) {
	if maxId != "" {
		query = query + " max_id:" + maxId
	}

	req, err := newRequest(fmt.Sprintf(ajaxSearchURL, url.PathEscape(query)))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Referer", "https://twitter.com/search/timeline")

	q := req.URL.Query()
	q.Add("f", "tweets")
	q.Add("include_available_features", "1")
	q.Add("include_entities", "1")
	q.Add("include_new_items_bar", "true")

	req.URL.RawQuery = q.Encode()

	htm, err := getHTMLFromJSON(req, "items_html")
	if err != nil {
		return nil, err
	}

	tweets, err := readTweetsFromHTML(htm)
	if err != nil {
		return nil, err
	}

	return tweets, nil
}

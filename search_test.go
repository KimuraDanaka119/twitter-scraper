package twitterscraper_test

import (
	"context"
	"strings"
	"testing"

	twitterscraper "github.com/n0madic/twitter-scraper"
)

func TestFetchSearchCursor(t *testing.T) {
	scraper := twitterscraper.New()
	maxTweetsNbr := 150
	tweetsNbr := 0
	nextCursor := ""
	for tweetsNbr < maxTweetsNbr {
		tweets, cursor, err := scraper.FetchSearchTweets("twitter", maxTweetsNbr, nextCursor)
		if err != nil {
			t.Fatal(err)
		}
		if strings.HasPrefix(cursor, "scroll:") {
			continue
		}
		if cursor == "" {
			t.Fatal("Expected search cursor is empty")
		}
		tweetsNbr += len(tweets)
		nextCursor = cursor
	}
}

func TestGetSearchProfiles(t *testing.T) {
	count := 0
	maxProfilesNbr := 150
	dupcheck := make(map[string]bool)
	scraper := twitterscraper.New().SetSearchMode(twitterscraper.SearchUsers)
	for profile := range scraper.SearchProfiles(context.Background(), "Twitter", maxProfilesNbr) {
		if profile.Error != nil {
			t.Error(profile.Error)
		} else {
			count++
			if profile.UserID == "" {
				t.Error("Expected UserID is empty")
			} else {
				if dupcheck[profile.UserID] {
					t.Errorf("Detect duplicated UserID: %s", profile.UserID)
				} else {
					dupcheck[profile.UserID] = true
				}
			}
		}
	}

	if count != maxProfilesNbr {
		t.Errorf("Expected profiles count=%v, got: %v", maxProfilesNbr, count)
	}
}
func TestGetSearchTweets(t *testing.T) {
	count := 0
	maxTweetsNbr := 150
	dupcheck := make(map[string]bool)
	scraper := twitterscraper.New().WithDelay(4)
	for tweet := range scraper.SearchTweets(context.Background(), "twitter", maxTweetsNbr) {
		if tweet.Error != nil {
			t.Error(tweet.Error)
		} else {
			count++
			if tweet.ID == "" {
				t.Error("Expected tweet ID is empty")
			} else {
				if dupcheck[tweet.ID] {
					t.Errorf("Detect duplicated tweet ID: %s", tweet.ID)
				} else {
					dupcheck[tweet.ID] = true
				}
			}
			if tweet.PermanentURL == "" {
				t.Error("Expected tweet PermanentURL is empty")
			}
			if tweet.IsRetweet {
				t.Error("Expected tweet IsRetweet is false")
			}
			if tweet.Text == "" {
				t.Error("Expected tweet Text is empty")
			}
		}
	}

	if count != maxTweetsNbr {
		t.Errorf("Expected tweets count=%v, got: %v", maxTweetsNbr, count)
	}
}

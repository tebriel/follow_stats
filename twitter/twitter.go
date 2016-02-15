package twitter

import (
	"github.com/ChimeraCoder/anaconda"
	"log"
	"net/url"
)

type TwitterCreds struct {
	ConsumerKey       string
	ConsumerSecret    string
	AccessToken       string
	AccessTokenSecret string
}

func Authenticate(creds TwitterCreds) *anaconda.TwitterApi {
	anaconda.SetConsumerKey(creds.ConsumerKey)
	anaconda.SetConsumerSecret(creds.ConsumerSecret)
	return anaconda.NewTwitterApi(creds.AccessToken, creds.AccessTokenSecret)
}

func GetTweets(api anaconda.TwitterApi, username string, is_verbose bool) []anaconda.Tweet {
	v := url.Values{}
	v.Set("screen_name", username)
	v.Set("include_rts", "1")
	v.Set("count", "100")
	tweets, err := api.GetUserTimeline(v)
	if err != nil {
		log.Fatal("Wasn't able to get user's timeline")
	}
	if is_verbose {
		log.Printf("Fetched %d tweets from user %s", len(tweets), username)
	}
	return tweets
}

func CalculateScore(tweets []anaconda.Tweet, is_verbose bool) float64 {
	num_tweets := float64(len(tweets))
	var num_rts, num_ats, num_plain float64
	for _, tweet := range tweets {
		if tweet.RetweetedStatus != nil {
			num_rts = num_rts + 1
		} else if tweet.InReplyToScreenName != "" {
			num_ats = num_ats + 1
		} else {
			num_plain = num_plain + 1
		}
	}

	result := num_plain / num_tweets

	if is_verbose {
		log.Printf("Saw %.0f @s, %.0f RTs, and %.0f plain tweets", num_ats, num_rts, num_plain)
		log.Printf("User's score calcualted as: %f", result)
	}

	return result
}

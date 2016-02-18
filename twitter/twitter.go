package twitter

import (
	"github.com/ChimeraCoder/anaconda"
	"github.com/golang/glog"
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

func ClosePlacePolygon(tweet *anaconda.Tweet) {
	if len(tweet.Place.BoundingBox.Coordinates) == 0 {
		tweet.Place.BoundingBox.Coordinates = [][][]float64{}
		return
	}

	start := tweet.Place.BoundingBox.Coordinates[0][0]
	tweet.Place.BoundingBox.Coordinates[0] = append(tweet.Place.BoundingBox.Coordinates[0], start)
}

func GetTweets(api anaconda.TwitterApi, username string) []anaconda.Tweet {
	v := url.Values{}
	v.Set("screen_name", username)
	v.Set("include_rts", "1")
	v.Set("count", "100")
	tweets, err := api.GetUserTimeline(v)
	if err != nil {
		glog.Fatal("Wasn't able to get user's timeline")
	}
	glog.V(2).Infof("Fetched %d tweets from user %s", len(tweets), username)

	glog.V(2).Info("Closing Polygons on tweets")
	// Close those polygons because ES can't handle a non-closed GeoJSON polygon
	for _, tweet := range tweets {
		ClosePlacePolygon(&tweet)
	}

	return tweets
}

func CalculateScore(tweets []anaconda.Tweet) float64 {
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

	glog.V(2).Infof("Saw %.0f @s, %.0f RTs, and %.0f plain tweets", num_ats, num_rts, num_plain)
	glog.V(2).Infof("User's score calcualted as: %f", result)

	return result
}

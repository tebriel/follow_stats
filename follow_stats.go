package main

import (
	"flag"
	"github.com/ChimeraCoder/anaconda"
	"github.com/gin-gonic/gin"
	"github.com/tebriel/follow_stats/twitter"
	"net/http"
	"os"
)

// Either pass these things in by cli or by env variables
func get_args() twitter.TwitterCreds {
	consumer_key := flag.String("consumer_key", "", "Twitter Consumer Key")
	consumer_secret := flag.String("consumer_secret", "", "Twitter Consumer Key")
	access_token := flag.String("access_token", "", "Twitter Access Token")
	access_token_secret := flag.String("access_token_secret", "", "Twitter Access Token Secret")

	flag.Parse()

	creds := twitter.TwitterCreds{
		ConsumerKey:       os.Getenv("CONSUMER_KEY"),
		ConsumerSecret:    os.Getenv("CONSUMER_SECRET"),
		AccessToken:       os.Getenv("ACCESS_TOKEN"),
		AccessTokenSecret: os.Getenv("ACCESS_TOKEN_SECRET"),
	}

	if *consumer_key != "" {
		creds.ConsumerKey = *consumer_key
	}
	if *consumer_secret != "" {
		creds.ConsumerSecret = *consumer_secret
	}
	if *access_token != "" {
		creds.AccessToken = *access_token
	}
	if *access_token_secret != "" {
		creds.AccessTokenSecret = *access_token_secret
	}

	return creds
}

func eval_user(api anaconda.TwitterApi) gin.HandlerFunc {
	return func(c *gin.Context) {
		name := c.Param("name")

		tweets := twitter.GetTweets(api, name)
		score := twitter.CalculateScore(tweets)

		c.String(http.StatusOK, "Score is: %f", score)
	}
}

func main() {
	creds := get_args()

	// Creates a gin router with default middleware:
	// logger and recovery (crash-free) middleware
	router := gin.Default()

	api := twitter.Authenticate(creds)

	// This handler will match /user/john but will not match neither /user/ or /user
	router.GET("/user/:name", eval_user(*api))

	// By default it serves on :8080 unless a
	// PORT environment variable was defined.
	router.Run()
	// router.Run.Run(":3000") for a hard coded port
}

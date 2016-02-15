package main

import (
	"./twitter"
	"github.com/ChimeraCoder/anaconda"
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
)

func eval_user(api anaconda.TwitterApi) gin.HandlerFunc {
	return func(c *gin.Context) {
		name := c.Param("name")

		tweets = twitter.GetTweets(api, name)
		score = twitter.CalculateScore(tweets)

		c.String(http.StatusOK, "Score is: %f", score)
	}
}

func main() {
	// Creates a gin router with default middleware:
	// logger and recovery (crash-free) middleware
	router := gin.Default()

	env := os.Environ()

	creds = twitter.TwitterCreds{
		ConsumerKey:       env["CONSUMER_KEY"],
		ConsumerSecret:    env["CONSUMER_SECRET"],
		AccessToken:       env["ACCESS_TOKEN"],
		AccessTokenSecret: env["ACCESS_TOKEN_SECRET"],
	}

	api := twitter.Authenticate(creds)

	// This handler will match /user/john but will not match neither /user/ or /user
	router.GET("/user/:name", eval_user(api))

	// By default it serves on :8080 unless a
	// PORT environment variable was defined.
	router.Run()
	// router.Run.Run(":3000") for a hard coded port
}

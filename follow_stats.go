package main

import (
	"flag"
	"github.com/ChimeraCoder/anaconda"
	"github.com/gin-gonic/gin"
	"github.com/golang/glog"
	"github.com/tebriel/follow_stats/cache"
	"github.com/tebriel/follow_stats/twitter"
	"net/http"
	"os"
)

type CliArgs struct {
	ElasticUrl   string
	TwitterCreds twitter.TwitterCreds
}

// Either pass these things in by cli or by env variables
func get_args() CliArgs {
	consumer_key := flag.String("consumer_key", "", "Twitter Consumer Key")
	consumer_secret := flag.String("consumer_secret", "", "Twitter Consumer Key")
	access_token := flag.String("access_token", "", "Twitter Access Token")
	access_token_secret := flag.String("access_token_secret", "", "Twitter Access Token Secret")
	es_url_flag := flag.String("es_url", "", "URL for elasticsearch instance")

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

	result := CliArgs{
		TwitterCreds: creds,
	}

	es_url := os.Getenv("ES_URL")
	if *es_url_flag != "" {
		result.ElasticUrl = *es_url_flag
	} else {
		result.ElasticUrl = es_url
	}

	return result
}

func eval_user(api anaconda.TwitterApi, es_url string) gin.HandlerFunc {
	return func(c *gin.Context) {
		name := c.Param("name")
		glog.V(2).Infof("Fetching stats for %s", name)

		tweets := twitter.GetTweets(api, name)
		score := twitter.CalculateScore(tweets)

		go cache.CacheTweets(tweets, es_url)
		go cache.CacheScore(name, score, es_url)

		c.String(http.StatusOK, "Score is: %.2f", score*100)
	}
}

func main() {
	cli_args := get_args()

	// Creates a gin router with default middleware:
	// logger and recovery (crash-free) middleware
	router := gin.Default()

	api := twitter.Authenticate(cli_args.TwitterCreds)

	// This handler will match /user/john but will not match neither /user/ or /user
	router.GET("/", func(c *gin.Context) { c.String(http.StatusOK, "OK") })
	router.GET("/nginx_status/", func(c *gin.Context) { c.String(http.StatusOK, "OK") })

	router.GET("/user/:name", eval_user(*api, cli_args.ElasticUrl))

	// By default it serves on :8080 unless a
	// PORT environment variable was defined.
	router.Run("0.0.0.0:8080")
	// router.Run.Run(":3000") for a hard coded port
}

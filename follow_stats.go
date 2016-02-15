package main

import (
	"flag"
	"github.com/ChimeraCoder/anaconda"
	"github.com/gin-gonic/gin"
	"github.com/tebriel/follow_stats/twitter"
	"log"
	"net/http"
	"os"
)

type CliArgs struct {
	VerboseOutput bool
	TwitterCreds  twitter.TwitterCreds
}

// Either pass these things in by cli or by env variables
func get_args() CliArgs {
	consumer_key := flag.String("consumer_key", "", "Twitter Consumer Key")
	consumer_secret := flag.String("consumer_secret", "", "Twitter Consumer Key")
	access_token := flag.String("access_token", "", "Twitter Access Token")
	access_token_secret := flag.String("access_token_secret", "", "Twitter Access Token Secret")
	is_verbose := flag.Bool("verbose", false, "Verbose log output")

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
		TwitterCreds:  creds,
		VerboseOutput: *is_verbose,
	}

	return result
}

func eval_user(api anaconda.TwitterApi, is_verbose bool) gin.HandlerFunc {
	return func(c *gin.Context) {
		name := c.Param("name")
		if is_verbose {
			log.Printf("Fetching stats for %s", name)
		}

		tweets := twitter.GetTweets(api, name, is_verbose)
		score := twitter.CalculateScore(tweets, is_verbose)

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
	router.GET("/user/:name", eval_user(*api, cli_args.VerboseOutput))

	// By default it serves on :8080 unless a
	// PORT environment variable was defined.
	router.Run()
	// router.Run.Run(":3000") for a hard coded port
}

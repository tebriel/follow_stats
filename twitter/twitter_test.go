package twitter

import (
	"github.com/ChimeraCoder/anaconda"
	"math"
	"testing"
)

func make_tweets() []anaconda.Tweet {
	var result = []anaconda.Tweet{
		anaconda.Tweet{InReplyToScreenName: "devnall"},
		anaconda.Tweet{RetweetedStatus: &anaconda.Tweet{}},
		anaconda.Tweet{},
	}

	return result
}

func TestCalculateScore(t *testing.T) {
	score := CalculateScore(make_tweets(), false)
	if score != 1.0/3.0 {
		t.Errorf("Expected 0.5 but got: %f", score)
	}
}

func TestCalculateScoreNoTweets(t *testing.T) {
	// Expect this when the user has no tweets, we can use NaN to alert the user that we can't
	// generate a score
	score := CalculateScore([]anaconda.Tweet{}, false)
	if !math.IsNaN(score) {
		t.Errorf("Expected NaN but got: %f", score)
	}
}

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

func TestClosePlacePolygon(t *testing.T) {
	tweet := anaconda.Tweet{InReplyToScreenName: "devnall"}
	coords := [][][]float64{{{0.0, 1.0}, {1.0, 2.0}}}

	tweet.Place.BoundingBox.Coordinates = coords
	ClosePlacePolygon(&tweet)
	if len(tweet.Place.BoundingBox.Coordinates[0]) != 3 {
		t.Errorf("Expected to have 3 sets of coordinates, only have %d", len(coords[0]))
		t.FailNow()
	}
	if !(coords[0][0][0] == coords[0][2][0] && coords[0][0][1] == coords[0][2][1]) {
		t.Errorf("Expected Last set to be the same as first set, instead they were: [%f,%f]", coords[0][2][0], coords[0][2][1])
		t.FailNow()
	}
}

func TestClosePlacePolygonNoCoords(t *testing.T) {
	// Shouldn't crash
	tweet := anaconda.Tweet{InReplyToScreenName: "devnall"}
	ClosePlacePolygon(&tweet)
}

func TestCalculateScore(t *testing.T) {
	score := CalculateScore(make_tweets())
	if score != 1.0/3.0 {
		t.Errorf("Expected 0.5 but got: %f", score)
	}
}

func TestCalculateScoreNoTweets(t *testing.T) {
	// Expect this when the user has no tweets, we can use NaN to alert the user that we can't
	// generate a score
	score := CalculateScore([]anaconda.Tweet{})
	if !math.IsNaN(score) {
		t.Errorf("Expected NaN but got: %f", score)
	}
}

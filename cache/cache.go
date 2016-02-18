package cache

import (
	"github.com/ChimeraCoder/anaconda"
	"github.com/golang/glog"
	"gopkg.in/olivere/elastic.v3"
	"time"
)

type StatsCalc struct {
	User    string  `json:"user"`
	Score   float64 `json:"score"`
	Created int64   `json:"created"`
}

func connect(url string) *elastic.Client {
	// Create a client
	glog.V(2).Infof("Connecting to ES at: %s", url)
	client, err := elastic.NewClient(elastic.SetURL(url))
	if err != nil {
		glog.Fatalf("Couldn't connect to ES: %s", err)
	}

	return client
}

func CacheScore(username string, score float64, es_url string) {
	client := connect(es_url)

	glog.V(2).Infof("Inserting %s:%.2f", username, score)
	doc := StatsCalc{User: username, Score: score, Created: time.Now().UTC().Unix()}
	_, err := client.Index().
		Index("scores").
		Type("statscalc").
		BodyJson(doc).
		Do()
	if err != nil {
		glog.Fatalf("Couldn't create doc: %s", err)
		panic(err)
	}
}

func CacheTweets(tweets []anaconda.Tweet, es_url string) {
	client := connect(es_url)
	glog.V(2).Infof("Building Bulk Insert for %d Tweets", len(tweets))

	bulkRequest := client.Bulk()
	for _, tweet := range tweets {
		action := elastic.NewBulkIndexRequest().Index("tweets").Type("tweet").Id(tweet.IdStr).Doc(tweet)
		bulkRequest.Add(action)
	}

	t := time.Now()
	glog.V(2).Info("Inserting bulk tweets into ES")
	_, err := bulkRequest.Do()
	if err != nil {
		glog.Fatalf("Couldn't bulk load documents: %s", err)
	}
	glog.V(2).Infof("Bulk Insert Complete, taking %ds", time.Since(t)/time.Second)
}

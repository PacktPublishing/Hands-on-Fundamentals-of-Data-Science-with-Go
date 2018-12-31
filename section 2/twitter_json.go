package main

import (
	"encoding/json"
	"fmt"
	"github.com/ChimeraCoder/anaconda"
	"io/ioutil"
	"net/url"
	"os"
)

var (
	consumerKey    = os.Getenv("CONSUMER_KEY_TWITTER")
	consumerSecret = os.Getenv("CONSUMER_SECRET_TWITTER")
	accessToken    = os.Getenv("ACCESS_KEY_TWITTER")
	accessSecret   = os.Getenv("ACCESS_SECRET_TWITTER")
)

type cleanTweet struct {
	Id       string
	Text     string
	Likes    int
	Retweets int
	Language string
	URL      string
}

var CleanTweets []cleanTweet

func getCleanTweet(tweet anaconda.Tweet) cleanTweet {
	var t = cleanTweet{tweet.IdStr, tweet.Text,
		tweet.FavoriteCount, tweet.RetweetCount,
		tweet.Lang,
		"www.twitter.com/i/web/status/" + tweet.IdStr}
	return t
}

func PrettyPrintTweet(tweet anaconda.Tweet) {
	t := getCleanTweet(tweet)
	tweetJSON, _ := json.MarshalIndent(t, "", "\t")
	fmt.Println(string(tweetJSON))
}

func SaveTweetsJSON(TweetsJSON []cleanTweet) error {
	tweetJSON, _ := json.MarshalIndent(TweetsJSON, "", "\t")
	err := ioutil.WriteFile("tweets.json", tweetJSON, 0644)
	if err != nil {
		return err
	}
	return nil
}

func LoadTweetsJSON() ([]cleanTweet, error) {
	fileData, err := ioutil.ReadFile("tweets.json")

	if err != nil {
		return CleanTweets, err
	}
	err = json.Unmarshal(fileData, &CleanTweets)
	if err != nil {
		return CleanTweets, err
	}
	return CleanTweets, nil
}

func main() {
	Tweets, err := LoadTweetsJSON()
	if err == nil {
		fmt.Println("Loading Tweets First..")
		for _, t := range Tweets {
			tweetJSON, _ := json.MarshalIndent(t, "", " \t")
			fmt.Println(string(tweetJSON))
		}
	}
	api := anaconda.NewTwitterApiWithCredentials(
		accessToken, accessSecret, consumerKey, consumerSecret)
	fmt.Println("Started the API...")

	searchResult, _ := api.GetSearch("deep learning",
		url.Values{"result_type": []string{"popular"}})

	fmt.Printf("Retrieved %v tweets\n",
		len(searchResult.Statuses))

	var TweetsForFile []cleanTweet
	for _, tweet := range searchResult.Statuses {
		if !tweet.Retweeted && tweet.FavoriteCount > 100 && tweet.RetweetCount > 50 {
			TweetsForFile = append(TweetsForFile, getCleanTweet(tweet))
		} else {
			fmt.Println("Skipping tweet")
			// PrettyPrintTweet(tweet)
		}
	}
	err = SaveTweetsJSON(TweetsForFile)
	if err != nil {
		fmt.Println("Error in saving tweets")
	} else {
		fmt.Println("Successfully saved popular tweets!")
	}
}

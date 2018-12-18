package main

import (
    "fmt"
    "github.com/ChimeraCoder/anaconda"
    "net/url"
    "os"
)

var (
    consumerKey    = os.Getenv("CONSUMER_KEY_TWITTER")
    consumerSecret = os.Getenv("CONSUMER_SECRET_TWITTER")
    accessToken    = os.Getenv("ACCESS_KEY_TWITTER")
    accessSecret   = os.Getenv("ACCESS_SECRET_TWITTER")
)

func main() {
    api := anaconda.NewTwitterApiWithCredentials(
        accessToken, accessSecret, consumerKey, consumerSecret)
    fmt.Println("Started the API...")

    searchResult, _ := api.GetSearch("deep learning",
        url.Values{"result_type": []string{"popular"}})

    fmt.Printf("Retrieved %v tweets\n",
        len(searchResult.Statuses))

    for _, tweet := range searchResult.Statuses {
        if tweet.FavoriteCount > 5000 && tweet.RetweetCount > 2000 {
            _, err := api.Retweet(tweet.Id, false)
            if err != nil {
                fmt.Println("Error in Retweeting")
                continue
            }
        } else {
            fmt.Printf("Skipping tweet with %v retweets and %v likes\n",
                tweet.RetweetCount, tweet.FavoriteCount)
        }
    }
}

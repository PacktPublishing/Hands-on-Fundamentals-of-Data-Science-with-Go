package main

import (
    "encoding/csv"
    "encoding/json"
    "fmt"
    "github.com/ChimeraCoder/anaconda"
    "net/url"
    "os"
    "strconv"
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
    var t = cleanTweet{tweet.IdStr, tweet.Text, tweet.FavoriteCount, tweet.RetweetCount,
        tweet.Lang, "www.twitter.com/i/web/status/" + tweet.IdStr}
    return t
}

func PrettyPrintTweet(tweet anaconda.Tweet) {
    t := getCleanTweet(tweet)
    tweetJSON, _ := json.MarshalIndent(t, "", "\t")
    fmt.Println(string(tweetJSON))
}

func SaveTweetsCSV(tweets []cleanTweet) error {
    file, err := os.Create("tweets.csv")

    defer file.Close()
    if err != nil {
        panic(err)
    }

    w := csv.NewWriter(file)
    defer w.Flush()

    err = w.Write([]string{"Index", "Id", "Likes",
        "Retweets", "Language", "URL", "Text"})
    if err != nil {
        return err
    }

    for idx, tweet := range tweets {
        stringData := []string{strconv.Itoa(idx), tweet.Id,
            strconv.Itoa(tweet.Likes),
            strconv.Itoa(tweet.Likes),
            tweet.Language, tweet.URL, tweet.Text}
        err = w.Write(stringData)
        if err != nil {
            return err
        }
    }
    return nil
}

func LoadTweetsCSV() ([]cleanTweet, error) {
    csvFile, err := os.Open("tweets.csv")

    defer csvFile.Close()

    if err != nil {
        return CleanTweets, err
    }

    lines, err := csv.NewReader(csvFile).ReadAll()
    if err != nil {
        return CleanTweets, err
    }

    for idx, line := range lines {
        if idx == 0 {
            continue
        }
        likeCounts, _ := strconv.Atoi(line[3])
        retweetCounts, _ := strconv.Atoi(line[4])
        lineData := cleanTweet{
            line[1], line[2], likeCounts, retweetCounts, line[5], line[6],
        }
        CleanTweets = append(CleanTweets, lineData)
    }
    return CleanTweets, err
}

func main() {
    Tweets, err := LoadTweetsCSV()
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
    err = SaveTweetsCSV(TweetsForFile)
    if err != nil {
        fmt.Println("Error in saving tweets")
    } else {
        fmt.Println("Successfully saved popular tweets!")
    }
}

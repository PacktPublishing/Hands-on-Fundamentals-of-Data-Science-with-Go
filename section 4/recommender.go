package main

import (
    "encoding/csv"
    "fmt"
    "math"
    "os"
    "sort"
    "strconv"
    "strings"
)

type User struct {
    UserId      int
    Ratings     map[string]float64 // movie title : Rating
    LikedMovies []Movie            // movies rated 4.5 or more
}

type Movie struct {
    MovieId       int
    Title         string
    Genres        []string // string representation of all genres
    FeatureVector []int    // vector representation of all genres
}

// movie Title : Movie
var Movies = make(map[string]Movie)

// movie Id: title
var MovieIdsToTitle = make(map[int]string)

// Slice of movies for content recommendation
var MoviesSlice []Movie

// user ID : User
var Users = make(map[int]User)

// Genre name : Index (to vectorize movie features)
var GenresIndex = make(map[string]int)

// Starting index for genres map
var defaultIdx = 0

// Have a function to load movies.csv data
// movieId, title, genres (pipe separated) movies.csv
func LoadMoviesCSV() (map[string]Movie, error) {
    csvFile, err := os.Open("ml-latest-small/movies.csv")

    defer csvFile.Close()

    if err != nil {
        return Movies, err
    }

    lines, err := csv.NewReader(csvFile).ReadAll()
    if err != nil {
        return Movies, err
    }

    // Find all genres in data and create an index
    for idx, line := range lines {
        if idx == 0 {
            continue
        }
        genresPipeSeparated := line[2]
        genres := strings.Split(genresPipeSeparated, "|")

        for _, g := range genres {
            if g != "(no genres listed)" {
                if _, ok := GenresIndex[g]; !ok {
                    GenresIndex[g] = defaultIdx
                    defaultIdx++
                }
            }
        }
    }

    // Go through all data and vectorize movie genres
    for idx, line := range lines {
        if idx == 0 {
            continue
        }
        id, _ := strconv.Atoi(line[0])
        title := line[1]
        genresPipeSeparated := line[2]
        genres := strings.Split(genresPipeSeparated, "|")
        var vector = make([]int, len(GenresIndex))
        for _, g := range genres {
            vector[GenresIndex[g]] = 1
        }
        lineData := Movie{id, title, genres, vector}
        Movies[title] = lineData
        MovieIdsToTitle[id] = title
        MoviesSlice = append(MoviesSlice, lineData)
    }
    return Movies, err
}

// Have a function to load ratings.csv data
// userId, movieId, rating, timestamp
func LoadUserRatingsCSV() (map[int]User, error) {
    csvFile, err := os.Open("ml-latest-small/ratings.csv")

    defer csvFile.Close()

    if err != nil {
        return Users, err
    }

    lines, err := csv.NewReader(csvFile).ReadAll()
    if err != nil {
        return Users, err
    }

    var ratingsMap = make(map[string]float64)
    var nextUserIdStr string

    // Get all unique users in ratings.csv
    for idx, line := range lines {
        if idx == 0 {
            continue
        }
        userIdStr := line[0]
        if idx < len(lines)-1 {
            nextUserIdStr = lines[idx+1][0]
        } else {
            nextUserIdStr = ""
        }

        userId, _ := strconv.Atoi(userIdStr)
        movieId, _ := strconv.Atoi(line[1])
        rating, _ := strconv.ParseFloat(line[2], 64) // String to float64

        var currUser User

        // Reset maps, and store this user once this user's
        // info is getting over
        if userIdStr != nextUserIdStr {
            movieTitle := MovieIdsToTitle[movieId]
            ratingsMap[movieTitle] = rating

            currUser.Ratings = ratingsMap
            currUser.UserId = userId

            for k, v := range ratingsMap {
                if v >= 4.5 {
                    currUser.LikedMovies = append(currUser.LikedMovies, Movies[k])
                }
            }
            Users[userId] = currUser
            // Reset the ratingsMap for next user
            ratingsMap = make(map[string]float64)
        } else {
            movieTitle := MovieIdsToTitle[movieId]
            // Keep adding the ratings for the same user
            ratingsMap[movieTitle] = rating
        }
    }
    fmt.Println("Total Unique Users: ", len(Users))
    return Users, err
}

func LoadData() (map[string]Movie, map[int]User) {
    allMovies, err := LoadMoviesCSV()
    if err != nil {
        panic(err)
    }
    fmt.Println("Total movies in dataset:", len(allMovies))

    allUsers, err := LoadUserRatingsCSV()
    if err != nil {
        panic(err)
    }
    return allMovies, allUsers
}

// given two vectors, calculate cosine similarity between them
func calcCosineSim(v1, v2 []int) float64 {
    var numerator float64
    var v1Magnitude float64
    var v2Magnitude float64
    var denominator float64
    var cosineSim float64

    for i := range v1 {
        numerator = numerator + float64(v1[i]*v2[i])
        v1Magnitude = v1Magnitude + float64(v1[i]*v1[i])
        v2Magnitude = v2Magnitude + float64(v2[i]*v2[i])
    }
    denominator = math.Sqrt(v1Magnitude) * math.Sqrt(v2Magnitude)

    if denominator == 0.0 {
        return 0.0
    }
    cosineSim = numerator / denominator
    return cosineSim
}

// given a movie title, recommend similar movie titles
func GetMovieRecommendations(title string) []Movie {
    vector := Movies[title].FeatureVector

    sort.Slice(MoviesSlice, func(i, j int) bool {
        cosineSim1 := calcCosineSim(MoviesSlice[i].FeatureVector,
            vector)
        cosineSim2 := calcCosineSim(MoviesSlice[j].FeatureVector,
            vector)

        return cosineSim1 > cosineSim2 // descending order
    })
    return MoviesSlice
}

func main() {
    allMovies, allUsers := LoadData()

    thisUserId := 2 // user ID

    total := 0
    for _, m := range allUsers[thisUserId].LikedMovies {
        sortedMovies := GetMovieRecommendations(m.Title)
        for i, n := range sortedMovies {
            _, ok := allUsers[thisUserId].Ratings[n.Title]
            if total == 20 {
                break
            }
            if !ok {
                fmt.Printf("Because you liked %s,"+
                    "we recommend %s\n", m.Title, n.Title)
                total = total + 1
            }
            if i == 1 {
                break
            }
        }

    }
    fmt.Println("---------------------")
    fmt.Println("Movie Details...")

    fmt.Println(allMovies["Toy Story (1995)"])

    userMovie := "Toy Story (1995)"
    sortedMovies := GetMovieRecommendations(userMovie)
    fmt.Println("---------------------")

    fmt.Printf("Recommended movies for %s:\n", userMovie)

    total = 0
    for _, m := range sortedMovies {
        if m.Title == userMovie {
            continue
        }
        cSine := calcCosineSim(Movies[userMovie].FeatureVector, Movies[m.Title].FeatureVector)
        total = total + 1
        fmt.Println(total, m.Title, fmt.Sprintf("%.4f", cSine))
        if total == 15 {
            break
        }
    }
}

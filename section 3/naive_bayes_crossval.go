package main

import (
    "fmt"
    . "github.com/cdipaolo/goml/base"
    "github.com/cdipaolo/goml/text"
    "io/ioutil"
    "math/rand"
    "regexp"
    "strings"
)

const (
    TRAIN_DATA_PATH          = "aclImdb/train/"
    TEST_DATA_PATH           = "aclImdb/test/"
    POSITIVE_TRAIN_DATA_PATH = TRAIN_DATA_PATH + "/pos/"
    NEGATIVE_TRAIN_DATA_PATH = TRAIN_DATA_PATH + "/neg/"
    POSITIVE_TEST_DATA_PATH  = TEST_DATA_PATH + "/pos/"
    NEGATIVE_TEST_DATA_PATH  = TEST_DATA_PATH + "/neg/"
)

// Struct used in goml:
// type TextDatapoint struct {
//     X string `json:"x"`
//     Y uint8  `json:"y"`
// }

func preProcess(text string) string {
    // Find all chars that are not alphabets
    reg := regexp.MustCompile("[^a-zA-Z]+")

    // Replace those chars with spaces
    text = reg.ReplaceAllString(text, " ")

    // Lower case
    text = strings.ToLower(text)

    // Tokenize on whitespace, while removing excess whitespace
    tokens := strings.Fields(text)

    // Join the tokens back to string
    return strings.Join(tokens, " ")
}

// Read data from IMDB dataset
func ReadData(dir string) []string {
    fileInfo, err := ioutil.ReadDir(dir)
    if err != nil {
        panic(err)
    }
    var dataStrings []string
    for _, file := range fileInfo {
        bytes, err := ioutil.ReadFile(dir + file.Name())
        if err != nil {
            panic(err)
        }
        dataStrings = append(dataStrings, string(bytes))
    }
    return dataStrings
}

// Prepare data in TextDatapoint struct for text classifier
func GetDataForClassifier() ([]TextDatapoint, []TextDatapoint) {
    var TrainData []TextDatapoint
    var TestData []TextDatapoint

    positiveTrainDataStrings := ReadData(POSITIVE_TRAIN_DATA_PATH)
    negativeTrainDataStrings := ReadData(NEGATIVE_TRAIN_DATA_PATH)
    positiveTestDataStrings := ReadData(POSITIVE_TEST_DATA_PATH)
    negativeTestDataStrings := ReadData(NEGATIVE_TEST_DATA_PATH)

    for _, str := range positiveTrainDataStrings {
        d := TextDatapoint{preProcess(str), 1}
        TrainData = append(TrainData, d)
    }
    for _, str := range negativeTrainDataStrings {
        d := TextDatapoint{preProcess(str), 0}
        TrainData = append(TrainData, d)
    }
    for _, str := range positiveTestDataStrings {
        d := TextDatapoint{preProcess(str), 1}
        TestData = append(TestData, d)
    }
    for _, str := range negativeTestDataStrings {
        d := TextDatapoint{preProcess(str), 0}
        TestData = append(TestData, d)
    }
    // If you want to randomize differently on every run
    //rand.Seed(time.Now().UnixNano())

    // Shuffling is really important here
    rand.Shuffle(len(TrainData), func(i, j int) {
        TrainData[i], TrainData[j] = TrainData[j], TrainData[i]
    })
    return TrainData, TestData
}

func Avg(nums []float64) float64 {
    var sum float64
    for _, num := range nums {
        sum = sum + num
    }
    return sum / float64(len(nums))
}

func main() {
    // Get training and test data set
    TrainData, TestData := GetDataForClassifier()

    // Cross validation data set
    var CrossValData [][]TextDatapoint

    k := len(TrainData) * 1 / 5
    fmt.Println("Items in each fold are:", k)
    for idx := 0; idx < len(TrainData); idx += k {
        CrossValData = append(CrossValData, TrainData[idx:idx+k])
        fmt.Println(len(TrainData[idx:idx+k]), len(CrossValData))
    }

    // USE CROSS VALIDATION TO TRAIN CLASSIFIER 5 TIMES
    var mistakes int
    var TrainAccuracies []float64
    var ValAccuracies []float64
    var model *text.NaiveBayes

    for idx, Piece := range CrossValData {
        // Keep the current piece for validation
        // Collect other 4 pieces for training data
        TrainDataPieces := []TextDatapoint{}
        for i := 0; i < len(CrossValData); i++ {
            if i != idx {
                TrainDataPieces = append(
                    TrainDataPieces, CrossValData[i]...)
            }
        }

        // Use concurrency. Create the channels
        stream := make(chan TextDatapoint, 100)
        errors := make(chan error)

        model = text.NewNaiveBayes(stream, 2, OnlyWordsAndNumbers)
        go model.OnlineLearn(errors)

        // Train the model
        for _, data := range TrainDataPieces {
            stream <- data
        }

        close(stream)

        for {
            err, _ := <-errors
            if err != nil {
                fmt.Printf("Error passed: %v", err)
            } else {
                break
            }
        }

        mistakes = 0
        // calculate mistakes in training data
        for _, t := range TrainDataPieces {
            class := model.Predict(t.X)
            if class != t.Y {
                mistakes += 1
            }
        }

        total := float64(len(TrainDataPieces))
        err := float64(mistakes)
        accuracy := (total - err) / total * 100
        TrainAccuracies = append(TrainAccuracies, accuracy)
        fmt.Printf("Train accuracy %v is %v%%\n", idx, accuracy)

        // calculate mistakes in validation data
        mistakes = 0
        for _, p := range Piece {
            class := model.Predict(p.X)
            if class != p.Y {
                mistakes += 1
            }
        }
        total = float64(len(Piece))
        err = float64(mistakes)
        accuracy = (total - err) / total * 100
        ValAccuracies = append(ValAccuracies, accuracy)
        fmt.Printf("Val accuracy for %v is %v%%\n", idx, accuracy)

    }
    fmt.Printf("Average Train accuracy is %v%%\n",
        Avg(TrainAccuracies))
    fmt.Printf("Average Validation accuracy is %v%%\n",
        Avg(ValAccuracies))

    //now we can predict any new statement
    s := "brilliant. amazing. awesome"
    class := model.Predict(s)
    fmt.Println("Class Predicted", class, "for", s)

    //now we can predict for some test data as well
    //for a positive test
    fmt.Println("Positive Test Example..")
    class = model.Predict(TestData[0].X)
    fmt.Println("Class Predicted", class, "for", TestData[0].X)
    fmt.Println(model.Probability(TestData[0].X))

    // and for a negative test
    fmt.Println("Negative Test Example..")
    class = model.Predict(TestData[12500].X)
    fmt.Println("Class Predicted", class, "for", TestData[12500].X)
    fmt.Println(model.Probability(TestData[12500].X))

}

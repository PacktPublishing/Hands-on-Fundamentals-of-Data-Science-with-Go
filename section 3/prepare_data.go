package main

import (
    "fmt"
    . "github.com/cdipaolo/goml/base"
    "io/ioutil"
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
        d := TextDatapoint{str, 1}
        TrainData = append(TrainData, d)
    }
    for _, str := range negativeTrainDataStrings {
        d := TextDatapoint{str, 0}
        TrainData = append(TrainData, d)
    }
    for _, str := range positiveTestDataStrings {
        d := TextDatapoint{str, 1}
        TestData = append(TrainData, d)
    }
    for _, str := range negativeTestDataStrings {
        d := TextDatapoint{str, 0}
        TestData = append(TrainData, d)
    }
    return TrainData, TestData
}

func main() {
    TrainData, _ := GetDataForClassifier()
    var CrossValData [][]TextDatapoint

    k := len(TrainData) * 1 / 5
    fmt.Println("Items in each fold are:", k)
    for idx := 0; idx < len(TrainData); idx += k {
        CrossValData = append(CrossValData, TrainData[idx:idx+k])
        fmt.Println(len(TrainData[idx:idx+k]), len(CrossValData))
    }
}

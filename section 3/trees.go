package main

import (
    "encoding/csv"
    "fmt"
    "github.com/sjwhitworth/golearn/base"
    "github.com/sjwhitworth/golearn/evaluation"
    "github.com/sjwhitworth/golearn/trees"
    "io/ioutil"
    "math/rand"
    "os"
    "regexp"
    "strconv"
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

// adding a struct to handle the data for trees
// similar to what we did with NaiveBayes
type Data struct {
    X string
    Y string
}

// collection of some of the stopwords that can be used
// to clean up our vocabulary
var STOPWORDS = []string{"i", "me", "my", "myself", "we", "our", "ours",
    "ourselves", "you", "your", "yours",
    "yourself", "yourselves", "he", "him", "his", "himself",
    "she", "her", "hers", "herself", "it", "its", "itself", "they",
    "them", "their", "theirs", "themselves", "what", "which", "who",
    "whom", "this", "that", "these", "those", "am", "is", "are", "was",
    "were", "be", "been", "being", "have", "has", "had", "having", "do",
    "does", "did", "doing", "a", "an", "the", "and", "but", "if", "or",
    "because", "as", "until", "while", "of", "at", "by", "for", "with",
    "about", "against", "between", "into", "through", "during", "before",
    "after", "above", "below", "to", "from", "up", "down", "in", "out",
    "on", "off", "over", "under", "again", "further", "then", "once",
    "here", "there", "when", "where", "why", "how", "all", "any", "both",
    "each", "few", "more", "most", "other", "some", "such", "no", "nor",
    "not", "only", "own", "same", "so", "than", "too", "very", "s", "t",
    "can", "will", "just", "don", "should", "now", "br", "could"}

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

// Prepare data in Data struct for text classifier
func GetDataForClassifier() ([]Data, []Data) {
    var TrainData []Data
    var TestData []Data

    positiveTrainDataStrings := ReadData(POSITIVE_TRAIN_DATA_PATH)
    negativeTrainDataStrings := ReadData(NEGATIVE_TRAIN_DATA_PATH)
    positiveTestDataStrings := ReadData(POSITIVE_TEST_DATA_PATH)
    negativeTestDataStrings := ReadData(NEGATIVE_TEST_DATA_PATH)

    for _, str := range positiveTrainDataStrings {
        d := Data{preProcess(str), "1"}
        TrainData = append(TrainData, d)
    }
    for _, str := range negativeTrainDataStrings {
        d := Data{preProcess(str), "0"}
        TrainData = append(TrainData, d)
    }
    for _, str := range positiveTestDataStrings {
        d := Data{preProcess(str), "1"}
        TestData = append(TestData, d)
    }
    for _, str := range negativeTestDataStrings {
        d := Data{preProcess(str), "0"}
        TestData = append(TestData, d)
    }

    rand.Shuffle(len(TrainData), func(i, j int) {
        TrainData[i], TrainData[j] = TrainData[j], TrainData[i]
    })
    return TrainData, TestData
}

func SaveStringsCSV(AllData []Data, IdxToVocab map[int]string, VocabToIdx map[string]int, fname string) error {
    file, err := os.Create(fname)
    defer file.Close()
    if err != nil {
        panic(err)
    }

    w := csv.NewWriter(file)
    defer w.Flush()

    // prepare the header of the csv
    var firstLine []string
    for i := 0; i < len(IdxToVocab); i++ {
        firstLine = append(firstLine, IdxToVocab[i])
    }
    firstLine = append(firstLine, "label")
    err = w.Write(firstLine)
    if err != nil {
        return err
    }

    // write rest of the training reviews data row by row
    for _, str := range AllData {
        // Last column for label
        var lineData = make([]string, len(VocabToIdx)+1)

        for idx := range lineData {
            lineData[idx] = "0"
        }
        // get word counts in this review
        var reviewWordCount = make(map[string]int)
        // break review into tokens
        tokens := strings.Fields(str.X)
        // loop through tokens and increment count
        for _, token := range tokens {
            _, ok := VocabToIdx[token]
            if ok {
                reviewWordCount[token] += 1
            }
        }
        // add the token counts to each row in CSV
        for w, c := range reviewWordCount {
            val, ok := VocabToIdx[w]
            if ok {
                lineData[val] = strconv.Itoa(c)
            }
        }
        // add the label to end of the column in CSV
        lineData[len(VocabToIdx)] = str.Y

        err = w.Write(lineData)
        if err != nil {
            return err
        }
    }
    fmt.Println("Writing to CSV Complete...")
    return nil
}

func containsString(slice []string, val string) bool {
    for _, item := range slice {
        if item == val {
            return true
        }
    }
    return false
}

func PrepareDataForTrees(TrainData []Data) (map[int]string, map[string]int) {
    var vocabToCount = make(map[string]int) // Vocab to count mapping

    // Loop through entire TrainData and prepare vocab counts
    for _, d := range TrainData {
        tokens := strings.Fields(d.X)
        for _, word := range tokens {
            if !containsString(STOPWORDS, word) {
                _, ok := vocabToCount[word]
                if ok {
                    vocabToCount[word] = vocabToCount[word] + 1
                } else {
                    vocabToCount[word] = 1
                }
            }
        }
    }
    // create the two maps required to make CSV for golearn
    var IdxToVocab = make(map[int]string)
    var VocabToIdx = make(map[string]int)

    i := 0
    for word, count := range vocabToCount {
        if count > 2000 {
            VocabToIdx[word] = i
            IdxToVocab[i] = word
            i++
        }
    }
    return IdxToVocab, VocabToIdx
}

func main() {
    TrainDataStrings, _ := GetDataForClassifier()
    IdxToVocab, VocabToIdx := PrepareDataForTrees(TrainDataStrings)
    fmt.Println("Total length of vocabulary used for trees:", len(VocabToIdx))
    err := SaveStringsCSV(TrainDataStrings,
        IdxToVocab, VocabToIdx, "imdb_train_trees.csv")
    if err != nil {
        panic(err)
    }

    var tree base.Classifier

    rand.Seed(123145)

    // Load in the IMDB review dataset
    imdbData, err := base.ParseCSVToInstances(
        "imdb_train_trees.csv", true) // true is for "hasHeaders"
    if err != nil {
        panic(err)
    }

    fmt.Println("Starting to train the decision tree using information gain...")
    // Create a 80-20 training-validation split
    trainData, valData := base.InstancesTrainTestSplit(
        imdbData, 0.80)
    // TREE USING INFORMATION GAIN RATIO
    tree = trees.NewID3DecisionTreeFromRule(0.6,
        new(trees.InformationGainRatioRuleGenerator))
    // This parameter controls train-prune split
    // Using Information Gain Ratio to build the tree

    // Train the ID3 tree
    err = tree.Fit(trainData)
    if err != nil {
        panic(err)
    }

    // Generate train data predictions
    trainDataPredictions, err := tree.Predict(trainData)
    if err != nil {
        panic(err)
    }
    // Generate validation data predictions
    predictions, err := tree.Predict(valData)
    if err != nil {
        panic(err)
    }

    // Evaluate
    fmt.Println("ID3 Performance (information gain ratio)")
    cfTrain, err := evaluation.GetConfusionMatrix(trainData, trainDataPredictions)

    if err != nil {
        panic(err)
    }
    cf, err := evaluation.GetConfusionMatrix(valData, predictions)

    if err != nil {
        panic(err)
    }
    fmt.Printf("Training accuracy: %v%%\n",
        evaluation.GetAccuracy(cfTrain)*100)
    fmt.Printf("Validation accuracy: %v%%\n",
        evaluation.GetAccuracy(cf)*100)

    fmt.Println(evaluation.GetSummary(cf))

    // TREE USING GINI COEFFICIENT
    fmt.Println("Starting to train the decision tree using gini coefficient...")

    tree = trees.NewID3DecisionTreeFromRule(0.6,
        new(trees.GiniCoefficientRuleGenerator))
    // Using Gini Coeffiecient to build the tree

    // Train the ID3 tree
    err = tree.Fit(trainData)
    if err != nil {
        panic(err)
    }

    // Generate train data predictions
    trainDataPredictions, err = tree.Predict(trainData)
    if err != nil {
        panic(err)
    }
    // Generate validation data predictions
    predictions, err = tree.Predict(valData)
    if err != nil {
        panic(err)
    }

    // Evaluate
    fmt.Println("ID3 Performance (gini index generator)")
    cfTrain, err = evaluation.GetConfusionMatrix(trainData,
        trainDataPredictions)

    if err != nil {
        panic(err)
    }
    cf, err = evaluation.GetConfusionMatrix(valData, predictions)

    if err != nil {
        panic(err)
    }
    fmt.Printf("Training accuracy: %v%%\n",
        evaluation.GetAccuracy(cfTrain)*100)
    fmt.Printf("Validation accuracy: %v%%\n",
        evaluation.GetAccuracy(cf)*100)

    fmt.Println(evaluation.GetSummary(cf))

}

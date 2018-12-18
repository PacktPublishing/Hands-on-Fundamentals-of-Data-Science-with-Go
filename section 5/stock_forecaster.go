package main

import (
    "encoding/csv"
    "fmt"
    "github.com/sajari/regression"
    "os"
    "strconv"
)

func trainModel(examplesToTry int) (*regression.Regression, error) {
    r := new(regression.Regression)
    r.SetObserved("Closing Price")
    r.SetVar(0, "MACD Signal")
    r.SetVar(1, "Middle Band")
    r.SetVar(2, "Z Score")

    csvFile, err := os.Open("AAPL_training_data.csv")

    defer csvFile.Close()

    if err != nil {
        return r, err
    }

    lines, err := csv.NewReader(csvFile).ReadAll()
    if err != nil {
        return r, err
    }

    for idx, line := range lines[:examplesToTry] {
        // removing initial data as it is used to set up MACD, signals and Z score
        if idx <= 35 {
            continue
        }
        closingPrice, _ := strconv.ParseFloat(lines[idx+1][1], 64) // Predict next day price value
        macdSignal, _ := strconv.ParseFloat(line[5], 64)
        middleBand, _ := strconv.ParseFloat(line[6], 64)
        zScore, _ := strconv.ParseFloat(line[10], 64)

        r.Train(regression.DataPoint(closingPrice,
            []float64{macdSignal, middleBand, zScore}))
    }

    err = r.Run()

    fmt.Printf("Regression formula:\n%v\n", r.Formula)
    fmt.Printf("Regression:\n%s\n", r)

    return r, err
}

func main() {
    r, err := trainModel(500)
    if err != nil {
        panic(err)
    }
    fmt.Println("R2 of this model is", r.R2)
}

package main

import (
    "encoding/csv"
    "fmt"
    "github.com/sajari/regression"
    "math"
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
    fmt.Printf("Training model with %v data points\n",
        examplesToTry)
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

func testModel(r *regression.Regression, startIdx int, endIdx int) float64 {
    csvFile, err := os.Open("AAPL_training_data.csv")

    defer csvFile.Close()

    if err != nil {
        panic(err)
    }

    lines, err := csv.NewReader(csvFile).ReadAll()
    if err != nil {
        panic(err)
    }
    n := float64(endIdx - startIdx) // total data points in float
    var squaredErrors float64
    for idx, line := range lines[startIdx:endIdx] {
        macdSignal, _ := strconv.ParseFloat(line[5], 64)
        middleBand, _ := strconv.ParseFloat(line[6], 64)
        zScore, _ := strconv.ParseFloat(line[10], 64)
        val, _ := r.Predict([]float64{macdSignal, middleBand,
            zScore})
        actual, _ := strconv.ParseFloat(lines[startIdx+idx+1][1],
            64) // Actual to predicted is 1 index ahead
        squaredErrors = squaredErrors + (val-actual)*(val-actual)
    }
    return math.Sqrt(squaredErrors / n)
}

func main() {
    var rmse []float64
    var r2 []float64

    // Roll forward cross-validation for stock market dataset
    r, err := trainModel(1000)
    r2 = append(r2, r.R2)

    if err != nil {
        panic(err)
    }

    rmse = append(rmse, testModel(r, 1000, 1100))

    r, err = trainModel(1100)
    r2 = append(r2, r.R2)

    if err != nil {
        panic(err)
    }

    rmse = append(rmse, testModel(r, 1100, 1200))

    r, err = trainModel(1200)
    r2 = append(r2, r.R2)

    if err != nil {
        panic(err)
    }

    rmse = append(rmse, testModel(r, 1200, 1300))

    fmt.Println("All RMSEs are", rmse)
    fmt.Println("All R2s are", r2)

    // 7/2/13   59.784286   60.16810321 61.47931346 -1.31121025 -0.790538176    60.2745712  65.06090627 55.48823613 2.393167534 -0.204868733
    val, _ := r.Predict([]float64{-0.790538176, 60.2745712, -0.204868733}) // Predicted 59.77, Actual 60.114285
    fmt.Println("Predicted value for test data point is", val)

    // 12/29/16  116.730003  114.5145495 113.0951307 1.419418752 0.841754384 114.2615002 120.0888164 108.434184  2.913658109 0.847217727
    val, _ = r.Predict([]float64{0.841754384, 114.2615002, 0.847217727}) // Predicted 116.13, Actual 115.82
    fmt.Println("Predicted value for test data point is", val)

}

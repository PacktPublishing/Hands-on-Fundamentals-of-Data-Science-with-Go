package main

import (
    "fmt"
    "github.com/sajari/regression"
)

func main() {
    r := new(regression.Regression)
    r.SetObserved("Temperature in Celcius")
    r.SetVar(0, "Humidity")
    r.SetVar(1, "Wind Speed")

    var observed = []float64{9.47, 9.36, 9.38, 8.29, 8.76,
        9.22, 7.73, 8.77, 10.82, 13.77, 16.02,
        17.14, 17.80, 17.33, 18.88, 18.91, 15.39, 15.55,
        14.26, 13.14, 11.55, 11.18, 10.12,
        10.20, 10.42, 9.91, 11.18, 7.16, 6.11, 6.79}

    var humidity = []float64{0.89, 0.86, 0.89, 0.83, 0.83,
        0.85, 0.95, 0.89, 0.82, 0.72, 0.67, 0.54, 0.55, 0.51,
        0.47, 0.46, 0.6, 0.63, 0.69, 0.7, 0.77, 0.76, 0.79,
        0.77, 0.62, 0.66, 0.8, 0.79, 0.82, 0.83}

    var windSpeed = []float64{14.1197, 14.2646, 3.9284,
        14.1036, 11.0446, 13.9587, 12.3648, 14.1519, 11.3183,
        12.5258, 17.5651, 19.7869, 21.9443, 20.6885,
        15.3755, 10.4006, 14.4095, 11.1573, 8.5169, 7.6314,
        7.3899, 4.9266, 6.6493, 3.9284, 16.9855, 17.2109,
        10.8192, 11.0768, 6.6493, 13.0088}

    // Add data points to the regression model
    for i := 0; i < len(observed); i++ {
        r.Train(regression.DataPoint(observed[i],
            []float64{humidity[i], windSpeed[i]}))
    }
    r.Run() // Run the regression model

    fmt.Printf("Regression formula:\n%v\n", r.Formula)
    fmt.Printf("Regression:\n%s\n", r)

    // Actual 15.09
    val, _ := r.Predict([]float64{0.61, 17.549})
    fmt.Println("Predicted value for test data point is", val)
    fmt.Println("R2 value for this regression is", r.R2)
}

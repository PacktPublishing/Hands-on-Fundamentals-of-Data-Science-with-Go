package main

import (
    "fmt"
    "math"
)

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

func main() {
    a := []int{1, 0, 1, 1, 0}
    b := []int{1, 0, 1, 1, 0}

    fmt.Println("Cosine Sim of a and b", calcCosineSim(a, b))

    c := []int{0, 1, 0, 0, 0}
    d := []int{1, 0, 1, 1, 0}
    fmt.Println("Cosine Sim of c and d", calcCosineSim(c, d))

    e := []int{0, 1, 1, 1, 0}
    f := []int{1, 0, 1, 1, 0}
    fmt.Println("Cosine Sim of e and f", calcCosineSim(e, f))

}

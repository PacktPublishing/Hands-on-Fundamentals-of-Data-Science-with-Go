package main

import (
    "fmt"
    "github.com/DannyBen/quandl"
    "gonum.org/v1/plot"
    "gonum.org/v1/plot/plotter"
    "gonum.org/v1/plot/plotutil"
    "gonum.org/v1/plot/vg"
    "math/rand"
    "os"
)

var (
    quandlAPIKey = os.Getenv("QUANDL_API_KEY")
)

func main() {
    rand.Seed(int64(0))

    p, err := plot.New()
    if err != nil {
        panic(err)
    }

    p.Title.Text = "Amazon vs Apple Stock Price 2012-2017"
    p.X.Label.Text = "Time"
    p.Y.Label.Text = "Price"

    quandl.APIKey = quandlAPIKey
    v := quandl.Options{}
    v.Set("start_date", "2012-01-01")
    v.Set("end_date", "2017-10-01")
    v.Set("column_index", "4")
    data, _ := quandl.GetSymbol("WIKI/AMZN", v)
    pts := make(plotter.XYs, len(data.Data))

    for i, item := range data.Data {
        fmt.Println(item)
        fmt.Println(i, item[0], item[1])
        pts[len(data.Data)-1-i].X = float64(len(data.Data) - 1 - i) // type conversion
        pts[len(data.Data)-1-i].Y = item[1].(float64)               // type assertion
    }

    data2, _ := quandl.GetSymbol("WIKI/AAPL", v)
    pts2 := make(plotter.XYs, len(data2.Data))

    for i, item := range data2.Data {
        fmt.Println(item)
        fmt.Println(i, item[0], item[1])
        pts2[len(data2.Data)-1-i].X = float64(len(data2.Data) - 1 - i) // type conversion
        pts2[len(data2.Data)-1-i].Y = item[1].(float64)                // type assertion
    }

    err = plotutil.AddLinePoints(p,
        "AMZN", pts, "AAPL", pts2)

    if err != nil {
        panic(err)
    }

    err = p.Save(6.5*vg.Inch, 4*vg.Inch, "amznvsaapl.jpg")
    if err != nil {
        panic(err)
    }

}

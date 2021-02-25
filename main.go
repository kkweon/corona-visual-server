package main

import (
	"math/rand"
	"net/http"
	"os"

	"github.com/go-echarts/go-echarts/v2/charts"
	"github.com/go-echarts/go-echarts/v2/opts"
)

// func generateBarItems(data []int) []opts.BarData {
// 	items := make([]opts.BarData, len(data))
// 	for _, d := range data {
// 		items = append(items, opts.BarData{Value: d})
// 	}
// 	return items
// }

func generateBarItems() []opts.BarData {
	items := make([]opts.BarData, 0)
	for i := 0; i < 7; i++ {
		items = append(items, opts.BarData{Value: rand.Intn(300)})
	}
	return items
}

func weeklyHandler(w http.ResponseWriter, r *http.Request) {
	// namesX := []string{"Mon", "Tue", "Wed", "Thr", "Fri", "Sat", "Sun"}
	// fmt.Println("weeklyHandler")
	// bar := charts.NewBar()

	// // set some global options like Title/Legend/ToolTip or anything else
	// bar.SetGlobalOptions(charts.WithTitleOpts(opts.Title{
	// 	Title:    "My first bar chart generated by go-echarts",
	// 	Subtitle: "It's extremely easy to use, right?",
	// }))

	// bar.SetXAxis([]string{"Mon", "Tue", "Wed", "Thr", "Fri", "Sat", "Sun"}).
	// 	AddSeries("3 weeks ago", generateBarItems([]int{10, 20, 30, 40, 50, 60, 70})).
	// 	AddSeries("2 weeks ago", generateBarItems([]int{15, 20, 32, 10, 55, 60, 79})).
	// 	AddSeries("1 weeks ago", generateBarItems([]int{25, 10, 32, 20, 35, 69, 99}))
	// f, err := os.Create("bar.html")
	// if err != nil {
	// 	log.Println(err)
	// 	return
	// }
	// bar.Render(f)

	// create a new bar instance
	bar := charts.NewBar()
	// set some global options like Title/Legend/ToolTip or anything else
	bar.SetGlobalOptions(
		charts.WithTitleOpts(opts.Title{
			Title:    "My first bar chart generated by go-echarts",
			Subtitle: "It's extremely easy to use, right?",
		}),
		charts.WithLegendOpts(opts.Legend{Show: true}),
	)

	// Put data into instance
	bar.SetXAxis([]string{"Mon", "Tue", "Wed", "Thu", "Fri", "Sat", "Sun"}).
		AddSeries("Category A", generateBarItems()).
		AddSeries("Category B", generateBarItems()).
		AddSeries("Category C", generateBarItems())
	// Where the magic happens
	f, _ := os.Create("bar.html")
	bar.Render(f)

	htmlFile := "./bar.html"
	http.ServeFile(w, r, htmlFile)
}

const port = ":8081"

func main() {
	http.HandleFunc("/", weeklyHandler)
	http.ListenAndServe(port, nil)
}

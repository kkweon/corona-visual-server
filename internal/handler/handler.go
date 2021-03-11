package handler

import (
	"corona-visual-server/internal/config"
	"corona-visual-server/internal/fetcher"
	"corona-visual-server/internal/model"
	"github.com/sirupsen/logrus"
	"net/http"
	"strings"
	"time"

	"github.com/go-echarts/go-echarts/v2/charts"
	"github.com/go-echarts/go-echarts/v2/opts"
)

var weekdays = []string{"Mon", "Tue", "Wed", "Thu", "Fri", "Sat", "Sun"}

// Handler is the main handler.
type Handler struct {
	fetcher *fetcher.Fetcher
	config  *config.Config
}

// New returns the Handler.
func New(config *config.Config, fetcher *fetcher.Fetcher) Handler {
	return Handler{
		fetcher: fetcher,
		config:  config,
	}
}

// GetWeeklyHandler handles weekly request.
func (h *Handler) GetWeeklyHandler(w http.ResponseWriter, _ *http.Request) {
	logrus.Info("weeklyHandler")
	data, err := h.fetcher.GetCoronaData(time.Now())
	if err != nil {
		logrus.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// create a new bar instance
	bar := charts.NewBar()
	// set some global options like Title/Legend/ToolTip or anything else
	bar.SetGlobalOptions(
		charts.WithTitleOpts(opts.Title{
			Title:    "Covid confirmed person data comparison",
			Subtitle: "3 Weeks comparison of each weekday",
			Left:     "5%",
		}),
		charts.WithLegendOpts(opts.Legend{
			Show: true,
			Left: "48%",
			Top:  "7%",
		}),
	)

	// Put data into the bar instance
	bar.SetXAxis(h.getWeeklyAxis(data[0])).
		AddSeries("3rd weeks ago", generateWeeklyBarItems(data[:7])).
		AddSeries("2nd weeks ago", generateWeeklyBarItems(data[7:14])).
		AddSeries("1st weeks ago", generateWeeklyBarItems(data[14:])).
		SetSeriesOptions(charts.WithLabelOpts(opts.Label{
			Show:     true,
			Position: "top",
		}),
		)
	err = bar.Render(w)
	if err != nil {
		logrus.Error(err)
	}
}

// getWeeklyAxis finds the starting weekday of the xAxis
func (h *Handler) getWeeklyAxis(data model.CoronaDailyData) []string {
	t, err := time.Parse(h.config.DateFormat, data.Date)
	if err != nil {
		logrus.Error(err)
		return weekdays
	}
	wDay := t.Weekday().String()
	logrus.Infof("weekday start: %v", wDay)

	var idx int
	for i, d := range weekdays {
		if strings.Contains(wDay, d) {
			idx = i
		}
	}
	result := append(weekdays[idx:], weekdays[:idx]...)
	return result
}

func generateWeeklyBarItems(data []model.CoronaDailyData) []opts.BarData {
	items := make([]opts.BarData, 0)
	for i := 0; i < 7; i++ {
		items = append(items, opts.BarData{Value: data[i].AddCount})
	}
	return items
}

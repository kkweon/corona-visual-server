package fetcher

import (
	"corona-visual-server/internal/config"
	"corona-visual-server/internal/model"
	"encoding/xml"
	"fmt"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
	"time"
)

type httpClient interface {
	Do(req *http.Request) (*http.Response, error)
}

// Fetcher fetches Corona data.
type Fetcher struct {
	config *config.Config
	client httpClient
}

// New returns Fetcher
func New(config *config.Config, client httpClient) Fetcher {
	return Fetcher{
		config: config,
		client: client,
	}
}

func (f *Fetcher) get3WeeksRange(cTime time.Time) (string, string) {
	endDate := cTime.Format(f.config.DateFormat)
	startDate := cTime.AddDate(0, 0, -23).Format(f.config.DateFormat) // I need 21 days, but I have 23 days just in case
	logrus.Infof("startDate %v, endDate %v", startDate, endDate)
	return startDate, endDate
}

// GetCoronaData returns CoronaData
// TODO: This function should return CoronaStruct instead of []byte
func (f *Fetcher) GetCoronaData(cTime time.Time) ([]model.CoronaDailyData, error) {
	// make request with query https://stackoverflow.com/a/30657518/6513756
	logrus.Info("GetCoronaData")

	b, err := f.requestOpenAPI(cTime)
	if err != nil {
		return nil, err
	}

	logrus.Info("GetCoronaData success")

	data, err := f.getCoronaDailyDataFromResponse(b)
	if err != nil {
		return nil, err
	}

	return data, nil
}

func (f *Fetcher) getCoronaDailyDataFromResponse(b []byte) ([]model.CoronaDailyData, error) {
	var resp model.Response
	if err := xml.Unmarshal(b, &resp); err != nil {
		return nil, err
	}

	var data []model.CoronaDailyData
	for i, item := range resp.Body.Items.Item {
		if i == len(resp.Body.Items.Item)-1 {
			continue
		}

		data = append(data, model.CoronaDailyData{
			Date:     item.StateDt.AddDate(0, 0, -1),
			AddCount: getAddCount(item, resp.Body.Items.Item[i+1]),
		})
	}

	if data == nil || len(data) < 21 {
		return nil, fmt.Errorf(
			"err = %v\nlen([]CoronaDailyData) has less than 21 days of data",
			data)
	}

	// data is in reverse order
	// so, here get 21 latest data (D0, D-1, D-2, D-3, ..., D-20).
	data = data[:21]

	// reverse so it becomes (D-20, D-19, ..., D0)
	for i, j := 0, len(data)-1; i < j; i, j = i+1, j-1 {
		data[i], data[j] = data[j], data[i]
	}

	return data, nil
}

func (f *Fetcher) requestOpenAPI(time time.Time) ([]byte, error) {
	req, err := http.NewRequest("GET", f.config.OpenAPIURL, nil)
	if err != nil {
		return nil, err
	}

	startDate, endDate := f.get3WeeksRange(time)
	q := req.URL.Query()
	q.Add("serviceKey", f.config.ServiceKey)
	q.Add("pageNo", "1")
	q.Add("numOfRows", "25") // I will have max 23 days result
	q.Add("startCreateDt", startDate)
	q.Add("endCreateDt", endDate)

	req.URL.RawQuery = q.Encode() // this make added query to attached AND URL encoding

	resp, err := f.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// response
	if resp.StatusCode != http.StatusOK {
		return nil, err
	}

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return b, nil
}

func getAddCount(today model.Item, yday model.Item) int64 {
	return today.CareCnt + today.ClearCnt + today.DeathCnt - yday.CareCnt - yday.ClearCnt - yday.DeathCnt
}

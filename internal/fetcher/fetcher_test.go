package fetcher

import (
	"bytes"
	"corona-visual-server/internal/config"
	"corona-visual-server/internal/model"
	"io/ioutil"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

type httpMockClient struct {
	DoFunc func(*http.Request) (*http.Response, error)
}

func (h *httpMockClient) Do(req *http.Request) (*http.Response, error) {
	return h.DoFunc(req)
}

func MustParseYYYYMMDD(yyyymmdd string) time.Time {
	t, _ := time.ParseInLocation("20060102", yyyymmdd, config.SeoulTZ)
	return t
}

func TestFetcher_GetCoronaData(t *testing.T) {
	sampleXML, err := ioutil.ReadFile("sample_data/sample_response.xml")
	now := time.Now()
	if err != nil {
		assert.Error(t, err, "Unable to open sample_response.xml. Maybe in a wrong directory?")
	}

	type fields struct {
		config *config.Config
		client httpClient
	}

	tests := []struct {
		name    string
		fields  fields
		want    []model.CoronaDailyData
		wantErr bool
	}{
		{
			name: "[positive] it should return CoronaData",
			fields: fields{
				config: &config.Config{
					DateFormat: "20060102",
				},
				client: &httpMockClient{
					DoFunc: func(request *http.Request) (*http.Response, error) {
						return &http.Response{
							StatusCode: 200,
							Body:       ioutil.NopCloser(bytes.NewReader(sampleXML)),
						}, nil
					},
				},
			},
			want: []model.CoronaDailyData{
				{MustParseYYYYMMDD("20200218"), -1},
				{MustParseYYYYMMDD("20200219"), -1},
				{MustParseYYYYMMDD("20200220"), -1},
				{MustParseYYYYMMDD("20200221"), -1},
				{MustParseYYYYMMDD("20200222"), -1},
				{MustParseYYYYMMDD("20200223"), -1},
				{MustParseYYYYMMDD("20200224"), -1},
				{MustParseYYYYMMDD("20200225"), -1},
				{MustParseYYYYMMDD("20200226"), -1},
				{MustParseYYYYMMDD("20200227"), -1},
				{MustParseYYYYMMDD("20200228"), -1},
				{MustParseYYYYMMDD("20200229"), -1},
				{MustParseYYYYMMDD("20200301"), -1},
				{MustParseYYYYMMDD("20200302"), -1},
				{MustParseYYYYMMDD("20200303"), 602},
				{MustParseYYYYMMDD("20200304"), 555},
				{MustParseYYYYMMDD("20200305"), 415},
				{MustParseYYYYMMDD("20200306"), 503},
				{MustParseYYYYMMDD("20200307"), 489},
				{MustParseYYYYMMDD("20200308"), 386},
				{MustParseYYYYMMDD("20200309"), 295},
			},
			wantErr: false,
		},
		{
			name: "[negative] it should return error when there are less than 21 data",
			fields: fields{
				config: &config.Config{
					DateFormat: "20060102",
				},
				client: &httpMockClient{
					DoFunc: func(request *http.Request) (*http.Response, error) {
						return &http.Response{
							StatusCode: 200,
							Body:       ioutil.NopCloser(bytes.NewReader([]byte("<response></response>"))),
						}, nil
					},
				},
			},
			want:    nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := &Fetcher{
				config: tt.fields.config,
				client: tt.fields.client,
			}
			got, err := f.GetCoronaData(now)

			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

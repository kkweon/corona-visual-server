package model

import (
	"corona-visual-server/internal/config"
	"encoding/xml"
	"fmt"
	"github.com/sirupsen/logrus"
	"time"
)

// Response represents the result of research.swtch.com
type Response struct {
	XMLName xml.Name `xml:"response"`
	Text    string   `xml:",chardata"`
	Header  struct {
		Text       string `xml:",chardata"`
		ResultCode string `xml:"resultCode"`
		ResultMsg  string `xml:"resultMsg"`
	} `xml:"header"`
	Body struct {
		Text  string `xml:",chardata"`
		Items struct {
			Text string `xml:",chardata"`
			Item []Item  `xml:"item"`
		} `xml:"items"`
		NumOfRows  string `xml:"numOfRows"`
		PageNo     string `xml:"pageNo"`
		TotalCount string `xml:"totalCount"`
	} `xml:"body"`
}

// Item represents an individual item of Response.
type Item struct {
	Text           string        `xml:",chardata"`
	AccDefRate     float64       `xml:"accDefRate"`
	AccExamCnt     int64         `xml:"accExamCnt"`
	AccExamCompCnt int64         `xml:"accExamCompCnt"`
	CareCnt        int64         `xml:"careCnt"`
	ClearCnt       int64         `xml:"clearCnt"`
	CreateDt       timeInSeoulTZ `xml:"createDt"`
	DeathCnt       int64         `xml:"deathCnt"`
	DecideCnt      int64         `xml:"decideCnt"`
	ExamCnt        int64         `xml:"examCnt"`
	ResutlNegCnt   int64         `xml:"resutlNegCnt"`
	Seq            int64         `xml:"seq"`
	StateDt        timeInSeoulTZ `xml:"stateDt"`
	StateTime      string        `xml:"stateTime"`
	UpdateDt       timeInSeoulTZ `xml:"updateDt"`
}

type timeInSeoulTZ struct {
	time.Time
}

func (c *timeInSeoulTZ) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	var text string
	err := d.DecodeElement(&text, &start)
	if err != nil {
		return err
	}

	t, err := time.ParseInLocation(config.LongTimeFormat, text, config.SeoulTZ)
	if err == nil {
		*c = timeInSeoulTZ{t}
		return nil
	}

	t, err2 := time.ParseInLocation(config.ShortTimeFormat, text, config.SeoulTZ)
	if err2 == nil {
		*c = timeInSeoulTZ{t}
		return nil
	}

	logrus.WithFields(logrus.Fields{
		"err": err,
		"err2": err2,
	}).Errorf("failed to UnmarshalXML(%s)", text)

	return fmt.Errorf("failed to UnmarshalXML(%s)", text)
}

// CoronaDailyData is a single data point.
type CoronaDailyData struct {
	Date     time.Time
	AddCount int64
}

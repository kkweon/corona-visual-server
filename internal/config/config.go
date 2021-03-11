package config

import (
	"github.com/sirupsen/logrus"
	"time"
)

// Config represents Application global config.
type Config struct {
	OpenAPIURL string
	DateFormat string
	ServiceKey string
	Port       string
}

// SeoulTZ is the timezone for Seoul.
var SeoulTZ *time.Location

func init() {
	var err error
	 SeoulTZ, err = time.LoadLocation("Asia/Seoul")

	 if err != nil {
	 	logrus.WithError(err).Panic("unable to load timezone Asia/Seoul")
	 }
}

const (
	LongTimeFormat = "2006-01-02 15:04:05.999999999"
	ShortTimeFormat = "20060102"
)

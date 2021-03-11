package main

import (
	"corona-visual-server/internal/config"
	"corona-visual-server/internal/fetcher"
	"corona-visual-server/internal/handler"
	"github.com/sirupsen/logrus"
	"net/http"
	"os"
)

func main() {
	serviceKey := os.Getenv("SERVICE_KEY")
	if serviceKey == "" {
		logrus.Fatal("$SERVICE_KEY is not set")
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8081"
		logrus.Info("$PORT is not set, so port set to ", port)
	}

	cfg := config.Config{
		OpenAPIURL: openAPIURL,
		DateFormat: dateFormat,
		ServiceKey: serviceKey,
		Port:       port,
	}

	f := fetcher.New(&cfg, netClient)
	h := handler.New(&cfg, &f)

	http.HandleFunc("/", h.GetWeeklyHandler)
	logrus.Fatal(http.ListenAndServe(":"+cfg.Port, nil))
}

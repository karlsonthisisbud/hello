package main

import (
	"github.com/i25959341/sku-aggregator/internal/bootstrap"
	"github.com/sirupsen/logrus"
)

// main Initialises the App
func main() {
	err := bootstrap.New()
	if err != nil {
		logrus.Panic(err)
	}
}

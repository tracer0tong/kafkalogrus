package main

import (
	"time"
	"github.com/sirupsen/logrus"
	"kafkalogrus"
)

func main() {
	var err error
	var hook *kafkalogrus.KafkaLogrusHook

	// Create a new hook
	if hook, err = kafkalogrus.NewKafkaLogrusHook(
		"kh",
		logrus.AllLevels,
		logrus.JSONFormatter{},
		[]string{"127.0.0.1:9092"},
		"test",
		true); err != nil {
		panic(err)
	}

	// Create a new logrus.Logger
	logger := logrus.New()

	// Add hook to logger
	logger.Hooks.Add(hook)

	// Add topics
	l := logger.WithField("topic", "nondefaulttopic")

	// Send message to logger
	l.Debug("This must not be logged")
	l.Info("This is an Info msg")
	l.Warn("This is a Warn msg")
	l.Error("This is an Error msg")

	// Ensure log messages were written to Kafka
	time.Sleep(time.Second)
}

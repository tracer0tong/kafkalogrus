package kafkalogrus

import (
	"testing"
	"time"

	"github.com/sirupsen/logrus"
)

func TestKafkaHook(t *testing.T) {
	// Create a new KafkaHook
	hook, err := NewKafkaLogrusHook(
		"kh",
		[]logrus.Level{logrus.InfoLevel, logrus.WarnLevel, logrus.ErrorLevel},
		&logrus.JSONFormatter{},
		[]string{"127.0.0.1:9092"},
		"test",
		true)

	if err != nil {
		t.Errorf("Can not create KafkaHook: %v\n", err)
	}

	// Create a new logrus.Logger
	logger := logrus.New()

	// Add hook to logger
	logger.Hooks.Add(hook)

	t.Logf("logger: %v", logger)
	t.Logf("logger.Out: %v", logger.Out)
	t.Logf("logger.Formatter: %v", logger.Formatter)
	t.Logf("logger.Hooks: %v", logger.Hooks)
	t.Logf("logger.Level: %v", logger.Level)

	// Add topics
	l := logger.WithField("topic", "nondefaulttopic")

	l.Debug("This must not be logged")

	l.Info("This is an Info msg")

	l.Warn("This is a Warn msg")

	l.Error("This is an Error msg")

	// Ensure log messages were written to Kafka
	time.Sleep(time.Second)
}

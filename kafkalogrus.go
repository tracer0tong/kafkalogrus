package kafkalogrus

import (
	"errors"
	"log"
	"os"
	"time"

	"github.com/Shopify/sarama"
	"github.com/sirupsen/logrus"
	"runtime"
)

// KafkaLogrusHook
type KafkaLogrusHook struct {
	id             string
	defaultTopic   string
	injectHostname bool
	injectCaller   bool
	hostname       string
	levels         []logrus.Level
	formatter      logrus.Formatter
	producer       sarama.AsyncProducer
}

// Create a new KafkaHook
func NewKafkaLogrusHook(id string,
	levels []logrus.Level,
	formatter logrus.Formatter,
	brokers []string,
	defaultTopic string,
	injectHostname bool,
	injectCaller bool) (*KafkaLogrusHook, error) {
	var err error
	var producer sarama.AsyncProducer
	kafkaConfig := sarama.NewConfig()
	kafkaConfig.Producer.RequiredAcks = sarama.WaitForLocal       // Only wait for the leader to ack
	kafkaConfig.Producer.Compression = sarama.CompressionSnappy   // Compress messages
	kafkaConfig.Producer.Flush.Frequency = 500 * time.Millisecond // Flush batches every 500ms

	if producer, err = sarama.NewAsyncProducer(brokers, kafkaConfig); err != nil {
		return nil, err
	}

	go func() {
		for err := range producer.Errors() {
			log.Printf("Failed to send log entry to Kafka: %v\n", err)
		}
	}()

	var hostname string
	if hostname, err = os.Hostname(); err != nil {
		hostname = "localhost"
	}

	hook := &KafkaLogrusHook{
		id,
		defaultTopic,
		injectHostname,
		injectCaller,
		hostname,
		levels,
		formatter,
		producer,
	}

	return hook, nil
}

func (hook *KafkaLogrusHook) Id() string {
	return hook.id
}

func (hook *KafkaLogrusHook) Levels() []logrus.Level {
	return hook.levels
}

type CallerDesc struct {
	filename string
	line     int
	callee   string
}

func (hook *KafkaLogrusHook) GetCallerInfo() *CallerDesc {
	pc, file, line, ok := runtime.Caller(5)
	details := runtime.FuncForPC(pc)
	if ok && details != nil {
		return &CallerDesc{file, line, details.Name()}
	}
	return nil
}

func (hook *KafkaLogrusHook) Fire(entry *logrus.Entry) error {
	var partitionKey sarama.ByteEncoder
	var b []byte
	var err error

	t, _ := entry.Data["time"].(time.Time)
	if b, err = t.MarshalBinary(); err != nil {
		return err
	}
	partitionKey = sarama.ByteEncoder(b)

	if hook.injectHostname {
		if _, ok := entry.Data["host"]; !ok {
			entry.Data["host"] = hook.hostname
		}
	}

	if hook.injectCaller {
		if _, ok := entry.Data["func"]; !ok {
			entry.Data["func"] = hook.GetCallerInfo()
		}
	}

	if b, err = hook.formatter.Format(entry); err != nil {
		return err
	}
	value := sarama.ByteEncoder(b)

	topic := hook.defaultTopic
	if tsRaw, ok := entry.Data["topic"]; ok {
		if ts, ok := tsRaw.(string); !ok {
			return errors.New("Incorrect topic filed type (should be string)")
		} else {
			topic = ts
		}
	}
	hook.producer.Input() <- &sarama.ProducerMessage{
		Key:   partitionKey,
		Topic: topic,
		Value: value,
	}
	return nil
}

# Kafka Logrus hook

Modified version of Kafka hook for Logrus.
Just need hook for my projects and I couldn't see that somebody actively supports any related projects.

A [logrus.Hook](https://godoc.org/github.com/sirupsen/logrus#Hook) which sends a single
log entry to kafka topics simultaneously.

## How to use

### Import package

```Go
import kl "github.com/tracer0tong/kafkalogrus"
```

### Create a hook (KafkaHook)

```Go
NewKafkaLogrusHook(id string, levels []logrus.Level, formatter logrus.Formatter, brokers []string, defaultTopic string, injectHostname bool) (*KafkaHook, error)
```

- id: Hook Id
- levels: [logrus.Levels](https://godoc.org/github.com/sirupsen/logrus#Level) supported by the hook
- formatter: [logrus.Formatter](https://godoc.org/github.com/sirupsen/logrus#Formatter) used by the hook
- brokers: Kafka brokers

For example:

```Go
hook, err := kl.NewKafkaHook(
        "klh",
        []logrus.Level{logrus.InfoLevel, logrus.WarnLevel, logrus.ErrorLevel},
        &logrus.JSONFormatter{},
        []string{"192.168.60.5:9092", "192.168.60.6:9092", "192.168.60.7:9092"},
        "test",
        true
    )
```

### Create a [logrus.Logger](https://godoc.org/github.com/sirupsen/logrus#Logger)

For example:

```Go
logger := logrus.New()
```

### Add hook to logger

```Go
logger.Hooks.Add(hook)
```

### Modify topic

```Go
l := logger.WithField("topic", "nondefaulttopic")
```

The field name must be ***topic***.


### Send messages to logger

For example:

```Go
l.Debug("This must not be logged")
l.Info("This is an Info msg")
l.Warn("This is a Warn msg")
l.Error("This is an Error msg")
```

#### Complete examples

[https://github.com/tracer0tong/kafkalogrus/tree/master/examples](https://github.com/tracer0tong/kafkalogrus/tree/master/examples)

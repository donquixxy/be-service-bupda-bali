package utilities

import (
	"fmt"

	"github.com/sirupsen/logrus"
	"github.com/tensuqiuwulu/be-service-bupda-bali/config"
)

func NewLogger(configurationLog config.Log) *logrus.Logger {
	logger := logrus.New()
	logger.SetFormatter(&logrus.JSONFormatter{})

	if configurationLog.Level == "trace" {
		logger.SetLevel(logrus.TraceLevel)
	} else if configurationLog.Level == "debug" {
		logger.SetLevel(logrus.DebugLevel)
	} else if configurationLog.Level == "info" {
		logger.SetLevel(logrus.InfoLevel)
	} else if configurationLog.Level == "warn" {
		logger.SetLevel(logrus.WarnLevel)
	} else if configurationLog.Level == "error" {
		logger.SetLevel(logrus.ErrorLevel)
	} else if configurationLog.Level == "fatal" {
		logger.SetLevel(logrus.FatalLevel)
	} else if configurationLog.Level == "panic" {
		logger.SetLevel(logrus.PanicLevel)
	}

	for _, level := range configurationLog.Levels {
		if level == "trace" {
			logger.AddHook(&TraceHook{})
		} else if level == "debug" {
			logger.AddHook(&DebugHook{})
		} else if level == "info" {
			logger.AddHook(&InfoHook{})
		} else if level == "warn" {
			logger.AddHook(&WarnHook{})
		} else if level == "error" {
			logger.AddHook(&ErrorHook{})
		} else if level == "fatal" {
			logger.AddHook(&FatalHook{})
		} else if level == "panic" {
			logger.AddHook(&PanicHook{})
		}
	}

	return logger
}

type TraceHook struct {
}

func (hook *TraceHook) Levels() []logrus.Level {
	return []logrus.Level{logrus.TraceLevel}
}

func (hook *TraceHook) Fire(entry *logrus.Entry) error {
	fmt.Println("message level trace:", entry.Level, entry.Message)
	return nil
}

type DebugHook struct {
}

func (hook *DebugHook) Levels() []logrus.Level {
	return []logrus.Level{logrus.DebugLevel}
}

func (hook *DebugHook) Fire(entry *logrus.Entry) error {
	fmt.Println("message level debug:", entry.Level, entry.Message)
	return nil
}

type InfoHook struct {
}

func (hook *InfoHook) Levels() []logrus.Level {
	return []logrus.Level{logrus.InfoLevel}
}

func (hook *InfoHook) Fire(entry *logrus.Entry) error {
	fmt.Println("message level info:", entry.Level, entry.Message)
	return nil
}

type WarnHook struct {
}

func (hook *WarnHook) Levels() []logrus.Level {
	return []logrus.Level{logrus.WarnLevel}
}

func (hook *WarnHook) Fire(entry *logrus.Entry) error {
	fmt.Println("message level warn:", entry.Level, entry.Message)
	return nil
}

type ErrorHook struct {
}

func (hook *ErrorHook) Levels() []logrus.Level {
	return []logrus.Level{logrus.ErrorLevel}
}

func (hook *ErrorHook) Fire(entry *logrus.Entry) error {
	fmt.Println("message level error:", entry.Level, entry.Message)
	return nil
}

type FatalHook struct {
}

func (hook *FatalHook) Levels() []logrus.Level {
	return []logrus.Level{logrus.FatalLevel}
}

func (hook *FatalHook) Fire(entry *logrus.Entry) error {
	fmt.Println("message level fatal:", entry.Level, entry.Message)
	return nil
}

type PanicHook struct {
}

func (hook *PanicHook) Levels() []logrus.Level {
	return []logrus.Level{logrus.PanicLevel}
}

func (hook *PanicHook) Fire(entry *logrus.Entry) error {
	fmt.Println("message level panic:", entry.Level, entry.Message)
	return nil
}

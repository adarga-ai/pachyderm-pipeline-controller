// Copyright 2021 Adarga Limited
// SPDX-License-Identifier: Apache-2.0

package config

import (
	"os"

	"github.com/spotahome/kooper/v2/log"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Want to use zap, but the SugaredLogger isn't quite compatible
// with kooper's Logger object so we need to create our own wrapper

type ZapLogger struct {
	*zap.SugaredLogger
}

func (s ZapLogger) Warningf(template string, args ...interface{}) {
	s.Warnf(template, args...)
}

func (s ZapLogger) WithKV(kv log.KV) log.Logger {
	ss := s.SugaredLogger
	for k, v := range kv {
		ss = ss.With(k, v)
	}
	return ZapLogger{ss}
}

// mklogger initialises a structured logger; it returns the core zap
// object and a version wrapped to support the kooper Logger interface
func GetLogger() (*zap.Logger, *ZapLogger) {
	// Initialize logger. JSON formatted error messages, outputted to stdout
	logcfg := zap.NewProductionEncoderConfig()
	logcfg.TimeKey = "time"
	logcfg.EncodeTime = zapcore.ISO8601TimeEncoder
	logger := zap.New(zapcore.NewCore(zapcore.NewJSONEncoder(logcfg), os.Stdout, zap.InfoLevel))
	return logger, &ZapLogger{logger.Sugar()}
}

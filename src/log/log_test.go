// Copyright 2000-2018 by ChinanetCenter Corporation. All rights reserved.

package log

import (
	"os"
	"testing"
)

func TestNewLogger(t *testing.T) {
	logger := New(os.Stdout)
	logger.SetLevel(LvlDebug)
	logger.Debug("debug")
	logger.Info("info")
	logger.Error("error")
	logger.SetLevel(LvlError)
	logger.Debug("debug")
	logger.Info("info")
	logger.Error("error")
}

func TestDefaultLogger(t *testing.T) {
	SetLevel(LvlDebug)
	Debug("debug")
	Info("info")
	Error("error")
	SetLevel(LvlError)
	Debug("debug")
	Info("info")
	Error("error")
}

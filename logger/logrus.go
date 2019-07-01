// Copyright 2019 Drone.IO Inc. All rights reserved.
// Use of this source code is governed by the Parity Public License
// that can be found in the LICENSE file.

package logger

import "github.com/sirupsen/logrus"

// Logrus returns a Logger that wraps a logrus.Logger.
func Logrus(logrus *logrus.Logger) Logger {
	return &wrapLogrus{logrus}
}

type wrapLogrus struct {
	*logrus.Logger
}

func (w *wrapLogrus) WithError(err error) Logger {
	return &wrapLogrus{w.Logger.WithError(err).Logger}
}

func (w *wrapLogrus) WithField(key string, value interface{}) Logger {
	return &wrapLogrus{w.Logger.WithField(key, value).Logger}
}

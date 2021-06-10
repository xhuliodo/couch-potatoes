package logger

import (
	"time"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

type ErrorLogger struct {
	Logger *logrus.Logger
}

func (l *ErrorLogger) Log(err error) {
	entry := &ErrorLoggerEntry{Logger: logrus.NewEntry(l.Logger)}
	logFields := logrus.Fields{}

	logFields["ts"] = time.Now().UTC().Format(time.RFC1123)
	
	stack := errors.WithStack(err)
	logFields["stack"] = stack

	cause := errors.Cause(err)

	entry.Logger.WithFields(logFields).Error(cause)
}

type ErrorLoggerEntry struct {
	Logger logrus.FieldLogger
}

// func (l *ErrorLoggerEntry) Panic(v interface{}, stack []byte) {
// 	l.Logger = l.Logger.WithFields(logrus.Fields{
// 		"stack": string(stack),
// 		"panic": fmt.Sprintf("%+v", v),
// 	})
// }

// func GetLogEntry(r *http.Request) logrus.FieldLogger {
// 	entry := middleware.GetLogEntry(r).(*ErrorLoggerEntry)
// 	return entry.Logger
// }

// func LogEntrySetField(r *http.Request, key string, value interface{}) {
// 	if entry, ok := r.Context().Value(middleware.LogEntryCtxKey).(*ErrorLoggerEntry); ok {
// 		entry.Logger = entry.Logger.WithField(key, value)
// 	}
// }

// func LogEntrySetFields(r *http.Request, fields map[string]interface{}) {
// 	if entry, ok := r.Context().Value(middleware.LogEntryCtxKey).(*ErrorLoggerEntry); ok {
// 		entry.Logger = entry.Logger.WithFields(fields)
// 	}
// }

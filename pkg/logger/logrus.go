package logger

import (
	"github.com/sirupsen/logrus"
)

type logrusLogger struct {
	log      *logrus.Logger
	name     string
	hostName string
}

// Option to config log
type Option func(*logrusLogger)

// WithServiceName option for service name attribute
func WithServiceName(serviceName string) func(*logrusLogger) {
	return func(lg *logrusLogger) {
		lg.name = serviceName
	}
}

// WithHostName options for host name attribute
func WithHostName(hostName string) func(*logrusLogger) {
	return func(lg *logrusLogger) {
		lg.hostName = hostName
	}
}

// Operation attribute log operation name
type Operation string

// Op generate operation with name
func Op(val string) Operation {
	return Operation(val)
}

// Event attribute help log can make different when print log
type Event string

// E generate event with name
func E(op string) Event {
	return Event(op)
}

// NewJSONLogger initializes the json standard logger by logrus lib
func NewJSONLogger(opts ...Option) Log {
	var baseLogger = logrus.New()
	var logrusLogger = &logrusLogger{
		log: baseLogger,
	}
	for idx := range opts {
		opt := opts[idx]
		opt(logrusLogger)
	}

	logrusLogger.log.Formatter = &logrus.JSONFormatter{}
	return logrusLogger
}

func (l *logrusLogger) prepareFields() *logrus.Entry {
	return l.log.WithFields(
		logrus.Fields{
			"name":     l.name,
			"hostname": l.hostName,
		},
	)
}
func (l *logrusLogger) Log(vals ...interface{}) error {
	fields, vals := detachFields(vals)
	l.prepareFields().WithFields(fields).Info(vals...)
	return nil
}

func (l *logrusLogger) Debug(vals ...interface{}) error {
	fields, vals := detachFields(vals)
	l.prepareFields().WithFields(fields).Debug(vals...)
	return nil
}

func (l *logrusLogger) Info(vals ...interface{}) error {
	fields, vals := detachFields(vals)
	l.prepareFields().WithFields(fields).Info(vals...)
	return nil
}

func (l *logrusLogger) Warn(vals ...interface{}) error {
	fields, vals := detachFields(vals)
	l.prepareFields().WithFields(fields).Warning(vals...)
	return nil
}

func (l *logrusLogger) Error(vals ...interface{}) error {
	fields, vals := detachFields(vals)
	l.prepareFields().WithFields(fields).Error(vals...)
	return nil
}

func (l *logrusLogger) Errorf(str string, vals ...interface{}) error {
	l.prepareFields().Errorf(str, vals...)
	return nil
}

func detachFields(vals []interface{}) (logrus.Fields, []interface{}) {
	fields := logrus.Fields{}
	others := []interface{}{}
	for idx := range vals {
		arg := vals[idx]
		switch arg := arg.(type) {
		case Event:
			fields["event"] = arg
		case Operation:
			fields["op"] = arg
		default:
			others = append(others, arg)
		}
	}
	return fields, others
}

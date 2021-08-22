package logger

import (
	"os"
	"sync"

	"github.com/natefinch/lumberjack"
	"github.com/sirupsen/logrus"
	"motherbear/backend/constants"
)

// Event stores messages to log later, from our standard interface
type Event struct {
	id      int
	message string
}

// StandardLogger enforces specific log message formats
type StandardLogger struct {
	*logrus.Logger
}

var instance *StandardLogger
var instanceSymptom *StandardLogger
var once sync.Once

// We used the Singletone pattern.
// NewLogger initializes the standard logger
func init() {
	once.Do(func() {
		var baseLogger = logrus.New()
		var symptomLogger = logrus.New()
		instance = &StandardLogger{baseLogger}
		instance.Formatter = &logrus.JSONFormatter{}

		instanceSymptom = &StandardLogger{symptomLogger}
		instanceSymptom.Formatter = &logrus.JSONFormatter{}

		instance.Level = logrus.InfoLevel
		instanceSymptom.Level = logrus.InfoLevel

		if _, err := os.Stat("./log"); os.IsNotExist(err) {
			os.MkdirAll("./log", os.ModePerm)
		}
		lumberjackLogrotate := &lumberjack.Logger{
			Filename:   "log/isaac_server.log",
			MaxSize:    100,  // Max megabytes before log is rotated
			MaxBackups: 30, // Max number of old log files to keep
			MaxAge:     30, // Max number of days to retain log files
			Compress:   true,
		}
		instance.SetOutput(lumberjackLogrotate)

		lumberjackSymptomrotate := &lumberjack.Logger{
			Filename:   "log/peer_symptom.log",
			MaxSize:    100,  // Max megabytes before log is rotated
			MaxBackups: 30, // Max number of old log files to keep
			MaxAge:     30, // Max number of days to retain log files
			Compress:   true,
		}
		instanceSymptom.SetOutput(lumberjackSymptomrotate)
	})
}

// Logger returns the instance of StandardLogger.
func Logger() *StandardLogger {
	return instance
}

// Declare variables to store log messages as new Events
var (
	defaultMessage         = Event{1, "%s"}
	invalidArgMessage      = Event{2, "Invalid arg: %s"}
	invalidArgValueMessage = Event{3, "Invalid value for argument: %s: %v"}
	missingArgMessage      = Event{4, "Missing arg: %s"}
)

func SetLogLevel(logLevel int) {
	if logLevel == 1 {
		instance.Level = logrus.InfoLevel
	} else {
		instance.Level = logrus.DebugLevel
	}
}

// InvalidArg is a standard error message
func InvalidArg(argumentName string, argumentIDOptional ...string) {
	argumentID := constants.LoggerServerUser
	if len(argumentIDOptional) > 0 {
		argumentID = argumentIDOptional[0]
	}
	instance.WithFields(logrus.Fields{
		"userid": argumentID,
	}).Errorf(invalidArgMessage.message, argumentName)
}

// InvalidArgValue is a standard error message
func InvalidArgValue(argumentName string, argumentValue string, argumentIDOptional ...string) {
	argumentID := constants.LoggerServerUser
	if len(argumentIDOptional) > 0 {
		argumentID = argumentIDOptional[0]
	}
	instance.WithFields(logrus.Fields{
		"userid": argumentID,
	}).Errorf(invalidArgValueMessage.message, argumentName, argumentValue)
}

// MissingArg is a standard error message
func MissingArg(argumentName string, argumentIDOptional ...string) {
	argumentID := constants.LoggerServerUser
	if len(argumentIDOptional) > 0 {
		argumentID = argumentIDOptional[0]
	}
	instance.WithFields(logrus.Fields{
		"userid": argumentID,
	}).Errorf(missingArgMessage.message, argumentName)
}

// Info is a standard info message
func Info(argumentData string, argumentIDOptional ...string) {
	argumentID := constants.LoggerServerUser
	if len(argumentIDOptional) > 0 {
		argumentID = argumentIDOptional[0]
	}
	instance.WithFields(logrus.Fields{
		"userid": argumentID,
	}).Infof(defaultMessage.message, argumentData)
}

// Info is a standard info message
func Infof(format string, args ...interface{}) {
	argumentID := constants.LoggerServerUser
	instance.WithFields(logrus.Fields{
		"userid": argumentID,
	}).Infof(format, args...)
}

// Warn is a standard warn message
func Warn(argumentData string, argumentIDOptional ...string) {
	argumentID := constants.LoggerServerUser
	if len(argumentIDOptional) > 0 {
		argumentID = argumentIDOptional[0]
	}
	instance.WithFields(logrus.Fields{
		"userid": argumentID,
	}).Warnf(defaultMessage.message, argumentData)
}

// Error is a standard error message
func Error(argumentData string, argumentIDOptional ...string) {
	argumentID := constants.LoggerServerUser
	if len(argumentIDOptional) > 0 {
		argumentID = argumentIDOptional[0]
	}
	instance.WithFields(logrus.Fields{
		"userid": argumentID,
	}).Errorf(defaultMessage.message, argumentData)
}

// Fatal is a standard fatal message
func Fatal(argumentData string, argumentIDOptional ...string) {
	argumentID := constants.LoggerServerUser
	if len(argumentIDOptional) > 0 {
		argumentID = argumentIDOptional[0]
	}
	instance.WithFields(logrus.Fields{
		"userid": argumentID,
	}).Fatalf(defaultMessage.message, argumentData)
}

// Errorf is a standard error message
func Errorf(format string, args ...interface{}) {
	argumentID := constants.LoggerServerUser
	instance.WithFields(logrus.Fields{
		"userid": argumentID,
	}).Errorf(format, args...)
}

// Panicln is a standard panic message
func Panicln(args ...interface{}) {
	argumentID := constants.LoggerServerUser
	instance.WithFields(logrus.Fields{
		"userid": argumentID,
	}).Panicln(args...)
}

// Fatal is a standard fatal message
func Fatalln(args ...interface{}) {
	argumentID := constants.LoggerServerUser
	instance.WithFields(logrus.Fields{
		"userid": argumentID,
	}).Fatalln(args...)
}

// Debug is a standard info message
func Debug(argumentData string, argumentIDOptional ...string) {
	argumentID := constants.LoggerServerUser
	if len(argumentIDOptional) > 0 {
		argumentID = argumentIDOptional[0]
	}
	instance.WithFields(logrus.Fields{
		"userid": argumentID,
	}).Debug(argumentData)
}

// Info is a standard info message
func Debugf(format string, args ...interface{}) {
	argumentID := constants.LoggerServerUser
	instance.WithFields(logrus.Fields{
		"userid": argumentID,
	}).Debugf(format, args...)
}

// peer symptom is a standard info message
func Symptomf(target string, symptom string, format string, args ...interface{}) {
	argumentID := target
	instanceSymptom.WithFields(logrus.Fields{
		"channel": argumentID,
		"symptom": symptom,
	}).Errorf(format, args...)
}
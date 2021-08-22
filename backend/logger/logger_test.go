package logger

import (
	"github.com/sirupsen/logrus"
	"os"
	"testing"

	"gopkg.in/go-playground/assert.v1"
)

func TestLogger(t *testing.T) {
	fpLog, err := os.OpenFile("log/isaac_server.log", os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0666)
	if err != nil {
		panic(err)
	}
	Logger().SetOutput(fpLog)
	InvalidArgValue("client", "nil", "bkmoon")
	Errorf("%s : %s", "12345", "sdfsdg")

	fi, err := fpLog.Stat()
	if err != nil {
		t.Error("Not created log file")
	}

	assert.NotEqual(t, 0, fi.Size())
}

func TestLogLevel(t *testing.T) {
	fpLog, err := os.OpenFile("log/isaac_server.log", os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0666)
	if err != nil {
		panic(err)
	}
	Logger().Level = logrus.InfoLevel
	Logger().SetOutput(fpLog)

	Debug("Debugger data")
	Debugf("%s : %s", "12345", "sdfsdg")

	fi, err := fpLog.Stat()
	if err != nil {
		t.Error("Not created log file")
	}

	assert.NotEqual(t, 0, fi.Size())
}

func TestSymptomLevel(t *testing.T) {
	fpLog, err := os.OpenFile("log/peer_symptom.log", os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0666)
	if err != nil {
		panic(err)
	}
	Logger().Level = logrus.ErrorLevel
	Logger().SetOutput(fpLog)

	Symptomf("channel", "Slow Response", "[%s]response time slowly [%f] sec", "node1", 10.0)

	fi, err := fpLog.Stat()
	if err != nil {
		t.Error("Not created log file")
	}

	assert.NotEqual(t, 0, fi.Size())
}
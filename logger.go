package linkcxo

import (
	"os"

	"github.com/sirupsen/logrus"
)

func InitializeLogger(logLevel string) {
	logrus.SetFormatter(&logrus.TextFormatter{})
	logrus.SetReportCaller(true)

	ll, err := logrus.ParseLevel(logLevel)
	if err != nil {
		ll = logrus.InfoLevel
	}
	logrus.SetLevel(ll)
	logrus.AddHook(newExtraFieldHook(os.Getenv("APP_VERSION")))
}
func newExtraFieldHook(appVersion string) *ExtraFieldHook {
	return &ExtraFieldHook{
		AppVersion: appVersion,
		Pid:        os.Getpid(),
	}
}
func (h *ExtraFieldHook) Levels() []logrus.Level {
	return logrus.AllLevels
}
func (h *ExtraFieldHook) Fire(entry *logrus.Entry) error {
	entry.Data["appVersion"] = h.AppVersion
	entry.Data["pid"] = h.Pid
	return nil
}

type ExtraFieldHook struct {
	AppVersion string
	Pid        int
}
